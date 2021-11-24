terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name = "test_policy_export"
  cloudlet_code = "ER"
  description = "Testing\\ exported policy"
  group_id = "12345"
  match_rule_format = "1.0"
  match_rules = data.akamai_cloudlets_edge_redirector_match_rule.match_rules_er.json
}

data "akamai_cloudlets_edge_redirector_match_rule" "match_rules_er" {
  match_rules {
    name = "\\r2"
    start = 0
    end = 0
    matches {
      match_type = "header"
      match_value = "value\\"
      match_operator = "equals"
      case_sensitive = false
      negate = false
      check_ips = ""
      object_match_value {
        name = "ER\\"
        type = "object"
        name_case_sensitive = false
        name_has_wildcard = false
        options {
          value = ["\\y"]
          value_has_wildcard = true
          value_case_sensitive = false
          value_escaped = false
        }
      }
    }
    use_relative_url = "copy_scheme_hostname"
    status_code = 301
    redirect_url = "/\\ddd"
    match_url = "abc.\\com"
    use_incoming_query_string = false
    use_incoming_scheme_and_host = true
  }
}
