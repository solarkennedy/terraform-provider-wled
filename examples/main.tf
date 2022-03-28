terraform {
  required_providers {
    wled = {
      source  = "solarkennedy/wled"
    }
  }
}

resource wled_settings "zbench2" {
  host = "10.0.2.119"
  ui_description = "zBench5"
}


output "zbench2_settings" {
  value = wled_settings.zbench2
}