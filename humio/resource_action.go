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
	"net/http"
	"reflect"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func resourceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActionCreate,
		ReadContext:   resourceActionRead,
		UpdateContext: resourceActionUpdate,
		DeleteContext: resourceActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					humio.ActionTypeEmail,
					humio.ActionTypeHumioRepo,
					humio.ActionTypeOpsGenie,
					humio.ActionTypePagerDuty,
					humio.ActionTypeSlack,
					humio.ActionTypeSlackPostMessage,
					humio.ActionTypeVictorOps,
					humio.ActionTypeWebhook,
				}, false)),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"recipients": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
									v := val.(string)
									if len(v) > 254 || !rxEmail.MatchString(v) {
										return diag.FromErr(fmt.Errorf("%q must be a valid email, got: %s", key, v))
									}
									return nil
								},
							},
						},
						"subject_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"humiorepo": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingest_token": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"opsgenie": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_url": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "https://api.opsgenie.com",
							ValidateDiagFunc: validateURL,
						},
						"genie_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"pagerduty": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routing_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"critical",
								"error",
								"warning",
								"info",
							}, false)),
						},
					},
				},
			},
			"slack": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
			"slackpostmessage": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Required: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fields": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"use_proxy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"victorops": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
			"webhook": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_template": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"actionID\": \"{alert_action_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }",
						},
						"headers": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								http.MethodGet,
								http.MethodPost,
								http.MethodPut,
							}, false)),
						},
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
		},
	}
}

func resourceActionCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	action, err := actionFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain action from resource data: %s", err)
	}

	a, err := client.(*humio.Client).Actions().Add(
		d.Get("repository").(string),
		&action,
	)
	if err != nil {
		return diag.Errorf("could not create action: %s", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository").(string), a.Name))

	return resourceActionRead(ctx, d, client)
}

func resourceActionRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	parts := parseRepositoryAndID(d.Id())
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return diag.Errorf("error importing humio_action. Please make sure the ID is in the form REPOSITORYNAME+NOTIFIERID (i.e. myRepoName+12345678901234567890123456789012")
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

	action, err := client.(*humio.Client).Actions().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil || reflect.DeepEqual(*action, humio.Action{}) {
		return diag.Errorf("could not get action: %s", err)
	}
	return resourceDataFromAction(action, d)
}

func resourceDataFromAction(a *humio.Action, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("action_id", a.ID)
	if err != nil {
		return diag.Errorf("could not set action_id for action: %s", err)
	}
	err = d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("could not set name for action: %s", err)
	}
	err = d.Set("type", a.Type)
	if err != nil {
		return diag.Errorf("could not set type for action: %s", err)
	}

	switch a.Type {
	case humio.ActionTypeEmail:
		if err := d.Set("email", emailFromAction(a)); err != nil {
			return diag.Errorf("error setting email settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeHumioRepo:
		if err := d.Set("humiorepo", humiorepoFromAction(a)); err != nil {
			return diag.Errorf("error setting humiorepo settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeOpsGenie:
		if err := d.Set("opsgenie", opsgenieFromAction(a)); err != nil {
			return diag.Errorf("error setting opsgenie settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypePagerDuty:
		if err := d.Set("pagerduty", pagerdutyFromAction(a)); err != nil {
			return diag.Errorf("error setting pagerduty settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeSlack:
		if err := d.Set("slack", slackFromAction(a)); err != nil {
			return diag.Errorf("error setting slack settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeSlackPostMessage:
		if err := d.Set("slackpostmessage", slackpostmessageFromAction(a)); err != nil {
			return diag.Errorf("error setting slackpostmessage settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeVictorOps:
		if err := d.Set("victorops", victoropsFromAction(a)); err != nil {
			return diag.Errorf("error setting victorops settings for resource %s: %s", d.Id(), err)
		}
	case humio.ActionTypeWebhook:
		if err := d.Set("webhook", webhookFromAction(a)); err != nil {
			return diag.Errorf("error setting webhook settings for resource %s: %s", d.Id(), err)
		}
	default:
		return diag.Errorf("unsupported action type: %s", a.Type)
	}

	return nil
}

func resourceActionUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	action, err := actionFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain action from resource data: %s", err)
	}

	_, err = client.(*humio.Client).Actions().Update(
		d.Get("repository").(string),
		&action,
	)
	if err != nil {
		return diag.Errorf("could not update action: %s", err)
	}

	return resourceActionRead(ctx, d, client)
}

func resourceActionDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	action, err := actionFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain action from resource data: %s", err)
	}

	err = client.(*humio.Client).Actions().Delete(
		d.Get("repository").(string),
		action.Name,
	)
	if err != nil {
		return diag.Errorf("could not delete action: %s", err)
	}
	return nil
}

// actionFromResourceData returns a humio.Action based on either the new change or the current state depending on update bool.
func actionFromResourceData(d *schema.ResourceData) (humio.Action, error) {
	action := humio.Action{
		Type: d.Get("type").(string),
		ID:   d.Get("action_id").(string),
		Name: d.Get("name").(string),
	}

	switch d.Get("type") {
	case humio.ActionTypeEmail:
		properties := getActionPropertiesFromResourceData(d, "email", "recipients")
		var recipients []string
		for _, recipient := range properties[0]["recipients"].([]interface{}) {
			recipients = append(recipients, recipient.(string))
		}
		action.EmailAction = humio.EmailAction{
			Recipients:      recipients,
			BodyTemplate:    properties[0]["body_template"].(string),
			SubjectTemplate: properties[0]["subject_template"].(string),
		}
	case humio.ActionTypeHumioRepo:
		properties := getActionPropertiesFromResourceData(d, "humiorepo", "ingest_token")
		action.HumioRepoAction = humio.HumioRepoAction{
			IngestToken: properties[0]["ingest_token"].(string),
		}
	case humio.ActionTypeOpsGenie:
		properties := getActionPropertiesFromResourceData(d, "opsgenie", "genie_key")
		action.OpsGenieAction = humio.OpsGenieAction{
			ApiUrl:   properties[0]["api_url"].(string),
			GenieKey: properties[0]["genie_key"].(string),
		}
	case humio.ActionTypePagerDuty:
		properties := getActionPropertiesFromResourceData(d, "pagerduty", "routing_key")
		action.PagerDutyAction = humio.PagerDutyAction{
			RoutingKey: properties[0]["routing_key"].(string),
			Severity:   properties[0]["severity"].(string),
		}
	case humio.ActionTypeSlack:
		properties := getActionPropertiesFromResourceData(d, "slack", "url")
		fields := []humio.SlackFieldEntryInput{}
		for fieldName, value := range properties[0]["fields"].(map[string]interface{}) {
			fields = append(fields, humio.SlackFieldEntryInput{
				FieldName: fieldName,
				Value:     value.(string),
			})
		}
		action.SlackAction = humio.SlackAction{
			Url:    properties[0]["url"].(string),
			Fields: fields,
		}
	case humio.ActionTypeSlackPostMessage:
		properties := getActionPropertiesFromResourceData(d, "slackpostmessage", "api_token")
		fields := []humio.SlackFieldEntryInput{}
		for fieldName, value := range properties[0]["fields"].(map[string]interface{}) {
			fields = append(fields, humio.SlackFieldEntryInput{
				FieldName: fieldName,
				Value:     value.(string),
			})
		}
		channels := []string{}
		for _, channel := range properties[0]["channels"].([]interface{}) {
			channels = append(channels, channel.(string))
		}
		action.SlackPostMessageAction = humio.SlackPostMessageAction{
			ApiToken: properties[0]["api_token"].(string),
			Channels: channels,
			Fields:   fields,
			UseProxy: properties[0]["use_proxy"].(bool),
		}
	case humio.ActionTypeVictorOps:
		properties := getActionPropertiesFromResourceData(d, "victorops", "notify_url")
		action.VictorOpsAction = humio.VictorOpsAction{
			MessageType: properties[0]["message_type"].(string),
			NotifyUrl:   properties[0]["notify_url"].(string),
		}
	case humio.ActionTypeWebhook:
		properties := getActionPropertiesFromResourceData(d, "webhook", "url")
		headers := []humio.HttpHeaderEntryInput{}
		for header, value := range properties[0]["headers"].(map[string]interface{}) {
			headers = append(headers, humio.HttpHeaderEntryInput{
				Header: header,
				Value:  value.(string),
			})
		}
		action.WebhookAction = humio.WebhookAction{
			BodyTemplate: properties[0]["body_template"].(string),
			Headers:      headers,
			Method:       properties[0]["method"].(string),
			Url:          properties[0]["url"].(string),
		}
	default:
		return humio.Action{}, fmt.Errorf("unsupported action type: %s", d.Get("type"))
	}

	return action, nil
}

// getActionPropertiesFromResourceData returns the first non-empty set of action properties related to a given action.
// We do this as a workaround for an issue where we get a list longer than 1 which should not happen given MaxItems is
// set to 1 in the schema definition.
func getActionPropertiesFromResourceData(d *schema.ResourceData, actionName, requiredPropertyName string) []tfMap {
	_, newProperties := d.GetChange(actionName)
	newPropertiesList := newProperties.(*schema.Set).List()
	if len(newPropertiesList) == 0 {
		properties := d.Get(actionName).(*schema.Set).List()[0]
		return []tfMap{properties.(tfMap)}
	}
	for idx := range newPropertiesList {
		if newPropertiesList[idx].(tfMap)[requiredPropertyName] != "" {
			return []tfMap{newPropertiesList[idx].(tfMap)}
		}
	}

	return []tfMap{}
}

func emailFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	s["recipients"] = a.EmailAction.Recipients
	s["body_template"] = a.EmailAction.BodyTemplate
	s["subject_template"] = a.EmailAction.SubjectTemplate
	return []tfMap{s}
}

func humiorepoFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	s["ingest_token"] = a.HumioRepoAction.IngestToken
	return []tfMap{s}
}

func opsgenieFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	s["api_url"] = a.OpsGenieAction.ApiUrl
	s["genie_key"] = a.OpsGenieAction.GenieKey
	return []tfMap{s}
}

func pagerdutyFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	s["routing_key"] = a.PagerDutyAction.RoutingKey
	s["severity"] = a.PagerDutyAction.Severity
	return []tfMap{s}
}

func slackFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	fields := make(map[string]string)
	for _, field := range a.SlackAction.Fields {
		fields[field.FieldName] = field.Value
	}
	s["fields"] = fields
	s["url"] = a.SlackAction.Url
	return []tfMap{s}
}

func slackpostmessageFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	fields := make(map[string]string)
	for _, field := range a.SlackPostMessageAction.Fields {
		fields[field.FieldName] = field.Value
	}
	s["api_token"] = a.SlackPostMessageAction.ApiToken
	s["channels"] = a.SlackPostMessageAction.Channels
	s["fields"] = fields
	s["use_proxy"] = a.SlackPostMessageAction.UseProxy
	return []tfMap{s}
}

func victoropsFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	s["message_type"] = a.VictorOpsAction.MessageType
	s["notify_url"] = a.VictorOpsAction.NotifyUrl
	return []tfMap{s}
}

func webhookFromAction(a *humio.Action) []tfMap {
	s := tfMap{}
	headers := make(map[string]string)
	for _, pair := range a.WebhookAction.Headers {
		headers[pair.Header] = pair.Value
	}
	s["body_template"] = a.WebhookAction.BodyTemplate
	s["headers"] = headers
	s["method"] = a.WebhookAction.Method
	s["url"] = a.WebhookAction.Url
	return []tfMap{s}
}
