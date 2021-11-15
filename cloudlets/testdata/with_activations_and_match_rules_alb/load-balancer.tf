resource "akamai_cloudlets_application_load_balancer" "test_origin" {
  origin_id = "test_origin"
  description = "test description"
  balancing_type = "WEIGHTED"

  data_centers {
    latitude = 102.78108
    longitude = -116.07064
    continent = "NA"
    country = "US"
    origin_id = "test_origin"
    percent = 10
    cloud_service = true
    liveness_hosts = ["tf1.test", "tf2.test"]
    hostname = "test-hostname"
    state_or_province = "MA"
    city = "Boston"
    cloud_server_host_header_override = false
  }

  liveness_settings {
    port = 1234
    protocol = "HTTP"
    path = "/status"
    host_header = "header"
    additional_headers = {
      abc = "123"
    }
    interval = 10
    peer_certificate_verification = false
    request_string = "test_request_string"
    response_string = "test_response_string"
    status_3xx_failure = false
    status_4xx_failure = false
    status_5xx_failure = false
    timeout = 60
  }
}

resource "akamai_cloudlets_application_load_balancer_activation" "test_origin_PRODUCTION" {
  origin_id = "test_origin"
  network = "PRODUCTION"
  version = 2
}

resource "akamai_cloudlets_application_load_balancer_activation" "test_origin_STAGING" {
  origin_id = "test_origin"
  network = "STAGING"
  version = 2
}
