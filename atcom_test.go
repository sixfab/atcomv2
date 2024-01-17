package atcom

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/tarm/serial"
)

type SerialShell struct {
	mocked map[string]interface{}
}

var (
	mockSerialPort = &serial.Port{}
	mockWrite      = 2
	mockSleepFunc  = func(d time.Duration) {}
	mockRead       = 2
)

var serialShell = &SerialShell{
	mocked: map[string]interface{}{
		"openPort": map[string]interface{}{
			"resp": mockSerialPort,
			"err":  nil,
		},
		"Write": map[string]interface{}{
			"resp": mockWrite,
			"err":  nil,
		},
		"Read": map[string]interface{}{
			"resp": mockRead,
			"err":  nil,
		},
	},
}

func (t *SerialShell) Patch(cmd string, resp interface{}, err error) {
	t.mocked[cmd] = map[string]interface{}{"resp": resp, "err": err}
}

func (t *SerialShell) OpenPort(c *serial.Config) (*serial.Port, error) {

	for mocked_name := range t.mocked {

		if mocked_name == "openPort" {
			response, _ := t.mocked[mocked_name].(map[string]interface{})

			if response["err"] == nil {
				return response["resp"].(*serial.Port), nil
			}
			return response["resp"].(*serial.Port), response["err"].(error)
		}
	}
	return nil, nil
}

func (t *SerialShell) Write(port *serial.Port, command []byte) (n int, err error) {
	for mocked_name := range t.mocked {

		if mocked_name == "Write" {
			response, _ := t.mocked[mocked_name].(map[string]interface{})

			if response["err"] == nil {
				return response["resp"].(int), nil
			}
			return response["resp"].(int), response["err"].(error)
		}

	}
	return 5, nil
}

func (t *SerialShell) Close(port *serial.Port) (err error) {
	return nil
}

func (t *SerialShell) Read(port *serial.Port, buffer []byte) (n int, err error) {
	for mocked_name := range t.mocked {

		if mocked_name == "Read" {
			response, _ := t.mocked[mocked_name].(map[string]interface{})

			if response["err"] == nil {
				return response["resp"].(int), nil
			}
			return response["resp"].(int), response["err"].(error)
		}

	}
	return 5, nil
}

func TestOpen(t *testing.T) {

	t.Run("Should return port with nil input parameter", func(t *testing.T) {

		at := NewAtcom(serialShell, nil, nil)
		port, err := at.open(nil)

		expectedPort := &serial.Port{}

		if !reflect.DeepEqual(port, expectedPort) {
			t.Errorf("Expected %p, but got %p", expectedPort, port)
		}
		if err != nil {
			t.Errorf("Expected nil error, but got %v", err)
		}
	})

	t.Run("Should return port with input parameter", func(t *testing.T) {

		at := NewAtcom(serialShell, nil, nil)

		arg := map[string]interface{}{
			"port": "/dev/ttyUSB2",
			"baud": 115200,
		}
		port, err := at.open(arg)

		expectedPort := &serial.Port{}

		if !reflect.DeepEqual(port, expectedPort) {
			t.Errorf("Expected %p, but got %p", expectedPort, port)
		}
		if err != nil {
			t.Errorf("Expected nil error, but got %v", err)
		}
	})

}

func TestSendAT(t *testing.T) {

	t.Run("Should return error for open function", func(t *testing.T) {

		commandName := "openPort"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Serial port error")

		serialShell.Patch(commandName, mockSerialPort, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil, nil)
		_, err := at.SendAT("AT", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error for write function", func(t *testing.T) {

		commandName := "Write"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Write function error")

		serialShell.Patch(commandName, mockWrite, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil, nil)
		_, err := at.SendAT("AT", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error for Read ", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Read Error")

		serialShell.Patch(commandName, mockRead, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil, mockSleepFunc)
		_, err := at.SendAT("AT", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

}
