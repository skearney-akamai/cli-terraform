terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_imaging_policy_set" "policyset" {
  name        = "some policy set"
  region      = "EMEA"
  type        = "VIDEO"
  contract_id = "ctr_123"
}

data "akamai_imaging_policy_video" "data_policy__auto" {
  policy {
    breakpoints {

      widths = [320, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
  }
}

data "akamai_imaging_policy_video" "data_policy_test_policy_video" {
  policy {
    breakpoints {

      widths = [420, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
  }
}

resource "akamai_imaging_policy_video" "policy__auto" {
  policy_id              = ".auto"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_video.data_policy__auto.json
}

resource "akamai_imaging_policy_video" "policy_test_policy_video" {
  policy_id              = "test_policy_video"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_video.data_policy_test_policy_video.json
}
