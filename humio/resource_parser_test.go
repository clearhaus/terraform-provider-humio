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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	humio "github.com/humio/terraform-provider-humio/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccParserRequiredFields(t *testing.T) {
	config := parserEmpty
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`The argument "repository" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`)},
	}, nil)
}

func TestAccParserInvalidInputs(t *testing.T) {
	config := parserInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "repository"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "tag_fields"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "test_data"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "parser_script"`)},
	}, nil)
}

func TestAccParserBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: parserBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_parser.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_parser.test", "name", "parser-test"),
				resource.TestCheckResourceAttr("humio_parser.test", "parser_script", ""),
				resource.TestCheckNoResourceAttr("humio_parser.test", "tag_fields"),
				resource.TestCheckNoResourceAttr("humio_parser.test", "test_data"),
			),
		},
	}, testAccCheckParserDestroy)
}

func TestAccParserBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: parserBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_parser.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_parser.test", "name", "parser-test"),
				resource.TestCheckResourceAttr("humio_parser.test", "parser_script", ""),
				resource.TestCheckNoResourceAttr("humio_parser.test", "tag_fields"),
				resource.TestCheckNoResourceAttr("humio_parser.test", "test_data"),
			),
		},
		{
			Config: parserFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_parser.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_parser.test", "name", "parser-test"),
				resource.TestCheckResourceAttr("humio_parser.test", "parser_script", "parser script here"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.0", "json"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.1", "test"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.0", "data1"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.1", "data2"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: parserFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_parser.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_parser.test", "name", "parser-test"),
				resource.TestCheckResourceAttr("humio_parser.test", "parser_script", "parser script here"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.0", "json"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.1", "test"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.0", "data1"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.1", "data2"),
			),
		},
	}, testAccCheckParserDestroy)
}

func TestAccParserFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: parserFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_parser.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_parser.test", "name", "parser-test"),
				resource.TestCheckResourceAttr("humio_parser.test", "parser_script", "parser script here"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.0", "json"),
				resource.TestCheckResourceAttr("humio_parser.test", "tag_fields.1", "test"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.#", "2"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.0", "data1"),
				resource.TestCheckResourceAttr("humio_parser.test", "test_data.1", "data2"),
			),
		},
	}, testAccCheckParserDestroy)
}

func testAccCheckParserDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_parser" {
			continue
		}
		parts := parseRepositoryAndID(rs.Primary.ID)
		resp, err := conn.Parsers().Get(parts[0], parts[1])
		emptyParser := humio.Parser{
			Name:        "",
			Script:      "",
			FieldsToTag: nil,
		}
		if err == nil {
			if !reflect.DeepEqual(*resp, emptyParser) {
				return fmt.Errorf("parsers still exist for id %s: %#+v", rs.Primary.ID, *resp)
			}
		}
	}
	return nil
}

const parserEmpty = `
resource "humio_parser" "test" {}
`

const parserInvalidInputs = `
resource "humio_parser" "test" {
    repository    = ["invalid"]
    name          = ["invalid"]
    tag_fields    = "invalid"
    test_data     = "invalid"
    parser_script = ["invalid"]
}
`

const parserBasic = `
resource "humio_parser" "test" {
    repository = "sandbox"
    name       = "parser-test"
}
`

const parserFull = `
resource "humio_parser" "test" {
    repository    = "sandbox"
    name          = "parser-test"
    parser_script = "parser script here"
    tag_fields    = ["json","test"]
    test_data     = ["data1","data2"]
}
`

var wantParser = humio.Parser{
	Name:        "test-parser",
	TestCases:   nil,
	Script:      "",
	FieldsToTag: nil,
}

func TestEncodeDecodeParserResource(t *testing.T) {
	res := resourceParser()
	data := res.TestResourceData()
	resourceDataFromParser(&wantParser, data)
	got, err := parserFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantParser, got) {
		t.Error(cmp.Diff(wantParser, got))
	}
}
