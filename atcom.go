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
	serial SerialModel
	shell  ShellModel
}

// Serial Implementation for normal usage
type Serial struct{}

// Serial interface
type SerialModel interface {
	OpenPort(c *serial.Config) (*serial.Port, error)
	Write(port *serial.Port, command []byte) (n int, err error)
	Close(port *serial.Port) (err error)
	Read(port *serial.Port, buffer []byte) (n int, err error)
}

// Serial implements Serial interface
func (s *Serial) OpenPort(c *serial.Config) (*serial.Port, error) {
	return serial.OpenPort(c)
}

func (s *Serial) Write(port *serial.Port, command []byte) (n int, err error) {
	return port.Write([]byte(command))
}

func (s *Serial) Close(port *serial.Port) (err error) {
	return port.Close()
}

func (s *Serial) Read(port *serial.Port, buffer []byte) (n int, err error) {
	return port.Read(buffer)
}

// Shell Implementation for normal usage
type Shell struct{}

// Shell interface
type ShellModel interface {
	Command(name string, arg ...string) (string, error)
}

// RealShell implements Shell interface
func (s *Shell) Command(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	return string(output), err
}

// NewAtcom creates a new Atcom instance with default serial and shell implementations
func NewAtcom(s SerialModel, sh ShellModel) *Atcom {

	if s == nil {
		s = &Serial{}
	}

	if sh == nil {
		sh = &Shell{}
	}

	return &Atcom{
		serial: s,
		shell:  sh,
	}
}

// Function to open serial port
func (t *Atcom) open(portname string, baudrate int) (port *serial.Port, err error) {

	if baudrate == 0 {
		baudrate = 115200
	}

	if portname == "" {
		return nil, errors.New("serialport is required")
	}

	config := &serial.Config{
		Name:        portname,
		Baud:        baudrate,
		ReadTimeout: time.Millisecond * 100,
	}

	return t.serial.OpenPort(config)
}

// SendAT sends AT command to modem and returns response
func (t *Atcom) SendAT(c *ATCommand) *ATCommand {

	command := c.Command
	lineEnd := c.LineEnd
	timeout := c.Timeout
	desired := c.Desired
	fault := c.Fault
	portname := c.SerialAttr.Port
	baudrate := c.SerialAttr.Baud
	responseChan := c.ResponseChan

	serialPort, err := t.open(portname, baudrate)

	if err != nil {
		c.Error = err
		return c
	}

	defer t.serial.Close(serialPort)

	if lineEnd {
		command += "\r\n"
	}

	_, err = t.serial.Write(serialPort, []byte(command))
	if err != nil {
		c.Error = err
		return c
	}

	data := make([]string, 0)
	timeoutDuration := time.Duration(timeout) * time.Second

	found := make(chan error)

	ctxScan, cancelScan := context.WithCancel(context.Background())
	defer cancelScan()

	go func(ctx context.Context) {
		response := ""
		buf := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				close(found)
				return
			default:
				time.Sleep(time.Millisecond * 5)
				n, err := t.serial.Read(serialPort, buf)
				if err != nil {
					if err.Error() == "EOF" {
						continue
					}

					found <- err
					return
				}

				// Send real-time responses through the channel if ResponseChan is set
				if responseChan != nil {
					if n > 0 {
						response = string(buf[:n])
					}

					lines := strings.Split(response, "\r\n")

					for _, line := range lines {
						line = strings.TrimSpace(line)
						line = strings.Trim(line, "\r")
						line = strings.Trim(line, "\n")

						if line != "" {
							data = append(data, line)
							c.ResponseChan <- line
						}

						if strings.Contains(line, "ERROR") {
							found <- errors.New(line)
							break
						}

						if strings.Contains(line, "OK") {
							found <- nil
							break
						}
					}
					// Read responses continuously until timeout is reached
					continue
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
						ok := false
						for _, desiredStr := range desired {
							if strings.Contains(response, desiredStr) {
								ok = true
								found <- nil
								return
							}
						}
						for _, faultStr := range fault {
							if strings.Contains(response, faultStr) {
								found <- errors.New("faulty response detected")
								return
							}
						}

						if !ok {
							found <- errors.New("desired response not found")
							return
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
		}
	}(ctxScan)

	timeoutCh := time.After(timeoutDuration)

	for {
		select {
		case err := <-found:
			c.Response = data
			c.Error = err
			return c
		case <-timeoutCh:
			cancelScan()
			c.Response = data
			if c.ResponseChan == nil {
				c.Error = errors.New("timeout")
			}
			return c
		}
	}
}
