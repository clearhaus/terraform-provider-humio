package humio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	humio "github.com/clearhaus/terraform-provider-humio/internal/api"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"retention": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_in_days": {
							Type:             schema.TypeFloat,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.FloatBetween(1, 365)),
						},
					},
				},
			},
		},
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	repository, err := repositoryFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain repository from resource data: %s", err)
	}

	err = client.(*humio.Client).Repositories().Create(
		repository.Name,
	)
	if err != nil {
		return diag.Errorf("could not create repository: %s", err)
	}

	err = client.(*humio.Client).Repositories().UpdateDescription(
		repository.Name,
		repository.Description,
	)
	if err != nil {
		return diag.Errorf("could not set description for repository: %s", err)
	}

	err = client.(*humio.Client).Repositories().UpdateTimeBasedRetention(
		repository.Name,
		repository.RetentionDays,
	)
	if err != nil {
		return diag.Errorf("could not set time based retention for repository: %s", err)
	}

	d.SetId(repository.Name)

	return resourceRepositoryRead(ctx, d, client)
}

func resourceRepositoryRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	repo, err := client.(*humio.Client).Repositories().Get(d.Id())
	if err != nil {
		diag.Errorf("could not get repository: %s", err)
	}
	return resourceDataFromRepository(&repo, d)
}

func resourceDataFromRepository(a *humio.Repository, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
	}
	err = d.Set("description", a.Description)
	if err != nil {
		return diag.Errorf("error setting description for resource %s: %s", d.Id(), err)
	}
	// Only set retention in state if the config has a retention block
	// This prevents drift when the API returns a default retention value
	if _, ok := d.GetOk("retention"); ok {
		if err := d.Set("retention", retentionFromRepository(a)); err != nil {
			return diag.Errorf("error setting retention settings for resource %s: %s", d.Id(), err)
		}
	}
	return nil
}

func retentionFromRepository(a *humio.Repository) []tfMap {
	s := tfMap{}
	s["time_in_days"] = a.RetentionDays
	return []tfMap{s}
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	repository, err := repositoryFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain repository from resource data: %s", err)
	}

	err = client.(*humio.Client).Repositories().UpdateDescription(
		repository.Name,
		repository.Description,
	)
	if err != nil {
		return diag.Errorf("could not update description for repository: %s", err)
	}
	err = client.(*humio.Client).Repositories().UpdateTimeBasedRetention(
		repository.Name,
		repository.RetentionDays,
	)
	if err != nil {
		return diag.Errorf("could not update time based retention for repository: %s", err)
	}

	return resourceRepositoryRead(ctx, d, client)
}

func resourceRepositoryDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	repository, err := repositoryFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain repository from resource data: %s", err)
	}

	deleteReason := "Deleted by Terraform"
	err = client.(*humio.Client).Repositories().Delete(
		repository.Name,
		deleteReason,
	)
	if err != nil {
		return diag.Errorf("could not delete repository: %s", err)
	}
	return nil
}

func repositoryFromResourceData(d *schema.ResourceData) (humio.Repository, error) {
	var retentionDays float64
	if rawRetention, ok := d.GetOk("retention"); ok {
		retentionDays = rawRetention.(*schema.Set).List()[0].(tfMap)["time_in_days"].(float64)
	}

	return humio.Repository{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		RetentionDays:          retentionDays,
	}, nil
}
