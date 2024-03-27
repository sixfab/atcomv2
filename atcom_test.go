package atcom

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/tarm/serial"
)

type MockSerial struct {
	mocked map[string]interface{}
}

var (
	mockSerialPort = &serial.Port{}
	mockWrite      = 2
	mockSleepFunc  = func(d time.Duration) {}
	mockRead       = 2
	mockBuffer     = ""
)

var mockSerial = &MockSerial{
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

func (t *MockSerial) Patch(cmd string, resp interface{}, err error) {
	t.mocked[cmd] = map[string]interface{}{"resp": resp, "err": err}
}

func (t *MockSerial) OpenPort(c *serial.Config) (*serial.Port, error) {

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

func (t *MockSerial) Write(port *serial.Port, command []byte) (n int, err error) {
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

func (t *MockSerial) Close(port *serial.Port) (err error) {
	return nil
}

func (t *MockSerial) Read(port *serial.Port, buf []byte) (n int, err error) {

	for mocked_name := range t.mocked {

		if mocked_name == "Read" {
			response, _ := t.mocked[mocked_name].(map[string]interface{})

			copy(buf, []byte(mockBuffer))

			if response["err"] == nil {
				return response["resp"].(int), nil
			}
			return response["resp"].(int), response["err"].(error)
		}

	}
	return 5, nil
}

func TestNewAtcom(t *testing.T) {

	t.Run("Should return pointer of atcom structure", func(t *testing.T) {
	})
}

func TestOpen(t *testing.T) {

	t.Run("Should return port with nil input parameter", func(t *testing.T) {

		at := NewAtcom(mockSerial, nil, nil)
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

		at := NewAtcom(mockSerial, nil, nil)

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
		mockedDefault := mockSerial.mocked[commandName]
		error := errors.New("Serial port error")

		mockSerial.Patch(commandName, mockSerialPort, error)
		defer func() { mockSerial.mocked[commandName] = mockedDefault }()

		at := NewAtcom(mockSerial, nil, nil)
		_, err := at.SendAT("ATE1", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error for Write function", func(t *testing.T) {

		commandName := "Write"
		mockedDefault := mockSerial.mocked[commandName]
		error := errors.New("Write function error")

		mockSerial.Patch(commandName, mockWrite, error)
		defer func() { mockSerial.mocked[commandName] = mockedDefault }()

		at := NewAtcom(mockSerial, nil, nil)
		_, err := at.SendAT("ATE1", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := mockSerial.mocked[commandName]
		error := errors.New("modem error")

		mockSerial.Patch(commandName, mockRead, error)
		defer func() { mockSerial.mocked[commandName] = mockedDefault }()

		at := NewAtcom(mockSerial, nil, mockSleepFunc)
		_, err := at.SendAT("AT#ECMD=0", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return successfull response", func(t *testing.T) {

		parameters := []struct {
			command          string
			mock_response    int
			mock_buffer      string
			desired_response []string
		}{
			{"ATE0", 6, "\r\nOK\r\n", []string{"OK"}},
			{"ATE1", 11, "ATE1\r\r\nOK\r\n", []string{"ATE1", "OK"}},
			{"AT+COPS?", 27, "AT+COPS?\r\r\n+COPS: 0\r\n\r\nOK\r\n", []string{"AT+COPS?", "+COPS: 0", "OK"}},
		}

		for _, tt := range parameters {

			t.Run(tt.command, func(t *testing.T) {

				commandName := "Read"
				mockedDefault := mockSerial.mocked[commandName]
				mockSerial.Patch(commandName, tt.mock_response, nil)

				mockBuffer = tt.mock_buffer
				defer func() {
					mockSerial.mocked[commandName] = mockedDefault
					mockBuffer = ""
				}()

				at := NewAtcom(mockSerial, nil, mockSleepFunc)

				response, _ := at.SendAT(tt.command, nil)

				if !reflect.DeepEqual(response, tt.desired_response) {
					t.Errorf("Expected %s, but got %s", tt.desired_response, response)
				}
			})
		}

	})
	t.Run("Should return timeout", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := mockSerial.mocked[commandName]

		mockSerial.Patch(commandName, 316, nil)
		defer func() { mockSerial.mocked[commandName] = mockedDefault }()

		at := NewAtcom(mockSerial, nil, mockSleepFunc)

		args := map[string]interface{}{
			"timeout": 1,
		}

		_, err := at.SendAT("AT+WRONGATCOMMAND", args)

		error := errors.New("timeout")
		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return response ", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := mockSerial.mocked[commandName]

		error := errors.New("timeout")
		mockSerial.Patch(commandName, 0, error)
		defer func() { mockSerial.mocked[commandName] = mockedDefault }()

		at := NewAtcom(mockSerial, nil, mockSleepFunc)
		_, err := at.SendAT("AT+COPS=?", nil)

		if err.Error() != error.Error() {
			t.Errorf("Expected error %v, but got %v", error, err)
		}
	})

	t.Run("Should return response for desired argument", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := mockSerial.mocked[commandName]

		mockSerial.Patch(commandName, 27, nil)

		mockBuffer = "AT+COPS?\r\r\n+COPS: 0\r\n\r\nOK\r\n"

		defer func() {
			mockSerial.mocked[commandName] = mockedDefault
			mockBuffer = ""
		}()

		at := NewAtcom(mockSerial, nil, mockSleepFunc)
		arg := map[string]interface{}{"desired": []string{"+COPS"}}
		res, err := at.SendAT("AT+COPS?", arg)

		if err != nil {
			t.Errorf("Expected nil error, but got %v", err)
		}

		expectedResult := []string{"AT+COPS?", "+COPS: 0", "OK"}
		if !reflect.DeepEqual(res, expectedResult) {
			t.Errorf("Expected %s, but got %s", expectedResult, res)
		}
	})

	t.Run("Should return response for fault argument", func(t *testing.T) {

		commandName := "Read"
		mockedDefault := mockSerial.mocked[commandName]

		mockSerial.Patch(commandName, 27, nil)

		mockBuffer = "AT+COPS?\r\r\n+COPS: 0\r\n\r\nOK\r\n"

		defer func() {
			mockSerial.mocked[commandName] = mockedDefault
			mockBuffer = ""
		}()

		at := NewAtcom(mockSerial, nil, mockSleepFunc)
		arg := map[string]interface{}{"fault": []string{"+COPS"}}
		res, err := at.SendAT("AT+COPS?", arg)

		if err != nil {
			t.Errorf("Expected nil error, but got %v", err)
		}

		expectedResult := []string{"AT+COPS?", "+COPS: 0", "OK"}
		if !reflect.DeepEqual(res, expectedResult) {
			t.Errorf("Expected %s, but got %s", expectedResult, res)
		}
	})

}
