package serviceuser

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_service_user", func(r *config.Resource) {
		r.ShortGroup = "rbac"
		r.Kind = "ServiceUser"

		r.References["default_role"] = config.Reference{
			TerraformName: "snowflake_account_role",
		}

		r.References["default_warehouse"] = config.Reference{
			TerraformName: "snowflake_warehouse",
		}

		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				// These values use internal default values, that cannot be set directly, avoid copying them to the spec via late-init
				"mins_to_unlock",
				"disabled",
				// On initial creation this is set to the snowflake default of "ignore", but the actual valid key is "IGNORE", if late init copies the lower case to spec, we will forever reconcile loop
				"unsupported_ddl_action",
			},
		}
	})
}
