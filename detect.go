package atcom

import (
	"errors"
	"strings"
)

func (t *Atcom) getAvailablePorts() (availablePorts []map[string]string, err error) {
	output, err := t.shell.Command("bash", "-c", "/usr/bin/find /sys/bus/usb/devices/usb*/ -name dev")

	if err != nil {
		return nil, err
	}

	ports := make([]string, 0)

	for _, port := range strings.Split(output, "\n") {
		if strings.HasSuffix(port, "/dev") {
			port = strings.TrimSuffix(port, "/dev")
			ports = append(ports, port)
		}
	}

	for _, port := range ports {
		output, err := t.shell.Command("bash", "-c", "udevadm info -q property --export -p "+port)

		if err != nil {
			return nil, err
		}

		deviceDetails := strings.Split(output, "\n")

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

func (t *Atcom) findModem(smodems []SupportedModem) (SupportedModem, error) {
	output, err := t.shell.Command("lsusb")

	if err != nil {
		return SupportedModem{}, err
	}

	for _, modem := range smodems {
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, modem.vid) && strings.Contains(line, modem.pid) {
				return modem, nil
			}
		}
	}

	return SupportedModem{}, errors.New("no supported modem found")
}

func (t *Atcom) DecidePort() (map[string]string, error) {
	modem, err := t.findModem(supportedModems)

	if err != nil {
		return nil, err
	}

	ports, err := t.getAvailablePorts()

	if err != nil {
		return nil, err
	}

	attr := DefaultSerialAttr()

	for _, port := range ports {
		if port["vendor_id"] == modem.vid &&
			port["product_id"] == modem.pid &&
			port["interface"] == modem.ifs {

			// set port and baudrate on atcom instance
			attr.Port = port["port"]
			t.SerialAttr.Port = attr.Port
			t.SerialAttr.Baud = attr.Baud

			detectedModem := map[string]string{
				"port":   port["port"],
				"vid":    port["vendor_id"],
				"pid":    port["product_id"],
				"vendor": port["vendor"],
				"model":  port["model"],
			}

			return detectedModem, nil
		}
	}

	return nil, nil
}
