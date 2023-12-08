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
		result, err := open(args)

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

func TestSendAT(t *testing.T) {

	command := "command"
	args := map[string]interface{}{
		"desired": nil,
		"fault":   nil,
		"timeout": 5,
		"end":     true,
	}

	t.Run("should return error on calling open function", func(t *testing.T) {

		_, result := SendAT(command, args)
		_, expectedResult := open(args)

		if result != expectedResult {
			t.Errorf("Expected %s, got %s", expectedResult, result)
		}
	})

	args = map[string]interface{}{
		"desired": nil,
		"fault":   nil,
		"timeout": 5,
		"end":     true,
		"port":    "/dev/ttyUSB2",
		"baud":    115200,
	}

	t.Run("should return error on calling serialPort.Write function", func(t *testing.T) {

		_, result := SendAT(command, args)

		serialPort, _ := open(args)
		_, expectedResult := serialPort.Write([]byte(command))

		if result != expectedResult {
			t.Errorf("Expected %s, got %s", expectedResult, result)
		}
	})
}
