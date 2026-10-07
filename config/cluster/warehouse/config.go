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
				// When unset, the upstream provider reports sentinel values (-1 / "default") in the state, which fail its own validation once late initialized into the spec
				"auto_suspend",
				"auto_resume",
				"enable_query_acceleration",
				"query_acceleration_max_scale_factor",
			},
		}
	})
}
