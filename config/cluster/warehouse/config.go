package warehouse

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_warehouse", func(r *config.Resource) {
		r.ShortGroup = "warehouse"

		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				"max_concurrency_level",
				"statement_queued_timeout_in_seconds",
				"statement_timeout_in_seconds",
			},
		}
	})
}
