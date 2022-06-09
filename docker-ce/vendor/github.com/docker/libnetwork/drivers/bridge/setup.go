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
package bridge

type setupStep func(*networkConfiguration, *bridgeInterface) error

type bridgeSetup struct {
	config *networkConfiguration
	bridge *bridgeInterface
	steps  []setupStep
}

func newBridgeSetup(c *networkConfiguration, i *bridgeInterface) *bridgeSetup {
	return &bridgeSetup{config: c, bridge: i}
}

func (b *bridgeSetup) apply() error {
	for _, fn := range b.steps {
		if err := fn(b.config, b.bridge); err != nil {
			return err
		}
	}
	return nil
}

func (b *bridgeSetup) queueStep(step setupStep) {
	b.steps = append(b.steps, step)
}
