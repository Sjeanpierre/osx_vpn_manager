package main

//todo use keychain to store psks instead of plaintext config file
//todo write out IP-up template with all host for routing instead of using route command

import (
	"os/exec"
	"log"
	"fmt"
	"github.com/lextoumbourou/goodhosts"
	"os"
	"time"
	"regexp"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

var (
	vpnProfileFilePath = path.Join(resourcePath, "vpn_profiles.json")
	managedName = "osx_managed_vpn"
	managedHost = "managedvpn.local"
	managedPSK = "osx_managed_psk"
	managedUserName = "osx_managed_un"
	managedPW = "osx_managed_pw"
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
	connectionRegex = regexp.MustCompile(`^Connected`)
	existingHostRegex = regexp.MustCompile(strings.Join([]string{managedHost, "$"}, ""))
	vpcUIDRegex = regexp.MustCompile(`^vpc-`)
	vpcIndexRegex = regexp.MustCompile(`\d?`)
	weirdRouteExitCodeRegex = regexp.MustCompile(`exit status 64`)
	vpnProfileFields = []string{"ID#","Name","Username"}
	sameConnection bool
)

type vpnProfile struct {
	Name     string `json:"name"`
	Psk      string `json:"psk"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func loadprofileFile() []vpnProfile {
	file, e := ioutil.ReadFile(vpnProfileFilePath)
	if e != nil {
		fmt.Printf("Could not: %v\n", e)
		os.Exit(1)
	}
	var profiles []vpnProfile
	json.Unmarshal(file, &profiles)
	return profiles
}

func printVPNProfiles() {
	vpnProfiles := loadprofileFile()
	consoleTable := tablewriter.NewWriter(os.Stdout)
	consoleTable.SetHeader(vpnProfileFields)
	for index, vpnProfile := range vpnProfiles {
		row := []string{
			strconv.Itoa(index),
			vpnProfile.Name,
			vpnProfile.UserName,
		}
		consoleTable.Append(row)
	}
	consoleTable.Render()
}

func selectVPNProfileDetails(profileName string) vpnProfile {
	var selectedProfile vpnProfile
	profiles := loadprofileFile()
	for _, profile := range profiles {
		if profile.Name == profileName {
			selectedProfile = profile
		}
	}
	if selectedProfile.Name == "" {
		log.Fatalf("VPN Profile: %s not found in %s", profileName, vpnProfileFilePath)
	}
	return selectedProfile
}

func createManagedVPN() {
	cmd := exec.Command(macvpnCMD, macvpnArgs...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created %s VPN configuration", managedName)
}

func updateManagedVPNHost(vpnHost vpnInstance) {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatal("Could not read hostfile")
	}
	if hosts.Has(vpnHost.PublicIP, managedHost) {
		sameConnection = true
		return
	}
	removeExistingHost()
	addManagedVPNHost(vpnHost)
}

func needsDisconnection() bool{
	var disconnect bool
	if connectionEstablished() {
		if sameConnection {
			disconnect = false
		} else {
			disconnect = true
		}
	}

	return disconnect
}

func disconnectExistingConnection() {
	if needsDisconnection() {
		fmt.Println("Disconnecting existing managed VPN connection")
		cmd := exec.Command("scutil",
			"--nc",
			"stop",
			managedName,
		)
		err := cmd.Run()
		if err != nil {
			log.Fatal("Could not stop previous managed VPN connection")
		}
	}
}

func addManagedVPNHost(vpnHost vpnInstance) {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatal("Could not read hostfile")
	}
	hosts.Add(vpnHost.PublicIP, managedHost)
	if err := hosts.Flush(); err != nil {
		log.Fatalf("Error writing host entry", err)
	}
}

func removeExistingHost() {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatal("Could not read hostfile")
	}
	for _, hostLine := range hosts.Lines {
		if existingHostRegex.MatchString(hostLine.Raw) {
			fmt.Printf("Removing `%s` from hostfile\n", hostLine.Raw)
			hosts.Remove(hostLine.IP, hostLine.Hosts[0])
		}
	}
	if err := hosts.Flush(); err != nil {
		log.Fatalf("Error writing host entry", err)
	}
}

func establishManagedVPNConnection(vpnDetails vpnProfile, vpnHost *vpnInstance) {
	cmd := exec.Command("scutil",
		"--nc",
		"start",
		managedName,
		"--user",
		vpnDetails.UserName,
		"--password",
		vpnDetails.PassWord,
		"--secret",
		vpnDetails.Psk)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Could not connect to vpn via scutil: %s", err)
	}
	i := 0
	print("connecting...")
	for {
		print(".")
		if connectionEstablished() {
			print("\n")
			updateRouting(*vpnHost)
			fmt.Printf("VPN connection to %s established!!", vpnHost.Name)
			break
		} else if i < 20 {
			i++
			time.Sleep(500 * time.Millisecond)
		} else {
			log.Fatal("Could not set route, timed after 10 seconds waiting for VPN connection")
			break
		}
	}
}

func verifyManagedVPNConnection() bool {
	cmd := exec.Command("scutil",
		"--nc",
		"show",
		managedName,
	)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func setupManagedVPNConnection() {
	if verifyManagedVPNConnection() {
		return
	}
	log.Printf("Managed VPN `%s` not found, creating...\n", managedName)
	createManagedVPN()
	if verifyManagedVPNConnection() {
		fmt.Println("Managed VPN settings applied, please rerun last command\n")
		os.Exit(0)
	}
	log.Fatal("Could not setup managed VPN connection\n")
}

func connectionEstablished() bool {
	output, err := exec.Command("scutil", "--nc", "status", managedName).Output()
	if err != nil {
		log.Fatal(err)
	}
	if connectionRegex.MatchString(string(output)) {
		return true
	}
	return false
}

func updateRouting(vpnHost vpnInstance) {
	print("updating route table\n")
	_, err := exec.Command("route","-v", "add", "-net", vpnHost.VpcCidr, "-interface ppp0").Output()
	if err != nil {
		if weirdRouteExitCodeRegex.MatchString(err.Error()) {
			return
		}
		log.Fatalf("Could not update route table after VPN connection: %s", err)

	}
}

func selectVPNHost(identifier string) vpnInstance{
	vpnHostsList := readHostsJSONFile()
	if vpcUIDRegex.MatchString(identifier) {
		fmt.Println("Connecting to VPN by UID")
		for _,host := range vpnHostsList {
			if host.VpcID == identifier {
				return host
			}
		}
	}
	if vpcIndexRegex.MatchString(identifier) {
		fmt.Println("Connecting to VPN by ID#")
		for index,host := range vpnHostsList {
			if strconv.Itoa(index) == identifier  {
				return host
			}
		}
	}
	fmt.Println("Connecting to VPN by instnace Name")
	for _,host := range vpnHostsList {
		if host.Name == identifier  {
			return host
		}
	}
	log.Fatal("Could not find VPN with provided identifier")
	return vpnInstance{}
}

func startConnection(vpnIdentifier string, profileName string) {
	setupManagedVPNConnection()
	vpnHost := selectVPNHost(vpnIdentifier)
	updateManagedVPNHost(vpnHost)
	disconnectExistingConnection()
	profile := selectVPNProfileDetails(profileName)
	establishManagedVPNConnection(profile, &vpnHost)
}