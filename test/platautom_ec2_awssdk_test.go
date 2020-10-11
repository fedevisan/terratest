package test

import (
	"fmt"
	"testing"

	awsy "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/aws"
	awsx "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-example using Terratest.
func TestEC2PlatAutom(t *testing.T) {
	t.Parallel()

	// Give this EC2 Instance a unique ID for a name tag so we can distinguish it from any other EC2 Instance running
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-example-%s", random.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	//awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsRegion := "us-east-1"

	// website::tag::1::Configure Terraform setting path to Terraform code, EC2 instance name, and AWS Region.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"instance_name": expectedName,
		},

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

	aws.AddTagsToResource(t, awsRegion, instanceID, map[string]string{"testing": "testing-tag-value"})

	// Look up the tags for the given Instance ID
	instanceTags := aws.GetTagsForEc2Instance(t, awsRegion, instanceID)

	// website::tag::3::Check if the EC2 instance with a given tag and name is set.
	testingTag, containsTestingTag := instanceTags["testing"]
	assert.True(t, containsTestingTag)
	assert.Equal(t, "testing-tag-value", testingTag)

	// Verify that our expected name tag is one of the tags
	nameTag, containsNameTag := instanceTags["Name"]
	assert.True(t, containsNameTag)
	assert.Equal(t, expectedName, nameTag)

	svc := ec2.New(session.New(&awsy.Config{Region: awsy.String("us-east-1"), Credentials: credentials.NewSharedCredentials("", "default")}))

	input := &ec2.DescribeVolumesInput{}
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
	fmt.Println("Results WITHOUT Filter:")
	fmt.Println(result)
	fmt.Println("End of Results WITHOUT Filter:")
	/*
	   inputn := &ec2.DescribeVolumesInput{}
	   resultn, err := svc.DescribeVolumes(inputn)
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

	   fmt.Println(resultn)
	*/
	inputx := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			awsy.String("vpc-f9599e84"),
		},
	}

	resultx, err := svc.DescribeVpcs(inputx)
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
	fmt.Println("Results VPC vpc-f9599e84:")
	fmt.Println(resultx)
	fmt.Println("End Of Results VPC vpc-f9599e84:")

	// Describe Instance Atrributes
	inputz := &ec2.DescribeInstanceAttributeInput{
		Attribute:  awsy.String("instanceType"),
		InstanceId: awsy.String(instanceID),
	}

	resultz, err := svc.DescribeInstanceAttribute(inputz)
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

	fmt.Println(resultz)

	fmt.Println("-------------Instance ID------------------------")
	fmt.Println(instanceID)
	fmt.Println("------------------------------------------------")

	fmt.Println("---------------Public IP------------------------")
	fmt.Println(awsx.GetPublicIpOfEc2Instance(t, instanceID, awsRegion))
	fmt.Println("------------------------------------------------")

	fmt.Println("-----------Private Hostname---------------------")
	fmt.Println(awsx.GetPrivateHostnameOfEc2Instance(t, instanceID, awsRegion))
	fmt.Println("------------------------------------------------")

	fmt.Println("---------------TAGS-----------------------------")
	//fmt.Println(awsx.GetTagsForEc2Instance(t, instanceID, awsRegion))
	fmt.Println("------------------------------------------------")

}
