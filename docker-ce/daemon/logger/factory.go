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
package logger

import (
	"fmt"
	"sync"

	containertypes "github.com/docker/docker/api/types/container"
	units "github.com/docker/go-units"
	"github.com/pkg/errors"
)

// Creator builds a logging driver instance with given context.
type Creator func(Info) (Logger, error)

// LogOptValidator checks the options specific to the underlying
// logging implementation.
type LogOptValidator func(cfg map[string]string) error

type logdriverFactory struct {
	registry     map[string]Creator
	optValidator map[string]LogOptValidator
	m            sync.Mutex
}

func (lf *logdriverFactory) register(name string, c Creator) error {
	if lf.driverRegistered(name) {
		return fmt.Errorf("logger: log driver named '%s' is already registered", name)
	}

	lf.m.Lock()
	lf.registry[name] = c
	lf.m.Unlock()
	return nil
}

func (lf *logdriverFactory) driverRegistered(name string) bool {
	lf.m.Lock()
	_, ok := lf.registry[name]
	lf.m.Unlock()
	return ok
}

func (lf *logdriverFactory) registerLogOptValidator(name string, l LogOptValidator) error {
	lf.m.Lock()
	defer lf.m.Unlock()

	if _, ok := lf.optValidator[name]; ok {
		return fmt.Errorf("logger: log validator named '%s' is already registered", name)
	}
	lf.optValidator[name] = l
	return nil
}

func (lf *logdriverFactory) get(name string) (Creator, error) {
	lf.m.Lock()
	defer lf.m.Unlock()

	c, ok := lf.registry[name]
	if !ok {
		return c, fmt.Errorf("logger: no log driver named '%s' is registered", name)
	}
	return c, nil
}

func (lf *logdriverFactory) getLogOptValidator(name string) LogOptValidator {
	lf.m.Lock()
	defer lf.m.Unlock()

	c, _ := lf.optValidator[name]
	return c
}

var factory = &logdriverFactory{registry: make(map[string]Creator), optValidator: make(map[string]LogOptValidator)} // global factory instance

// RegisterLogDriver registers the given logging driver builder with given logging
// driver name.
func RegisterLogDriver(name string, c Creator) error {
	return factory.register(name, c)
}

// RegisterLogOptValidator registers the logging option validator with
// the given logging driver name.
func RegisterLogOptValidator(name string, l LogOptValidator) error {
	return factory.registerLogOptValidator(name, l)
}

// GetLogDriver provides the logging driver builder for a logging driver name.
func GetLogDriver(name string) (Creator, error) {
	return factory.get(name)
}

var builtInLogOpts = map[string]bool{
	"mode":            true,
	"max-buffer-size": true,
}

// ValidateLogOpts checks the options for the given log driver. The
// options supported are specific to the LogDriver implementation.
func ValidateLogOpts(name string, cfg map[string]string) error {
	if name == "none" {
		return nil
	}

	switch containertypes.LogMode(cfg["mode"]) {
	case containertypes.LogModeBlocking, containertypes.LogModeNonBlock, containertypes.LogModeUnset:
	default:
		return fmt.Errorf("logger: logging mode not supported: %s", cfg["mode"])
	}

	if s, ok := cfg["max-buffer-size"]; ok {
		if containertypes.LogMode(cfg["mode"]) != containertypes.LogModeNonBlock {
			return fmt.Errorf("logger: max-buffer-size option is only supported with 'mode=%s'", containertypes.LogModeNonBlock)
		}
		if _, err := units.RAMInBytes(s); err != nil {
			return errors.Wrap(err, "error parsing option max-buffer-size")
		}
	}

	if !factory.driverRegistered(name) {
		return fmt.Errorf("logger: no log driver named '%s' is registered", name)
	}

	filteredOpts := make(map[string]string, len(builtInLogOpts))
	for k, v := range cfg {
		if !builtInLogOpts[k] {
			filteredOpts[k] = v
		}
	}

	validator := factory.getLogOptValidator(name)
	if validator != nil {
		return validator(filteredOpts)
	}
	return nil
}
