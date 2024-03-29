package wled

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"wled_settings": resourceWLEDSettings(),
			"wled_preset":   resourceWLEDPreset(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
