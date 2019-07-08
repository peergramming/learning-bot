package utils

import (
	"errors"
	"os/exec"
)

// SystemPackagesCheck makes sure that the required packages are installed
// on the system.
func SystemPackagesCheck() (err error) {
	// Check java
	_, err = exec.Command("java", "-version").Output()
	if err != nil {
		return errors.New("OpenJDK or Oracle Java is not installed on the system")
	}

	// Check unzip utility
	_, err = exec.Command("unzip").Output()
	if err != nil {
		return errors.New("unzip utility is not installed on the system")
	}
	return nil
}
