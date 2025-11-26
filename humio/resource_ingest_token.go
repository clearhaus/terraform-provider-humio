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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

func resourceIngestToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIngestTokenCreate,
		ReadContext:   resourceIngestTokenRead,
		UpdateContext: resourceIngestTokenUpdate,
		DeleteContext: resourceIngestTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parser": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceIngestTokenCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	ingestToken, err := ingestTokenFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	_, err = client.(*humio.Client).IngestTokens().Add(
		d.Get("repository").(string),
		ingestToken.Name,
		ingestToken.AssignedParser,
	)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to create token",
			Detail:        fmt.Sprintf("could not create ingest token: %s", err),
			AttributePath: nil,
		}}
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository"), d.Get("name")))

	return resourceIngestTokenRead(ctx, d, client)
}

func resourceIngestTokenRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		parts := parseRepositoryAndID(d.Id())
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return diag.Errorf("error importing humio_ingest_token. Please make sure the ID is in the form REPOSITORYNAME+INGESTTOKENNAMENAME (i.e. myRepoName+myIngestTokenName")
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

	ingestToken, err := client.(*humio.Client).IngestTokens().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil {
		return diag.Errorf("could not get ingest token: %s", err)
	}
	return resourceDataFromIngestToken(ingestToken, d)
}

func resourceDataFromIngestToken(a *humio.IngestToken, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
	}
	err = d.Set("token", a.Token)
	if err != nil {
		return diag.Errorf("error setting token for resource %s: %s", d.Id(), err)
	}
	err = d.Set("parser", a.AssignedParser)
	if err != nil {
		return diag.Errorf("error setting parser for resource %s: %s", d.Id(), err)
	}
	return nil
}

func resourceIngestTokenUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	ingestToken, err := ingestTokenFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	_, err = client.(*humio.Client).IngestTokens().Update(
		d.Get("repository").(string),
		ingestToken.Name,
		ingestToken.AssignedParser,
	)
	if err != nil {
		return diag.Errorf("could not update ingest token: %s", err)
	}
	return resourceIngestTokenRead(ctx, d, client)
}

func resourceIngestTokenDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	ingestToken, err := ingestTokenFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain alert from resource data: %s", err)
	}

	err = client.(*humio.Client).IngestTokens().Remove(
		d.Get("repository").(string),
		ingestToken.Name,
	)
	if err != nil {
		return diag.Errorf("could not delete ingest token: %s", err)
	}
	return nil
}

func ingestTokenFromResourceData(d *schema.ResourceData) (humio.IngestToken, error) {
	return humio.IngestToken{
		Name:           d.Get("name").(string),
		Token:          d.Get("token").(string),
		AssignedParser: d.Get("parser").(string),
	}, nil
}
