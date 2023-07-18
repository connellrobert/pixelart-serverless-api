todo
- [ ] Create route for returning POST-able s3 presigned urls for clients to upload images
- [X] Scheduler should set cloudwatch alarm associated with the db table to OK when putting objects in
- [ ] oracle should pull images from s3 prior to sending openai request
- [ ] oracle should accept only b64 from openai to store image in s3. Eliminate option from request logic stream
- [ ] oracle should delete sqs messages upon consumption
- [ ] Create terraform plan file local save for proper caching of backend infra proposed state. 
- [ ] state should be stored in aws, a backend should be created for it.
- [ ] Convert hardcoded mapping for scheduler to aws paramater store configured in terraform
- [ ] Add tracing configuration to terraform and xray sampling logic to functions
- [ ] Add status function to poll the analytics table for completed results

bugs
- [ ] status errors when there are no attempts and only pulls the first. It should check for nil values and find the latest entry in the attempts
- [ ] Poller still sends an sqs message even when erroring out. It's causing the result function to overwrite the db object and cause the entire request to be nil.