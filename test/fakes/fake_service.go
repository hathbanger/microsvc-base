// Code generated by counterfeiter. DO NOT EDIT.
package microsvcfakes

import (
	"context"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hathbanger/microsvc-base/pkg/microsvc"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

type FakeService struct {
	FooStub        func(context.Context, models.FooRequest) (models.FooResponse, error)
	fooMutex       sync.RWMutex
	fooArgsForCall []struct {
		arg1 context.Context
		arg2 models.FooRequest
	}
	fooReturns struct {
		result1 models.FooResponse
		result2 error
	}
	fooReturnsOnCall map[int]struct {
		result1 models.FooResponse
		result2 error
	}
	HealthStub        func() bool
	healthMutex       sync.RWMutex
	healthArgsForCall []struct {
	}
	healthReturns struct {
		result1 bool
	}
	healthReturnsOnCall map[int]struct {
		result1 bool
	}
	ServiceDiscoveryStub        func(string, string) (*api.Client, *api.AgentServiceRegistration, error)
	serviceDiscoveryMutex       sync.RWMutex
	serviceDiscoveryArgsForCall []struct {
		arg1 string
		arg2 string
	}
	serviceDiscoveryReturns struct {
		result1 *api.Client
		result2 *api.AgentServiceRegistration
		result3 error
	}
	serviceDiscoveryReturnsOnCall map[int]struct {
		result1 *api.Client
		result2 *api.AgentServiceRegistration
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeService) Foo(arg1 context.Context, arg2 models.FooRequest) (models.FooResponse, error) {
	fake.fooMutex.Lock()
	ret, specificReturn := fake.fooReturnsOnCall[len(fake.fooArgsForCall)]
	fake.fooArgsForCall = append(fake.fooArgsForCall, struct {
		arg1 context.Context
		arg2 models.FooRequest
	}{arg1, arg2})
	fake.recordInvocation("Foo", []interface{}{arg1, arg2})
	fake.fooMutex.Unlock()
	if fake.FooStub != nil {
		return fake.FooStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.fooReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) FooCallCount() int {
	fake.fooMutex.RLock()
	defer fake.fooMutex.RUnlock()
	return len(fake.fooArgsForCall)
}

func (fake *FakeService) FooCalls(stub func(context.Context, models.FooRequest) (models.FooResponse, error)) {
	fake.fooMutex.Lock()
	defer fake.fooMutex.Unlock()
	fake.FooStub = stub
}

func (fake *FakeService) FooArgsForCall(i int) (context.Context, models.FooRequest) {
	fake.fooMutex.RLock()
	defer fake.fooMutex.RUnlock()
	argsForCall := fake.fooArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeService) FooReturns(result1 models.FooResponse, result2 error) {
	fake.fooMutex.Lock()
	defer fake.fooMutex.Unlock()
	fake.FooStub = nil
	fake.fooReturns = struct {
		result1 models.FooResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeService) FooReturnsOnCall(i int, result1 models.FooResponse, result2 error) {
	fake.fooMutex.Lock()
	defer fake.fooMutex.Unlock()
	fake.FooStub = nil
	if fake.fooReturnsOnCall == nil {
		fake.fooReturnsOnCall = make(map[int]struct {
			result1 models.FooResponse
			result2 error
		})
	}
	fake.fooReturnsOnCall[i] = struct {
		result1 models.FooResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeService) Health() bool {
	fake.healthMutex.Lock()
	ret, specificReturn := fake.healthReturnsOnCall[len(fake.healthArgsForCall)]
	fake.healthArgsForCall = append(fake.healthArgsForCall, struct {
	}{})
	fake.recordInvocation("Health", []interface{}{})
	fake.healthMutex.Unlock()
	if fake.HealthStub != nil {
		return fake.HealthStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.healthReturns
	return fakeReturns.result1
}

func (fake *FakeService) HealthCallCount() int {
	fake.healthMutex.RLock()
	defer fake.healthMutex.RUnlock()
	return len(fake.healthArgsForCall)
}

func (fake *FakeService) HealthCalls(stub func() bool) {
	fake.healthMutex.Lock()
	defer fake.healthMutex.Unlock()
	fake.HealthStub = stub
}

func (fake *FakeService) HealthReturns(result1 bool) {
	fake.healthMutex.Lock()
	defer fake.healthMutex.Unlock()
	fake.HealthStub = nil
	fake.healthReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeService) HealthReturnsOnCall(i int, result1 bool) {
	fake.healthMutex.Lock()
	defer fake.healthMutex.Unlock()
	fake.HealthStub = nil
	if fake.healthReturnsOnCall == nil {
		fake.healthReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.healthReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeService) ServiceDiscovery(arg1 string, arg2 string) (*api.Client, *api.AgentServiceRegistration, error) {
	fake.serviceDiscoveryMutex.Lock()
	ret, specificReturn := fake.serviceDiscoveryReturnsOnCall[len(fake.serviceDiscoveryArgsForCall)]
	fake.serviceDiscoveryArgsForCall = append(fake.serviceDiscoveryArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ServiceDiscovery", []interface{}{arg1, arg2})
	fake.serviceDiscoveryMutex.Unlock()
	if fake.ServiceDiscoveryStub != nil {
		return fake.ServiceDiscoveryStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.serviceDiscoveryReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeService) ServiceDiscoveryCallCount() int {
	fake.serviceDiscoveryMutex.RLock()
	defer fake.serviceDiscoveryMutex.RUnlock()
	return len(fake.serviceDiscoveryArgsForCall)
}

func (fake *FakeService) ServiceDiscoveryCalls(stub func(string, string) (*api.Client, *api.AgentServiceRegistration, error)) {
	fake.serviceDiscoveryMutex.Lock()
	defer fake.serviceDiscoveryMutex.Unlock()
	fake.ServiceDiscoveryStub = stub
}

func (fake *FakeService) ServiceDiscoveryArgsForCall(i int) (string, string) {
	fake.serviceDiscoveryMutex.RLock()
	defer fake.serviceDiscoveryMutex.RUnlock()
	argsForCall := fake.serviceDiscoveryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeService) ServiceDiscoveryReturns(result1 *api.Client, result2 *api.AgentServiceRegistration, result3 error) {
	fake.serviceDiscoveryMutex.Lock()
	defer fake.serviceDiscoveryMutex.Unlock()
	fake.ServiceDiscoveryStub = nil
	fake.serviceDiscoveryReturns = struct {
		result1 *api.Client
		result2 *api.AgentServiceRegistration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeService) ServiceDiscoveryReturnsOnCall(i int, result1 *api.Client, result2 *api.AgentServiceRegistration, result3 error) {
	fake.serviceDiscoveryMutex.Lock()
	defer fake.serviceDiscoveryMutex.Unlock()
	fake.ServiceDiscoveryStub = nil
	if fake.serviceDiscoveryReturnsOnCall == nil {
		fake.serviceDiscoveryReturnsOnCall = make(map[int]struct {
			result1 *api.Client
			result2 *api.AgentServiceRegistration
			result3 error
		})
	}
	fake.serviceDiscoveryReturnsOnCall[i] = struct {
		result1 *api.Client
		result2 *api.AgentServiceRegistration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.fooMutex.RLock()
	defer fake.fooMutex.RUnlock()
	fake.healthMutex.RLock()
	defer fake.healthMutex.RUnlock()
	fake.serviceDiscoveryMutex.RLock()
	defer fake.serviceDiscoveryMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeService) recordInvocation(key string, args []interface{}) {
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

var _ microsvc.Service = new(FakeService)
