package atcom

import (
	"os/exec"
	"strings"
)

func getAvailablePorts() (availablePorts []map[string]string, err error) {
	cmd := exec.Command("bash", "-c", "/usr/bin/find /sys/bus/usb/devices/usb*/ -name dev")
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	outputStr := string(output)

	ports := make([]string, 0)

	for _, port := range strings.Split(outputStr, "\n") {
		if strings.HasSuffix(port, "/dev") {
			port = strings.TrimSuffix(port, "/dev")
			ports = append(ports, port)
		}
	}

	for _, port := range ports {
		cmd := exec.Command("bash", "-c", "udevadm info -q property --export -p "+port)
		output, err := cmd.Output()

		if err != nil {
			return nil, err
		}

		outputStr := string(output)

		deviceDetails := strings.Split(outputStr, "\n")

		portDetails := make(map[string]string)

		for _, line := range deviceDetails {
			switch {
			case strings.HasPrefix(line, "DEVNAME="):
				portDetails["port"] = strings.Trim(line[8:], "'")
			case strings.HasPrefix(line, "ID_VENDOR="):
				portDetails["vendor"] = strings.Trim(line[10:], "'")
			case strings.HasPrefix(line, "ID_VENDOR_ID="):
				portDetails["vendor_id"] = strings.Trim(line[13:], "'")
			case strings.HasPrefix(line, "ID_MODEL="):
				portDetails["model"] = strings.Trim(line[9:], "'")
			case strings.HasPrefix(line, "ID_MODEL_FROM_DATABASE="):
				portDetails["model_from_database"] = strings.Trim(line[23:], "'")
			case strings.HasPrefix(line, "ID_MODEL_ID="):
				portDetails["product_id"] = strings.Trim(line[12:], "'")
			case strings.HasPrefix(line, "ID_USB_INTERFACE_NUM="):
				portDetails["interface"] = "if" + strings.Trim(line[21:], "'")
			case strings.HasPrefix(line, "ID_USB_VENDOR_ID="):
				portDetails["ID_USB_VENDOR_ID"] = strings.Trim(line[17:], "'")
			case strings.HasPrefix(line, "ID_USB_MODEL_ID="):
				portDetails["ID_USB_MODEL_ID"] = strings.Trim(line[16:], "'")
			}
		}

		if !strings.Contains(portDetails["port"], "bus") {
			availablePorts = append(availablePorts, portDetails)
		}
	}
	return availablePorts, nil
}

func findModems() (SupportedModem, error) {
	cmd := exec.Command("lsusb")
	output, err := cmd.Output()

	if err != nil {
		return SupportedModem{}, err
	}

	outputStr := string(output)

	for _, modem := range supportedModems {
		if outputStr == "" {
			return modem, nil
		}

		for _, line := range strings.Split(outputStr, "\n") {
			if strings.Contains(line, modem.vid) && strings.Contains(line, modem.pid) {
				return modem, nil
			}
		}
	}
	return SupportedModem{}, nil
}

func DecidePort() (map[string]string, error) {
	modem, err := findModems()

	if err != nil {
		return nil, err
	}

	ports, err := getAvailablePorts()

	if err != nil {
		return nil, err
	}

	for _, port := range ports {
		if port["vendor_id"] == modem.vid &&
			port["product_id"] == modem.pid &&
			port["interface"] == modem.ifs {

			detectedModem := map[string]string{
				"port":       port["port"],
				"vendor_id":  port["vendor_id"],
				"product_id": port["product_id"],
			}

			return detectedModem, nil
		}
	}

	return nil, nil
}
