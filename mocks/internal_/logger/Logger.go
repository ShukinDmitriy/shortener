// Code generated by mockery v2.43.2. DO NOT EDIT.

package logger

import (
	mock "github.com/stretchr/testify/mock"
	zapcore "go.uber.org/zap/zapcore"
)

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

type Logger_Expecter struct {
	mock *mock.Mock
}

func (_m *Logger) EXPECT() *Logger_Expecter {
	return &Logger_Expecter{mock: &_m.Mock}
}

// Error provides a mock function with given fields: msg, fields
func (_m *Logger) Error(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type Logger_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Error(msg interface{}, fields ...interface{}) *Logger_Error_Call {
	return &Logger_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Error_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Error_Call) Return() *Logger_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Error_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Fatal provides a mock function with given fields: msg, fields
func (_m *Logger) Fatal(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Fatal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fatal'
type Logger_Fatal_Call struct {
	*mock.Call
}

// Fatal is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Fatal(msg interface{}, fields ...interface{}) *Logger_Fatal_Call {
	return &Logger_Fatal_Call{Call: _e.mock.On("Fatal",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Fatal_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Fatal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Fatal_Call) Return() *Logger_Fatal_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Fatal_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Fatal_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: msg, fields
func (_m *Logger) Info(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type Logger_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Info(msg interface{}, fields ...interface{}) *Logger_Info_Call {
	return &Logger_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Info_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Info_Call) Return() *Logger_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Info_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Info_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: msg, fields
func (_m *Logger) Warn(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Logger_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type Logger_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - msg string
//   - fields ...zapcore.Field
func (_e *Logger_Expecter) Warn(msg interface{}, fields ...interface{}) *Logger_Warn_Call {
	return &Logger_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{msg}, fields...)...)}
}

func (_c *Logger_Warn_Call) Run(run func(msg string, fields ...zapcore.Field)) *Logger_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]zapcore.Field, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(zapcore.Field)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *Logger_Warn_Call) Return() *Logger_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *Logger_Warn_Call) RunAndReturn(run func(string, ...zapcore.Field)) *Logger_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// NewLogger creates a new instance of Logger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
