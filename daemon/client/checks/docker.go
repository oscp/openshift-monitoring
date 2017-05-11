package checks

import (
	"bytes"
	"os/exec"
	"log"
)

func CheckDockerPool(okSize int) (bool, string) {
	isOk := false
	var msg string

	cmd := exec.Command("lvs -o data_percent,metadata_percent,LV_NAME --noheadings --units G --nosuffix | grep docker-pool")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		isOk = false
		log.Println("error while parsing docker pool size:", err)
		msg = "Could not parse docker pool size: " + err.Error()
		return isOk, msg
	}

	isOk = isLvsSizeOk(out.String(), okSize)
	if (!isOk) {
		msg = "Docker pool size is above: " + string(okSize)
	}
	return isOk, msg
}

func isLvsSizeOk(stdOut string, okSize int) bool {
	// Example: arr = [" ", "3.34", "   ", "0.64", "   docker-pool"];
	log.Println("LVS: ", stdOut)

	return true
}
