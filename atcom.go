/*
Created by: Yasin Kaya (selengalp), yasinkaya.121@gmail.com, 2023

Copyright (c) 2023 Sixfab Inc.
*/
package atcom

import (
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
	var desired []string = nil
	var fault []string = nil
	var timeout int = 5

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

	data := make([]string, 0)
	timeoutDuration := time.Duration(timeout) * time.Second

	found := make(chan error)
	defer close(found)

	ctxScan, cancelScan := context.WithCancel(context.Background())
	defer cancelScan()

	go func(ctx context.Context) {
		response := ""
		buf := make([]byte, 1024)

		for {
			time.Sleep(time.Millisecond * 5)
			n, err := serialPort.Read(buf)
			if err != nil {
				if err.Error() == "EOF" {
					continue
				}
				found <- err
				return
			}
			if n > 0 {
				response += string(buf[:n])
			}

			if strings.Contains(response, "\r\nOK\r\n") {
				lines := strings.Split(response, "\r\n")

				for _, line := range lines {
					line = strings.TrimSpace(line)
					line = strings.Trim(line, "\r")
					line = strings.Trim(line, "\n")

					if line != "" {
						data = append(data, line)
					}

					if line == "OK" {
						break
					}
				}

				// check desired and fault existed in response
				if desired != nil || fault != nil {
					for _, desiredStr := range desired {
						if strings.Contains(response, desiredStr) {
							found <- nil
							return
						}

						for _, faultStr := range fault {
							if strings.Contains(response, faultStr) {
								found <- errors.New("faulty response detected")
								return
							}
						}
					}
				} else {
					found <- nil
					return
				}

				found <- nil
				return
			} else if strings.Contains(response, "\r\nERROR\r\n") {
				found <- errors.New("modem error")
				return
			}
		}
	}(ctxScan)

	timeoutCh := time.After(timeoutDuration)

	for {
		select {
		case err := <-found:
			return data, err
		case <-timeoutCh:
			cancelScan()
			return data, errors.New("timeout")
		}
	}
}
