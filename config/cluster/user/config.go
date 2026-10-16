package user

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/user"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_user", func(r *config.Resource) {
		r.Kind = "ClusterUser"
	})
}
