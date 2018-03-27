package checks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"os"
	"regexp"
)

func CheckOpenFileCount() error {
	log.Println("Checking open files")

	out, err := exec.Command("bash", "-c", "cat /proc/sys/fs/file-nr | cut -f1").Output()
	if err != nil {
		msg := "Could not evaluate open file count: " + err.Error()
		log.Println(msg)
		return errors.New(msg)
	}

	nr, err := strconv.Atoi(strings.TrimSpace(string(out)))

	if err != nil {
		return errors.New("Could not parse output to integer: " + string(out))
	}

	if nr < 200000 {
		return nil
	} else {
		return errors.New("Open files are higher than 200'000 files!")
	}
}

func CheckGlusterStatus() error {
	log.Println("Checking gluster status with gstatus")

	out, err := exec.Command("bash", "-c", "gstatus -o json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 16") || strings.Contains(err.Error(), "exit status 1") || strings.Contains(err.Error(), "exit status 12") {
			// Other gluster server did the same check the same time
			// Try again 5 seconds
			time.Sleep(5 * time.Second)
			out, err = exec.Command("bash", "-c", "gstatus -o json").Output()
			if err != nil {
				msg := "Could not check gstatus output. Tryed 2 times. Error: " + err.Error()
				log.Println(msg)
				return errors.New(msg)
			}
		} else {
			msg := "Could not check gstatus output: " + err.Error()
			log.Println(msg)
			return errors.New(msg)
		}
	}

	// Sample JSON
	// 2017-03-27 12:34:17.626544 {"brick_count": 4, "bricks_active": 4, "glfs_version": "3.7.9", "node_count": 2, "nodes_active": 2, "over_commit": "No", "product_name": "Red Hat Gluster Storage Server 3.1 Update 3", "raw_capacity": 214639312896, "sh_active": 2, "sh_enabled": 2, "snapshot_count": 0, "status": "healthy", "usable_capacity": 107319656448, "used_capacity": 11712278528, "volume_count": 2, "volume_summary": [{"snapshot_count": 0, "state": "up", "usable_capacity": 53659828224, "used_capacity": 34619392, "volume_name": "vol_fast_registry"}, {"snapshot_count": 0, "state": "up", "usable_capacity": 53659828224, "used_capacity": 5821519872, "volume_name": "vol_slow_openshift-infra"}]}
	res := string(out)[27:(len(string(out)))]

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		msg := "Error decoding gstatus output: " + res
		log.Println(msg)
		return errors.New(msg)
	}

	if dat["status"] != "healthy" {
		return errors.New("Status of GlusterFS is not healthy")
	}

	return nil
}

func CheckVGSizes(okSize int) error {
	log.Println("Checking VG free size")

	out, err := exec.Command("bash", "-c", "vgs -o vg_free,vg_size,VG_NAME --noheadings --units G --nosuffix | grep -v crash").Output()
	if err != nil {
		msg := "Could not evaluate VG sizes: " + err.Error()
		log.Println(msg)
		return errors.New(msg)
	}

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if len(l) > 0 {
			isOk := isVgSizeOk(l, okSize)

			log.Println("Checking VG size: ", l)

			if !isOk {
				return fmt.Errorf("VG size is below: %v | %v", strconv.Itoa(okSize), l)
			}
		}
	}

	return nil
}

func CheckLVPoolSizes(okSize int) error {
	log.Println("Checking LV pool used size")

	out, err := exec.Command("bash", "-c", "lvs -o data_percent,metadata_percent,LV_NAME --noheadings --units G --nosuffix | grep pool").Output()
	if err != nil {
		msg := "Could not evaluate LV pool size: " + err.Error()
		log.Println(msg)
		return errors.New(msg)
	}

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if len(l) > 0 {
			isOk := isLvsSizeOk(l, okSize)

			log.Println("Checking LV Pool: ", l)

			if !isOk {
				return fmt.Errorf("LV pool size is above: %v | %v", strconv.Itoa(okSize), l)
			}
		}
	}

	return nil
}

func CheckMountPointSizes(okSize int) error {
	mounts := os.Getenv("MOUNTPOINTS_TO_CHECK")

	if mounts == "" {
		return nil
	}

	mountList := strings.Split(mounts, ",")

	for _, m := range mountList {
		log.Printf("Checking free disk size of %v.", m)

		out, err := exec.Command("bash", "-c", "df --output=target,pcent | grep -w "+m).Output()
		if err != nil {
			msg := "Could not evaluate df of mount point: " + m + ". err: " + err.Error()
			log.Println(msg)
			return errors.New(msg)
		}

		// Example: /gluster/fast_registry                               8%
		num := regexp.MustCompile(`\d+%`)
		usages := num.FindAllString(string(out), 1)
		log.Println(m, usages)
		if len(usages) != 1 {
			return errors.New("Could not parse output to integer: " + string(out))
		}
		usageInt, err := strconv.Atoi(strings.Replace(usages[0], "%", "", 1))
		if err != nil {
			return errors.New("Could not parse output to integer: " + string(out))
		}

		if usageInt > okSize {
			msg := fmt.Sprintf("Usage %% of volume %v is bigger than treshold. Is: %v%% - treshold: %v%%", m, usageInt, okSize)
			log.Println(msg)
			return errors.New(msg)
		}
	}

	return nil
}

func CheckIfGlusterdIsRunning() error {
	log.Print("Checking if glusterd is running")

	out, err := exec.Command("bash", "-c", "systemctl status glusterd").Output()
	if err != nil {
		msg := "Could not run 'systemctl status glusterd'. err: " + err.Error()
		log.Println(msg)
		return errors.New(msg)
	}

	if !strings.Contains(string(out), "active (running)") {
		return fmt.Errorf("Glusterd seems not to be running! Output: %v", string(out))
	}

	return nil
}
