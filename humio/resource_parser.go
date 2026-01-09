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
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

func resourceParser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParserCreate,
		ReadContext:   resourceParserRead,
		UpdateContext: resourceParserUpdate,
		DeleteContext: resourceParserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"test_data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"parser_script": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceParserCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	parser, err := parserFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain parser from resource data: %s", err)
	}

	_, err = client.(*humio.Client).Parsers().Add(
		d.Get("repository").(string),
		&parser,
		false,
	)
	if err != nil {
		return diag.Errorf("could not create parser: %s", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository"), d.Get("name")))

	return resourceParserRead(ctx, d, client)
}

func resourceParserRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		parts := parseRepositoryAndID(d.Id())
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return diag.Errorf("error importing humio_parser. Please make sure the ID is in the form REPOSITORYNAME+PARSERNAME (i.e. myRepoName+myParserName")
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

	parser, err := client.(*humio.Client).Parsers().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil || reflect.DeepEqual(*parser, humio.Parser{TestCases: []humio.ParserTestCase{}}) {
		return diag.Errorf("could not get parser: %s", err)
	}
	return resourceDataFromParser(parser, d)
}

func resourceDataFromParser(a *humio.Parser, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
	}
	err = d.Set("parser_script", a.Script)
	if err != nil {
		return diag.Errorf("error setting parser_script for resource %s: %s", d.Id(), err)
	}
	err = d.Set("tag_fields", a.FieldsToTag)
	if err != nil {
		return diag.Errorf("error setting tag_fields for resource %s: %s", d.Id(), err)
	}
	var tests []string
	for _, testCase := range a.TestCases {
		tests = append(tests, testCase.Event.RawString)
	}
	err = d.Set("test_data", tests)
	if err != nil {
		return diag.Errorf("error setting test_data for resource %s: %s", d.Id(), err)
	}
	return nil
}

func resourceParserUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	parser, err := parserFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain parser from resource data: %s", err)
	}

	// Get existing parser to obtain its ID
	existingParser, err := client.(*humio.Client).Parsers().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil {
		return diag.Errorf("could not get existing parser for update: %s", err)
	}
	parser.ID = existingParser.ID

	_, err = client.(*humio.Client).Parsers().Update(
		d.Get("repository").(string),
		&parser,
	)
	if err != nil {
		return diag.Errorf("could not update parser: %s", err)
	}
	return resourceParserRead(ctx, d, client)
}

func parserFromResourceData(d *schema.ResourceData) (humio.Parser, error) {
	testData := convertInterfaceListToStringSlice(d.Get("test_data").([]interface{}))
	var testCases []humio.ParserTestCase
	for _, testString := range testData {
		testCases = append(testCases, humio.ParserTestCase{
			Event: humio.ParserTestEvent{
				RawString: testString,
			},
		})
	}

	return humio.Parser{
		Name:        d.Get("name").(string),
		Script:      d.Get("parser_script").(string),
		FieldsToTag: convertInterfaceListToStringSlice(d.Get("tag_fields").([]interface{})),
		TestCases:   testCases,
	}, nil
}

func resourceParserDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	parser, err := parserFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain parser from resource data: %s", err)
	}

	err = client.(*humio.Client).Parsers().Delete(
		d.Get("repository").(string),
		parser.Name,
	)
	if err != nil {
		return diag.Errorf("could not delete parser: %s", err)
	}
	return nil
}
