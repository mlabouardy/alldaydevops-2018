package main

import (
	"context"
	"encoding/json"
	"log"

	. "./utils"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

type Request struct {
	Records []struct {
		Body          string `json:"body"`
		ReceiptHandle string `json:"receiptHandle"`
	} `json:"Records"`
}

func HandleRequest(ctx context.Context, r Request) error {
	log.Println(r)

	var cluster Cluster
	json.Unmarshal([]byte(r.Records[0].Body), &cluster)

	log.Println(cluster)

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	master := cluster.Instances[0]
	workers := cluster.Instances[1:]

	err = SwarmInit(master.ID)
	if err != nil {
		return err
	}

	token, err := SwarmWorkerToken(master.ID)
	if err != nil {
		return err
	}

	log.Println(token)

	_, err = DeployVisualizer(master.ID)
	if err != nil {
		return err
	}

	for _, worker := range workers {
		err = SwarmJoinWorker([]string{worker.ID}, master.IP, token)
		if err != nil {
			return err
		}
	}

	err = UpdateDynamoDB(cfg, cluster.ID)
	if err != nil {
		return err
	}

	err = DeleteSQSMessage(cfg, r.Records[0].ReceiptHandle)
	if err != nil {
		return err
	}

	err = NotifySlack(cluster.Name)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
