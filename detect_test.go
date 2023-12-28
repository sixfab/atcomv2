package atcom

import (
	"errors"
	"strings"
	"testing"
)

type MockShell struct{}

func (t *MockShell) Command(name string, arg ...string) (string, error) {

	if name == "lsusb" {
		mockOutput := `
			Bus 001 Device 002: ID 2c7c:0125 Quectel Wireless Solutions Co., Ltd. EC25 LTE modem
			Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
			Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
			Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
			`

		return mockOutput, nil
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

		return strings.Join(mockOutput, "\n"), nil
	}

	if name == "bash" &&
		arg[0] == "-c" &&
		arg[1] == "udevadm info -q property --export -p /sys/bus/usb/devices/usb1/1-2/1-2:1.2/ttyUSB2/tty/ttyUSB2" {

		mockOutput := []string{
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

		return strings.Join(mockOutput, "\n"), nil
	}
	return "", nil
}

func TestGetAvailablePorts(t *testing.T) {
	at := NewAtcom(nil, &MockShell{})

	availablePorts, err := at.getAvailablePorts()

	if err != nil {
		t.Error(err)
	}

	ok := strings.Contains(availablePorts[1]["port"], "/dev/ttyUSB2")

	if !ok {
		t.Error("GetAvailablePorts: Failed")
	}
}
