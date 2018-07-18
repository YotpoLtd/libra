// Nomad Client configuration
nomad {
  address = "https://nomad-ui.us-east-1.yotpo.xyz"
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

backend "other-backend" {
  kind     = "graphite"
  host     = "https://my-grafana.hosted-metrics.grafana.net"
  username = "api_key"
}

job "api-resque-orderproductreviewrequestsenderjobs" {
  group "orderproductreviewrequestsenderjobs" {

    min_count = 1
    max_count = 360

//    rule "increase worker count" {
//      backend = "influxtest"
//      database_name = "yotpo_api_production"
//      measurement_name = "resque_backlog"
//      selector = "last"
//      field = "backlog"
//      where_clause = {
//        "queue" = "OrderProductReviewRequestSenderJobs"
//      }
//      time_period = "60s"
//      comparison = "above"
//      comparison_value = 0.9
//      cron = "*/30 * * * * *"
//      action       = "increase_count"
//      action_value = 32
//    }

    rule "decrease worker count" {
      backend = "influxtest"
      database_name = "yotpo_api_production"
      measurement_name = "resque_backlog"
      selector = "min"
      field = "backlog"
      where_clause = {
        "queue" = "OrderProductReviewRequestSenderJobs"
      }
      time_period = "60s"
      comparison = "below"
      comparison_value = 99999999999.0
      cron = "*/5 * * * * *"
      action       = "decrease_count"
      action_value = 32
    }
  }
}

