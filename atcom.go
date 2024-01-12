/*
Created by: Yasin Kaya (selengalp), yasinkaya.121@gmail.com, 2023

Copyright (c) 2023 Sixfab Inc.
*/
package atcom

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type Atcom struct {
	serial Serial
	shell  Shell
}

// Serial Implementation for normal usage
type RealSerial struct {
}

// Serial interface
type Serial interface {
	OpenPort(c *serial.Config) (*serial.Port, error)
	Write(port *serial.Port, command []byte) (n int, err error)
	Close(port *serial.Port) (err error)
	Read(port *serial.Port, buffer []byte) (n int, err error)
}

// RealSerial implements Serial interface
func (s *RealSerial) OpenPort(c *serial.Config) (*serial.Port, error) {
	return serial.OpenPort(c)
}

func (s *RealSerial) Write(port *serial.Port, command []byte) (n int, err error) {
	return port.Write([]byte(command))
}

func (s *RealSerial) Close(port *serial.Port) (err error) {
	return port.Close()
}

func (s *RealSerial) Read(port *serial.Port, buffer []byte) (n int, err error) {
	return port.Read(buffer)
}

// Shell Implementation for normal usage
type RealShell struct{}

// Shell interface
type Shell interface {
	Command(name string, arg ...string) (string, error)
}

// RealShell implements Shell interface
func (s *RealShell) Command(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	return string(output), err
}

// NewAtcom creates a new Atcom instance with default serial and shell implementations
func NewAtcom(s Serial, sh Shell) *Atcom {

	if s == nil {
		s = &RealSerial{}
	}

	if sh == nil {
		sh = &RealShell{}
	}

	return &Atcom{
		serial: s,
		shell:  sh,
	}
}

// Function to open serial port
func (t *Atcom) open(args map[string]interface{}) (port *serial.Port, err error) {

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

	return t.serial.OpenPort(config)
}

// SendAT sends AT command to modem and returns response
func (t *Atcom) SendAT(command string, args map[string]interface{}) ([]string, error) {

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

	serialPort, err := t.open(args)

	if err != nil {
		return nil, err
	}

	defer t.serial.Close(serialPort)

	if lineEnd {
		command += "\r\n"
	}

	_, err = t.serial.Write(serialPort, []byte(command))
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
			n, err := t.serial.Read(serialPort, buf)
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
