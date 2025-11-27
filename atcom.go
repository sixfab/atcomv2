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
	urc := c.Urc

	serialPort, err := t.open(portname, baudrate)

	if err != nil {
		c.Error = err
		return c
	}

	defer t.serial.Close(serialPort)

	if lineEnd {
		command += "\r\n"
	}

	// If urc is true, do not send command to serial port.
	if !urc {
		_, err = t.serial.Write(serialPort, []byte(command))
	}

	if err != nil {
		c.Error = err
		return c
	}

	data := make([]string, 0)
	responseBuffer := ""
	timeoutDuration := time.Duration(timeout) * time.Second

	found := make(chan error)

	ctxScan, cancelScan := context.WithCancel(context.Background())
	defer cancelScan()

	go func(ctx context.Context) {
		response := ""
		buf := make([]byte, 1024)
		leftFromLastRead := make([]byte, 1024)
		nLeftBytes := 0

		for {
			select {
			case <-ctx.Done():
				close(found)
				return
			default:
				n, err := t.serial.Read(serialPort, buf)
				if err != nil {
					if err.Error() == "EOF" {
						continue
					}

					found <- err
					return
				}

				// if no data, continue
				if n == 0 {
					continue
				}

				response = string(leftFromLastRead[:nLeftBytes]) + string(buf[:n]) // prepend left bytes from last read
				nLeftBytes = 0                                                     // reset left bytes count

				// replace \r with \n for uniformity
				response = strings.ReplaceAll(response, "\r", "\n")
				lines := strings.Split(response, "\n")
				responseBuffer += response

				if len(lines) == 1 && n > 0 {
					// save all bytes to leftFromLastRead
					copy(leftFromLastRead, []byte(response))
					nLeftBytes = n
					continue
				}

				// split response to lines and trim spaces
				for index, line := range lines {
					line = strings.TrimSpace(line)

					if line == "" {
						if index == len(lines)-1 {
							// save left bytes for next read
							copy(leftFromLastRead, []byte(line))
							nLeftBytes = len(line)

							if nLeftBytes > 0 {
								// remove last index of data
								if len(data) > 0 {
									data = data[:len(data)-1]
								}
							}
							break
						}
						continue
					}

					data = append(data, line)

					// send line to response channel if exists
					if responseChan != nil {
						c.ResponseChan <- line
					}
				}

				if responseChan == nil { // if responseChan is not existed
					for _, line := range data {
						if line == "OK" {
							// check desired and fault existed in response
							if desired != nil || fault != nil {
								ok := false
								for _, desiredStr := range desired {
									if strings.Contains(responseBuffer, desiredStr) {
										ok = true
										found <- nil
										return
									}
								}
								for _, faultStr := range fault {
									if strings.Contains(responseBuffer, faultStr) {
										ok = true
										found <- errors.New("faulty response detected")
										return
									}
								}

								if !ok {
									found <- errors.New("desired or fault response not found")
									return
								}
							} else {
								found <- nil
								return
							}
						}
					}
				} else { // if responseChan is existed
					// check desired and fault existed in response
					if desired != nil || fault != nil {
						for _, desiredStr := range desired {
							if strings.Contains(responseBuffer, desiredStr) {
								found <- nil
								return
							}
						}
						for _, faultStr := range fault {
							if strings.Contains(responseBuffer, faultStr) {
								found <- errors.New("faulty response detected")
								return
							}
						}
					}
				}

				// check "ERROR" existed in response
				if strings.Contains(responseBuffer, "ERROR") {
					found <- errors.New(responseBuffer)
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
			c.Error = errors.New("timeout")
			return c
		}
	}
}
