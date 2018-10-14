package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rs/xid"
)

type Cluster struct {
	ID        string
	Name      string
	Instances []Instance
}

type Instance struct {
	Name string
	IP   string
	ID   string
}

type CommandOutput struct {
	Status ssm.CommandInvocationStatus
	Output string
}

type DBItem struct {
	ID            string
	Name          string
	Size          int
	ClusterStatus string
	ManagerIp     string
	ManagerID     string
}

func setupInfrastructure(cfg aws.Config, name string, count int64) (Cluster, error) {
	svc := ec2.New(cfg)
	req := svc.RunInstancesRequest(&ec2.RunInstancesInput{
		ImageId:      aws.String(os.Getenv("AMI")),
		InstanceType: ec2.InstanceTypeT1Micro,
		KeyName:      aws.String(os.Getenv("KEYPAIR")),
		MaxCount:     &count,
		MinCount:     &count,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(os.Getenv("SSM_ROLE_NAME")),
		},
		SecurityGroupIds: []string{os.Getenv("SECURITY_GROUP_ID")},
		TagSpecifications: []ec2.TagSpecification{
			ec2.TagSpecification{
				ResourceType: ec2.ResourceTypeInstance,
				Tags: []ec2.Tag{
					ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String(name),
					},
				},
			},
		},
	})
	result, err := req.Send()
	if err != nil {
		return Cluster{}, err
	}

	instances := make([]Instance, 0, len(result.Instances))
	for index, instance := range result.Instances {
		instances = append(instances, Instance{
			Name: fmt.Sprintf("node-%d", (index + 1)),
			IP:   *instance.PrivateIpAddress,
			ID:   *instance.InstanceId,
		})
	}
	return Cluster{
		ID:        xid.New().String(),
		Name:      name,
		Instances: instances,
	}, nil
}

func pushToSQS(cfg aws.Config, cluster Cluster) error {
	data, _ := json.Marshal(cluster)
	raw := string(data)
	svc := sqs.New(cfg)
	req := svc.SendMessageRequest(&sqs.SendMessageInput{
		MessageBody: &raw,
		QueueUrl:    aws.String(os.Getenv("SQS_URL")),
	})
	_, err := req.Send()
	if err != nil {
		return err
	}
	return nil
}

func insertToDynamoDB(cfg aws.Config, cluster Cluster) error {
	item := DBItem{
		ID:            cluster.ID,
		Name:          cluster.Name,
		Size:          len(cluster.Instances),
		ClusterStatus: "Pending",
		ManagerIp:     cluster.Instances[0].IP,
		ManagerID:     cluster.Instances[0].ID,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	svc := dynamodb.New(cfg)
	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      av,
	})
	_, err = req.Send()
	if err != nil {
		return err
	}
	return nil
}

func DeployInfrastructure(count int64, name string) error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	cluster, err := setupInfrastructure(cfg, name, count)
	if err != nil {
		log.Println(err)
		return err
	}

	err = pushToSQS(cfg, cluster)
	if err != nil {
		log.Println(err)
		return err
	}

	err = insertToDynamoDB(cfg, cluster)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func executeCommand(instanceIds []string, command string) (string, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return "", err
	}
	svc := ssm.New(cfg)
	req := svc.SendCommandRequest(&ssm.SendCommandInput{
		InstanceIds:  instanceIds,
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]string{
			"commands": []string{command},
		},
	})
	result, err := req.Send()
	if err != nil {
		return "", err
	}
	return *result.Command.CommandId, nil
}

func getCommandOutput(instanceId string, commandId string) (CommandOutput, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return CommandOutput{}, err
	}
	svc := ssm.New(cfg)
	req := svc.GetCommandInvocationRequest(&ssm.GetCommandInvocationInput{
		CommandId:  &commandId,
		InstanceId: &instanceId,
	})
	result, err := req.Send()
	if err != nil {
		return CommandOutput{}, err
	}
	return CommandOutput{
		Status: result.Status,
		Output: *result.StandardOutputContent,
	}, nil
}

func deployContainer(instanceId string) (string, error) {
	commandId, err := executeCommand([]string{instanceId}, `docker service create -d --name alldaydevops2018 -p 80:3000 mlabouardy/alldaydevops2018`)
	if err != nil {
		return "", err
	}
	time.Sleep(time.Second * 5)
	output, err := getCommandOutput(instanceId, commandId)
	if err != nil {
		return "", err
	}
	if output.Status != ssm.CommandInvocationStatusSuccess {
		log.Println(output)
		return "", errors.New("Cannot deploy the container")
	}
	return strings.TrimSpace(output.Output), nil
}

func DeployApplication(name string) error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.New(cfg)
	req := svc.ScanRequest(&dynamodb.ScanInput{
		TableName:        aws.String(os.Getenv("TABLE_NAME")),
		FilterExpression: aws.String("#c = :cluster"),
		ExpressionAttributeNames: map[string]string{
			"#c": "Name",
		},
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":cluster": {
				S: aws.String(name),
			},
		},
		Limit: aws.Int64(1),
	})
	res, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}

	output, err := deployContainer(*res.Items[0]["ManagerID"].S)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(output)
	return nil
}
