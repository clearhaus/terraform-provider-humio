package humio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	user, err := client.(*humio.Client).Users().GetCurrent()
	if err != nil {
		return diag.Errorf("could not get current user: %s", err)
	}

	d.SetId(user.ID)
	if err := d.Set("username", user.Username); err != nil {
		return diag.Errorf("error setting username: %s", err)
	}
	if err := d.Set("full_name", user.FullName); err != nil {
		return diag.Errorf("error setting full_name: %s", err)
	}
	if err := d.Set("email", user.Email); err != nil {
		return diag.Errorf("error setting email: %s", err)
	}
	if err := d.Set("is_root", user.IsRoot); err != nil {
		return diag.Errorf("error setting is_root: %s", err)
	}

	return nil
}
