package mocks

import "github.com/stretchr/testify/mock"

type MockShell struct {
	mock.Mock
}

func (m *MockShell) Command(name string, arg ...string) (string, error) {
	ret := m.Called(name, arg)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, ...string) string); ok {
		r0 = rf(name, arg...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(name, arg...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
