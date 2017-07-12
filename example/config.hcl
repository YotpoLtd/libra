// Nomad Client configuration
nomad {
  address = "http://localhost:4646"
}

backend "test-backend" {
  kind     = "cloudwatch"
  region   = "us-east-1"
  aws_access_key_id = ""
  aws_secret_access_key = ""
}

backend "other-backend" {
  kind     = "graphite"
  host     = "https://tsdb-56767-ua.hosted-metrics.grafana.net"
  username = "api_key"  
}

// Scale for the job "content-themes"
job "nginx" {
  // For group "populator"
  group "nginx" {
    // (required) The minimum nuber of tasks to run for this job
    min_count = 1

    // (required) The maximum number of tasks to run for this job
    max_count = 3

    // Scale by a rule
    // Again the name does not matter for execution
    rule "cloudwatch asg cpu usage upper bound" {
      // (required) What backend to use, this will define which configuration
      // is valid and which checks you can execute
      backend = "test-backend"

      // (required) The check type, could be anything implemented by the backend
      //
      // Example for cloudwatch
      //   - "CPUUtilization"
      dimension_name = "AutoScalingGroupName"
      dimension_value = "infra-2-us-east-1-consul-httpapi-asg"

      metric_namespace = "AWS/EC2"
      metric_name = "CPUUtilization"

      // (required) The comparison do do, this supports the basic match operations like
      //
      //   - above (>)
      //   - below (<)
      //   - equal (==)
      //   - not_equal (!=)
      //   - above_or_equal (>=)
      //   - below_or_equal (<=)
      comparison = "above_or_equal"

      // (required) The value to compare to, this should be a float or integer
      comparison_value = 90.0

      // (optional) You can define how often this rule should be checked, by default it will me checked every minute
      cron = "* * * * *"

      action       = "increase_count"
      action_value = 1
    }

    rule "cloudwatch asg cpu usage lower bound" {
      backend          = "test-backend"
      dimension_name   = "AutoScalingGroupName"
      dimension_value  = "infra-2-us-east-1-consul-httpapi-asg"
      metric_namespace = "AWS/EC2"
      metric_name      = "CPUUtilization"
      comparison       = "below"
      comparison_value = 20.0
      cron             = "* * * * *"
      action           = "decrease_count"
      action_value     = 1
    }
    
    rule "graphite ouath lower bound" {
      backend          = "other-backend"
      metric_name      = "stats.prod.panama.Fitbit.oauth2.api_request._api_request.count"
      comparison       = "below"
      comparison_value = 20.0
      cron             = "* * * * *"
      action           = "decrease_count"
      action_value     = 1
    }
  }
}
