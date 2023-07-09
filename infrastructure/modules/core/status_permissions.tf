
resource "aws_iam_role_policy_attachment" "status_role_policy" {
    role = module.status_function.function_iam_role_name
    policy_arn = aws_iam_policy.status_policy.arn
}

resource "aws_iam_policy" "status_policy" {
    name = "status_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "dynamodb:GetItem"
            ],
            "Effect": "Allow",
            "Resource": [
                "${aws_dynamodb_table.analytics_table.arn}"
            ]
        },
        {
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Effect": "Allow",
            "Resource": [
                "*"
            ]
        },
        {
            "Action": "xray:*",
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
EOF
}