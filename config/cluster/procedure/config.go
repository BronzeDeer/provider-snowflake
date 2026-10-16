package procedure

import (
	namespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/procedure"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	namespaced.Configure(p)
	p.AddResourceConfigurator("snowflake_procedure_python", func(r *config.Resource) {
		r.Kind = "ClusterProcedurePython"
	})
}
