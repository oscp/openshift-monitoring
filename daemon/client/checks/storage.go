package checks

import (
	"os/exec"
	"log"
	"strings"
	"encoding/json"
	"strconv"
)

func CheckOpenFileCount() (bool, string) {
	isOk := false
	var msg string
	out, err := exec.Command("bash", "-c", "cat /proc/sys/fs/file-nr | cut -f1").Output()
	if err != nil {
		msg = "Could not evaluate open file count: " + err.Error()
		log.Println(msg)
		return isOk, msg
	}

	nr, err := strconv.Atoi(strings.TrimSpace(string(out)))

	if (err != nil) {
		msg = "Could not parse output to integer: " + string(out)
		return isOk, msg
	}

	if (nr < 200000) {
		isOk = true
	}

	if (!isOk) {
		msg = "Open files are higher than 200'000 files!"
	}
	return isOk, msg
}

func CheckGlusterStatus() (bool, string) {
	var msg string
	out, err := exec.Command("bash", "-c", "gstatus -abw -o json").Output()
	if err != nil {
		msg = "Could not check gstatus output: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	// Sample JSON
	// 2017-03-27 12:34:17.626544 {"brick_count": 4, "bricks_active": 4, "glfs_version": "3.7.9", "node_count": 2, "nodes_active": 2, "over_commit": "No", "product_name": "Red Hat Gluster Storage Server 3.1 Update 3", "raw_capacity": 214639312896, "sh_active": 2, "sh_enabled": 2, "snapshot_count": 0, "status": "healthy", "usable_capacity": 107319656448, "used_capacity": 11712278528, "volume_count": 2, "volume_summary": [{"snapshot_count": 0, "state": "up", "usable_capacity": 53659828224, "used_capacity": 34619392, "volume_name": "vol_fast_registry"}, {"snapshot_count": 0, "state": "up", "usable_capacity": 53659828224, "used_capacity": 5821519872, "volume_name": "vol_slow_openshift-infra"}]}
	res := string(out)[27:(len(string(out)))]

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		msg = "Error decoding gstatus output: " + res
		log.Println(msg)
		return false, msg
	}

	if (dat["status"] != "healthy") {
		return false, "Status of GlusterFS is not healthy"
	}

	return true, ""
}