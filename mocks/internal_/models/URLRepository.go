// Code generated by mockery v2.43.2. DO NOT EDIT.

package models

import (
	context "context"

	models "github.com/ShukinDmitriy/shortener/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// URLRepository is an autogenerated mock type for the URLRepository type
type URLRepository struct {
	mock.Mock
}

type URLRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *URLRepository) EXPECT() *URLRepository_Expecter {
	return &URLRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, events
func (_m *URLRepository) Delete(ctx context.Context, events []models.DeleteRequestBatch) error {
	ret := _m.Called(ctx, events)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []models.DeleteRequestBatch) error); ok {
		r0 = rf(ctx, events)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type URLRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - events []models.DeleteRequestBatch
func (_e *URLRepository_Expecter) Delete(ctx interface{}, events interface{}) *URLRepository_Delete_Call {
	return &URLRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, events)}
}

func (_c *URLRepository_Delete_Call) Run(run func(ctx context.Context, events []models.DeleteRequestBatch)) *URLRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]models.DeleteRequestBatch))
	})
	return _c
}

func (_c *URLRepository_Delete_Call) Return(_a0 error) *URLRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLRepository_Delete_Call) RunAndReturn(run func(context.Context, []models.DeleteRequestBatch) error) *URLRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: shortKey
func (_m *URLRepository) Get(shortKey string) (models.Event, bool) {
	ret := _m.Called(shortKey)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 models.Event
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (models.Event, bool)); ok {
		return rf(shortKey)
	}
	if rf, ok := ret.Get(0).(func(string) models.Event); ok {
		r0 = rf(shortKey)
	} else {
		r0 = ret.Get(0).(models.Event)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(shortKey)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// URLRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type URLRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - shortKey string
func (_e *URLRepository_Expecter) Get(shortKey interface{}) *URLRepository_Get_Call {
	return &URLRepository_Get_Call{Call: _e.mock.On("Get", shortKey)}
}

func (_c *URLRepository_Get_Call) Run(run func(shortKey string)) *URLRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *URLRepository_Get_Call) Return(_a0 models.Event, _a1 bool) *URLRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLRepository_Get_Call) RunAndReturn(run func(string) (models.Event, bool)) *URLRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetEventsByUserID provides a mock function with given fields: ctx, userID
func (_m *URLRepository) GetEventsByUserID(ctx context.Context, userID string) []*models.Event {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByUserID")
	}

	var r0 []*models.Event
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Event); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Event)
		}
	}

	return r0
}

// URLRepository_GetEventsByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEventsByUserID'
type URLRepository_GetEventsByUserID_Call struct {
	*mock.Call
}

// GetEventsByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *URLRepository_Expecter) GetEventsByUserID(ctx interface{}, userID interface{}) *URLRepository_GetEventsByUserID_Call {
	return &URLRepository_GetEventsByUserID_Call{Call: _e.mock.On("GetEventsByUserID", ctx, userID)}
}

func (_c *URLRepository_GetEventsByUserID_Call) Run(run func(ctx context.Context, userID string)) *URLRepository_GetEventsByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *URLRepository_GetEventsByUserID_Call) Return(_a0 []*models.Event) *URLRepository_GetEventsByUserID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLRepository_GetEventsByUserID_Call) RunAndReturn(run func(context.Context, string) []*models.Event) *URLRepository_GetEventsByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// GetShortKeyByOriginalURL provides a mock function with given fields: originalURL
func (_m *URLRepository) GetShortKeyByOriginalURL(originalURL string) (string, bool) {
	ret := _m.Called(originalURL)

	if len(ret) == 0 {
		panic("no return value specified for GetShortKeyByOriginalURL")
	}

	var r0 string
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (string, bool)); ok {
		return rf(originalURL)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(originalURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(originalURL)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// URLRepository_GetShortKeyByOriginalURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetShortKeyByOriginalURL'
type URLRepository_GetShortKeyByOriginalURL_Call struct {
	*mock.Call
}

// GetShortKeyByOriginalURL is a helper method to define mock.On call
//   - originalURL string
func (_e *URLRepository_Expecter) GetShortKeyByOriginalURL(originalURL interface{}) *URLRepository_GetShortKeyByOriginalURL_Call {
	return &URLRepository_GetShortKeyByOriginalURL_Call{Call: _e.mock.On("GetShortKeyByOriginalURL", originalURL)}
}

func (_c *URLRepository_GetShortKeyByOriginalURL_Call) Run(run func(originalURL string)) *URLRepository_GetShortKeyByOriginalURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *URLRepository_GetShortKeyByOriginalURL_Call) Return(_a0 string, _a1 bool) *URLRepository_GetShortKeyByOriginalURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLRepository_GetShortKeyByOriginalURL_Call) RunAndReturn(run func(string) (string, bool)) *URLRepository_GetShortKeyByOriginalURL_Call {
	_c.Call.Return(run)
	return _c
}

// Initialize provides a mock function with given fields:
func (_m *URLRepository) Initialize() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Initialize")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLRepository_Initialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Initialize'
type URLRepository_Initialize_Call struct {
	*mock.Call
}

// Initialize is a helper method to define mock.On call
func (_e *URLRepository_Expecter) Initialize() *URLRepository_Initialize_Call {
	return &URLRepository_Initialize_Call{Call: _e.mock.On("Initialize")}
}

func (_c *URLRepository_Initialize_Call) Run(run func()) *URLRepository_Initialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *URLRepository_Initialize_Call) Return(_a0 error) *URLRepository_Initialize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLRepository_Initialize_Call) RunAndReturn(run func() error) *URLRepository_Initialize_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, events
func (_m *URLRepository) Save(ctx context.Context, events []*models.Event) error {
	ret := _m.Called(ctx, events)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*models.Event) error); ok {
		r0 = rf(ctx, events)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type URLRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - events []*models.Event
func (_e *URLRepository_Expecter) Save(ctx interface{}, events interface{}) *URLRepository_Save_Call {
	return &URLRepository_Save_Call{Call: _e.mock.On("Save", ctx, events)}
}

func (_c *URLRepository_Save_Call) Run(run func(ctx context.Context, events []*models.Event)) *URLRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*models.Event))
	})
	return _c
}

func (_c *URLRepository_Save_Call) Return(_a0 error) *URLRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLRepository_Save_Call) RunAndReturn(run func(context.Context, []*models.Event) error) *URLRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewURLRepository creates a new instance of URLRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLRepository {
	mock := &URLRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
