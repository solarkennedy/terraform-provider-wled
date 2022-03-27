package wled

import (
	"context"
	"strconv"
	"time"

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
	err := client.SetSettings(new_settings)
	if err != nil {
		return diag.FromErr(err)
	}
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

	if err := d.Set("ui_description", settings.Description); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func resourceWLEDSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceWLEDSettingsRead(ctx, d, m)
}

func resourceWLEDSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
