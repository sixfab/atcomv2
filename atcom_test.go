package atcom

import (
	"testing"
	"time"

	"github.com/tarm/serial"
)

func TestOpen(t *testing.T) {

	args := map[string]interface{}{
		"port": "/dev/ttyUSB2",
		"baud": 115200,
	}

	t.Run("open function success", func(t *testing.T) {
		result, _ := open(args)

		config := &serial.Config{
			Name:        "/dev/ttyUSB2",
			Baud:        115200,
			ReadTimeout: time.Millisecond * 100,
		}

		expectedResult, err := serial.OpenPort(config)

		if result != expectedResult {
			t.Errorf(" %s error message", err)
		}
	})
}

/*
// MockSerialPort is a mock implementation of the serial.Port interface
type MockSerialPort struct {
	mock.Mock
}

func (m *MockSerialPort) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockSerialPort) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockSerialPort) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestSendAT(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		args           map[string]interface{}
		expectedResp   []string
		expectedErrStr string
	}{
		{
			name:         "Normal case with expected response",
			command:      "AT",
			args:         map[string]interface{}{"desired": []string{"expected_response"}, "fault": []string{"fault_response"}, "timeout": 5, "lineEnd": true},
			expectedResp: []string{"expected_response"},
		},
		{
			name:           "Timeout case",
			command:        "AT",
			args:           map[string]interface{}{"desired": nil, "fault": nil, "timeout": 5, "lineEnd": true},
			expectedResp:   []string{},
			expectedErrStr: "timeout",
		},
		{
			name:           "Faulty response case",
			command:        "AT",
			args:           map[string]interface{}{"desired": []string{"expected_response"}, "fault": []string{"fault_response"}, "timeout": 5, "lineEnd": true},
			expectedResp:   []string{"fault_response"},
			expectedErrStr: "faulty response detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSerialPort := new(MockSerialPort)
			mockSerialPort.On("Write", []byte(tt.command+"\r\n")).Return(3, nil)
			mockSerialPort.On("Read", mock.Anything).Return(0, nil).Times(3) // Simulating 3 reads
			mockSerialPort.On("Close").Return(nil)

			response, err := SendAT(tt.command, tt.args)

			assert.Equal(t, tt.expectedResp, response, "Response does not match")

			if tt.expectedErrStr != "" {
				assert.NotNil(t, err, "Expected an error")
				assert.EqualError(t, err, tt.expectedErrStr, "Error message does not match")
			} else {
				assert.Nil(t, err, "Expected no error")
			}
		})
	}
}*/
