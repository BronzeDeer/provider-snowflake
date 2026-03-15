package schema

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_schema", func(r *config.Resource) {
		r.ShortGroup = "schema"
		r.References["database"] = config.Reference{
			TerraformName: "snowflake_database",
		}
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				"with_managed_access",
				"is_transient",
				// Even though these are object level parameters, modifying them requires the executing user to have been given the rights to manipulate them at the account level , which most least-privilege users won't have
				// If we late initialize the fields with the default values retrieved from the state we will trigger a reconciliation attempt from the provider trying to set the default value as concrete value, which then fails with the cryptic error "Insufficient privileges to operate on Account 'ACCOUNT'."
				// By disabling late initialization for these fields we ensure that absence of fields in the CR is treated correctly as "do not try to modify externally derived/default value, even if it changes externally"
				"trace_level",
				"data_retention_time_in_days",
				"log_level",
				"max_data_extension_time_in_days",
				"storage_serialization_policy",
				"suspend_task_after_num_failures",
				"user_task_managed_initial_warehouse_size",
				"user_task_minimum_trigger_interval_in_seconds",
				"user_task_timeout_ms",
			},
		}
	})
}
