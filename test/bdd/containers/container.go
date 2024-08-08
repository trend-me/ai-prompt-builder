package containers

import (
	"fmt"
	"os/exec"
)

var errChan chan error

func Up() {
	errChan = make(chan error, 1)
	// Stop any running containers
	cmd := exec.Command("docker-compose", "-f", "../containers/docker-compose.yml", "down", "--volumes")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

	go func() {
		// Start the containers
		fmt.Println("Starting containers")
		cmd = exec.Command("docker-compose", "-f", "../containers/docker-compose.yml", "up")
		stdout, err = cmd.Output()
		fmt.Println(string(stdout))
		if err != nil {
			fmt.Println(err.Error())
			errChan <- err
			return
		}

		errChan <- nil
	}()
}

func Down() error {
	errChan <- nil
	cmd := exec.Command("docker-compose", "-f", "../containers/docker-compose.yml", "down", "--volumes")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(stdout))
	return nil
}
