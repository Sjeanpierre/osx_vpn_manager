####Command line tool wirtten in Go to facilitate the configuration and use of l2tp/ipsec vpn connections on Mac OSX.

## Requirements
- [Configured AWS credentials with the ability to list EC2 instnaces](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs)

## Installation

If you have [Homebrew](http://brew.sh) installed, you can simply start a Terminal and run:

```bash
brew install Sjeanpierre/tools/osx_vpn_manager
```

## Usage
#### profile add - configure vpn profiles, which consist of username, password, and pre-shared key values
<img width="468" alt="1__sudo" src="https://cloud.githubusercontent.com/assets/673382/17197819/20fcc8a4-543e-11e6-8d79-26859362ac57.png">

#### profile list - list configured vpn profiles, which consist of username, password, and pre-shared key values
```
sudo vpn profile list
+-----+--------+-------------+
| ID# |  NAME  |  USERNAME   |
+-----+--------+-------------+
|   0 | prod   | sjeanpierre |
|   1 | dev    | jstevenson  |
+-----+--------+-------------+
```
#### host refresh - download details about vpn instnaces in AWS
```
sudo vpn host refresh
fetching vpc details for region: us-west-1
fetching vpc details for region: us-west-2
...
fetching instances with tag vpn in: us-west-1
fetching instances with tag vpn in: us-west-2
...
```
#### host list - list instnaces from AWS which contain the vpn substring in their name. (more host sources coming soon?)
```
$ sudo vpn host list
+-----+--------------+----------------------------------------+-------------+----------------+-----------------+
| ID# |    VPC ID    |                VPN NAME                | ENVIRONMENT |   PUBLIC IP    |    VPC CIDR     |
+-----+--------------+----------------------------------------+-------------+----------------+-----------------+
|   0 | vpc-xxxxxxxx | us-preprod-data-services-vpn           | preprod     | 59.xxx.xx.11   | 10.183.24.0/23  |
|   1 | vpc-xxxxxxxx | global-accts-prod-app-vpn              | staging     | 59.x.xx.104    | 10.183.22.0/23  |
|   2 | vpc-xxxxxxxx | xxxxxxx-libreswan-vpn                  | staging     | 59.xx.xx.250   | 10.181.208.0/24 |
|   3 | vpc-xxxxxxxx | us-preprod-apps-vpn                    | preprod     | 59.xxx.xx.54   | 10.183.26.0/23  |
|   4 | vpc-xxxxxxxx | global-accts-preprod-data-services-vpn | preprod     | 59.x.x.47      | 10.183.20.0/23  |
|   5 | vpc-xxxxxxxx | global-accts-preprod-apps-vpn          | preprod     | 59.x.xx.111    | 10.183.20.0/23  |
|   6 | vpc-xxxxxxxx | xxxxxxx-libreswan-vpn                  | staging     | 59.xxx.xxx.164 | 10.181.208.0/24 |
|   7 | vpc-xxxxxxxx | us-prod-mso-vpn                        | preprod     | 59.xxx.xx.95   | 10.183.22.0/23  |
|   8 | vpc-xxxxxxxx | us-preprod-xxx-vpn                     | preprod     | 59.xxx.xx.0    | 10.183.28.0/23  |
|   9 | vpc-xxxxxxxx | us-prod-data-services-vpn              | preprod     | 59.xxx.x.19    | 10.183.28.0/23  |
|  10 | vpc-xxxxxxxx | global-xxxxx-preprod-apps-vpn          | preprod     | 59.x.xxx.241   | 10.183.22.0/23  |
----------------------------------------------------------------------------------------------------------------
```
#### connect - Connect to vpn host from host list using ID#,VPC ID, or instnace name. Supply profile name using -p flag or setting VPN_PROFILE environment variable
```
sudo vpn connect -p prod vpc-xxxxxxxx
Connecting to VPN by ID#
connecting.........
updating route table
VPN connection to us-preprod-data-services-vpn established!!
```
