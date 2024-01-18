package atcom

import (
	"errors"
	"fmt"
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
			//"buffer": "",
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

			//buffer = response["buffer"].(string)

			if response["err"] == nil {
				return response["resp"].(int), nil
			}
			return response["resp"].(int), response["err"].(error)
		}

	}
	return 5, nil
}

func TestNewAtcom(t *testing.T) {

	t.Run("Should return NewAtcom for", func(t *testing.T) {

		parameters := []struct {
			name    string
			serial  Serial
			shell   Shell
			sleep   func(d time.Duration)
			desired *Atcom
		}{
			{"no mocking", nil, nil, nil, &Atcom{&RealSerial{}, &RealShell{}, time.Sleep}},
			{"mocking serial", &SerialShell{}, nil, nil, &Atcom{&SerialShell{}, &RealShell{}, time.Sleep}},
			{"mocking shell", nil, &MockShell{}, nil, &Atcom{&RealSerial{}, &MockShell{}, time.Sleep}},
			{"mocking Sleep", nil, nil, mockSleepFunc, &Atcom{&RealSerial{}, &RealShell{}, mockSleepFunc}},
			{"mocking all", &SerialShell{}, &MockShell{}, mockSleepFunc, &Atcom{&SerialShell{}, &MockShell{}, mockSleepFunc}},
		}

		for _, tt := range parameters {

			t.Run(tt.name, func(t *testing.T) {

				response := NewAtcom(tt.serial, tt.shell, tt.sleep)

				if !reflect.DeepEqual(response.serial, tt.desired.serial) {
					t.Errorf("Expected %v, but got %v", response.serial, tt.desired.serial)
				}
				if !reflect.DeepEqual(response.shell, tt.desired.shell) {
					t.Errorf("Expected %v, but got %v", response.shell, tt.desired.shell)
				}
				if fmt.Sprintf("%p", response.Sleep) != fmt.Sprintf("%p", tt.desired.Sleep) {
					t.Errorf("Expected %p, but got %p", tt.desired.Sleep, response.Sleep)
				}
			})
		}
	})
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
		_, err := at.SendAT("ATE1", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error for Write function", func(t *testing.T) {

		commandName := "Write"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Write function error")

		serialShell.Patch(commandName, mockWrite, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil, nil)
		_, err := at.SendAT("ATE1", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error for Read function", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := serialShell.mocked[commandName]
		error := errors.New("Read Error")

		serialShell.Patch(commandName, mockRead, error)
		defer func() { serialShell.mocked[commandName] = mockedDefault }()

		at := NewAtcom(serialShell, nil, mockSleepFunc)
		_, err := at.SendAT("ATE1", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})
	/*
			t.Run("Should return successfull response", func(t *testing.T) {

				parameters := []struct {
					command          string
					mock             int
					desired_response string
				}{
					{"ATE1", 316, `[+COPS: (1,"AVEA","AVEA","28603",7),(1,"AVEA","AVEA","28603",2),(1,"TR TURKCELL","TCELL","28601",0),(1,"Vodafone Turkiye","VF-TR","28602",0),(1,"Vodafone Turkiye","VF-TR","28602",7),(1,"Vodafone Turkiye","VF-TR","28602",2),(1,"TR TURKCELL","TCELL","28601",7),(1,"TR TURKCELL","TCELL","28601",2),,(0-4),(0-2) OK]`},
				}

				for _, tt := range parameters {

					t.Run(tt.command, func(t *testing.T) {

						commandName := "Read"
						mockedDefault := serialShell.mocked[commandName]

						serialShell.Patch(commandName, tt.mock, nil)
						defer func() { serialShell.mocked[commandName] = mockedDefault }()

						at := NewAtcom(serialShell, nil, mockSleepFunc)
						response, _ := at.SendAT(tt.command, nil)

						if !reflect.DeepEqual(response, tt.desired_response) {
							t.Errorf("Expected %v, but got %v", tt.desired_response, response)
						}
					})
				}

			})

		t.Run("Should return timeout", func(t *testing.T) {

			commandName := "Read"
			mockedDefault := serialShell.mocked[commandName]

			serialShell.Patch(commandName, 316, nil)
			defer func() { serialShell.mocked[commandName] = mockedDefault }()

			at := NewAtcom(serialShell, nil, mockSleepFunc)
			_, err := at.SendAT("AT+COPS=?", nil)

			error := errors.New("timeout")
			if err.Error() != error.Error() {
				t.Errorf("Expected error %v, but got %v", error, err)
			}
		})*/

}
