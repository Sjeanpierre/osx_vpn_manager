package main

//references
//http://apple.stackexchange.com/questions/128297/how-to-create-a-vpn-connection-via-terminal
//https://developer.apple.com/legacy/library/documentation/Darwin/Reference/ManPages/man8/scutil.8.html
//https://github.com/halo/macosvpn

//todo read config/profile files for connection details
//todo read vpn host list
//todo create vpn skeleton
//todo establish vpn connection
//todo configure routing

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"regexp"
	"os"
	"log"
	"path"
)

var (
	resourcePath = path.Join(os.Getenv("HOME") ,".vpn_host_manager")
	connect = kingpin.Command("connect", "Connect to a VPN")
	profile = connect.Flag("profile", "profile name.").Required().Short('p').String()
	list = kingpin.Command("list", "List stuff")
	hosts = list.Command("hosts", "List available vpn hosts")
	profileList = list.Command("profiles", "List available vpn profiles")
	listRegex = regexp.MustCompile(`^list`)
	refreshRegex = regexp.MustCompile(`^refresh`)
	refresh = kingpin.Command("refresh", "Refreshes resources")
)

func listVpnHosts() {
	printHosts()
}

func listVpnProfiles() {
	fmt.Print("profile list")
}

func listFunctions(listMethod string) {
	switch listMethod {
	case "list hosts":
		listVpnHosts()
	case "list profiles":
		listVpnProfiles()
	}
}

func refreshFunctions(refreshMethod string) {
	if refreshMethod == "refresh" {
		refreshHosts()
	}

}

func setupDirectories() {
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		error := os.Mkdir(resourcePath,0644)
		if error != nil {
			log.Fatalf("encountered error during setup, %s", error)
		}
	}
}

func main() {
        setupDirectories()
	parsed := kingpin.Parse()
	switch {
	case listRegex.MatchString(parsed):
		listFunctions(parsed)
	case refreshRegex.MatchString(parsed):
		refreshFunctions(parsed)
	}
}
