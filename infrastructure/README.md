todo
- [X] create deployment s3 bucket
- [X] create module for operations
- [X] zip executables into archives and store in s3 bucket
- [X] create oracle and poll function and use the s3 object
- [X] create dynamodb table and sqs queue and connect to functions
- [X] create cloudwatch composite alarm and sns topic/subscription, connect to poll function
- [X] create apiGW and connect to scheduler function
- [X] create analytics table and connect to scheduler and result
- [X] connect result function to modules
- [X] create iam roles for functions
- [X] create secret manager with openai key with secret passed in earthly