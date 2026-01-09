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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	humio "github.com/clearhaus/terraform-provider-humio/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIngestTokenRequiredFields(t *testing.T) {
	config := ingestTokenEmpty
	accTestCase(t, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "repository" is required, but no definition was found.`),
		},
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
		},
	}, nil)
}

func TestAccIngestTokenInvalidInputs(t *testing.T) {
	config := ingestTokenInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "repository"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "parser"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "token"`)},
	}, nil)
}

func TestAccIngestTokenBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: ingestTokenBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_ingest_token.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "name", "ingest-token-test"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "parser", ""),
				resource.TestCheckResourceAttrSet("humio_ingest_token.test", "token"),
			),
		},
	}, testAccCheckIngestTokenDestroy)
}

func TestAccIngestTokenBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: ingestTokenBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_ingest_token.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "name", "ingest-token-test"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "parser", ""),
				resource.TestCheckResourceAttrSet("humio_ingest_token.test", "token"),
			),
		},
		{
			Config: ingestTokenFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_ingest_token.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "name", "ingest-token-test"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "parser", "json"),
				resource.TestCheckResourceAttrSet("humio_ingest_token.test", "token"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: ingestTokenFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_ingest_token.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "name", "ingest-token-test"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "parser", "json"),
				resource.TestCheckResourceAttrSet("humio_ingest_token.test", "token"),
			),
		},
	}, testAccCheckIngestTokenDestroy)
}

func TestAccIngestTokenFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: ingestTokenFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_ingest_token.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "name", "ingest-token-test"),
				resource.TestCheckResourceAttr("humio_ingest_token.test", "parser", "json"),
				resource.TestCheckResourceAttrSet("humio_ingest_token.test", "token"),
			),
		},
	}, testAccCheckIngestTokenDestroy)
}

func testAccCheckIngestTokenDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_ingest_token" {
			continue
		}
		// TODO: Use rs.Primary.ID to figure out if ingest token exists, and not just list all ingest tokens.
		resp, err := conn.IngestTokens().List("sandbox")
		if err == nil {
			if len(resp) > 1 { // by default there is an ingest token called "default"
				return fmt.Errorf("ingest tokens still exist: %#+v", resp)
			}
		}
	}
	return nil
}

const ingestTokenEmpty = `
resource "humio_ingest_token" "test" {}
`

const ingestTokenInvalidInputs = `
resource "humio_ingest_token" "test" {
	repository = ["invalid"]
	name       = ["invalid"]
	parser     = ["invalid"]
	token      = ["invalid"]
}
`

const ingestTokenBasic = `
resource "humio_ingest_token" "test" {
	repository = "sandbox"
	name       = "ingest-token-test"
}
`

const ingestTokenFull = `
resource "humio_ingest_token" "test" {
	repository = "sandbox"
	name       = "ingest-token-test"
	parser     = "json"
}
`

var wantIngestToken = humio.IngestToken{
	Name:           "testing-shipper",
	Token:          "",
	AssignedParser: "json",
}

func TestEncodeDecodeIngestTokenResource(t *testing.T) {
	res := resourceIngestToken()
	data := res.TestResourceData()
	resourceDataFromIngestToken(&wantIngestToken, data)
	got, err := ingestTokenFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantIngestToken, got) {
		t.Error(cmp.Diff(wantIngestToken, got))
	}
}
