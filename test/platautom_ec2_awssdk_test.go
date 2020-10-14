package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	//  awsx "github.com/gruntwork-io/terratest/modules/aws"
	//  This "awsx" is the package "github.com/gruntwork-io/terratest/modules/aws". Finally, it is not necessary use "awsx." because
	//  the Objects of "github.com/aws/aws-sdk-go/aws" are identified by "awssdk."

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-example using Terratest.

// Define Struct for JSON Output Parsing

type Vol struct {
	Type    string   `json:"type,omitempty"`
	Volumes []Volume `json:"Volumes,omitempty"`
}

type Volume struct {
	Attachments      []Attachment `json:"Attachments,omitempty"`
	Size             int          `json:"Size,omitempty"`
	AvailabilityZone string       `json:"AvailabilityZone,omitempty"`
}

type Attachment struct {
	Device string `json:"Device,omitempty"`
}

func TestEC2PlatAutom(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	//awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsRegion := "us-east-1"

	// website::tag::1::Configure Terraform setting path to Terraform code, EC2 instance name, and AWS Region.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-example",

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	// website::tag::4::At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2::Run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceID := terraform.Output(t, terraformOptions, "instance_id")

	svc := ec2.New(session.New(&awssdk.Config{Region: awssdk.String("us-east-1"), Credentials: credentials.NewSharedCredentials("", "default")}))

	//----------- Describe Volumes attached to this specific instance ----------------------
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: awssdk.String("attachment.instance-id"),
				Values: []*string{
					awssdk.String(instanceID),
				},
			},
		},
	}
	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	//fmt.Println(result)

	// ------------------ Convert to JSON Format the AWS Output of Volume Description --------------------------
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Print the JSON Marshal Output
	os.Stdout.Write(b)

	// Parse JSON Output and Get the Volume Size of Instance created in this execution
	m := []byte(b)

	r := bytes.NewReader(m)
	decoder := json.NewDecoder(r)

	val := &Vol{}
	error := decoder.Decode(val)

	if error != nil {
		log.Fatal(error)
	}

	// If you want to read a response body
	// decoder := json.NewDecoder(res.Body)
	// err := decoder.Decode(val)

	// ------------ Print Volume Size and Compare. Volumes is a slice so you must loop over it. -----------------------------------
	for _, s := range val.Volumes {

		fmt.Println("\n-------------Volume Size------------------------")
		fmt.Println(s.Size)
		fmt.Println("-------------------------------------------------")
		/*    for _, a := range s.Attachments {
		      fmt.Println(a.Device)
		  }*/
		//Verify that the volume size is 8
		assert.Equal(t, int(10), s.Size)

	}

}
