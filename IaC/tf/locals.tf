locals {
  mappings = {
    ContainerMetaARM = {
      AMIMeta = {
        RepoReleaseVersion = "2.0"
        Owner = "amazon"
        AMIVersion = "2.0.20221103"
        AMIName = "amzn2-ami-hvm-2.0.20221103.3-arm64-gp2"
      }
    }
    AWSEBAWSRegionArch2AMIBase = {
      eu-central-1 = {
        pv = ""
        graphics = ""
        gpu = ""
        hvm = "ami-034a7d2833172671e"
      }
    }
    AWSEBOptions = {
      options = {
        OptionDefinitionOverrideEnabled = True
        platformAssetsUrl = "https://elasticbeanstalk-platform-assets-eu-central-1.s3.eu-central-1.amazonaws.com/stalks/eb_docker_amazon_linux_2_1.0.2238.0_20221124145420"
        CloudConfigOptions = "cloud_final_modules:
 - [scripts-user, always]"
        DefaultsScript = "/opt/elasticbeanstalk/config/private/containercommandsupport/container-defaults.sh"
        LeaderTestScript = "/opt/elasticbeanstalk/config/private/containercommandsupport/leader-test.sh"
        AWSEBHealthdGroupId = "53a31d24-b111-4e10-a89b-635b650f44db"
        HaltStartupCommandsOnFailure = "true"
        downloadSourceBundleScriptLocation = [
          "https://s3.eu-central-1.amazonaws.com/elasticbeanstalk-env-resources-eu-central-1/eb_patching_resources/download_source_bundle.py"
        ]
        DNSNameResource = "AWSEBEIP"
        UserDataScript = "https://elasticbeanstalk-platform-assets-eu-central-1.s3.eu-central-1.amazonaws.com/stalks/eb_docker_amazon_linux_2_1.0.2238.0_20221124145420/lib/UserDataScript.sh"
        DefaultSSHPort = "22"
        LaunchType = "Migration"
        ProxyServer = "nginx"
        FastVersionDeployment = "true"
        AWSEBHealthdEndpoint = "https://elasticbeanstalk-health.eu-central-1.amazonaws.com"
        ServiceRole = "arn:aws:iam::452360699504:role/aws-elasticbeanstalk-service-role"
        EnvironmentType = "SingleInstance"
        EmbeddedConfigsetsEnabled = "true"
        EBSNSTopicArn = "arn:aws:sns:eu-central-1:452360699504:ElasticBeanstalkNotifications-Environment-p-trck-production"
        HealthdProxyLogLocation = ""
        CommandBasedLeaderElection = "true"
        nodeploymentOnStartup = "true"
        ebpatchscripturl = [
          "https://s3.eu-central-1.amazonaws.com/elasticbeanstalk-env-resources-eu-central-1/eb_patching_resources/patch_instance.py"
        ]
      }
    }
    AWSEBAWSRegionArch2AMI = {
      eu-central-1 = {
        pv = ""
        graphics = ""
        gpu = ""
        hvm = "ami-0c42688b7fc01d745"
      }
    }
    EnvironmentInfoTasks = {
      systemtail = {
        LocationPrefix = "resources/environments/logs/systemtail/"
        AutoClean = "true"
        CommandName = "CMD-SystemTailLogs"
      }
      tail = {
        LocationPrefix = "resources/environments/logs/tail/"
        AutoClean = "true"
        CommandName = "CMD-TailLogs"
      }
      publish = {
        LocationPrefix = "resources/environments/logs/publish/"
      }
      bundle = {
        LocationPrefix = "resources/environments/logs/bundle/"
        AutoClean = "true"
        CommandName = "CMD-BundleLogs"
      }
    }
    AWSEBAWSRegionArch2AMIBaseARM = {
      eu-central-1 = {
        pv = ""
        graphics = ""
        gpu = ""
        hvm = "ami-0dabdd53883de82df"
      }
    }
    AWSEBAWSRegionArch2AMIARM = {
      eu-central-1 = {
        pv = ""
        graphics = ""
        gpu = ""
        hvm = "ami-0a1a040a26881ef27"
      }
    }
    ContainerMeta = {
      AMIMeta = {
        RepoReleaseVersion = "2.0"
        Owner = "amazon"
        AMIVersion = "2.0.20221103"
        AMIName = "amzn2-ami-hvm-2.0.20221103.3-x86_64-gp2"
      }
    }
  }
  stack_id = uuidv5("dns", "cf")
}

