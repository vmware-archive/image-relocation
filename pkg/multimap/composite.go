/*
 * Copyright (c) 2019-Present Pivotal Software, Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package multimap

import "fmt"

// Composite is an interface to a consistent collection of mappings from string to string. The mappings are consistent
// in the sense that no two mappings map a given value to distinct results.
type Composite interface {
	// Add updates the composite mapping by adding a mapping with the given name. If a mapping with the given name
	// already exists, it is replaced providing it is consistent with all other mappings. If the named mapping is not
	// consistent with all other mappings, the mapping is removed from the composite mapping and error is returned.
	Add(name string, mapping map[string]string) error

	// Delete updates the composite mapping by removing a mapping with the given name. If there is no mapping with the
	// given name, an error is returned.
	Delete(name string) error

	// Map applies the composite mapping to the given value and returns the mapped value. If the given value is not
	// in the domain of the composite mapping, the given value is returned. In other words, the default mapping is
	// the identity function.
	Map(string) string
}

type errCh chan error

type addOp struct {
	name    string
	mapping map[string]string
	errCh   errCh
}

type deleteOp struct {
	name  string
	errCh errCh
}

type mapOp struct {
	value    string
	resultCh chan string
}

type composite struct {
	// a consistent collection of mappings
	mappings map[string]map[string]string

	// the composition of all the mappings
	composite map[string]string

	addCh    chan *addOp
	deleteCh chan *deleteOp
	mapCh    chan *mapOp
	stopCh   <-chan struct{}
}

func New(stopCh <-chan struct{}) Composite {
	c := &composite{
		mappings:  make(map[string]map[string]string),
		composite: make(map[string]string),

		addCh:    make(chan *addOp),
		deleteCh: make(chan *deleteOp),
		mapCh:    make(chan *mapOp),
		stopCh:   stopCh,
	}

	go c.monitor()

	return c
}

func (c *composite) Add(name string, mapping map[string]string) error {
	errCh := make(chan error)
	c.addCh <- &addOp{
		name:    name,
		mapping: mapping,
		errCh:   errCh,
	}
	return <-errCh
}

func (c *composite) Delete(name string) error {
	errCh := make(chan error)
	c.deleteCh <- &deleteOp{
		name:  name,
		errCh: errCh,
	}
	return <-errCh
}

func (c *composite) Map(value string) string {
	resultCh := make(chan string)
	c.mapCh <- &mapOp{
		value:    value,
		resultCh: resultCh,
	}
	return <-resultCh
}

func (c *composite) monitor() {
	for {
		select {
		case addOp := <-c.addCh:
			addOp.errCh <- c.add(addOp.name, addOp.mapping)

		case deleteOp := <-c.deleteCh:
			deleteOp.errCh <- c.delete(deleteOp.name)

		case mapOp := <-c.mapCh:
			mapOp.resultCh <- c.doMap(mapOp.value)

		case <-c.stopCh:
			close(c.addCh)
			close(c.deleteCh)
			close(c.mapCh)
			return
		}
	}
}

func (c *composite) add(name string, mapping map[string]string) error {
	_ = c.delete(name) // name may not be present, so ignore any error
	if err := c.checkConsistency(mapping); err != nil {
		return err
	}

	// save a copy of mapping
	c.mappings[name] = make(map[string]string, len(mapping))
	for k, v := range mapping {
		c.mappings[name][k] = v
	}

	c.merge()
	return nil
}

func (c *composite) delete(name string) error {
	if _, ok := c.mappings[name]; !ok {
		return fmt.Errorf("mapping not found: %s", name)
	}
	delete(c.mappings, name)
	c.merge()
	return nil
}

func (c *composite) doMap(value string) string {
	if result, ok := c.composite[value]; ok {
		return result
	}
	return value
}

func (c *composite) merge() {
	c.composite = make(map[string]string)
	for _, m := range c.mappings {
		for k, v := range m {
			c.composite[k] = v
		}
	}
}

func (c *composite) checkConsistency(mapping map[string]string) error {
	for k, v := range mapping {
		if w, ok := c.composite[k]; ok && v != w {
			for n, m := range c.mappings {
				if w, ok := m[k]; ok && v != w {
					return fmt.Errorf("inconsistent mapping: %s maps %q to %q but given mapping maps %q to %q", n, k, w, k, v)
				}
			}
		}
	}
	return nil
}
