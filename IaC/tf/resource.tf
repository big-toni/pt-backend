resource "aws_security_group" "awseb_security_group" {
  description = "SecurityGroup for ElasticBeanstalk environment."
  ingress = [
    {
      cidr_blocks = "0.0.0.0/0"
      from_port = var.instance_port
      to_port = var.instance_port
      protocol = "tcp"
    }
  ]
}

resource "aws_cloudformation_stack" "awseb_update_wait_condition_eo_l_zf_m" {
  timeout_in_minutes = "43200"
  // CF Property(Count) = "1"
  // CF Property(Handle) = aws_cloudformation_stack.awseb_update_wait_condition_handle_eo_l_zf_m.id
}

resource "aws_cloudformation_stack" "awseb_instance_launch_wait_handle" {}

resource "aws_cloudformation_stack" "awseb_update_wait_condition_handle_eo_l_zf_m" {}

resource "aws_quicksight_group" "awseb_auto_scaling_group" {
  // CF Property(MinSize) = 1
  // CF Property(CapacityRebalance) = False
  // CF Property(AvailabilityZones) = [
  //   "eu-central-1a",
  //   "eu-central-1c",
  //   "eu-central-1b"
  // ]
  // CF Property(Cooldown) = "360"
  // CF Property(LaunchTemplate) = {
  //   Version = aws_launch_template.awsebec2_launch_template.latest_version
  //   LaunchTemplateId = aws_launch_template.awsebec2_launch_template.arn
  // }
  // CF Property(MaxSize) = 1
  // CF Property(Tags) = [
  //   {
  //     Value = var.awseb_environment_name
  //     Key = "elasticbeanstalk:environment-name"
  //     PropagateAtLaunch = True
  //   },
  //   {
  //     Value = var.awseb_environment_name
  //     Key = "Name"
  //     PropagateAtLaunch = True
  //   },
  //   {
  //     Value = var.awseb_environment_id
  //     Key = "elasticbeanstalk:environment-id"
  //     PropagateAtLaunch = True
  //   }
  // ]
}

resource "aws_cloudformation_stack" "awseb_instance_launch_wait_condition" {
  timeout_in_minutes = "1200"
  // CF Property(Count) = 1
  // CF Property(Handle) = aws_cloudformation_stack.awseb_instance_launch_wait_handle.id
}

resource "aws_eip" "awsebeip" {}

resource "aws_launch_template" "awsebec2_launch_template" {
  user_data = {
    SecurityGroups = [
      aws_security_group.awseb_security_group.arn
    ]
    MetadataOptions = {
      HttpPutResponseHopLimit = 2
      HttpTokens = "required"
    }
    UserData = base64encode(join("", ["Content-Type: multipart/mixed; boundary="===============5189065377222898407=="", "
", "MIME-Version: 1.0", "
", "", "
", "--===============5189065377222898407==", "
", "Content-Type: text/cloud-config; charset="us-ascii"", "
", "MIME-Version: 1.0", "
", "Content-Transfer-Encoding: 7bit", "
", "Content-Disposition: attachment; filename="cloud-config.txt"", "
", "", "
", "#cloud-config", "
", "repo_upgrade: none", "
", "repo_releasever: ", local.mappings["ContainerMeta"]["AMIMeta"]["RepoReleaseVersion"], "
", local.mappings["AWSEBOptions"]["options"]["CloudConfigOptions"], "
", "", "
", "--===============5189065377222898407==", "
", "Content-Type: text/x-shellscript; charset="us-ascii"", "
", "MIME-Version: 1.0", "
", "Content-Transfer-Encoding: 7bit", "
", "Content-Disposition: attachment; filename="user-data.txt"", "

", "#!/bin/bash", "
", "exec > >(tee -a /var/log/eb-cfn-init.log|logger -t [eb-cfn-init] -s 2>/dev/console) 2>&1", "
", "echo [`date -u +"%Y-%m-%dT%H:%M:%SZ"`] Started EB User Data", "
", "set -x", "
", "

", "function sleep_delay ", "
", "{", "
", "  if (( $SLEEP_TIME < $SLEEP_TIME_MAX )); then ", "
", "    echo Sleeping $SLEEP_TIME", "
", "    sleep $SLEEP_TIME  ", "
", "    SLEEP_TIME=$(($SLEEP_TIME * 2)) ", "
", "  else ", "
", "    echo Sleeping $SLEEP_TIME_MAX  ", "
", "    sleep $SLEEP_TIME_MAX  ", "
", "  fi", "
", "}", "

", "# Executing bootstrap script", "
", "SLEEP_TIME=2", "
", "SLEEP_TIME_MAX=3600", "
", "while true; do ", "
", "  curl ", local.mappings["AWSEBOptions"]["options"]["UserDataScript"], " > /tmp/ebbootstrap.sh ", "
", "  RESULT=$?", "
", "  if [[ "$RESULT" -ne 0 ]]; then ", "
", "    sleep_delay ", "
", "  else", "
", "    /bin/bash /tmp/ebbootstrap.sh ", "    '", aws_cloudformation_stack.awseb_instance_launch_wait_handle.id, "'", "    '", local.stack_id, "'", "    '", local.mappings["AWSEBOptions"]["options"]["AWSEBHealthdGroupId"], "'", "    '", local.mappings["AWSEBOptions"]["options"]["AWSEBHealthdEndpoint"], "'", "    '", local.mappings["AWSEBOptions"]["options"]["HealthdProxyLogLocation"], "'", "    '", local.mappings["AWSEBOptions"]["options"]["platformAssetsUrl"], "'", "    '", data.aws_region.current.name, "'", "
", "    RESULT=$?", "
", "    if [[ "$RESULT" -ne 0 ]]; then ", "
", "      sleep_delay ", "
", "    else ", "
", "      exit 0  ", "
", "    fi ", "
", "  fi ", "
", "done", "

", "--===============5189065377222898407==-- "]))
    ImageId = local.mappings["AWSEBAWSRegionArch2AMI"][data.aws_region.current.name]["hvm"]
    IamInstanceProfile = {
      Name = "aws-elasticbeanstalk-ec2-role"
    }
    InstanceType = var.instance_type
    Monitoring = {
      Enabled = False
    }
  }
}

resource "aws_cloudformation_stack" "awseb_beanstalk_metadata" {}

