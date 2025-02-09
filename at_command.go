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

	Command      string
	Response     []string
	Processed    []string
	Error        error
	ResponseChan chan string

	Desired []string
	Fault   []string
	Timeout int
	LineEnd bool
}

func NewATCommand(command string) *ATCommand {
	return &ATCommand{
		SerialAttr:   DefaultSerialAttr(),
		Command:      command,
		Desired:      nil,
		Fault:        nil,
		Timeout:      5,
		LineEnd:      true,
		ResponseChan: nil,
	}
}

func (atc *ATCommand) GetMeaningfulPart(prefix string) error {

	if atc.Response == nil {
		return fmt.Errorf("no response")
	}

	var firstRow int = 0
	var lastRow int = 0
	var data []string

	// decide start and end of meaningful part of response
	for index, line := range atc.Response {
		line = strings.TrimSpace(line)

		// decide echo line if exists
		if strings.HasPrefix(line, atc.Command) {
			firstRow = index + 1
			continue
		}

		// skip modem +X lines when prefix is empty
		if prefix == "" {
			if strings.HasPrefix(line, "+") {
				firstRow++
				continue
			}
		}

		// decide last line of response
		if line == "OK" {
			lastRow = index
			break
		}
	}

	if lastRow == 0 {
		return fmt.Errorf("no ok response")
	}

	if prefix == "" {
		data = atc.Response[firstRow:lastRow]
	} else {
		for _, line := range atc.Response[firstRow:lastRow] {
			if strings.HasPrefix(line, prefix) {
				line = strings.Trim(line, prefix)
				data = append(data, line)
			}
		}
	}

	atc.Processed = data
	return nil
}
