//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cbft

import (
	"fmt"
	"io"
	"testing"
)

type TestDest struct{}

func (s *TestDest) Close() error {
	return nil
}

func (s *TestDest) OnDataUpdate(partition string,
	key []byte, seq uint64, val []byte) error {
	return nil
}

func (s *TestDest) OnDataDelete(partition string,
	key []byte, seq uint64) error {
	return nil
}

func (s *TestDest) OnSnapshotStart(partition string,
	snapStart, snapEnd uint64) error {
	return nil
}

func (s *TestDest) SetOpaque(partition string,
	value []byte) error {
	return nil
}

func (s *TestDest) GetOpaque(partition string) (
	value []byte, lastSeq uint64, err error) {
	return nil, 0, nil
}

func (s *TestDest) Rollback(partition string,
	rollbackSeq uint64) error {
	return nil
}

func (s *TestDest) ConsistencyWait(partition string,
	consistencyLevel string,
	consistencySeq uint64,
	cancelCh <-chan bool) error {
	return nil
}

func (t *TestDest) Count(pindex *PIndex,
	cancelCh <-chan bool) (uint64, error) {
	return 0, nil
}

func (t *TestDest) Query(pindex *PIndex, req []byte, res io.Writer,
	cancelCh <-chan bool) error {
	return nil
}

func (t *TestDest) Stats(w io.Writer) error {
	return nil
}

func TestBasicPartitionFunc(t *testing.T) {
	dest := &TestDest{}
	dest2 := &TestDest{}
	s, err := BasicPartitionFunc("", nil, map[string]Dest{"": dest})
	if err != nil || s != dest {
		t.Errorf("expected BasicPartitionFunc to work")
	}
	s, err = BasicPartitionFunc("foo", nil, map[string]Dest{"": dest})
	if err != nil || s != dest {
		t.Errorf("expected BasicPartitionFunc to hit the catch-all dest")
	}
	s, err = BasicPartitionFunc("", nil, map[string]Dest{"foo": dest})
	if err == nil || s == dest {
		t.Errorf("expected BasicPartitionFunc to not work")
	}
	s, err = BasicPartitionFunc("foo", nil, map[string]Dest{"foo": dest})
	if err != nil || s != dest {
		t.Errorf("expected BasicPartitionFunc to work on partition hit")
	}
	s, err = BasicPartitionFunc("foo", nil, map[string]Dest{"foo": dest, "": dest2})
	if err != nil || s != dest {
		t.Errorf("expected BasicPartitionFunc to work on partition hit")
	}
}

type ErrorOnlyDestProvider struct{}

func (dp *ErrorOnlyDestProvider) Dest(partition string) (Dest, error) {
	return nil, fmt.Errorf("always error for testing")
}

func (dp *ErrorOnlyDestProvider) Count(pindex *PIndex,
	cancelCh <-chan bool) (uint64, error) {
	return 0, fmt.Errorf("always error for testing")
}

func (dp *ErrorOnlyDestProvider) Query(pindex *PIndex, req []byte, res io.Writer,
	cancelCh <-chan bool) error {
	return fmt.Errorf("always error for testing")
}

func (dp *ErrorOnlyDestProvider) Stats(io.Writer) error {
	return fmt.Errorf("always error for testing")
}

func (dp *ErrorOnlyDestProvider) Close() error {
	return fmt.Errorf("always error for testing")
}

func TestErrorOnlyDestProviderWithDestForwarder(t *testing.T) {
	df := &DestForwarder{&ErrorOnlyDestProvider{}}
	if df.OnDataUpdate("", nil, 0, nil) == nil {
		t.Errorf("expected err")
	}
	if df.OnDataDelete("", nil, 0) == nil {
		t.Errorf("expected err")
	}
	if df.OnSnapshotStart("", 0, 0) == nil {
		t.Errorf("expected err")
	}
	if df.SetOpaque("", nil) == nil {
		t.Errorf("expected err")
	}
	value, lastSeq, err := df.GetOpaque("")
	if err == nil || value != nil || lastSeq != 0 {
		t.Errorf("expected err")
	}
	if df.Rollback("", 0) == nil {
		t.Errorf("expected err")
	}
	if df.ConsistencyWait("", "", 0, nil) == nil {
		t.Errorf("expected err")
	}
}
