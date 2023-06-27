
// Attach a policy to the lambda role
resource "aws_iam_role_policy_attachment" "oracle_role_policy" {
  role       = var.oracle_role_name
  policy_arn = aws_iam_policy.oracle_policy.arn
}

// Create an IAM policy with basic lambda execution permissions
resource "aws_iam_policy" "oracle_policy" {
  name = "oracle_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sqs:DeleteMessage",
        "sqs:ReceiveMessage",
        "sqs:GetQueueUrl",
        "sqs:GetQueueAttributes",
        "sqs:ListQueues"
      ],
      "Effect": "Allow",
      "Resource": "${var.queue_arn}"
    },
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

// Attach a policy to the lambda role
resource "aws_iam_role_policy_attachment" "poll_role_policy" {
  role       = var.poll_role_name
  policy_arn = aws_iam_policy.poll_policy.arn
}


// Create an IAM policy with basic lambda execution permissions
resource "aws_iam_policy" "poll_policy" {
  name = "poll_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sqs:PutMessage",
        "sqs:GetQueueUrl",
        "sqs:GetQueueAttributes",
        "sqs:ListQueues"
      ],
      "Effect": "Allow",
      "Resource": "${var.queue_arn}"
    },
    {
      "Action": [
        "dynamodb:DeleteItem",
        "dynamodb:GetItem",
        "dynamodb:Scan",
        "dynamodb:Query",
        "dynamodb:DescribeTable",
        "dynamodb:BatchGetItem",
        "dynamodb:ListTables"
      ],
      "Effect": "Allow",
      "Resource": "${var.table_arn}"
    },
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}