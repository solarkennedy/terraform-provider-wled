package wled

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solarkennedy/terraform-provider-wled/wled_client"
)

func resourceWLEDSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWLEDSettingsCreate,
		ReadContext:   resourceWLEDSettingsRead,
		UpdateContext: resourceWLEDSettingsUpdate,
		DeleteContext: resourceWLEDSettingsDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ui_description": {
				Type:     schema.TypeString,
				Optional: true},
		},
	}
}

func resourceWLEDSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	new_settings := wled_client.WLEDSettings{}
	host := d.Get("host").(string)
	client := wled_client.NewWLEDClient(host)

	ui_description := d.Get("ui_description")
	if ui_description != nil {
		new_settings.Description = ui_description.(string)
	}
	log.Printf("[DEBUG] Read raw new settings %+v", new_settings)
	err := client.SetSettings(new_settings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("wled_settings_" + host)
	resourceWLEDSettingsRead(ctx, d, m)
	return diags
}

func resourceWLEDSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	host := d.Get("host").(string)
	client := wled_client.NewWLEDClient(host)
	settings, err := client.GetSettings()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Read raw settings %+v", settings)

	ui_description := d.Get("ui_description")
	if ui_description != nil {
		log.Printf("[DEBUG] Setting the 'ui_description' in state to %s", settings.Description)
		if err := d.Set("ui_description", settings.Description); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceWLEDSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceWLEDSettingsCreate(ctx, d, m)
}

func resourceWLEDSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
