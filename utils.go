package atcom

import "os/exec"

func RunCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	return string(output), err
}
