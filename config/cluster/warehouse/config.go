package warehouse

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/warehouse"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_warehouse", func(r *config.Resource) {
		r.Kind = "ClusterWarehouse"
	})
}
