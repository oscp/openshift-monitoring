package checks

import (
	"os/exec"
	"log"
	"strconv"
)

func CheckDockerPool(okSize int) (bool, string) {
	var msg string
	out, err := exec.Command("bash", "-c", "lvs -o data_percent,metadata_percent,LV_NAME --noheadings --units G --nosuffix | grep docker-pool").Output()
	if err != nil {
		msg = "Could not parse docker pool size: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	isOk := isLvsSizeOk(string(out), okSize)
	if (!isOk) {
		msg = "Docker pool size is above: " + strconv.Itoa(okSize)
	}
	return isOk, msg
}
