// This file was generated by counterfeiter
package enforcerfakes

import (
	"sync"

	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/enforcer"
)

type FakeEnforcer struct {
	EnforceOnceStub        func() error
	enforceOnceMutex       sync.RWMutex
	enforceOnceArgsForCall []struct{}
	enforceOnceReturns     struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEnforcer) EnforceOnce() error {
	fake.enforceOnceMutex.Lock()
	fake.enforceOnceArgsForCall = append(fake.enforceOnceArgsForCall, struct{}{})
	fake.recordInvocation("EnforceOnce", []interface{}{})
	fake.enforceOnceMutex.Unlock()
	if fake.EnforceOnceStub != nil {
		return fake.EnforceOnceStub()
	} else {
		return fake.enforceOnceReturns.result1
	}
}

func (fake *FakeEnforcer) EnforceOnceCallCount() int {
	fake.enforceOnceMutex.RLock()
	defer fake.enforceOnceMutex.RUnlock()
	return len(fake.enforceOnceArgsForCall)
}

func (fake *FakeEnforcer) EnforceOnceReturns(result1 error) {
	fake.EnforceOnceStub = nil
	fake.enforceOnceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEnforcer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.enforceOnceMutex.RLock()
	defer fake.enforceOnceMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeEnforcer) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ enforcer.Enforcer = new(FakeEnforcer)
