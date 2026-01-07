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
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

// tfMap is a shorthand alias for convenience; Terraform uses this type a *lot*.
type tfMap = map[string]interface{}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: func(ctx context.Context, r *schema.ResourceData) (interface{}, diag.Diagnostics) {
			var diagnostics diag.Diagnostics
			addr := r.Get("addr").(string)
			url, err := url.Parse(addr)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			caBundlePEM, ok := r.GetOk("ca_certificate_pem")
			if ok {
				pem, _ := pem.Decode([]byte(caBundlePEM.(string)))
				if pem == nil {
					return nil, diag.FromErr(fmt.Errorf("ca_certificate_pem specified but no pem was found"))
				}
				return humio.NewClient(humio.Config{
					Address:          url,
					Token:            r.Get("api_token").(string),
					CACertificatePEM: caBundlePEM.(string),
				}), diagnostics
			}

			return humio.NewClient(humio.Config{
				Address: url,
				Token:   r.Get("api_token").(string),
			}), diagnostics
		},
		ResourcesMap: map[string]*schema.Resource{
			"humio_alert":        resourceAlert(),
			"humio_ingest_token": resourceIngestToken(),
			"humio_action":       resourceAction(),
			"humio_parser":       resourceParser(),
			"humio_repository":   resourceRepository(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"humio_user": dataSourceUser(),
		},
		Schema: map[string]*schema.Schema{
			"addr": {
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("HUMIO_ADDR", "https://cloud.humio.com/"),
				ValidateDiagFunc: validateURL,
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HUMIO_API_TOKEN", nil),
			},
			"ca_certificate_pem": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HUMIO_CA_CERTIFICATE_PEM", nil),
			},
		},
	}
}

func validateURL(val interface{}, key cty.Path) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	v := val.(string)
	u, err := url.Parse(v)
	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid URL",
			Detail:        fmt.Sprintf("%s is not a valid URL", v),
			AttributePath: key,
		})
	} else if u.Scheme == "" || u.Host == "" {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid URL",
			Detail:        fmt.Sprintf("%s must be an absolute URL", v),
			AttributePath: key,
		})
	} else if u.Scheme != "http" && u.Scheme != "https" {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid URL",
			Detail:        fmt.Sprintf("%s must begin with http or https", v),
			AttributePath: key,
		})
	}
	return diagnostics
}

func parseRepositoryAndID(fullIdentifier string) [2]string {
	var repository, id string
	parts := strings.SplitN(fullIdentifier, "+", 2)
	if len(parts) == 2 {
		repository = parts[0]
		id = parts[1]
	}
	return [2]string{repository, id}
}
