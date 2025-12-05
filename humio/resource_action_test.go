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
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	humio "github.com/clearhaus/terraform-provider-humio/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccActionRequiredFields(t *testing.T) {
	config := actionEmpty
	accTestCase(t, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "repository" is required, but no definition was found.`),
		},
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found.`),
		},
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
		},
	}, nil)
}

func TestAccActionInvalidInputs(t *testing.T) {
	config := actionInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "repository"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "type"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "email" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "humiorepo" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "opsgenie" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "pagerduty" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "slack" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "slackpostmessage" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "victorops" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "webhook" is not expected here`)},
	}, nil)
}

func TestAccActionInvalidEmailSettings(t *testing.T) {
	config := actionInvalidEmailSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "body_template"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "recipients"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "subject_template"`)},
	}, nil)
}

func TestAccActionInvalidHumioRepoSettings(t *testing.T) {
	config := actionInvalidHumioRepoSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "ingest_token"`)},
	}, nil)
}

func TestAccActionInvalidOpsGenieSettings(t *testing.T) {
	config := actionInvalidOpsGenieSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "api_url"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "genie_key"`)},
	}, nil)
}

func TestAccActionInvalidPagerDutySettings(t *testing.T) {
	config := actionInvalidPagerDutySettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "routing_key"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "severity"`)},
	}, nil)
}

func TestAccActionInvalidSlackSettings(t *testing.T) {
	config := actionInvalidSlackSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "fields"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "url"`)},
	}, nil)
}

func TestAccActionInvalidSlackPostMessageSettings(t *testing.T) {
	config := actionInvalidSlackPostMessageSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "api_token"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "channels"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "fields"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "use_proxy"`)},
	}, nil)
}

func TestAccActionInvalidVictorOpsSettings(t *testing.T) {
	config := actionInvalidVictorOpsSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "message_type"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "notify_url"`)},
	}, nil)
}

func TestAccActionInvalidWebHookSettings(t *testing.T) {
	config := actionInvalidWebHookSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "body_template"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "headers"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "method"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "url"`)},
	}, nil)
}

func TestAccActionEmailBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionEmailBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "EmailAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-email-test"),
				resource.TestCheckResourceAttr("humio_action.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.body_template", ""),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.subject_template", ""),

				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionEmailBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionEmailBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "EmailAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-email-test"),
				resource.TestCheckResourceAttr("humio_action.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.body_template", ""),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.subject_template", ""),

				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
		{
			Config: actionEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "EmailAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-email-test"),
				resource.TestCheckResourceAttr("humio_action.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: actionEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "EmailAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-email-test"),
				resource.TestCheckResourceAttr("humio_action.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.1", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionEmailFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "EmailAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-email-test"),
				resource.TestCheckResourceAttr("humio_action.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.recipients.1", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_action.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionHumioRepoFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionHumioRepoFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "HumioRepoAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-humiorepo-test"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.0.ingest_token", "secrettoken"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionOpsGenieBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionOpsGenieBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "OpsGenieAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.api_url", "https://api.opsgenie.com"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionOpsGenieBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionOpsGenieBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "OpsGenieAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.api_url", "https://api.opsgenie.com"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
		{
			Config: actionOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "OpsGenieAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: actionOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "OpsGenieAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionOpsGenieFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "OpsGenieAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionPagerDutyFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionPagerDutyFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "PagerDutyAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-pagerduty-test"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.0.routing_key", "secretroutingkey"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.0.severity", "critical"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slack-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slack-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
		{
			Config: actionSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slack-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: actionSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slack-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slack-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackPostMessageBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackPostMessageBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackPostMessageAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.use_proxy", "true"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackPostMessageBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackPostMessageBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackPostMessageAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.use_proxy", "true"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
		{
			Config: actionSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackPostMessageAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: actionSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackPostMessageAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionSlackPostMessageFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "SlackPostMessageAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionVictorOpsFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionVictorOpsFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "VictorOpsAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-victorops-test"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.0.message_type", "important"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.0.notify_url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionWebHookBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionWebHookBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "WebhookAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-webhook-test"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.body_template", "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"actionID\": \"{alert_action_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.%", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.Content-Type", "application/json"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.method", "POST"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionWebHookBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionWebHookBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "WebhookAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-webhook-test"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.body_template", "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"actionID\": \"{alert_action_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.%", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.Content-Type", "application/json"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.method", "POST"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
			),
		},
		{
			Config: actionWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "WebhookAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-webhook-test"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: actionWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "WebhookAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-webhook-test"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func TestAccActionWebHookFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: actionWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_action.test", "action_id"),
				resource.TestCheckResourceAttr("humio_action.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_action.test", "type", "WebhookAction"),
				resource.TestCheckResourceAttr("humio_action.test", "name", "action-webhook-test"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_action.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_action.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_action.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckActionDestroy)
}

func testAccCheckActionDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_action" {
			continue
		}

		parts := parseRepositoryAndID(rs.Primary.ID)
		resp, err := conn.Actions().Get(parts[0], parts[1])
		emptyAction := humio.Action{}
		if err == nil {
			if !reflect.DeepEqual(*resp, emptyAction) {
				return fmt.Errorf("action still exist for id %s: %#+v", rs.Primary.ID, *resp)
			}
		}
		if err != nil {
			if strings.HasPrefix(err.Error(), "could not find a action") {
				return nil
			}
			return fmt.Errorf("could not validate if notifers have been cleaned up: %s", err)
		}
	}
	return nil
}

const actionEmpty = `
resource "humio_action" "test" {}
`

const actionInvalidInputs = `
resource "humio_action" "test" {
    repository       = ["invalid"]
    type             = ["invalid"]
    name             = ["invalid"]
    email            = "invalid"
    humiorepo        = "invalid"
    opsgenie         = "invalid"
    pagerduty        = "invalid"
    slack            = "invalid"
    slackpostmessage = "invalid"
    victorops        = "invalid"
    webhook          = "invalid"
}
`

const actionInvalidEmailSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "EmailAction"
    name       = "action-invalid-email"
    email {
        body_template    = ["invalid"]
        recipients       = "invalid"
        subject_template = ["invalid"]
    }
}
`

const actionInvalidHumioRepoSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "HumioRepoAction"
    name       = "action-invalid-humiorepo"
    humiorepo {
        ingest_token = ["invalid"]
    }
}
`

const actionInvalidOpsGenieSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "OpsGenieAction"
    name       = "action-invalid-opsgenie"
    opsgenie {
        api_url   = ["invalid"]
        genie_key = ["invalid"]
    }
}
`

const actionInvalidPagerDutySettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "PagerDutyAction"
    name       = "action-invalid-pagerduty"
    pagerduty {
        routing_key = ["invalid"]
        severity    = ["invalid"]
    }
}
`

const actionInvalidSlackSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackAction"
    name       = "action-invalid-slack"
    slack {
        fields = "invalid"
        url    = ["invalid"]
    }
}
`

const actionInvalidSlackPostMessageSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackPostMessageAction"
    name       = "action-invalid-slackpostmessage"
    slackpostmessage {
        api_token = ["invalid"]
        channels  = "invalid"
        fields    = "invalid"
        use_proxy = ["invalid"]
    }
}
`

const actionInvalidVictorOpsSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "VictorOpsAction"
    name       = "action-invalid-victorops"
    victorops {
        message_type = ["invalid"]
        notify_url   = ["invalid"]
    }
}
`

const actionInvalidWebHookSettings = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "WebhookAction"
    name       = "action-invalid-webhook"
    webhook {
        body_template = ["invalid"]
        headers       = "invalid"
        method        = ["invalid"]
        url           = ["invalid"]
    }
}
`

const actionEmailBasic = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "EmailAction"
    name       = "action-email-test"
    email {
        recipients = ["test@example.org"]
    }
}
`

const actionEmailFull = `
resource "humio_action" "test" {
    repository  = "sandbox"
    type        = "EmailAction"
    name        = "action-email-test"
    email {
        body_template    = "this is the body"
        recipients       = ["test@example.org", "ops@example.org"]
        subject_template = "this is the subject"
    }
}
`

const actionHumioRepoFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "HumioRepoAction"
    name       = "action-humiorepo-test"
    humiorepo {
        ingest_token = "secrettoken"
    }
}
`

const actionOpsGenieBasic = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "OpsGenieAction"
    name       = "action-opsgenie-test"
    opsgenie {
        genie_key = "secretgeniekey"
    }
}
`

const actionOpsGenieFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "OpsGenieAction"
    name       = "action-opsgenie-test"
    opsgenie {
        api_url   = "https://127.0.0.1/iasjdojaoijdioajd"
        genie_key = "secretgeniekey"
    }
}
`

const actionPagerDutyFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "PagerDutyAction"
    name       = "action-pagerduty-test"
    pagerduty {
        routing_key = "secretroutingkey"
        severity    = "critical"
    }
}
`

const actionSlackBasic = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackAction"
    name       = "action-slack-test"
    slack {
        fields = {
            "Events String" = "{events_str}"
            "Query"         = "{query_string}"
            "Time Interval" = "{query_time_interval}"
        }
        url = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    }
}
`

const actionSlackFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackAction"
    name       = "action-slack-test"
    slack {
        fields = {
			"Link" = "{url}"
			"Query" = "{query_string}"
        }
        url = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    }
}
`

const actionSlackPostMessageBasic = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackPostMessageAction"
    name       = "action-slackpostmessage-test"
    slackpostmessage {
        api_token = "secretapitoken"
        channels  = ["#alerts","#ops"]
        fields = {
            "Events String" = "{events_str}"
            "Query"         = "{query_string}"
            "Time Interval" = "{query_time_interval}"
        }
    }
}
`

const actionSlackPostMessageFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "SlackPostMessageAction"
    name       = "action-slackpostmessage-test"
    slackpostmessage {
        api_token = "secretapitoken"
        channels  = ["#alerts","#ops"]
        fields = {
			"Link" = "{url}"
			"Query" = "{query_string}"
        }
        use_proxy = false
    }
}
`

const actionVictorOpsFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "VictorOpsAction"
    name       = "action-victorops-test"
    victorops {
        message_type = "important"
        notify_url   = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

const actionWebHookBasic = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "WebhookAction"
    name       = "action-webhook-test"
    webhook {
        headers = {
            "Content-Type" = "application/json"
        }
        url = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

const actionWebHookFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type       = "WebhookAction"
    name       = "action-webhook-test"
    webhook {
        body_template = "custom body"
        headers       = {
            "custom/header1" = "this1"
            custom2          = "this2"
        }
        method = "GET"
        url    = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

var wantEmailAction = humio.Action{
	ID:     "",
	Type: "EmailAction",
	Name:   "test-action",
	EmailAction: humio.EmailAction{
		Recipients:      []string{"test@example.org", "ops@example.org"},
		BodyTemplate:    "this is the subject",
		SubjectTemplate: "this is the body",
	},
}

var wantHumioRepoAction = humio.Action{
	ID:     "",
	Type: "HumioRepoAction",
	Name:   "test-action",
	HumioRepoAction: humio.HumioRepoAction{
		IngestToken: "12345678901234567890123456789012",
	},
}

var wantOpsGenieAction = humio.Action{
	ID:     "",
	Type: "OpsGenieAction",
	Name:   "test-action",
	OpsGenieAction: humio.OpsGenieAction{
		ApiUrl:   "https://example.org",
		GenieKey: "12345678901234567890123456789012",
	},
}

var wantPagerDutyAction = humio.Action{
	ID:     "",
	Type: "PagerDutyAction",
	Name:   "test-action",
	PagerDutyAction: humio.PagerDutyAction{
		RoutingKey: "12345678901234567890123456789012",
		Severity:   "critical",
	},
}

var wantSlackAction = humio.Action{
	ID:     "",
	Type: "SlackAction",
	Name:   "test-action",
	SlackAction: humio.SlackAction{
		Url: "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ",
		Fields: []humio.SlackFieldEntryInput{
			{"Link", "{url}" },
			{"Query", "{query_string}" },
		},
	},
}

var wantSlackPostMessageAction = humio.Action{
	ID:     "",
	Type: "SlackPostMessageAction",
	Name:   "test-action",
	SlackPostMessageAction: humio.SlackPostMessageAction{
		ApiToken: "12345678901234567890123456789012",
		Channels: []string{"#alerts", "ops"},
		Fields: []humio.SlackFieldEntryInput{
			{"Link", "{url}" },
			{"Query", "{query_string}" },
		},
		UseProxy: true,
	},
}

var wantVictorOpsAction = humio.Action{
	ID:     "",
	Type: "VictorOpsAction",
	Name:   "test-action",
	VictorOpsAction: humio.VictorOpsAction{
		MessageType: "12345678901234567890123456789012",
		NotifyUrl:   "https://example.org",
	},
}

var wantWebhookAction = humio.Action{
	ID:     "",
	Type: "WebhookAction",
	Name:   "test-action",
	WebhookAction: humio.WebhookAction{
		BodyTemplate: "12345678901234567890123456789012",
		Headers: []humio.HttpHeaderEntryInput{
			{"Token", "abcdefghij123456678"},
		},
		Method: "POST",
		Url:    "https://example.org",
	},
}

func TestEncodeDecodeEmailActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantEmailAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantEmailAction, got) {
		t.Error(cmp.Diff(wantEmailAction, got))
	}
}

func TestEncodeDecodeHumioRepoActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantHumioRepoAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantHumioRepoAction, got) {
		t.Error(cmp.Diff(wantHumioRepoAction, got))
	}
}

func TestEncodeDecodeOpsGenieActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantOpsGenieAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantOpsGenieAction, got) {
		t.Error(cmp.Diff(wantOpsGenieAction, got))
	}
}

func TestEncodeDecodePagerDutyActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantPagerDutyAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantPagerDutyAction, got) {
		t.Error(cmp.Diff(wantPagerDutyAction, got))
	}
}

func TestEncodeDecodeSlackActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantSlackAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantSlackAction, got) {
		t.Error(cmp.Diff(wantSlackAction, got))
	}
}

func TestEncodeDecodeSlackPostMessageActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantSlackPostMessageAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantSlackPostMessageAction, got) {
		t.Error(cmp.Diff(wantSlackPostMessageAction, got))
	}
}

func TestEncodeDecodeVictorOpsActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantVictorOpsAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantVictorOpsAction, got) {
		t.Error(cmp.Diff(wantVictorOpsAction, got))
	}
}

func TestEncodeDecodeWebhookActionResource(t *testing.T) {
	res := resourceAction()
	data := res.TestResourceData()
	resourceDataFromAction(&wantWebhookAction, data)
	got, err := actionFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantWebhookAction, got) {
		t.Error(cmp.Diff(wantWebhookAction, got))
	}
}
