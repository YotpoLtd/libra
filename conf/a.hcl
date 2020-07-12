// Nomad Client configuration
nomad {
  address = "http://10.101.5.124:4646"
}

# backend "cloudwatch" {
#   kind     = "cloudwatch"
#   region   = "us-east-1"
# }

backend "influxtest" {
  kind = "influxdb"
  addr = "http://influxdb.yotpo.com:8086"
  useragent = "libra"
}

backend "prometheus" {
  kind = "prometheus"
  host = "https://prometheus.us-east-1.yotpo.xyz"
}

job "ec2_events_staging" {
  group "ec2_events_staging" {
    min_count = 0
    max_count = 2

    rule "increase worker count" {
      backend = "prometheus"
      metric_name = "kafka_burrow_total_lag{cluster='application-msk',group='AccountPlatforms::Events::Consumers::ShopifyOrders'}"
      comparison = "above_or_equal"
      comparison_value = 1.0
      cron = "*/30 * * * * *"
      action       = "increase_count"
      action_value = 1
    }
  }
}

//job "api-resque-account" {
//  group "account" {
//    min_count = 0
//    max_count = 2
//
//    rule "increase worker count" {
//      backend = "influxtest"
//      database_name = "yotpo_api_production"
//      measurement_name = "resque_backlog"
//      selector = "last"
//      field = "backlog"
//      where_clause = {
//       "queue" = "Account"
//      }
//      time_period = "60s"
//      comparison = "above_or_equal"
//      comparison_value = 1.0
//      cron = "*/30 * * * * *"
//      action       = "increase_count"
//      action_value = 1
//    }
//
//    rule "decrease worker count" {
//      backend = "influxtest"
//      database_name = "yotpo_api_production"
//      measurement_name = "resque_backlog"
//      selector = "last"
//      field = "backlog"
//      where_clause = {
//       "queue" = "Account"
//      }
//      time_period = "60s"
//      comparison = "below"
//      comparison_value = 1.0
//      cron = "*/30 * * * * *"
//      action       = "decrease_count"
//      action_value = 1
//    }
//  }
//}