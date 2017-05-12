package checks

import (
	"os/exec"
	"log"
	"regexp"
	"strconv"
)

func CheckDockerPool(okSize int) (bool, string) {
	isOk := false
	var msg string
	out, err := exec.Command("bash", "-c", "lvs -o data_percent,metadata_percent,LV_NAME --noheadings --units G --nosuffix | grep docker-pool").Output()
	if err != nil {
		isOk = false
		msg = "Could not parse docker pool size: " + err.Error()
		log.Println(msg)
		return isOk, msg
	}

	isOk = isLvsSizeOk(string(out), okSize)
	if (!isOk) {
		msg = "Docker pool size is above: " + strconv.Itoa(okSize)
	}
	return isOk, msg
}

func isLvsSizeOk(stdOut string, okSize int) bool {
	// 42.10  8.86   docker-pool
	num := regexp.MustCompile("(\\d+\\.\\d+)")

	checksOk := 0
	for _, nr := range num.FindAllString(stdOut, -1) {
		i, err := strconv.ParseFloat(nr, 64)
		if (err != nil) {
			log.Print("Unable to parse int:", nr)
			return false
		}
		if (i < float64(okSize)) {
			checksOk++
		} else {
			log.Println("Docker pool size exceeded okSize:", i)
		}
	}

	return checksOk == 2
}
