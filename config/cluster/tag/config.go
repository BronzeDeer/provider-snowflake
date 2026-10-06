package tag

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_tag", func(r *config.Resource) {
		r.ShortGroup = "tag"
		r.References["database"] = config.Reference{
			TerraformName: "snowflake_database",
		}
		r.References["schema"] = config.Reference{
			TerraformName: "snowflake_schema",
		}
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
			},
		}
	})
}
