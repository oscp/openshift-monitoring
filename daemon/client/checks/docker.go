package checks

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func CheckDockerPool(okSize int) error {
	log.Println("Checking docker pool used size")

	out, err := exec.Command("bash", "-c", "lvs -o data_percent,metadata_percent,LV_NAME --noheadings --units G --nosuffix | grep docker-pool").Output()
	if err != nil {
		// ignore errors. grep exits with 1 if docker-pool is not found
		return nil
	}

	isOk := isLvsSizeOk(string(out), okSize)
	if !isOk {
		return fmt.Errorf("Docker pool size is above: %v", strconv.Itoa(okSize))
	}
	return nil
}
