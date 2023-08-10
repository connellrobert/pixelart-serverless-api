terraform {
    required_version = ">= 0.12"
    backend "s3" {
        bucket = "terraform-shitty-shit-backend-stuff"
        key = "pixelart/tf-backend"
        region = "us-east-1"
    }
    required_providers {
        aws = {
        source  = "hashicorp/aws"
        version = "~> 3.0"
        }
    }
}

provider "aws" {
    region = "us-east-1"
}
