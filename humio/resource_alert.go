// Copyright © 2020 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package humio

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertCreate,
		ReadContext:   resourceAlertRead,
		UpdateContext: resourceAlertUpdate,
		DeleteContext: resourceAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"alert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"throttle_time_millis": {
				Type:     schema.TypeInt,
				Required: true,
				// TODO: Figure out if we want to accept similar input as "start", if yes reuse the ValidateDiagFunc for "start" and rename this field.
			},
			"throttle_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: Add ValidateDiagFunc, we accept only digits followed by a unit, e.g. 5m, 24h, 2w
				// ValidateDiagFunc:
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"run_as_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"query_ownership_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					value := v.(string)
					if value == "Organization" || value == "User" || value == "" {
						return nil
					}
					return diag.Errorf("query_ownership_type must be 'User' or 'Organization' (case sensitive)")
				},
			},
		},
	}
}

func resourceAlertCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	alert, err := alertFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	_, err = client.(*humio.Client).Alerts().Add(
		d.Get("repository").(string),
		&alert,
	)
	if err != nil {
		return diag.Errorf("could not create alert: %s", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository"), d.Get("name")))

	return resourceAlertRead(ctx, d, client)
}

func resourceAlertRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		parts := parseRepositoryAndID(d.Id())
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return diag.Errorf("error importing humio_alert. Please make sure the ID is in the form REPOSITORYNAME+ALERTNAME (i.e. myRepoName+myAlertName")
		}
		err := d.Set("repository", parts[0])
		if err != nil {
			return diag.Errorf("error setting repository for resource %s: %s", d.Id(), err)
		}
		err = d.Set("name", parts[1])
		if err != nil {
			return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
		}
	}

	alert, err := client.(*humio.Client).Alerts().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil {
		return diag.Errorf("could not get alert: %s", err)
	}
	return resourceDataFromAlert(alert, d)
}

func resourceDataFromAlert(a *humio.Alert, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("alert_id", a.ID)
	if err != nil {
		return diag.Errorf("error setting alert_id for resource: %s, %s", d.Id(), err)
	}
	err = d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
	}
	err = d.Set("description", a.Description)
	if err != nil {
		return diag.Errorf("error setting description for resource %s: %s", d.Id(), err)
	}
	err = d.Set("throttle_time_millis", a.ThrottleTimeMillis)
	if err != nil {
		return diag.Errorf("error setting throttle_time_millis for resource %s: %s", d.Id(), err)
	}
	err = d.Set("throttle_field", a.ThrottleField)
	if err != nil {
		return diag.Errorf("error setting throttle_field for resource %s: %s", d.Id(), err)
	}
	err = d.Set("enabled", a.Enabled)
	if err != nil {
		return diag.Errorf("error setting enabled for resource %s: %s", d.Id(), err)
	}
	err = d.Set("actions", a.Actions)
	if err != nil {
		return diag.Errorf("error setting actions for resource %s: %s", d.Id(), err)
	}
	err = d.Set("labels", a.Labels)
	if err != nil {
		return diag.Errorf("error setting labels for resource %s: %s", d.Id(), err)
	}
	err = d.Set("query", a.QueryString)
	if err != nil {
		return diag.Errorf("error setting query for resource %s: %s", d.Id(), err)
	}
	err = d.Set("start", a.QueryStart)
	if err != nil {
		return diag.Errorf("error setting start for resource %s: %s", d.Id(), err)
	}
	err = d.Set("run_as_user_id", a.RunAsUserID)
	if err != nil {
		return diag.Errorf("error setting run_as_user_id for resource %s: %s", d.Id(), err)
	}
	err = d.Set("query_ownership_type", a.QueryOwnershipType)
	if err != nil {
		return diag.Errorf("error setting query_ownership_type for resource %s: %s", d.Id(), err)
	}
	return nil
}

func resourceAlertUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	alert, err := alertFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	_, err = client.(*humio.Client).Alerts().Update(
		d.Get("repository").(string),
		&alert,
	)
	if err != nil {
		return diag.Errorf("could not update alert: %s", err)
	}

	return resourceAlertRead(ctx, d, client)
}

func resourceAlertDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	alert, err := alertFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	err = client.(*humio.Client).Alerts().Delete(
		d.Get("repository").(string),
		alert.Name,
	)
	if err != nil {
		return diag.Errorf("could not delete alert: %s", err)
	}
	return nil
}

func alertFromResourceData(d *schema.ResourceData) (humio.Alert, error) {
	return humio.Alert{
		ID:                 d.Get("alert_id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		ThrottleTimeMillis: d.Get("throttle_time_millis").(int),
		ThrottleField:      d.Get("throttle_field").(string),
		Enabled:            d.Get("enabled").(bool),
		Actions:            convertInterfaceListToStringSlice(d.Get("actions").([]interface{})),
		Labels:             convertInterfaceListToStringSlice(d.Get("labels").([]interface{})),
		QueryString:        d.Get("query").(string),
		QueryStart:         d.Get("start").(string),
		RunAsUserID:        d.Get("run_as_user_id").(string),
		QueryOwnershipType: d.Get("query_ownership_type").(string),
	}, nil
}

func convertInterfaceListToStringSlice(s []interface{}) []string {
	var element []string
	for _, item := range s {
		value, _ := item.(string)
		element = append(element, value)
	}
	return element
}
