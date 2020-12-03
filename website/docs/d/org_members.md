---
layout: "buildkite"
page_title: "Buildkite: buildkite_org_members"
description: |-
  Get all memberships of an organization.
---

# buildkite_org_members

Use this data source to get all members of an organization.

## Example Usage

```hcl
data "buildkite_org_members" "main" {
}
```

## Attributes Reference

* `uuid` - the uuid of the organization member resource

* `member_id` - the id of the organization member resource

* `created_at` - the time at which the resource was created

* `role` - the organization role of the member

* `user_id` - the id of the user

* `user_name` - the name of the user

* `user_email` - the email of the user
