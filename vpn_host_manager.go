package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"strings"
	"encoding/json"
	"io/ioutil"
)

var awsRegions = []string{"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1", "sa-east-1"}

type vpnInstance struct {
	VpcID string `json:"vpc_id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	PublicIP    string `json:"public_ip"`
	VpcCidr     string `json:"vpc_cidr"`
}

func listVPCs() map[string]string {
	vpcList := make(map[string]string)
	for _, region := range awsRegions {
		fmt.Printf("fetching vpc details for region: %v\n", region)
		svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))
		params := &ec2.DescribeVpcsInput{}
		resp, err := svc.DescribeVpcs(params)
		if err != nil {
			fmt.Println("there was an error listing vpcs in", region, err.Error())
			log.Fatal(err.Error())
		}
		for _, vpc := range resp.Vpcs {
			vpcID := *vpc.VpcId
			vpcCIDR := *vpc.CidrBlock
			vpcList[vpcID] = vpcCIDR
		}

	}
	return vpcList
}

func listFilteredInstances(nameFilter string) []*ec2.Instance {
	var filteredInstances []*ec2.Instance
	for _, region := range awsRegions {
		svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))
		fmt.Printf("fetching instances with tag %v in: %v\n", nameFilter, region)
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []*string{
						aws.String(strings.Join([]string{"*", nameFilter, "*"}, "")),
					},
				},
			},
		}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			fmt.Println("there was an error listing instnaces in", region, err.Error())
			log.Fatal(err.Error())
		}
		for _, reservation := range resp.Reservations {
			for _, instances := range reservation.Instances {
				filteredInstances = append(filteredInstances, instances)
			}
		}
	}
	return filteredInstances
}

func extractTagValue(tagList []*ec2.Tag, lookup string) string {
	tagVale := ""
	for _, tag := range tagList {
		if *tag.Key == lookup {
			tagVale = *tag.Value
			break
		}
	}
	return tagVale
}

func listVpnInstnaces(vpcCidrs map[string]string) []vpnInstance {
	var vpnInstances []vpnInstance
	vpnInstanceList := listFilteredInstances("vpn")
	for _, instance := range vpnInstanceList {
		vpn := vpnInstance{
			VpcID: *instance.VpcId,
			VpcCidr: vpcCidrs[*instance.VpcId],
			Name: extractTagValue(instance.Tags, "Name"),
			Environment: extractTagValue(instance.Tags, "environment"),
			PublicIP: *instance.PublicIpAddress,
		}
		vpnInstances = append(vpnInstances, vpn)
	}
	return vpnInstances
}

func writevpnDetailFile(vpnList []vpnInstance) {
	vpnJSON, err := json.Marshal(vpnList)
	if err != nil {
		fmt.Println(err)
		return
	}
	ioutil.WriteFile("vpn_hosts.json", vpnJSON, 0644)
}