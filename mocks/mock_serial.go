package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tarm/serial"
)

type MockSerial struct {
	mock.Mock
}

func (m *MockSerial) OpenPort(c *serial.Config) (*serial.Port, error) {
	ret := m.Called(c)

	var r0 *serial.Port
	if rf, ok := ret.Get(0).(func(*serial.Config) *serial.Port); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(*serial.Port)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*serial.Config) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockSerial) Write(port *serial.Port, command []byte) (int, error) {
	ret := m.Called(port, command)

	var r0 int
	if rf, ok := ret.Get(0).(func(*serial.Port, []byte) int); ok {
		r0 = rf(port, command)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*serial.Port, []byte) error); ok {
		r1 = rf(port, command)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockSerial) Close(port *serial.Port) error {
	ret := m.Called(port)

	var r0 error
	if rf, ok := ret.Get(0).(func(*serial.Port) error); ok {
		r0 = rf(port)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *MockSerial) Read(port *serial.Port, buffer []byte) (int, error) {
	ret := m.Called(port, buffer)

	var r0 int
	if rf, ok := ret.Get(0).(func(*serial.Port, []byte) int); ok {
		r0 = rf(port, buffer)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*serial.Port, []byte) error); ok {
		r1 = rf(port, buffer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
