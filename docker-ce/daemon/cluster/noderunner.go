/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package cluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/executor/container"
	swarmapi "github.com/docker/swarmkit/api"
	swarmnode "github.com/docker/swarmkit/node"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// nodeRunner implements a manager for continuously running swarmkit node, restarting them with backoff delays if needed.
type nodeRunner struct {
	nodeState
	mu             sync.RWMutex
	done           chan struct{} // closed when swarmNode exits
	ready          chan struct{} // closed when swarmNode becomes active
	reconnectDelay time.Duration
	config         nodeStartConfig

	repeatedRun     bool
	cancelReconnect func()
	stopping        bool
	cluster         *Cluster // only for accessing config helpers, never call any methods. TODO: change to config struct
}

// nodeStartConfig holds configuration needed to start a new node. Exported
// fields of this structure are saved to disk in json. Unexported fields
// contain data that shouldn't be persisted between daemon reloads.
type nodeStartConfig struct {
	// LocalAddr is this machine's local IP or hostname, if specified.
	LocalAddr string
	// RemoteAddr is the address that was given to "swarm join". It is used
	// to find LocalAddr if necessary.
	RemoteAddr string
	// ListenAddr is the address we bind to, including a port.
	ListenAddr string
	// AdvertiseAddr is the address other nodes should connect to,
	// including a port.
	AdvertiseAddr   string
	joinAddr        string
	forceNewCluster bool
	joinToken       string
	lockKey         []byte
	autolock        bool
	availability    types.NodeAvailability
}

func (n *nodeRunner) Ready() chan error {
	c := make(chan error, 1)
	n.mu.RLock()
	ready, done := n.ready, n.done
	n.mu.RUnlock()
	go func() {
		select {
		case <-ready:
		case <-done:
		}
		select {
		case <-ready:
		default:
			n.mu.RLock()
			c <- n.err
			n.mu.RUnlock()
		}
		close(c)
	}()
	return c
}

func (n *nodeRunner) Start(conf nodeStartConfig) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.reconnectDelay = initialReconnectDelay

	return n.start(conf)
}

func (n *nodeRunner) start(conf nodeStartConfig) error {
	var control string
	if runtime.GOOS == "windows" {
		control = `\\.\pipe\` + controlSocket
	} else {
		control = filepath.Join(n.cluster.runtimeRoot, controlSocket)
	}

	// Hostname is not set here. Instead, it is obtained from
	// the node description that is reported periodically
	swarmnodeConfig := swarmnode.Config{
		ForceNewCluster:    conf.forceNewCluster,
		ListenControlAPI:   control,
		ListenRemoteAPI:    conf.ListenAddr,
		AdvertiseRemoteAPI: conf.AdvertiseAddr,
		JoinAddr:           conf.joinAddr,
		StateDir:           n.cluster.root,
		JoinToken:          conf.joinToken,
		Executor:           container.NewExecutor(n.cluster.config.Backend),
		HeartbeatTick:      1,
		ElectionTick:       3,
		UnlockKey:          conf.lockKey,
		AutoLockManagers:   conf.autolock,
		PluginGetter:       n.cluster.config.Backend.PluginGetter(),
	}
	if conf.availability != "" {
		avail, ok := swarmapi.NodeSpec_Availability_value[strings.ToUpper(string(conf.availability))]
		if !ok {
			return fmt.Errorf("invalid Availability: %q", conf.availability)
		}
		swarmnodeConfig.Availability = swarmapi.NodeSpec_Availability(avail)
	}
	node, err := swarmnode.New(&swarmnodeConfig)
	if err != nil {
		return err
	}
	if err := node.Start(context.Background()); err != nil {
		return err
	}

	n.done = make(chan struct{})
	n.ready = make(chan struct{})
	n.swarmNode = node
	n.config = conf
	savePersistentState(n.cluster.root, conf)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		n.handleNodeExit(node)
		cancel()
	}()

	go n.handleReadyEvent(ctx, node, n.ready)
	go n.handleControlSocketChange(ctx, node)

	return nil
}

func (n *nodeRunner) handleControlSocketChange(ctx context.Context, node *swarmnode.Node) {
	for conn := range node.ListenControlSocket(ctx) {
		n.mu.Lock()
		if n.grpcConn != conn {
			if conn == nil {
				n.controlClient = nil
				n.logsClient = nil
			} else {
				n.controlClient = swarmapi.NewControlClient(conn)
				n.logsClient = swarmapi.NewLogsClient(conn)
			}
		}
		n.grpcConn = conn
		n.mu.Unlock()
		n.cluster.configEvent <- struct{}{}
	}
}

func (n *nodeRunner) handleReadyEvent(ctx context.Context, node *swarmnode.Node, ready chan struct{}) {
	select {
	case <-node.Ready():
		n.mu.Lock()
		n.err = nil
		n.mu.Unlock()
		close(ready)
	case <-ctx.Done():
	}
	n.cluster.configEvent <- struct{}{}
}

func (n *nodeRunner) handleNodeExit(node *swarmnode.Node) {
	err := detectLockedError(node.Err(context.Background()))
	if err != nil {
		logrus.Errorf("cluster exited with error: %v", err)
	}
	n.mu.Lock()
	n.swarmNode = nil
	n.err = err
	close(n.done)
	select {
	case <-n.ready:
		n.enableReconnectWatcher()
	default:
		if n.repeatedRun {
			n.enableReconnectWatcher()
		}
	}
	n.repeatedRun = true
	n.mu.Unlock()
}

// Stop stops the current swarm node if it is running.
func (n *nodeRunner) Stop() error {
	n.mu.Lock()
	if n.cancelReconnect != nil { // between restarts
		n.cancelReconnect()
		n.cancelReconnect = nil
	}
	if n.swarmNode == nil {
		n.mu.Unlock()
		return nil
	}
	n.stopping = true
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := n.swarmNode.Stop(ctx); err != nil && !strings.Contains(err.Error(), "context canceled") {
		n.mu.Unlock()
		return err
	}
	n.mu.Unlock()
	<-n.done
	return nil
}

func (n *nodeRunner) State() nodeState {
	if n == nil {
		return nodeState{status: types.LocalNodeStateInactive}
	}
	n.mu.RLock()
	defer n.mu.RUnlock()

	ns := n.nodeState

	if ns.err != nil || n.cancelReconnect != nil {
		if errors.Cause(ns.err) == errSwarmLocked {
			ns.status = types.LocalNodeStateLocked
		} else {
			ns.status = types.LocalNodeStateError
		}
	} else {
		select {
		case <-n.ready:
			ns.status = types.LocalNodeStateActive
		default:
			ns.status = types.LocalNodeStatePending
		}
	}

	return ns
}

func (n *nodeRunner) enableReconnectWatcher() {
	if n.stopping {
		return
	}
	n.reconnectDelay *= 2
	if n.reconnectDelay > maxReconnectDelay {
		n.reconnectDelay = maxReconnectDelay
	}
	logrus.Warnf("Restarting swarm in %.2f seconds", n.reconnectDelay.Seconds())
	delayCtx, cancel := context.WithTimeout(context.Background(), n.reconnectDelay)
	n.cancelReconnect = cancel

	config := n.config
	go func() {
		<-delayCtx.Done()
		if delayCtx.Err() != context.DeadlineExceeded {
			return
		}
		n.mu.Lock()
		defer n.mu.Unlock()
		if n.stopping {
			return
		}
		config.RemoteAddr = n.cluster.getRemoteAddress()
		config.joinAddr = config.RemoteAddr
		if err := n.start(config); err != nil {
			n.err = err
		}
	}()
}

// nodeState represents information about the current state of the cluster and
// provides access to the grpc clients.
type nodeState struct {
	swarmNode       *swarmnode.Node
	grpcConn        *grpc.ClientConn
	controlClient   swarmapi.ControlClient
	logsClient      swarmapi.LogsClient
	status          types.LocalNodeState
	actualLocalAddr string
	err             error
}

// IsActiveManager returns true if node is a manager ready to accept control requests. It is safe to access the client properties if this returns true.
func (ns nodeState) IsActiveManager() bool {
	return ns.controlClient != nil
}

// IsManager returns true if node is a manager.
func (ns nodeState) IsManager() bool {
	return ns.swarmNode != nil && ns.swarmNode.Manager() != nil
}

// NodeID returns node's ID or empty string if node is inactive.
func (ns nodeState) NodeID() string {
	if ns.swarmNode != nil {
		return ns.swarmNode.NodeID()
	}
	return ""
}
