package main

import (
	"os/exec"
	"log"
	"fmt"
)

var (
	managedName string = "osx_managed_vpn"
	managedHost string = "managedvpn.local"
	managedPSK string = "dummy_psk"
	managedUserName = "dummy_un"
	managedPW = "dummy_pw"
	macvpnCMD = "macosvpn"
	macvpnArgs = []string{"create",
		"--l2tp",
		managedName,
		"--endpoint",
		managedHost,
		"--username",
		managedUserName,
		"--password",
		managedPW,
		"--shared-secret",
		managedPSK,
		"--split",
		"--force",
	}
)

func createManagedVPN() {
	cmd := exec.Command(macvpnCMD,macvpnArgs...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created %s VPN configuration",managedName)
}

func updateManagedVPNHost() {
	print("updating host\n")
	//todo
}

func establishManagedVPNConnection() {
	print("connecting to managed vpn\n")
	//todo
}

func updateRouting() {
	print("updating route table\n")
	//todo
}

func main() {
	createManagedVPN()
}
