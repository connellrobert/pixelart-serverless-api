todo
- [ ] Create route for returning POST-able s3 presigned urls for clients to upload images
- [X] Scheduler should set cloudwatch alarm associated with the db table to OK when putting objects in
- [ ] oracle should pull images from s3 prior to sending openai request
- [ ] oracle should accept only b64 from openai to store image in s3. Eliminate option from request logic stream
- [ ] oracle should delete sqs messages upon consumption
- [ ] Create terraform plan file local save for proper caching of backend infra proposed state. 
- [ ] state should be stored in aws, a backend should be created for it.