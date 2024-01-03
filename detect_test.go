package atcom

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var (
	nilError         error = nil
	mockDeviceOutput       = `
	Bus 001 Device 002: ID 2c7c:0125 Quectel Wireless Solutions Co., Ltd. EC25 LTE modem
	Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
	Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
	Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
	`
	mockUdevadmOutput = []string{
		"DEVNAME='/dev/ttyUSB2'",
		"MAJOR='188'",
		"MINOR='2'",
		"SUBSYSTEM='tty'",
		"USEC_INITIALIZED='166762587'",
		"ID_BUS='usb'",
		"ID_MODEL='EG25-G'",
		"ID_MODEL_ENC='EG25-G'",
		"ID_MODEL_ID='0125'",
		"ID_SERIAL='Quectel_EG25-G'",
		"ID_VENDOR='Quectel'",
		"ID_VENDOR_ENC='Quectel'",
		"ID_VENDOR_ID='2c7c'",
		"ID_REVISION='0318'",
		"ID_TYPE='generic'",
		"ID_USB_MODEL='EG25-G'",
		"ID_USB_MODEL_ENC='EG25-G'",
		"ID_USB_MODEL_ID='0125'",
		"ID_USB_SERIAL='Quectel_EG25-G'",
		"ID_USB_VENDOR='Quectel'",
		"ID_USB_VENDOR_ENC='Quectel'",
		"ID_USB_VENDOR_ID='2c7c'",
		"ID_USB_REVISION='0318'",
		"ID_USB_TYPE='generic'",
		"ID_USB_INTERFACES=':ffffff:ff0000:020600:0a0000:'",
		"ID_USB_INTERFACE_NUM='02'",
		"ID_USB_DRIVER='option'",
		"ID_USB_CLASS_FROM_DATABASE='Miscellaneous Device'",
		"ID_USB_PROTOCOL_FROM_DATABASE='Interface Association'",
		"ID_VENDOR_FROM_DATABASE='Quectel Wireless Solutions Co., Ltd.'",
		"ID_MODEL_FROM_DATABASE='EC25 LTE modem'",
	}
)

type MockShell struct{}

func (t *MockShell) Command(name string, arg ...string) (string, error) {

	if name == "lsusb" {
		mockOutput := mockDeviceOutput

		return mockOutput, nilError
	}

	if len(arg) < 2 {
		return "", errors.New("MockShell.Command: Length of arg is less than 2")
	}

	if name == "bash" && arg[0] == "-c" && arg[1] == "/usr/bin/find /sys/bus/usb/devices/usb*/ -name dev" {

		mockOutput := []string{
			"/sys/bus/usb/devices/usb1/dev",
			"/sys/bus/usb/devices/usb1/1-2/1-2:1.2/ttyUSB2/tty/ttyUSB2/dev",
			"/sys/bus/usb/devices/usb1/1-2/1-2:1.0/ttyUSB0/tty/ttyUSB0/dev",
			"/sys/bus/usb/devices/usb1/1-2/1-2:1.3/ttyUSB3/tty/ttyUSB3/dev",
			"/sys/bus/usb/devices/usb1/1-2/1-2:1.1/ttyUSB1/tty/ttyUSB1/dev",
		}

		return strings.Join(mockOutput, "\n"), nilError
	}

	if name == "bash" &&
		arg[0] == "-c" &&
		arg[1] == "udevadm info -q property --export -p /sys/bus/usb/devices/usb1/1-2/1-2:1.2/ttyUSB2/tty/ttyUSB2" {

		mockOutput := mockUdevadmOutput

		return strings.Join(mockOutput, "\n"), nilError
	}
	return "", nilError
}

func TestGetAvailablePorts(t *testing.T) {

	t.Run("Should return available ports", func(t *testing.T) {

		at := NewAtcom(nil, &MockShell{})
		availablePorts, err := at.getAvailablePorts()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		ok := strings.Contains(availablePorts[1]["port"], "/dev/ttyUSB2")

		if !ok {
			t.Error("Should return available ports: Failed")
		}
	})

	t.Run("Should return error for device", func(t *testing.T) {

		nilError = errors.New("device error")

		defer func() { nilError = nil }()

		at := NewAtcom(nil, &MockShell{})
		_, err := at.getAvailablePorts()

		if err.Error() != nilError.Error() {
			t.Errorf("Expected error %v, but got %v", nilError, err)
		}
	})

	t.Run("Should return error for udevadm", func(t *testing.T) {

		nilError = errors.New("udevadm error")
		defer func() { nilError = nil }()

		at := NewAtcom(nil, &MockShell{})
		_, err := at.getAvailablePorts()

		if err.Error() != nilError.Error() {
			t.Errorf("Expected error %v, but got %v", nilError, err)
		}
	})
}

func TestFindModem(t *testing.T) {

	t.Run("Should return error for lsusb", func(t *testing.T) {

		nilError = errors.New("lsusb error")
		defer func() { nilError = nil }()

		at := NewAtcom(nil, &MockShell{})
		_, err := at.findModem(supportedModems)

		if err.Error() != nilError.Error() {
			t.Errorf("Expected error %v, but got %v", nilError, err)
		}
	})

	t.Run("Should return modem", func(t *testing.T) {

		var tests = []struct {
			device string
			name   string
			want   SupportedModem
		}{
			{"Bus 001 Device 002: ID 2c7c:0125 Quectel Wireless Solutions Co., Ltd. EC25 LTE modem", "Quectel EC25", SupportedModem{"2c7c", "0125", "Quectel", "EC25", "if02"}},
			{"Bus 001 Device 002: ID 2c7c:0296 Quectel Wireless Solutions Co., Ltd. BG96 LTE modem", "Quectel BG96", SupportedModem{"2c7c", "0296", "Quectel", "BG96", "if02"}},
			{"Bus 001 Device 002: ID 1bc7:1201 Telit Wireless Solutions Co., Ltd. LE910Cx RMNET LTE modem", "Telit LE910Cx RMNET", SupportedModem{"1bc7", "1201", "Telit", "LE910Cx RMNET", "if04"}},
			{"Bus 001 Device 002: ID 1e2d:0069 Thales/Cinterion Wireless Solutions Co., Ltd. PLSx3 LTE modem", "Thales/Cinterion PLSx3", SupportedModem{"1e2d", "0069", "Thales/Cinterion", "PLSx3", "if04"}},
		}

		for _, tt := range tests {

			existingDeviceOutput := mockDeviceOutput
			mockDeviceOutput = tt.device
			defer func() { mockDeviceOutput = existingDeviceOutput }()

			at := NewAtcom(nil, &MockShell{})

			t.Run(tt.name, func(t *testing.T) {

				modem, err := at.findModem(supportedModems)
				if modem != tt.want {
					t.Errorf("Expected %s, but got %s", modem, tt.want)
				}

				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
			})
		}

	})

	t.Run("Should return no supported modem", func(t *testing.T) {

		existingDeviceOutput := mockDeviceOutput
		mockDeviceOutput = `
		Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
		Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		`
		defer func() { mockDeviceOutput = existingDeviceOutput }()

		at := NewAtcom(nil, &MockShell{})

		modem, err := at.findModem(supportedModems)

		if !reflect.DeepEqual(modem, SupportedModem{}) {
			t.Errorf("Expected %s, but got %s", SupportedModem{}, modem)
		}

		expected_error := errors.New("no supported modem found")

		if err.Error() != expected_error.Error() {
			t.Errorf("Expected %v, but got %v", expected_error, err)
		}
	})

}

func TestDecidePort(t *testing.T) {

	t.Run("Should return no supported modem", func(t *testing.T) {

		existingDeviceOutput := mockDeviceOutput
		mockDeviceOutput = `
		Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
		Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		`
		defer func() { mockDeviceOutput = existingDeviceOutput }()

		at := NewAtcom(nil, &MockShell{})
		_, err := at.DecidePort()

		expected_error := errors.New("no supported modem found")

		if err.Error() != expected_error.Error() {
			t.Errorf("Expected %v, but got %v", expected_error, err)
		}
	})

	t.Run("Should return no supported modem", func(t *testing.T) {

		existingDeviceOutput := mockDeviceOutput
		mockDeviceOutput = `
		Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
		Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
		`
		defer func() { mockDeviceOutput = existingDeviceOutput }()

		at := NewAtcom(nil, &MockShell{})
		_, err := at.DecidePort()

		expected_error := errors.New("no supported modem found")

		if err.Error() != expected_error.Error() {
			t.Errorf("Expected %v, but got %v", expected_error, err)
		}
	})
	t.Run("Should return detected modem", func(t *testing.T) {

		at := NewAtcom(nil, &MockShell{})
		detectedModem, err := at.DecidePort()

		expectedDetectedModem := map[string]string{
			"port":   "/dev/ttyUSB2",
			"vid":    "2c7c",
			"pid":    "0125",
			"vendor": "Quectel",
			"model":  "EG25-G",
		}

		if !reflect.DeepEqual(expectedDetectedModem, detectedModem) {
			t.Errorf("Expected %s, but got %s", expectedDetectedModem, detectedModem)
		}
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	})

	t.Run("Should return nil", func(t *testing.T) {

		existingmockUdevadmOutput := mockUdevadmOutput
		mockUdevadmOutput = []string{}

		defer func() { mockUdevadmOutput = existingmockUdevadmOutput }()

		at := NewAtcom(nil, &MockShell{})
		detectedModem, err := at.DecidePort()

		if detectedModem != nil {
			t.Errorf("Expected nil, but got %s", detectedModem)
		}
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	})
}
