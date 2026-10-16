package grantaccountrole

import (
	"context"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_grant_account_role", func(r *config.Resource) {
		r.ShortGroup = "rbac"
		r.Kind = "AccountRoleGrant"
	})
}

var ExternalName config.ExternalName = config.ExternalName{
	DisableNameInitializer: true,
	GetExternalNameFn:      config.IDAsExternalName,
	GetIDFn: func(ctx context.Context, externalName string, parameters, terraformProviderConfig map[string]any) (string, error) {
		if externalName != "" {
			return externalName, nil
		}

		roleName := parameters["role_name"]

		granteeType := ""
		grantee := ""

		if v, ok := parameters["parent_role_name"]; ok {
			granteeType = "ROLE"
			grantee = v.(string)
		} else if v, ok = parameters["user_name"]; ok {
			granteeType = "USER"
			grantee = v.(string)
		}

		if granteeType == "" {
			return "", fmt.Errorf("must set either 'parent_role_name' or 'user_name', but got neither")
		}

		return fmt.Sprintf("%s|%s|%s", roleName, granteeType, grantee), nil
	},
	SetIdentifierArgumentFn: func(base map[string]any, externalName string) {

	},
}
