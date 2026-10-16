package grantprivilegestoaccountrole

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/grant_privileges_to_account_role"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_grant_privileges_to_account_role", func(r *config.Resource) {
		r.Kind = "ClusterAccountRolePrivilegeGrant"
	})
}
