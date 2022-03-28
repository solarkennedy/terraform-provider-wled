package wled

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solarkennedy/terraform-provider-wled/wled_client"
)

func resourceWLEDPreset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWLEDPresetCreate,
		ReadContext:   resourceWLEDPresetRead,
		UpdateContext: resourceWLEDPresetUpdate,
		DeleteContext: resourceWLEDPresetDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host of the WLED device. Example: `wled.local`",
			},
			"preset_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Preset ID number. Must be >= 1",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional human readable name of the preset",
			},
		},
	}
}

func resourceWLEDPresetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	new_Preset := wled_client.WLEDPreset{}
	host := d.Get("host").(string)
	client := wled_client.NewWLEDClient(host)

	id := d.Get("preset_id")
	if id == nil {
		return diag.FromErr(fmt.Errorf("Unable to get the ID?"))
	}
	idStr := strconv.Itoa(id.(int))

	name, ok := d.GetOk("name")
	if name != nil && ok {
		new_Preset.Name = name.(string)
	}
	log.Printf("[DEBUG] Read raw new Preset %+v", new_Preset)
	err := client.SetPreset(wled_client.WLEDPresetID(idStr), new_Preset)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(idStr)
	resourceWLEDPresetRead(ctx, d, m)
	return diags
}

func resourceWLEDPresetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	host := d.Get("host").(string)
	client := wled_client.NewWLEDClient(host)

	id := d.Get("preset_id")
	if id == nil {
		return diag.FromErr(fmt.Errorf("Unable to get the ID?"))
	}
	idStr := strconv.Itoa(id.(int))

	preset, ok, err := client.GetPreset(wled_client.WLEDPresetID(idStr))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Read raw preset %+v", preset)

	if ok {
		if err := d.Set("name", preset.Name); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceWLEDPresetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceWLEDPresetCreate(ctx, d, m)
}

func resourceWLEDPresetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
