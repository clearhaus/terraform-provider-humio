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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	humio "github.com/humio/cli/api"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
			"link_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"silenced": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"throttle_time_millis": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notifiers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, client interface{}) error {
	alert, err := alertFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain alert from resource data: %v", err)
	}

	_, err = client.(*humio.Client).Alerts().Add(d.Get("repository").(string), &alert, false)
	if err != nil {
		return fmt.Errorf("could not create alert: %v", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository"), d.Get("name")))

	return resourceAlertRead(d, client)
}

func resourceAlertRead(d *schema.ResourceData, client interface{}) error {
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		parts := parseRepositoryAndName(d.Id())
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("Error Importing humio_alert. Please make sure the ID is in the form REPOSITORYNAME+ALERTNAME (i.e. myRepoName+myAlertName")
		}
		d.Set("repository", parts[0])
		d.Set("name", parts[1])
	}

	alert, err := client.(*humio.Client).Alerts().Get(d.Get("repository").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("could not get alert: %v", err)
	}
	resourceDataFromAlert(alert, d)
	return nil
}

func resourceDataFromAlert(a *humio.Alert, d *schema.ResourceData) error {
	d.Set("name", a.Name)
	d.Set("description", a.Description)
	d.Set("throttle_time_millis", a.ThrottleTimeMillis)
	d.Set("silenced", a.Silenced)
	d.Set("notifiers", a.Notifiers)
	d.Set("link_url", a.LinkURL)
	d.Set("labels", a.Labels)
	d.Set("query", a.Query.QueryString)
	d.Set("start", a.Query.Start)
	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, client interface{}) error {
	alert, err := alertFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain alert from resource data: %v", err)
	}

	_, err = client.(*humio.Client).Alerts().Add(d.Get("repository").(string), &alert, true)
	if err != nil {
		return fmt.Errorf("could not create alert: %v", err)
	}

	return resourceAlertRead(d, client)
}

func resourceAlertDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).Alerts().Delete(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete alert: %v", err)
	}
	return nil
}

func alertFromResourceData(d *schema.ResourceData, client interface{}) (humio.Alert, error) {
	return humio.Alert{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		ThrottleTimeMillis: d.Get("throttle_time_millis").(int),
		Silenced:           d.Get("silenced").(bool),
		Notifiers:          convertInterfaceListToStringSlice(d.Get("notifiers").([]interface{})),
		LinkURL:            d.Get("link_url").(string),
		Labels:             convertInterfaceListToStringSlice(d.Get("labels").([]interface{})),
		Query: humio.HumioQuery{
			QueryString: d.Get("query").(string),
			Start:       d.Get("start").(string),
			End:         "now",
			IsLive:      true,
		},
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
