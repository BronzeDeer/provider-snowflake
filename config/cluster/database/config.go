package database

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/database"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_database", func(r *config.Resource) {
		r.Kind = "ClusterDatabase"
	})
}
