package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Description: "PermissionSet data source in the Terraform provider mysso.",

		ReadContext: dataSourcePermissionSetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"session_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered permission set configuration",
			},
		},
	}
}

func dataSourcePermissionSetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := d.Set("rendered", "my permission set")
	if err != nil {
		return diag.Errorf("key is invalid or the value is not a correct type")
	}

	d.SetId("1")
	return nil
}
