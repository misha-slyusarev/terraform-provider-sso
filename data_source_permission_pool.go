package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/exp/maps"
)

func dataSourcePermissionPool() *schema.Resource {
	return &schema.Resource{
		Description: "PermissionPool data source in the Terraform provider mysso.",

		ReadContext: dataSourcePermissionPoolRead,

		Schema: map[string]*schema.Schema{
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"permission_set": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
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
						"tags": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"policy_attachments": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"based_on": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"relay_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"policy_attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePermissionPoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	var permissions []permissionSet
	var policyAttachments []interface{}
	var err error

	tflog.Info(ctx, " -->  Receive configured permissions  --")
	permissionSetData := d.Get("permission_set").(*schema.Set)

	tflog.Info(ctx, " -->  Aggregate permissions  --")
	if permissions, err = aggregatePermissions(permissionSetData); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't parse permission configuration",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, " -->  Receive global relay state --")
	if relayState, ok := d.GetOk("relay_state"); ok {
		for _, p := range permissions {
			if len(p["relay_state"].(string)) == 0 {
				p["relay_state"] = relayState
			}
		}
	}

	tflog.Info(ctx, " -->  Receive global tags --")
	if tags, ok := d.GetOk("tags"); ok {
		for _, p := range permissions {
			if p["tags"] == nil {
				p["tags"] = tags
			} else {
				ts := p["tags"].(map[string]interface{})
				maps.Copy(ts, tags.(map[string]interface{}))
				p["tags"] = ts
			}
		}
	}

	tflog.Info(ctx, " -->  Set permissions  --")
	if err = d.Set("permissions", permissions); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't parse permission configuration",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, " -->  Aggregate policy --")
	if policyAttachments, err = aggregatePolicy(ctx, permissionSetData); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't parse permission configuration",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, " -->  Set policy --")
	if err = d.Set("policy_attachments", policyAttachments); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't parse permission configuration",
			Detail:   err.Error(),
		})
	}

	d.SetId("1")
	return diags
}

func aggregatePermissions(ps *schema.Set) ([]permissionSet, error) {
	aggregated := make([]permissionSet, ps.Len())

	for i, v := range ps.List() {
		configuredPs, castOk := v.(map[string]interface{})
		if !castOk {
			return make([]permissionSet, 0), fmt.Errorf("unable to parse the permission_set block")
		}

		calculatedPs := make(map[string]interface{})
		calculatedPs["id"] = strconv.Itoa(i)

		if ps, ok := configuredPs["name"]; ok {
			calculatedPs["name"] = ps.(string)
		}
		if ps, ok := configuredPs["description"]; ok {
			calculatedPs["description"] = ps.(string)
		}
		if ps, ok := configuredPs["session_duration"]; ok {
			calculatedPs["session_duration"] = ps.(string)
		}
		if ps, ok := configuredPs["relay_state"]; ok {
			calculatedPs["relay_state"] = ps.(string)
		}
		if ps, ok := configuredPs["tags"]; ok {
			calculatedPs["tags"] = ps.(map[string]interface{})
		}
		aggregated[i] = calculatedPs
	}

	return aggregated, nil
}

func aggregatePolicy(ctx context.Context, ps *schema.Set) ([]interface{}, error) {
	var aggregated []interface{}

	basedOn := make(map[string][]string)

	for i, v := range ps.List() {
		configuredPs, castOk := v.(map[string]interface{})
		if !castOk {
			return make([]interface{}, 0), fmt.Errorf("unable to parse the permission_set block")
		}

		permissionSetName := configuredPs["name"].(string)
		permissionSetId := strconv.Itoa(i)
		policyArns := make([]string, 0)

		for j, policyArn := range configuredPs["policy_attachments"].([]interface{}) {
			plc := make(map[string]interface{})
			plc["id"] = permissionSetName + "-" + policyArn.(string) + "-" + strconv.Itoa(j)
			plc["policy_arn"] = policyArn
			plc["permission_set_id"] = permissionSetId
			tflog.Info(ctx, " -->  Aggregated before "+fmt.Sprintf("%v", aggregated)+" -- ")
			tflog.Info(ctx, " -->  Append "+plc["id"].(string)+" -- ")
			aggregated = append(aggregated, plc)
			tflog.Info(ctx, " -->  Aggregated after "+fmt.Sprintf("%v", aggregated)+" -- ")
			policyArns = append(policyArns, policyArn.(string))
		}

		basedOn[permissionSetName] = append(basedOn[permissionSetName], policyArns...)
	}

	// Add policies from the basedOn map if its key
	// was specified in the permission set based_on attribute
	for i, v := range ps.List() {
		configuredPs := v.(map[string]interface{})
		permissionSetName := configuredPs["name"].(string)
		configuredBasedOn := configuredPs["based_on"].([]interface{})
		permissionSetId := strconv.Itoa(i)

		if len(configuredBasedOn) > 0 {
			for _, basedOnName := range configuredBasedOn {
				if policyArns := basedOn[basedOnName.(string)]; len(policyArns) > 0 {
					for k, arn := range policyArns {
						plc := make(map[string]interface{})
						plc["id"] = permissionSetName + "-" + arn + "-" + strconv.Itoa(k)
						plc["policy_arn"] = arn
						plc["permission_set_id"] = permissionSetId
						aggregated = append(aggregated, plc)
					}
				}
			}
		}
	}

	return aggregated, nil
}

type permissionSet map[string]interface{}
