package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: nil,
		ResourcesMap: map[string]*schema.Resource{
			"mysso_permission_set": resourcePermissionSet(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mysso_permission_set":  dataSourcePermissionSet(),
			"mysso_permission_pool": dataSourcePermissionPool(),
		},
		ProviderMetaSchema: map[string]*schema.Schema{},
		ConfigureFunc:      providerConfigure,
		// ConfigureContextFunc: func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {},
		TerraformVersion: "",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	print("Provider configuration")
	return nil, nil
}
