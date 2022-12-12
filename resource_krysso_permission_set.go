package main

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionSetCreate,
		Read:   resourcePermissionSetRead,
		Update: resourcePermissionSetUpdate,
		Delete: resourcePermissionSetDelete,

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
		},
	}

	//TODO: add follwoing to the schema
	//"tags", "inline_policy", "policy_attachments"
}

func resourcePermissionSetCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePermissionSetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePermissionSetUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePermissionSetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
