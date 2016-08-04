package main

//todo allow for setting default profiles

import (
	"io/ioutil"
	"fmt"
	"os"
	"github.com/olekukonko/tablewriter"
	"path"
	"encoding/json"
	"strconv"
	"log"
	"regexp"
)

var (
	vpnProfileFields = []string{"ID #", "Name", "Username"}
	vpnProfileFilePath = path.Join(resourcePath, "vpn_profiles.json")
	noSuchFileErrRegexp = regexp.MustCompile(`no such file or directory`)
)

type vpnProfile struct {
	Name     string `json:"name"`
	Psk      string `json:"psk"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func loadProfileFile() []vpnProfile {
	file, e := ioutil.ReadFile(vpnProfileFilePath)
	if e != nil {
		if noSuchFileErrRegexp.MatchString(e.Error()) {
			return []vpnProfile{}
		}
		fmt.Printf("Could not: %v\n", e.Error())
		os.Exit(1)
	}
	var profiles []vpnProfile
	json.Unmarshal(file, &profiles)
	return profiles
}

func writeProfileFile(profileList []vpnProfile) {
	profileJSON, err := json.Marshal(profileList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Writing profile file to %s\n", vpnProfileFilePath)
	writeError := ioutil.WriteFile(vpnProfileFilePath, profileJSON, 0755)
	if writeError != nil {
		fmt.Print("Could not write profile file\n")
		log.Fatal(writeError)
	}
	fmt.Println("New profile saved!")
}

func printVPNProfileList() {
	vpnProfiles := loadProfileFile()
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
	vpnProfiles := loadProfileFile()
	var selectedProfile vpnProfile
	for _, profile := range vpnProfiles {
		if profile.Name == profileName {
			selectedProfile = profile
		}
	}
	if selectedProfile.Name == "" {
		log.Fatalf("VPN Profile: %s not found in %s", profileName, vpnProfileFilePath)
	}
	return selectedProfile
}

func detectDuplicateName(providedName string) {
	vpnProfiles := loadProfileFile()
	for _, profile := range vpnProfiles {
		if profile.Name == providedName {
			log.Fatalf("profile name %s: Already present, please select another name", providedName)
		}
	}
}

func detailCapture(attr string) string {
	var response string
	fmt.Printf("%s ", attr)
	_, err := fmt.Scanln(&response)
	if err != nil {
		if err.Error() == "unexpected newline" {
			return ""
		}
		log.Fatal(err)
	}
	return response
}

	var returnvar bool
func confirm() bool {
	confirmation := detailCapture("Save Profile? [y/n]:")
	switch confirmation {
	case "y":
		returnvar = true
	case "n":
		main()
	default:
		confirm(username, password, psk)
	}
	return returnvar
}

func addProfile(profileName string) {
	vpnProfiles := loadProfileFile()
	detectDuplicateName(profileName)
	fmt.Printf("Please enter the following values to configure VPN profile %s\n", profileName)
	username := detailCapture("USERNAME:")
	password := detailCapture("PASSWORD:")
	psk := detailCapture("PSK:")
	if confirm(username, password, psk) {
		var updated = []vpnProfile{}
		updated = append(vpnProfiles, vpnProfile{Name:profileName, UserName:username, PassWord:password, Psk:psk})
		writeProfileFile(updated)
	}
}

