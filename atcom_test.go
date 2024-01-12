package atcom

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tarm/serial"
)

type SerialShell struct {
	mocked map[string]interface{}
}

var mockSerialPort = &serial.Port{}

var serialShell = &SerialShell{
	mocked: map[string]interface{}{
		"openPort": map[string]interface{}{
			"resp": mockSerialPort,
			"err":  nil,
		},
	},
}

func (t *SerialShell) OpenPortPatch(resp *serial.Port, err error) {
	t.mocked["OpenPort"] = map[string]interface{}{"resp": resp, "err": err}
}

func (t *SerialShell) OpenPort(c *serial.Config) (*serial.Port, error) {

	for mocked_name := range t.mocked {
		response, _ := t.mocked[mocked_name].(map[string]interface{})

		if response["err"] == nil {
			return response["resp"].(*serial.Port), nil
		}
		return response["resp"].(*serial.Port), response["err"].(error)

	}
	return nil, nil
}

func (t *SerialShell) Write(port *serial.Port, command []byte) (n int, err error) {
	return 5, nil
}

func (t *SerialShell) Close(port *serial.Port) (err error) {
	return nil
}

func TestOpenPort(t *testing.T) {

	t.Run("Should return port with nil input parameter", func(t *testing.T) {

		at := NewAtcom(serialShell, nil)
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

		at := NewAtcom(serialShell, nil)

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

	t.Run("", func(t *testing.T) {

		commandName := "openPort"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Serial port error")

		serialShell.OpenPortPatch(mockSerialPort, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil)
		_, err := at.SendAT("", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})
}
