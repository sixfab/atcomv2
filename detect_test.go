package atcom

import (
	"os/exec"
	"testing"
)

func TestGetAvailablePorts(t *testing.T) {
	t.Run("Command error", func(t *testing.T) {
		// Test case for command failure
		_, err := getAvailablePorts()
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	t.Run("Port check", func(t *testing.T) {
		// Test case for port checking
		mockUsbPath := "/sys/bus/usb/devices/usb1/dev /sys/bus/usb/devices/usb2/dev"
		execCommand := exec.Command("mkdir", mockUsbPath)
		execCommand.Run()

		_, err := getAvailablePorts()
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		execCommand = exec.Command("rmdir", mockUsbPath)
		execCommand.Run()
	})
	/*
			t.Run("Valid ports", func(t *testing.T) {
				// Test case where everything works correctly
				mockUsbPath := "/sys/bus/usb/devices/usb1/dev /sys/bus/usb/devices/usb2/dev"
				execCommand := exec.Command("mkdir", mockUsbPath)
				execCommand.Run()

				mockUdevadmOutput := `DEVNAME='/dev/some_device'
		ID_VENDOR='Vendor'
		ID_VENDOR_ID='1234'
		ID_MODEL='Model'
		ID_MODEL_FROM_DATABASE='ModelDB'
		ID_MODEL_ID='5678'
		ID_USB_INTERFACE_NUM='if01'
		ID_USB_VENDOR_ID='abcd'
		ID_USB_MODEL_ID='efgh'`

				execCommand = exec.Command("echo", mockUdevadmOutput)
				execCommand.Run()

				expectedResult := []map[string]string{
					{
						"port":                "/dev/some_device",
						"vendor":              "Vendor",
						"vendor_id":           "1234",
						"model":               "Model",
						"model_from_database": "ModelDB",
						"product_id":          "5678",
						"interface":           "if01",
						"ID_USB_VENDOR_ID":    "abcd",
						"ID_USB_MODEL_ID":     "efgh",
					},
				}

				result, err := getAvailablePorts()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(result, expectedResult) {
					t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expectedResult)
				}
				execCommand = exec.Command("rmdir", mockUsbPath)
				execCommand.Run()
			})*/

}
