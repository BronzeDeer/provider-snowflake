package grantprivilegestoaccountrole

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_grant_privileges_to_account_role", func(r *config.Resource) {
		r.ShortGroup = "rbac"
		r.Kind = "AccountRolePrivilegeGrant"
		r.References["account_role_name"] = config.Reference{
			TerraformName: "snowflake_account_role",
		}
	})
}
