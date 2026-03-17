package accountrole

import "github.com/crossplane/upjet/v2/pkg/config"

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_account_role", func(r *config.Resource) {
		r.ShortGroup = "rbac"
		r.Kind = "AccountRole"

	})
}
