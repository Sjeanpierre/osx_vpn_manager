package main

//references
//http://apple.stackexchange.com/questions/128297/how-to-create-a-vpn-connection-via-terminal
//https://developer.apple.com/legacy/library/documentation/Darwin/Reference/ManPages/man8/scutil.8.html
//https://github.com/halo/macosvpn


import (
	"gopkg.in/alecthomas/kingpin.v2"
	"regexp"
	"os"
	"log"
	"path"
	"os/user"
	"fmt"
)

var (
	//Connection Commands
	connect = kingpin.Command("connect", "Connect to a VPN")
	profile = connect.Flag("profile", "profile name.").Required().Short('p').Envar("VPN_PROFILE").String()
	vpn = connect.Arg("vpn", "Identifier for VPN to be connected").Required().String()
	//Disconnect Commands
	_ = kingpin.Command("disconnect", "Disconnect current VPN connection")
	//Host Commands
	hosts = kingpin.Command("host", "Commands related to vpn hosts")
	_ = hosts.Command("list", "List vpn hosts")
	_ = hosts.Command("refresh", "Refreshes resources")
	//Profile Commands
	profiles = kingpin.Command("profile", "Commands related to VPN connection profiles")
	_ = profiles.Command("list", "List vpn connection profiles")
	addProfilecmd = profiles.Command("add", "Add new profile to existing set")
	newProfile = addProfilecmd.Arg("profile", "Name of profile to add").Required().String()
	//Command Regex Section
	connectRegex = regexp.MustCompile(`^connect`)
	hostCommadRegex = regexp.MustCompile(`^host`)
	profileCommandRegex = regexp.MustCompile(`^profile`)
	disconnectCommandRegex = regexp.MustCompile(`^disconnect`)
	//Global Vars
        cliVersion = "0.0.6"
	resourcePath = path.Join(os.Getenv("HOME"), ".vpn_host_manager")
	DEBUG = false
)

func permissionCheck() {
	cu, err := user.Current()
	if err != nil {
		log.Fatalln("Could not retrieve user information:",err.Error())
	}
	if cu.Uid != "0" {
		log.Fatal("Please rerun as root or with sudo")
	}
}

func listVpnHosts() {
	printVPNHostList()
}

func hostFunctions(hostMethod string) {
	switch hostMethod {
	case "host list":
		listVpnHosts()
	case "host refresh":
		refreshHosts()
	}
}

func profileFunctions(profileMethod string) {
	switch profileMethod {
	case "profile list":
		printVPNProfileList()
	case "profile add":
		addProfile(*newProfile)
	default:
		log.Fatalf("not sure what to do with command: %s", profileMethod)
	}

}

func connectVPN(profileName string, vpnIdentifier string) {
	startConnection(vpnIdentifier, profileName)
}

func disconnectVPN() {
	fmt.Println("😭  BYE!! 😭")
	disconnectConnection()
}

func setupDirectories() {
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		error := os.Mkdir(resourcePath, 0700)
		if error != nil {
			log.Fatalf("encountered error during setup, %s", error)
		}
	}
}

func setup() {
	permissionCheck()
	setupDirectories()
}

func main() {
	kingpin.Version(cliVersion)
	setup()
	parsedArg := kingpin.Parse()
	switch {
	case hostCommadRegex.MatchString(parsedArg):
		hostFunctions(parsedArg)
	case profileCommandRegex.MatchString(parsedArg):
		profileFunctions(parsedArg)
	case connectRegex.MatchString(parsedArg):
		connectVPN(*profile, *vpn)
	case disconnectCommandRegex.MatchString(parsedArg):
		disconnectVPN()
	default:
		//if we are in this error block it is because we have established
		//a command for the provided text, but have not specified a regex
		//for handling it.
		log.Fatalf("Command signature not recognized: %s", parsedArg)
	}
}
