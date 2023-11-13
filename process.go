package atcom

import (
	"errors"
	"strings"
)

func GetMeaningfulPart(response []string, command string, prefix string) (data []string, len int, err error) {

	var echoRow int = 0
	var lastRow int = 0

	for index, line := range response {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, command) {
			echoRow = index
			continue
		}

		if line == "OK" {
			lastRow = index
			break
		}
	}

	if lastRow == 0 {
		return nil, 0, errors.New("no ok response")
	}

	if prefix == "" {
		data = response[echoRow+1 : lastRow]
		len = lastRow - echoRow - 1
	} else {
		for _, line := range response[echoRow+1 : lastRow] {
			if strings.HasPrefix(line, prefix) {
				line = strings.Trim(line, prefix)
				data = append(data, line)
				len++
			}
		}
	}

	return data, len, nil
}
