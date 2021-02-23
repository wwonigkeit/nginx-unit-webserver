package cloud

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

// EC2CreateInstanceAPI defines the interface for the RunInstances and CreateTags functions.
// We use this interface to test the functions using a mocked service.
type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// MakeInstance creates an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a RunInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to RunInstances.
func MakeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

// MakeTags creates tags for an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateTagsOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateTags.
func MakeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

//BuildAwsInstance takes the machine image and creates the running virtual machine
func BuildAwsInstance(image string, machine string, c *websocket.Client, port int) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files

	//image := "nginx-unit-nodejs-20200209"
	//machine := "t2.nano"
	//port := 80
	//name := image

	fmt.Println(AWSCONFIG)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigFiles([]string{AWSCONFIG}))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	// First we need to fin the AMI id for the image we created

	inputImages := &ec2.DescribeImagesInput{

		Filters: []types.Filter{
			{
				Name: aws.String("name"),
				Values: []string{
					image,
				},
			},
		},
	}

	imageOutput, err := client.DescribeImages(context.TODO(), inputImages)

	if err != nil {
		fmt.Println("Got an error describing an image:")
		fmt.Println(err)
		return
	}

	fmt.Println("Found tagged image with ID " + *imageOutput.Images[0].ImageId)

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(*imageOutput.Images[0].ImageId),
		InstanceType: types.InstanceType(machine),
		MinCount:     1,
		MaxCount:     1,
	}

	result, err := MakeInstance(context.TODO(), client, input)

	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*result.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: &image,
			},
		},
	}

	_, err = MakeTags(context.TODO(), client, tagInput)
	if err != nil {
		fmt.Println("Got an error tagging the instance:")
		fmt.Println(err)
		return
	}

	//fmt.Println("Created tagged instance with ID " + *result.Instances[0].InstanceId)
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Created tagged instance with ID " + *result.Instances[0].InstanceId)}

	inputDescribeInstances := ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{*result.Instances[0].InstanceId},
	}
	now := time.Now()
	timestampstart := now.Unix()
	fmt.Println(timestampstart)

breakLoop:

	for {
		time.Sleep(5 * time.Second)
		resultDescribeInstance, err := client.DescribeInstanceStatus(context.TODO(), &inputDescribeInstances)
		if err != nil {
			fmt.Println("Got an error getting the instance information:")
			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: "Got an error getting the instance status output"}
			fmt.Println(err)
			//c.Pool.Broadcast <- websocket.Message{Type: 1, Body: err.Error}
			return
		}

		//fmt.Println(resultDescribeInstance.InstanceStatuses)

		if len(resultDescribeInstance.InstanceStatuses) == 0 {
			fmt.Println("No status for the deployed instance yet - waiting 5 seconds")
			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: "No status for the deployed instance yet - waiting 5 seconds"}
		} else if resultDescribeInstance.InstanceStatuses[0].InstanceState.Name == "running" {
			fmt.Println("Instance state is now: ", resultDescribeInstance.InstanceStatuses[0].InstanceState.Name)
			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Instance state is now: " + string(resultDescribeInstance.InstanceStatuses[0].InstanceState.Name))}

			inputGetConsole := ec2.GetConsoleOutputInput{
				InstanceId: result.Instances[0].InstanceId,
			}

			for {
				resultConsoleOutput, err := client.GetConsoleOutput(context.TODO(), &inputGetConsole)
				now := time.Now()
				timestampnow := now.Unix()

				if err != nil {
					c.Pool.Broadcast <- websocket.Message{Type: 1, Body: "Got an error getting the instance console output"}
					//c.Pool.Broadcast <- websocket.Message{Type: 1, Body: err.Error}
					fmt.Println("Got an error getting the instance console output:")
					fmt.Println(err)
					return
				} else if (resultConsoleOutput.Output == nil) && (timestampnow-timestampstart < 120) {
					fmt.Println("No console output for the deployed instance yet - waiting 5 seconds")
					c.Pool.Broadcast <- websocket.Message{Type: 1, Body: "No console output for the deployed instance yet - waiting 5 seconds"}
					time.Sleep(5 * time.Second)
				} else if timestampnow-timestampstart > 120 {
					c.Pool.Broadcast <- websocket.Message{Type: 1, Body: "Serial console log is delayed, but most likely the machine is already running - AWS is just terrible"}
					fmt.Println("Serial console log is delayed, but most likely the machine is already running - AWS is just terrible.")
					now := time.Now()
					timestampstart = now.Unix()

					inputDescribeInstances := ec2.DescribeInstancesInput{
						InstanceIds: []string{*result.Instances[0].InstanceId},
					}

					resultDescribeInstances, err := client.DescribeInstances(context.TODO(), &inputDescribeInstances)

					if err != nil {
						fmt.Println("Got an error getting the instance details output:")
						fmt.Println(err)
						return
					}

					extIP := resultDescribeInstances.Reservations[0].Instances[0].PublicIpAddress

					fmt.Println("External IP address: " + *extIP)

					mesg :=
						"Deployment of the nginx-unit-go instance on Google Cloud Platform has been completed.\n" +
							"The external IP adddress for the machine is:\n" +
							"\n" +
							"IP address: " + *extIP + "\n" +
							"\n" +
							"You can verify the configuration of the instance @\n" +
							"\n" +
							"\thttp://" + *extIP + ":" + strconv.Itoa(port) + "/\n" +
							"\n" +
							"if your application is available at this location. Alternatively the configuration for the Unit instance\n" +
							"can be changed or verified using the following URL:\n" +
							"\n" +
							"\thttp://" + *extIP + ":8080/config\n" +
							"\n"
					c.Pool.Broadcast <- websocket.Message{Type: 1, Body: mesg}
				} else {
					sDec, _ := b64.StdEncoding.DecodeString(*resultConsoleOutput.Output)
					fmt.Println(string(sDec))
					c.Pool.Broadcast <- websocket.Message{Type: 1, Body: string(sDec)}
					break breakLoop
				}
			}
		}
	}
}
