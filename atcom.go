/*
Created by: Yasin Kaya (selengalp), yasinkaya.121@gmail.com, 2023

Copyright (c) 2023 Sixfab Inc.
*/
package atcom

import (
	"bufio"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tarm/serial"
)

func open(args map[string]interface{}) (port *serial.Port, err error) {

	portname := "/dev/ttyUSB2"
	baudrate := 115200

	for key, value := range args {
		switch key {
		case "port":
			portname = value.(string)
		case "baud":
			baudrate = value.(int)
		}
	}

	config := &serial.Config{
		Name:        portname,
		Baud:        baudrate,
		ReadTimeout: time.Millisecond * 100,
	}

	return serial.OpenPort(config)
}

func SendAT(command string, args map[string]interface{}) ([]string, error) {

	var lineEnd bool = true
	var desired []string
	var fault []string
	var timeout int

	for key, value := range args {
		switch key {
		case "desired":
			desired = value.([]string)
		case "fault":
			fault = value.([]string)
		case "timeout":
			timeout = value.(int)
		case "lineEnd":
			lineEnd = value.(bool)
		}
	}

	serialPort, err := open(args)

	if err != nil {
		return nil, err
	}

	defer serialPort.Close()

	if lineEnd {
		command += "\r\n"
	}

	_, err = serialPort.Write([]byte(command))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(serialPort)
	response := make([]string, 0)
	timeoutDuration := time.Duration(timeout) * time.Second

	found := make(chan error)
	defer close(found)

	ctxScan, cancelScan := context.WithCancel(context.Background())
	defer cancelScan()

	go func(ctx context.Context) {
		for {
			if !scanner.Scan() {
				return
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\r")
			line = strings.Trim(line, "\n")

			if line != "" {
				response = append(response, line)
			}

			if line == "OK" {
				// make response string to check desired and fault
				data := ""
				for _, word := range response {
					if word != "OK" {
						data += word
					}
				}

				// check desired and fault existed in response
				if desired != nil || fault != nil {
					for _, desiredStr := range desired {
						if strings.Contains(data, desiredStr) {
							found <- nil
							return
						}
					}

					for _, faultStr := range fault {
						if strings.Contains(data, faultStr) {
							found <- nil
							return
						}
					}
				} else {
					return
				}
			} else if line == "ERROR" || strings.Contains(line, "+CME ERROR") {
				found <- errors.New("modem error")
				return
			}
		}
	}(ctxScan)

	timeoutCh := time.After(timeoutDuration)

	for {
		select {
		case err := <-found:
			return response, err
		case <-timeoutCh:
			cancelScan()
			return response, errors.New("timeout")
		}
	}
}
