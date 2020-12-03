package client

import (
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"log"
)

const (
	OrganizationMemberRoleMember = "MEMBER"
	OrganizationMemberRoleAdmin  = "ADMIN"
)

type paginatedResponse struct {
	Count    int
	PageInfo struct {
		HasNextPage bool
		EndCursor   string
	}
}

type organizationMembersResponse struct {
	Organization struct {
		Members struct {
			paginatedResponse
			Edges []struct {
				Node OrganizationMember
			}
		}
	}
}

type organizationMemberResponse struct {
	OrgMember OrganizationMember `json:"organizationMember"`
}

type organizationMemberUpdateResponse struct {
	OrganizationMemberUpdate struct {
		OrganizationMember OrganizationMember
	}
}

type organizationMemberDeleteResponse struct {
	DeletedOrganizationMemberID string `json:"deletedOrganizationMemberID"`
}

type OrganizationMember struct {
	Id        string `json:"id,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	User      User   `json:"user"`
}

type User struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (c *Client) GetOrganizationMembers() (*[]OrganizationMember, error) {
	req := graphql.NewRequest(`
query OrganizationMembers($orgSlug: ID!, $first: Int!, $after: String) {
	organization(slug: $orgSlug) {
		members(first: $first, after: $after, order: NAME) {
			count
			pageInfo {
				hasNextPage
				endCursor
			}
			edges {
				node {
					id
					createdAt
					uuid
					role
					user {
						id
						name
						email
						bot
					}
				}
			}
		}
	}
}
`)
	req.Var("orgSlug", c.orgSlug)
	req.Var("first", 200)

	var response organizationMembersResponse
	var result []OrganizationMember

	for {
		response = organizationMembersResponse{}
		if err := c.graphQLRequest(req, &response); err != nil {
			return nil, errors.Wrapf(err, "failed to get organization members")
		}
		for _, e := range response.Organization.Members.Edges {
			result = append(result, e.Node)
		}
		if response.Organization.Members.PageInfo.HasNextPage {
			req.Var("after", response.Organization.Members.PageInfo.EndCursor)
		} else {
			break
		}
	}

	return &result, nil
}

func (c *Client) GetOrganizationMember(uuid string) (*OrganizationMember, error) {
	log.Printf("[TRACE] Buildkite client GetOrganizationMember %s", uuid)

	req := graphql.NewRequest(`
query GetOrganizationMember($orgMemberSlug: ID!) {
  organizationMember(slug: $orgMemberSlug) {
    id
    uuid
    role
    createdAt
    user {
      id
      name
      email
    }
  }
}`)
	req.Var("orgMemberSlug", c.createOrgSlug(uuid))

	orgMemberResponse := organizationMemberResponse{}
	if err := c.graphQLRequest(req, &orgMemberResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to get organization member %s", uuid)
	}

	return &orgMemberResponse.OrgMember, nil
}

func (c *Client) UpdateOrganizationMember(orgMember *OrganizationMember) (*OrganizationMember, error) {
	log.Printf("[TRACE] Buildkite client UpdateOrganizationMember %s", orgMember.Id)

	req := graphql.NewRequest(`
mutation OrganizationMemberUpdateMutation($organizationMemberUpdateInput: OrganizationMemberUpdateInput!) {
  organizationMemberUpdate(input: $organizationMemberUpdateInput) {
    organizationMember {
      id
      uuid
      role
      createdAt
      user {
        id
        name
        email
      }
    }
  }
}
`)

	req.Var("organizationMemberUpdateInput", map[string]interface{}{
		"id":   orgMember.Id,
		"role": orgMember.Role,
	})

	orgMemberUpdateResponse := organizationMemberUpdateResponse{}
	if err := c.graphQLRequest(req, &orgMemberUpdateResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to update organization member %s", orgMember.Id)
	}

	return &orgMemberUpdateResponse.OrganizationMemberUpdate.OrganizationMember, nil
}

func (c *Client) DeleteOrganizationMember(id string) error {
	log.Printf("[TRACE] Buildkite client DeleteOrganizationMember %s", id)
	req := graphql.NewRequest(`
mutation OrganizationMemberDeleteMutation($organizationMemberDeleteInput: OrganizationMemberDeleteInput!) {
  organizationMemberDelete(input: $organizationMemberDeleteInput) {
    deletedOrganizationMemberID
  }
}
`)

	req.Var("organizationMemberDeleteInput", map[string]interface{}{
		"id": id,
	})

	orgMemberDeleteResponse := organizationMemberDeleteResponse{}
	if err := c.graphQLRequest(req, &orgMemberDeleteResponse); err != nil {
		return errors.Wrapf(err, "failed to delete organization member %s", id)
	}

	return nil
}
