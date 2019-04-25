// Code generated by counterfeiter. DO NOT EDIT.
package registryfakes

import (
	sync "sync"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	layout "github.com/google/go-containerregistry/pkg/v1/layout"
	registry "github.com/pivotal/image-relocation/pkg/registry"
)

type FakeLayoutPath struct {
	AppendImageStub        func(v1.Image, ...layout.Option) error
	appendImageMutex       sync.RWMutex
	appendImageArgsForCall []struct {
		arg1 v1.Image
		arg2 []layout.Option
	}
	appendImageReturns struct {
		result1 error
	}
	appendImageReturnsOnCall map[int]struct {
		result1 error
	}
	ImageIndexStub        func() (v1.ImageIndex, error)
	imageIndexMutex       sync.RWMutex
	imageIndexArgsForCall []struct {
	}
	imageIndexReturns struct {
		result1 v1.ImageIndex
		result2 error
	}
	imageIndexReturnsOnCall map[int]struct {
		result1 v1.ImageIndex
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLayoutPath) AppendImage(arg1 v1.Image, arg2 ...layout.Option) error {
	fake.appendImageMutex.Lock()
	ret, specificReturn := fake.appendImageReturnsOnCall[len(fake.appendImageArgsForCall)]
	fake.appendImageArgsForCall = append(fake.appendImageArgsForCall, struct {
		arg1 v1.Image
		arg2 []layout.Option
	}{arg1, arg2})
	fake.recordInvocation("AppendImage", []interface{}{arg1, arg2})
	fake.appendImageMutex.Unlock()
	if fake.AppendImageStub != nil {
		return fake.AppendImageStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.appendImageReturns
	return fakeReturns.result1
}

func (fake *FakeLayoutPath) AppendImageCallCount() int {
	fake.appendImageMutex.RLock()
	defer fake.appendImageMutex.RUnlock()
	return len(fake.appendImageArgsForCall)
}

func (fake *FakeLayoutPath) AppendImageCalls(stub func(v1.Image, ...layout.Option) error) {
	fake.appendImageMutex.Lock()
	defer fake.appendImageMutex.Unlock()
	fake.AppendImageStub = stub
}

func (fake *FakeLayoutPath) AppendImageArgsForCall(i int) (v1.Image, []layout.Option) {
	fake.appendImageMutex.RLock()
	defer fake.appendImageMutex.RUnlock()
	argsForCall := fake.appendImageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLayoutPath) AppendImageReturns(result1 error) {
	fake.appendImageMutex.Lock()
	defer fake.appendImageMutex.Unlock()
	fake.AppendImageStub = nil
	fake.appendImageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLayoutPath) AppendImageReturnsOnCall(i int, result1 error) {
	fake.appendImageMutex.Lock()
	defer fake.appendImageMutex.Unlock()
	fake.AppendImageStub = nil
	if fake.appendImageReturnsOnCall == nil {
		fake.appendImageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.appendImageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLayoutPath) ImageIndex() (v1.ImageIndex, error) {
	fake.imageIndexMutex.Lock()
	ret, specificReturn := fake.imageIndexReturnsOnCall[len(fake.imageIndexArgsForCall)]
	fake.imageIndexArgsForCall = append(fake.imageIndexArgsForCall, struct {
	}{})
	fake.recordInvocation("ImageIndex", []interface{}{})
	fake.imageIndexMutex.Unlock()
	if fake.ImageIndexStub != nil {
		return fake.ImageIndexStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.imageIndexReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLayoutPath) ImageIndexCallCount() int {
	fake.imageIndexMutex.RLock()
	defer fake.imageIndexMutex.RUnlock()
	return len(fake.imageIndexArgsForCall)
}

func (fake *FakeLayoutPath) ImageIndexCalls(stub func() (v1.ImageIndex, error)) {
	fake.imageIndexMutex.Lock()
	defer fake.imageIndexMutex.Unlock()
	fake.ImageIndexStub = stub
}

func (fake *FakeLayoutPath) ImageIndexReturns(result1 v1.ImageIndex, result2 error) {
	fake.imageIndexMutex.Lock()
	defer fake.imageIndexMutex.Unlock()
	fake.ImageIndexStub = nil
	fake.imageIndexReturns = struct {
		result1 v1.ImageIndex
		result2 error
	}{result1, result2}
}

func (fake *FakeLayoutPath) ImageIndexReturnsOnCall(i int, result1 v1.ImageIndex, result2 error) {
	fake.imageIndexMutex.Lock()
	defer fake.imageIndexMutex.Unlock()
	fake.ImageIndexStub = nil
	if fake.imageIndexReturnsOnCall == nil {
		fake.imageIndexReturnsOnCall = make(map[int]struct {
			result1 v1.ImageIndex
			result2 error
		})
	}
	fake.imageIndexReturnsOnCall[i] = struct {
		result1 v1.ImageIndex
		result2 error
	}{result1, result2}
}

func (fake *FakeLayoutPath) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.appendImageMutex.RLock()
	defer fake.appendImageMutex.RUnlock()
	fake.imageIndexMutex.RLock()
	defer fake.imageIndexMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLayoutPath) recordInvocation(key string, args []interface{}) {
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

var _ registry.LayoutPath = new(FakeLayoutPath)
