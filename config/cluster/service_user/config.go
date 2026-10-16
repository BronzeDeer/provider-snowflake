package serviceuser

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/service_user"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_service_user", func(r *config.Resource) {
		r.Kind = "ClusterServiceUser"
	})
}
