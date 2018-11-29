package main

import "os/exec"

const (

	cmdSerialNumber string = "ioreg -l | grep IOPlatformSerialNumber | awk -F '=' '{print $2}' | tr -d ' \n'"
	cmdCPU string = "/usr/sbin/sysctl -n machdep.cpu.brand_string"
	cmdDisk string = "system_profiler SPStorageDataType -detailLevel mini | grep Capacity | awk -F '(' '{print $1}' | awk -F ':' '{print $2}' | tr -d ' '"
	cmdMemory = "sysctl hw.memsize | awk -F ':' '{print $2}' | tr -d ' \n'"
	cmdGetAccount = "id -un"
)

func getAccount () string {
	data, _ := exec.Command ( "bash", "-c", cmdGetAccount ).Output()
	return string(data)
}

func getMemory () string {

	data, _ := exec.Command("bash", "-c", cmdMemory).Output()
	return string(data)
}


func getCPU () string {

	data, _ := exec.Command("bash", "-c", cmdCPU).Output()
	return string(data)
}

func getSerialNumber () string {

	sn, _ := exec.Command("bash", "-c", cmdSerialNumber).Output()
	return string(sn)
}

func getDisk () string {

	data, _ := exec.Command("bash", "-c", cmdDisk).Output()
	return string(data)
}
