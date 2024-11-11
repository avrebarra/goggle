// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package rpcserver_test

import (
	"context"
	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"sync"
)

// Ensure, that ToggleServiceMock does implement rpcserver.ToggleService.
// If this is not the case, regenerate this file with moq.
var _ rpcserver.ToggleService = &ToggleServiceMock{}

// ToggleServiceMock is a mock implementation of rpcserver.ToggleService.
//
//	func TestSomethingThatUsesToggleService(t *testing.T) {
//
//		// make and configure a mocked rpcserver.ToggleService
//		mockedToggleService := &ToggleServiceMock{
//			DoCreateToggleFunc: func(ctx context.Context, in domaintoggle.Toggle) (domaintoggle.Toggle, error) {
//				panic("mock out the DoCreateToggle method")
//			},
//			DoGetToggleFunc: func(ctx context.Context, id string) (domaintoggle.ToggleWithDetail, error) {
//				panic("mock out the DoGetToggle method")
//			},
//			DoListStrayTogglesFunc: func(ctx context.Context, in servicetoggle.ParamsDoListStrayToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
//				panic("mock out the DoListStrayToggles method")
//			},
//			DoListTogglesFunc: func(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
//				panic("mock out the DoListToggles method")
//			},
//			DoRemoveToggleFunc: func(ctx context.Context, id string) (domaintoggle.Toggle, error) {
//				panic("mock out the DoRemoveToggle method")
//			},
//			DoStatToggleFunc: func(ctx context.Context, id string) (domaintoggle.ToggleStat, error) {
//				panic("mock out the DoStatToggle method")
//			},
//			DoUpdateToggleFunc: func(ctx context.Context, id string, in domaintoggle.Toggle) (domaintoggle.Toggle, error) {
//				panic("mock out the DoUpdateToggle method")
//			},
//		}
//
//		// use mockedToggleService in code that requires rpcserver.ToggleService
//		// and then make assertions.
//
//	}
type ToggleServiceMock struct {
	// DoCreateToggleFunc mocks the DoCreateToggle method.
	DoCreateToggleFunc func(ctx context.Context, in domaintoggle.Toggle) (domaintoggle.Toggle, error)

	// DoGetToggleFunc mocks the DoGetToggle method.
	DoGetToggleFunc func(ctx context.Context, id string) (domaintoggle.ToggleWithDetail, error)

	// DoListStrayTogglesFunc mocks the DoListStrayToggles method.
	DoListStrayTogglesFunc func(ctx context.Context, in servicetoggle.ParamsDoListStrayToggles) ([]domaintoggle.ToggleWithDetail, int64, error)

	// DoListTogglesFunc mocks the DoListToggles method.
	DoListTogglesFunc func(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error)

	// DoRemoveToggleFunc mocks the DoRemoveToggle method.
	DoRemoveToggleFunc func(ctx context.Context, id string) (domaintoggle.Toggle, error)

	// DoStatToggleFunc mocks the DoStatToggle method.
	DoStatToggleFunc func(ctx context.Context, id string) (domaintoggle.ToggleStat, error)

	// DoUpdateToggleFunc mocks the DoUpdateToggle method.
	DoUpdateToggleFunc func(ctx context.Context, id string, in domaintoggle.Toggle) (domaintoggle.Toggle, error)

	// calls tracks calls to the methods.
	calls struct {
		// DoCreateToggle holds details about calls to the DoCreateToggle method.
		DoCreateToggle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In domaintoggle.Toggle
		}
		// DoGetToggle holds details about calls to the DoGetToggle method.
		DoGetToggle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// DoListStrayToggles holds details about calls to the DoListStrayToggles method.
		DoListStrayToggles []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In servicetoggle.ParamsDoListStrayToggles
		}
		// DoListToggles holds details about calls to the DoListToggles method.
		DoListToggles []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In servicetoggle.ParamsDoListToggles
		}
		// DoRemoveToggle holds details about calls to the DoRemoveToggle method.
		DoRemoveToggle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// DoStatToggle holds details about calls to the DoStatToggle method.
		DoStatToggle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// DoUpdateToggle holds details about calls to the DoUpdateToggle method.
		DoUpdateToggle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
			// In is the in argument value.
			In domaintoggle.Toggle
		}
	}
	lockDoCreateToggle     sync.RWMutex
	lockDoGetToggle        sync.RWMutex
	lockDoListStrayToggles sync.RWMutex
	lockDoListToggles      sync.RWMutex
	lockDoRemoveToggle     sync.RWMutex
	lockDoStatToggle       sync.RWMutex
	lockDoUpdateToggle     sync.RWMutex
}

// DoCreateToggle calls DoCreateToggleFunc.
func (mock *ToggleServiceMock) DoCreateToggle(ctx context.Context, in domaintoggle.Toggle) (domaintoggle.Toggle, error) {
	callInfo := struct {
		Ctx context.Context
		In  domaintoggle.Toggle
	}{
		Ctx: ctx,
		In:  in,
	}
	mock.lockDoCreateToggle.Lock()
	mock.calls.DoCreateToggle = append(mock.calls.DoCreateToggle, callInfo)
	mock.lockDoCreateToggle.Unlock()
	if mock.DoCreateToggleFunc == nil {
		var (
			outOut domaintoggle.Toggle
			errOut error
		)
		return outOut, errOut
	}
	return mock.DoCreateToggleFunc(ctx, in)
}

// DoCreateToggleCalls gets all the calls that were made to DoCreateToggle.
// Check the length with:
//
//	len(mockedToggleService.DoCreateToggleCalls())
func (mock *ToggleServiceMock) DoCreateToggleCalls() []struct {
	Ctx context.Context
	In  domaintoggle.Toggle
} {
	var calls []struct {
		Ctx context.Context
		In  domaintoggle.Toggle
	}
	mock.lockDoCreateToggle.RLock()
	calls = mock.calls.DoCreateToggle
	mock.lockDoCreateToggle.RUnlock()
	return calls
}

// DoGetToggle calls DoGetToggleFunc.
func (mock *ToggleServiceMock) DoGetToggle(ctx context.Context, id string) (domaintoggle.ToggleWithDetail, error) {
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDoGetToggle.Lock()
	mock.calls.DoGetToggle = append(mock.calls.DoGetToggle, callInfo)
	mock.lockDoGetToggle.Unlock()
	if mock.DoGetToggleFunc == nil {
		var (
			outOut domaintoggle.ToggleWithDetail
			errOut error
		)
		return outOut, errOut
	}
	return mock.DoGetToggleFunc(ctx, id)
}

// DoGetToggleCalls gets all the calls that were made to DoGetToggle.
// Check the length with:
//
//	len(mockedToggleService.DoGetToggleCalls())
func (mock *ToggleServiceMock) DoGetToggleCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockDoGetToggle.RLock()
	calls = mock.calls.DoGetToggle
	mock.lockDoGetToggle.RUnlock()
	return calls
}

// DoListStrayToggles calls DoListStrayTogglesFunc.
func (mock *ToggleServiceMock) DoListStrayToggles(ctx context.Context, in servicetoggle.ParamsDoListStrayToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
	callInfo := struct {
		Ctx context.Context
		In  servicetoggle.ParamsDoListStrayToggles
	}{
		Ctx: ctx,
		In:  in,
	}
	mock.lockDoListStrayToggles.Lock()
	mock.calls.DoListStrayToggles = append(mock.calls.DoListStrayToggles, callInfo)
	mock.lockDoListStrayToggles.Unlock()
	if mock.DoListStrayTogglesFunc == nil {
		var (
			outOut   []domaintoggle.ToggleWithDetail
			totalOut int64
			errOut   error
		)
		return outOut, totalOut, errOut
	}
	return mock.DoListStrayTogglesFunc(ctx, in)
}

// DoListStrayTogglesCalls gets all the calls that were made to DoListStrayToggles.
// Check the length with:
//
//	len(mockedToggleService.DoListStrayTogglesCalls())
func (mock *ToggleServiceMock) DoListStrayTogglesCalls() []struct {
	Ctx context.Context
	In  servicetoggle.ParamsDoListStrayToggles
} {
	var calls []struct {
		Ctx context.Context
		In  servicetoggle.ParamsDoListStrayToggles
	}
	mock.lockDoListStrayToggles.RLock()
	calls = mock.calls.DoListStrayToggles
	mock.lockDoListStrayToggles.RUnlock()
	return calls
}

// DoListToggles calls DoListTogglesFunc.
func (mock *ToggleServiceMock) DoListToggles(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
	callInfo := struct {
		Ctx context.Context
		In  servicetoggle.ParamsDoListToggles
	}{
		Ctx: ctx,
		In:  in,
	}
	mock.lockDoListToggles.Lock()
	mock.calls.DoListToggles = append(mock.calls.DoListToggles, callInfo)
	mock.lockDoListToggles.Unlock()
	if mock.DoListTogglesFunc == nil {
		var (
			outOut   []domaintoggle.ToggleWithDetail
			totalOut int64
			errOut   error
		)
		return outOut, totalOut, errOut
	}
	return mock.DoListTogglesFunc(ctx, in)
}

// DoListTogglesCalls gets all the calls that were made to DoListToggles.
// Check the length with:
//
//	len(mockedToggleService.DoListTogglesCalls())
func (mock *ToggleServiceMock) DoListTogglesCalls() []struct {
	Ctx context.Context
	In  servicetoggle.ParamsDoListToggles
} {
	var calls []struct {
		Ctx context.Context
		In  servicetoggle.ParamsDoListToggles
	}
	mock.lockDoListToggles.RLock()
	calls = mock.calls.DoListToggles
	mock.lockDoListToggles.RUnlock()
	return calls
}

// DoRemoveToggle calls DoRemoveToggleFunc.
func (mock *ToggleServiceMock) DoRemoveToggle(ctx context.Context, id string) (domaintoggle.Toggle, error) {
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDoRemoveToggle.Lock()
	mock.calls.DoRemoveToggle = append(mock.calls.DoRemoveToggle, callInfo)
	mock.lockDoRemoveToggle.Unlock()
	if mock.DoRemoveToggleFunc == nil {
		var (
			outOut domaintoggle.Toggle
			errOut error
		)
		return outOut, errOut
	}
	return mock.DoRemoveToggleFunc(ctx, id)
}

// DoRemoveToggleCalls gets all the calls that were made to DoRemoveToggle.
// Check the length with:
//
//	len(mockedToggleService.DoRemoveToggleCalls())
func (mock *ToggleServiceMock) DoRemoveToggleCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockDoRemoveToggle.RLock()
	calls = mock.calls.DoRemoveToggle
	mock.lockDoRemoveToggle.RUnlock()
	return calls
}

// DoStatToggle calls DoStatToggleFunc.
func (mock *ToggleServiceMock) DoStatToggle(ctx context.Context, id string) (domaintoggle.ToggleStat, error) {
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDoStatToggle.Lock()
	mock.calls.DoStatToggle = append(mock.calls.DoStatToggle, callInfo)
	mock.lockDoStatToggle.Unlock()
	if mock.DoStatToggleFunc == nil {
		var (
			outOut domaintoggle.ToggleStat
			errOut error
		)
		return outOut, errOut
	}
	return mock.DoStatToggleFunc(ctx, id)
}

// DoStatToggleCalls gets all the calls that were made to DoStatToggle.
// Check the length with:
//
//	len(mockedToggleService.DoStatToggleCalls())
func (mock *ToggleServiceMock) DoStatToggleCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockDoStatToggle.RLock()
	calls = mock.calls.DoStatToggle
	mock.lockDoStatToggle.RUnlock()
	return calls
}

// DoUpdateToggle calls DoUpdateToggleFunc.
func (mock *ToggleServiceMock) DoUpdateToggle(ctx context.Context, id string, in domaintoggle.Toggle) (domaintoggle.Toggle, error) {
	callInfo := struct {
		Ctx context.Context
		ID  string
		In  domaintoggle.Toggle
	}{
		Ctx: ctx,
		ID:  id,
		In:  in,
	}
	mock.lockDoUpdateToggle.Lock()
	mock.calls.DoUpdateToggle = append(mock.calls.DoUpdateToggle, callInfo)
	mock.lockDoUpdateToggle.Unlock()
	if mock.DoUpdateToggleFunc == nil {
		var (
			outOut domaintoggle.Toggle
			errOut error
		)
		return outOut, errOut
	}
	return mock.DoUpdateToggleFunc(ctx, id, in)
}

// DoUpdateToggleCalls gets all the calls that were made to DoUpdateToggle.
// Check the length with:
//
//	len(mockedToggleService.DoUpdateToggleCalls())
func (mock *ToggleServiceMock) DoUpdateToggleCalls() []struct {
	Ctx context.Context
	ID  string
	In  domaintoggle.Toggle
} {
	var calls []struct {
		Ctx context.Context
		ID  string
		In  domaintoggle.Toggle
	}
	mock.lockDoUpdateToggle.RLock()
	calls = mock.calls.DoUpdateToggle
	mock.lockDoUpdateToggle.RUnlock()
	return calls
}