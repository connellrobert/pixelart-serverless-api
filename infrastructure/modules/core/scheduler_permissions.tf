
resource "aws_iam_role_policy_attachment" "scheduler_role_policy" {
    role = module.scheduler_function.function_iam_role_name
    policy_arn = aws_iam_policy.scheduler_policy.arn
}

resource "aws_iam_policy" "scheduler_policy" {
    name = "scheduler_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "dynamodb:PutItem"
            ],
            "Effect": "Allow",
            "Resource": [
                "${aws_dynamodb_table.analytics_table.arn}",
                "${var.gi_table_arn}"
            ]
        },
        {
            "Action": [
                "cloudwatch:SetAlarmState"
            ],
            "Effect": "Allow",
            "Resource": [
                "${var.gi_empty_db_alarm_arn}"
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