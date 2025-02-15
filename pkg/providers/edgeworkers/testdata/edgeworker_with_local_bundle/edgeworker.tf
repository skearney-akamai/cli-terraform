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

resource "akamai_edgeworker" "edgeworker" {
  name             = "test_edgeworker"
  group_id         = 1
  resource_tier_id = 2
  local_bundle     = "testdata/bundle/sampleBundle.tgz"
}

resource "akamai_edgeworkers_activation" "edgeworker_activation" {
  edgeworker_id = 123
  network       = var.env
  version       = akamai_edgeworker.edgeworker.version
}
