terraform {
  required_providers {
    wled = {
      source  = "solarkennedy/wled"
    }
  }
}

resource wled_settings "zbench2" {
  host = "wled-zbench2.local"
}


output "zbench2_settings" {
  value = wled_settings.zbench2
}