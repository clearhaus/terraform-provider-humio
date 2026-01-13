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
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlertRequiredFields(t *testing.T) {
	config := alertEmpty
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`The argument "repository" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`The argument "throttle_time_millis" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`The argument "start" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`The argument "query" is required, but no definition was found.`)},
	}, nil)
}

func TestAccAlertInvalidInputs(t *testing.T) {
	config := alertInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "repository"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "throttle_time_millis"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "start"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "query"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "description"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "silenced"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "labels"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "actions"`)},
	}, nil)
}

func TestAccAlertBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: alertBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_alert.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_alert.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_alert.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_alert.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_alert.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_alert.test", "description", ""),
				resource.TestCheckResourceAttr("humio_alert.test", "silenced", "false"),
				resource.TestCheckNoResourceAttr("humio_alert.test", "actions"),
				resource.TestCheckNoResourceAttr("humio_alert.test", "labels"),
			),
		},
	}, testAccCheckAlertDestroy)
}

func TestAccAlertBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: alertBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_alert.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_alert.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_alert.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_alert.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_alert.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_alert.test", "description", ""),
				resource.TestCheckResourceAttr("humio_alert.test", "silenced", "false"),
				resource.TestCheckNoResourceAttr("humio_alert.test", "labels"),
				resource.TestCheckNoResourceAttr("humio_alert.test", "actions"),
			),
		},
		{
			Config: alertFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_alert.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_alert.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_alert.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_alert.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_alert.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_alert.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_alert.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_alert.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_alert.test", "actions.0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: alertFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_alert.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_alert.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_alert.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_alert.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_alert.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_alert.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_alert.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_alert.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_alert.test", "actions.0"),
			),
		},
	}, testAccCheckAlertDestroy)
}

func TestAccAlertFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: alertFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_alert.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_alert.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_alert.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_alert.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_alert.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_alert.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_alert.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_alert.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_alert.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_alert.test", "actions.0"),
			),
		},
	}, testAccCheckAlertDestroy)
}

func testAccCheckAlertDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_alert" {
			continue
		}
		// TODO: Use rs.Primary.ID to figure out if alert exists, and not just list all alerts.
		resp, err := conn.Alerts().List("sandbox")
		if err == nil {
			if len(resp) > 0 {
				return fmt.Errorf("alerts still exist: %#+v", resp)
			}
		}
	}
	return nil
}

const alertEmpty = `
resource "humio_alert" "test" {}
`

const alertInvalidInputs = `
resource "humio_alert" "test" {
	repository           = ["invalid"]
	name                 = ["invalid"]
	throttle_time_millis = "invalid"
	start                = ["invalid"]
	query                = ["invalid"]
	description          = ["invalid"]
	silenced             = "invalid"
	labels               = "invalid"
	actions            = "invalid"
}
`

const alertBasic = `
resource "humio_alert" "test" {
	repository           = "sandbox"
	name                 = "alert-test"
	throttle_time_millis = 3600000
	start                = "24h"
	query                = "loglevel=ERROR"
}
`

const alertFull = `
resource "humio_action" "test" {
    repository = "sandbox"
    type     = "SlackAction"
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

resource "humio_alert" "test" {
	repository           = "sandbox"
	name                 = "alert-test"
	throttle_time_millis = 3600000
	start                = "24h"
	query                = "loglevel=ERROR"
	description          = "some text"
	silenced             = true
	labels               = ["errors","important"]
	actions            = [humio_action.test.action_id]
}
`

var wantAlert = humio.Alert{
	ID:                 "",
	Name:               "over 1000 errors last 5 minutes",
	Description:        "errors occurred",
	ThrottleTimeMillis: 3600000,
	Enabled:            true,
	Actions:            []string{"action1", "action2"},
	Labels:             []string{"important", "error"},
	QueryString:        "loglevel=ERROR | count() > 1000",
	QueryStart:         "15m",
}

func TestEncodeDecodeAlertResource(t *testing.T) {
	res := resourceAlert()
	data := res.TestResourceData()
	resourceDataFromAlert(&wantAlert, data)
	got, err := alertFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantAlert, got) {
		t.Error(cmp.Diff(wantAlert, got))
	}
}
