package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tink-ab/terraform-provider-buildkite/buildkite/client"
	"github.com/tink-ab/terraform-provider-buildkite/buildkite/version"
)

const (
	userAgent = "terraform-provider-buildkite/"
)

func Provider() *schema.Provider {
	log.Printf("[DEBUG] Buildkite provider version %s", version.Version)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"buildkite_org_member":        resourceOrgMember(),
			"buildkite_pipeline":          resourcePipeline(),
			"buildkite_pipeline_schedule": resourcePipelineSchedule(),
			"buildkite_team":              resourceTeam(),
			"buildkite_team_member":       resourceTeamMember(),
			"buildkite_team_pipeline":     resourceTeamPipeline(),
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUILDKITE_ORGANIZATION", nil),
			},
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUILDKITE_API_TOKEN", nil),
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	orgName := d.Get("organization").(string)
	apiToken := d.Get("api_token").(string)

	return client.NewClient(orgName, apiToken, userAgent+version.Version), nil
}
