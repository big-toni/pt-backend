variable "instance_type_family" {
  description = "WebServer EC2 instance type family"
  type = string
}

variable "log_publication_control" {
  description = "If true customer service logs will be published to S3."
  type = string
  default = "false"
}

variable "instance_port" {
  description = "Listen Port"
  type = string
  default = "80"
}

variable "x_ray_enabled" {
  description = "Enables AWS X-Ray for your environment."
  type = string
  default = "false"
}

variable "awseb_environment_id" {
  type = string
}

variable "hooks_pkg_url" {
  description = "Hooks package URL"
  type = string
  default = "https://elasticbeanstalk-platform-assets-eu-central-1.s3.eu-central-1.amazonaws.com/stalks/eb_docker_amazon_linux_2_1.0.2238.0_20221124145420/lib/hooks.tar.gz"
}

variable "awseb_environment_name" {
  type = string
}

variable "awseb_referrer_id" {
  type = string
}

variable "app_source" {
  description = "The url of the application source."
  type = string
  default = "https://elasticbeanstalk-platform-assets-eu-central-1.s3.eu-central-1.amazonaws.com/stalks/eb_docker_amazon_linux_2_1.0.2238.0_20221124145420/sampleapp/EBSampleApp-Docker.zip"
}

variable "proxy_server" {
  description = "Specifies which proxy server to be used for client connections."
  type = string
  default = "nginx"
}

variable "environment_variables" {
  description = "Program environment variables."
  type = string
}

variable "awseb_agent_id" {
  type = string
}

variable "instance_type" {
  description = "WebServer EC2 instance type"
  type = string
}

variable "awseb_environment_bucket" {
  type = string
}

