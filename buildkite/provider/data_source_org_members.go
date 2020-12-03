package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/tink-ab/terraform-provider-buildkite/buildkite/client"
)

func dataSourceOrgMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrgMembersRead,
		Schema: map[string]*schema.Schema{
			"members": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ValidOrganizationMemberRole, false),
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrgMembersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	buildkiteClient := m.(*client.Client)
	ms, err := buildkiteClient.GetOrganizationMembers()
	if err != nil {
		return diag.FromErr(err)
	}
	members := make([]map[string]interface{}, 0)
	for _, t := range *ms {
		members = append(members, map[string]interface{}{
			"member_id":  t.Id,
			"uuid":       t.UUID,
			"role":       t.Role,
			"created_at": t.CreatedAt,
			"user_id":    t.User.Id,
			"user_name":  t.User.Name,
			"user_email": t.User.Email,
		})
	}
	if err := d.Set("members", members); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
