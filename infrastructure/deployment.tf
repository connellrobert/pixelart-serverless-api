// create a s3 resource
resource "random_string" "deployment_suffix" {
  length = 4
  special = false
  upper = false
  number = false
}

resource "random_string" "customer_data_suffix" {
  length = 4
  special = false
  upper = false
  number = false
}

resource "aws_s3_bucket" "deployment_bucket" {
  bucket = "${var.deployment_bucket_name}-${random_string.deployment_suffix.result}"
}

# resource "aws_s3_bucket_acl" "deployment_bucket_acl" {
#   bucket = aws_s3_bucket.deployment_bucket.id
#   acl    = "private" 

# }

# resource "aws_s3_bucket_ownership_controls" "deployment_bucket_ownership" {
#   bucket = aws_s3_bucket.deployment_bucket.id
#   rule {
#     object_ownership = "BucketOwnerPreferred"
#   }
# }

resource "aws_s3_bucket" "customer_data_bucket" {
    bucket = "${var.customer-data-bucket-name}-${random_string.customer_data_suffix.result}"
}


# resource "aws_s3_bucket_acl" "customer_data_acl" {
#   bucket = aws_s3_bucket.deployment_bucket.id
#   acl    = "private"

# }

# resource "aws_s3_bucket_ownership_controls" "customer_data_ownership" {
#   bucket = aws_s3_bucket.deployment_bucket.id
#   rule {
#     object_ownership = "BucketOwnerPreferred"
#   }
# }