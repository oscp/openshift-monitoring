package checks

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CheckBondNetworkInterface() error {
	log.Println("Checking bond0 interface")

	if _, err := os.Stat("/proc/net/bonding/bond0"); err == nil {
		// bond0 exists, execute check
		out, err := exec.Command("bash", "-c", "grep 'MII Status: up' /proc/net/bonding/bond0 | wc -l").Output()
		if err != nil {
			msg := "Could not evaluate bond0 status: " + err.Error()
			log.Println(msg)
			return errors.New(msg)
		}

		nr, err := strconv.Atoi(strings.TrimSpace(string(out)))
		if err != nil {
			return errors.New("Could not parse output to integer: " + string(out))
		}

		if nr != 3 {
			// 3 is the expected number of occurrences
			return errors.New("bond0 degraded: At least one interface is not 'UP'")
		}

	} else {
		log.Println("bond0 does not exist, skipping this check...")
	}

	return nil
}
