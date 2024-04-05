package atcom

import (
	"fmt"
	"strings"
)

type SerialAttr struct {
	Port string
	Baud int
}

func DefaultSerialAttr() SerialAttr {
	return SerialAttr{
		Port: "",
		Baud: 115200,
	}
}

type ATCommand struct {
	SerialAttr SerialAttr

	Command   string
	Response  []string
	Processed []string
	Error     error

	Desired []string
	Fault   []string
	Timeout int
	LineEnd bool
}

func NewATCommand(command string) *ATCommand {
	return &ATCommand{
		SerialAttr: DefaultSerialAttr(),
		Command:    command,
		Desired:    nil,
		Fault:      nil,
		Timeout:    5,
		LineEnd:    true,
	}
}

func (atc *ATCommand) GetMeaningfulPart(prefix string) error {

	if atc.Response == nil {
		return fmt.Errorf("no response")
	}

	if prefix == "" {
		atc.Processed = atc.Response
		return nil
	}

	var echoRow int = 0
	var lastRow int = 0
	var data []string

	for index, line := range atc.Response {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, atc.Command) {
			echoRow = index
			continue
		}

		if line == "OK" {
			lastRow = index
			break
		}
	}

	if lastRow == 0 {
		return fmt.Errorf("no ok response")
	}

	if prefix == "" {
		data = atc.Response[echoRow+1 : lastRow]
	} else {
		for _, line := range atc.Response[echoRow+1 : lastRow] {
			if strings.HasPrefix(line, prefix) {
				line = strings.Trim(line, prefix)
				data = append(data, line)
			}
		}
	}

	atc.Processed = data
	return nil
}
