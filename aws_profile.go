package main

import (
	"github.com/go-ini/ini"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"os"
	"path"
)

var (
	awsProfileNamesPath = path.Join(resourcePath, "aws_profile_names.json")
	awsCredentialFilePath = path.Join(os.Getenv("HOME"), ".aws", "credentials")
)

func existingProfiles() bool {
	if _, err := os.Stat(awsProfileNamesPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func awsProfiles() []string {
	if existingProfiles() {
		profiles, err := readAWSProfileFile()
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		return profiles
	}
	userQuestion := fmt.Sprint("Use AWS profiles? [y/n]:")
	if confirmUserSelection(userQuestion) {
		setupProfiles()
		profiles, err := readAWSProfileFile()
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		return profiles
	}
	return []string{"default"}
}

func writeAWSProfileFile(profileNames []string) {
	profileNamesJSON, err := json.Marshal(profileNames)
	if err != nil {
		fmt.Println(err)
		return
	}
	writeError := ioutil.WriteFile(awsProfileNamesPath, profileNamesJSON, 0755)
	if writeError != nil {
		fmt.Print("Could not write AWS profile names to config file\n")
		log.Fatal(writeError)
	}
}

func readAWSProfileFile() ([]string, error) {
	file, e := ioutil.ReadFile(awsProfileNamesPath)
	if e != nil {
		if noSuchFileErrRegexp.MatchString(e.Error()) {
			return []string{}, e
		}
		fmt.Printf("Could not: %v\n", e.Error())
		os.Exit(1)
	}
	var awsProfiles []string
	err := json.Unmarshal(file, &awsProfiles)
	if err != nil {
		log.Fatal("could not read AWS profile file")
	}
	return awsProfiles, nil
}

func detail4Capture(attr string) string {
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

func confirmUserSelection(userPrompt string) bool {
	var returnVar bool
	confirmation := detail4Capture(userPrompt)
	switch confirmation {
	case "y":
		returnVar = true
	case "n":
		returnVar = false
	default:
		confirmUserSelection(userPrompt)
	}
	return returnVar
}

func setupProfiles() {
	fmt.Println("Discovering AWS profile names from credentials file")
	cfg, err := ini.Load(awsCredentialFilePath)
	if err != nil {
		fmt.Printf("error reading AWS credential file: %s", err)
	}
	sections := cfg.SectionStrings()
	fmt.Println("Please select profiles AWS profiles to include")
	var addedProfiles []string
	for _, section := range sections {
		if section == ini.DEFAULT_SECTION {
			continue
		}
		userQuestion := fmt.Sprintf("Include profile [%s]? [y/n]:", section)
		if confirmUserSelection(userQuestion) {
			addedProfiles = append(addedProfiles, section)
		}
	}
	writeAWSProfileFile(addedProfiles)
}
