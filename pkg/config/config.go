package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Config contains the values about the different parsing mechanisms for
// Config names
type Config struct {
	automaticENIConfig bool
	eniConfigName      string
	eniConfigTagName   string
	awsSession         *session.Session
}

// New will return and instance of the ENIConfig object
func New(automaticENIConfig bool, eniconfigName string, eniconfigTagName string, awsSession *session.Session) Config {
	return Config{
		automaticENIConfig: automaticENIConfig,
		eniConfigName:      eniconfigName,
		eniConfigTagName:   eniconfigTagName,
		awsSession:         awsSession,
	}
}

// GetName will return the name of the Config using either the name in the
// base object OR will load it via the AWS EC2 DescribeTags API.
func (c Config) GetName(providerID string) (string, error) {
	if c.automaticENIConfig == false {
		return c.eniConfigName, nil
	}

	_, instanceID, err := getInstanceAZandIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}

	return c.GetENIConfigTag(instanceID)
}

// GetInstanceAZ will return the instance AZ from the provider ID
func (c Config) GetInstanceAZ(providerID string) (string, error) {
	instanceAZ, _, err := getInstanceAZandIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}
	return instanceAZ, nil
}

// GetSubnetAZ will return the AZ of the subnet from the subnet ID
func (c Config) GetSubnetAZ(subnetID string) (string, error) {
	svc := ec2.New(c.awsSession)

	input := &ec2.DescribeSubnetsInput{}
	input.SetSubnetIds([]*string{&subnetID})

	output, err := svc.DescribeSubnets(input)
	if err != nil {
		return "", fmt.Errorf("could not describe subnet with id '%s'", subnetID)
	}

	//
	if len(output.Subnets) == 0 {
		return "", fmt.Errorf("subnet with id '%s' not found", subnetID)
	}

	return *output.Subnets[0].AvailabilityZone, nil
}

// GetENIConfigTag allows you to get the ENIConfig name from the EC2 instances
// tag, customize this by using the ENIConfig.eniConfigTagName
func (c Config) GetENIConfigTag(instanceID string) (name string, err error) {
	svc := ec2.New(c.awsSession)

	filterID := &ec2.Filter{}
	filterID.SetName("resource-id")
	filterID.SetValues([]*string{aws.String(instanceID)})

	input := &ec2.DescribeTagsInput{}
	input.SetFilters([]*ec2.Filter{filterID})

	output, err := svc.DescribeTags(input)
	if err != nil {
		return "", fmt.Errorf("could not describe EC2 instance tags for instance id '%s'", instanceID)
	}

	for _, tag := range output.Tags {
		if *tag.Key == c.eniConfigTagName {
			name = *tag.Value
			break
		}
	}
	return name, nil
}

// Below are custom functions to parse the InstanceID and AZ from the ProviderID
// This can probably be refactored
func getInstanceAZandIDFromProviderID(providerID string) (string, string, error) {
	u, err := url.Parse(providerID)
	if err != nil {
		return "", "", fmt.Errorf("could not parse provider ID '%s'", providerID)
	}
	instanceAZ, instanceID := parsePath(u)
	return instanceAZ, instanceID, nil
}

func parsePath(u *url.URL) (az, instanceID string) {
	trimmed := trimLeftChar(u.Path)
	parts := strings.Split(trimmed, "/")
	return parts[0], parts[1]
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
