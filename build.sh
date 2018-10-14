#!/bin/bash

echo "Build deployment package for CreateInfrastructure function:"
GOOS=linux go build -o create/main create/*.go
zip -r create/deployment.zip create/main
rm create/main

echo "Build deployment package for ProvisionInfrastructure function:"
GOOS=linux go build -o provision/main provision/*.go
zip -r provision/deployment.zip provision/main
rm provision/main