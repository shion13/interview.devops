terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "eu-west-1"
}

locals {
  function_name = "bucket-checker"
}

# S3 bucket to store photos
resource "aws_s3_bucket" "photo" {
  bucket = "popsa-photo"
  
  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# IAM role for Lambda function
resource "aws_iam_role" "lambda_exec_role" {
  name = "${local.function_name}_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# Attach basic execution policy for Lambda
resource "aws_iam_role_policy_attachment" "lambda_policy_attach" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Lambda function definition that pulls Docker image from ECR
resource "aws_lambda_function" "photo-api" {
  function_name = "photo-api"
  role          = aws_iam_role.lambda_exec_role.arn

  package_type = "Image"
  
  image_uri = "123456789012.dkr.ecr.eu-west-1.amazonaws.com/devops-interview:latest"

  memory_size      = 1024
  timeout          = 30

  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# API Gateway to expose Lambda function via HTTPS
resource "aws_api_gateway_rest_api" "lambda_api" {
  name        = "${local.function_name}_api"
  description = "API Gateway exposing Lambda Docker function"

  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# API Gateway resource that maps to the path `/docker`
resource "aws_api_gateway_resource" "api_resource" {
  rest_api_id = aws_api_gateway_rest_api.lambda_api.id
  parent_id   = aws_api_gateway_rest_api.lambda_api.root_resource_id
  path_part   = "docker"
}

# POST method for API Gateway to invoke the Lambda function
resource "aws_api_gateway_method" "api_method" {
  rest_api_id   = aws_api_gateway_rest_api.lambda_api.id
  resource_id   = aws_api_gateway_resource.api_resource.id
  http_method   = "POST"
  authorization = "NONE"
}

# Integration of API Gateway with Lambda
resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.lambda_api.id
  resource_id = aws_api_gateway_resource.api_resource.id
  http_method = aws_api_gateway_method.api_method.http_method
  type        = "AWS_PROXY"
  integration_http_method = "POST"
  uri         = aws_lambda_function.docker_lambda.invoke_arn
}

# Deploy API Gateway
resource "aws_api_gateway_deployment" "api_deployment" {
  depends_on = [aws_api_gateway_integration.lambda_integration]
  rest_api_id = aws_api_gateway_rest_api.lambda_api.id
  stage_name  = "prod"
}

# HTTPS (TLS/SSL) enabled by default with API Gateway
resource "aws_api_gateway_stage" "https_stage" {
  deployment_id = aws_api_gateway_deployment.api_deployment.id
  rest_api_id = aws_api_gateway_rest_api.lambda_api.id
  stage_name  = aws_api_gateway_deployment.api_deployment.stage_name

  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# CloudWatch Log Group for Lambda function
resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = 14

  tags = {
    "CostCenter" = local.function_name  # Tag for cost tracking
  }
}

# Output API Gateway endpoint URL
output "api_url" {
  value       = "${aws_api_gateway_rest_api.lambda_api.execution_arn}/prod/docker"
  description = "URL of the API Gateway exposing the Lambda function"
}
