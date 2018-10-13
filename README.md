

# Lambda Functions

## Create Infrastructure 

### Environment Variables

- AMI
- KEYPAIR
- SSM_ROLE_NAME
- SECURITY_GROUP_ID
- SQS_URL
- TABLE_NAME

### IAM Policy

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": [
                "iam:PassRole",
                "dynamodb:PutItem",
                "dynamodb:Scan",
                "sqs:SendMessage",
                "sqs:SetQueueAttributes"
            ],
            "Resource": [
                "arn:aws:sqs:AWS_REGION:ACCOUNT_ID:clusters",
                "arn:aws:iam::ACCOUNT_ID:role/SwarmClusterSSMRole",
                "arn:aws:dynamodb:AWS_REGION:ACCOUNT_ID:table/clusters"
            ]
        },
        {
            "Sid": "2",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "ec2:CreateTags",
                "ec2:RunInstances",
                "logs:CreateLogGroup",
                "logs:PutLogEvents",
                "ssm:SendCommand",
                "ssm:GetCommandInvocation"
            ],
            "Resource": "*"
        }
    ]
}
```

## Provision Infrastructure

### Environment Variables

- SQS_URL
- TABLE_NAME

### IAM Policy

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": [
                "sqs:*",
                "dynamodb:UpdateItem"
            ],
            "Resource": [
                "arn:aws:sqs:AWS_REGION:ACCOUNT_ID:clusters",
                "arn:aws:dynamodb:AWS_REGION:ACCOUNT_ID:table/clusters"
            ]
        },
        {
            "Sid": "2",
            "Effect": "Allow",
            "Action": [
                "ssm:SendCommand",
                "logs:CreateLogStream",
                "logs:CreateLogGroup",
                "logs:PutLogEvents",
                "ssm:GetCommandInvocation"
            ],
            "Resource": "*"
        }
    ]
}
```

# Maintainers

Mohamed Labouardy <mohamed@labouardy.com>