package procedure

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_procedure_python", func(r *config.Resource) {
		r.ShortGroup = "procedure"
		r.Kind = "ProcedurePython"
		r.References["database"] = config.Reference{
			TerraformName: "snowflake_database",
		}
		r.References["schema"] = config.Reference{
			TerraformName: "snowflake_schema",
		}
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				"is_secure",
			},
		}
	})
}

var ExternalName config.ExternalName = config.ExternalName{
	DisableNameInitializer: true,
	GetExternalNameFn:      config.IDAsExternalName,
	GetIDFn: func(ctx context.Context, externalName string, parameters, terraformProviderConfig map[string]any) (string, error) {
		if externalName != "" {
			return externalName, nil
		}

		v, ok := parameters["database"]
		if !ok {
			return "", fmt.Errorf("must set 'database'")
		}
		db, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("'database' must be string")
		}

		v, ok = parameters["schema"]
		if !ok {
			return "", fmt.Errorf("must set 'schema'")
		}
		schema, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("'schema' must be string")
		}

		v, ok = parameters["name"]
		if !ok {
			return "", fmt.Errorf("must set 'name'")
		}
		name, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("'name' must be string")
		}

		arguments := make([](any), 0)
		v, ok = parameters["arguments"]
		if ok {
			arguments, ok = v.([](any))
			if !ok {
				return "", fmt.Errorf("'arguments' must be list")
			}
		}

		argumentTypes := make([]string, 0, len(arguments))
		for idx, val := range arguments {
			arg, ok := val.(map[string]any)
			if !ok {
				return "", fmt.Errorf("argument index %d in 'arguments' must be an object", idx)
			}

			argType, ok := arg["arg_data_type"]
			if !ok {
				return "", fmt.Errorf("argument index %d in 'arguments' is missing required field 'arg_data_type'", idx)
			}

			argTypeString, ok := argType.(string)
			if !ok {
				return "", fmt.Errorf("argument index %d in 'arguments': 'arg_data_type' is not castable to string", idx)
			}

			argumentTypes = append(argumentTypes, strings.ToLower(argTypeString))
		}

		// terraform import snowflake_procedure_python.example '"<database_name>"."<schema_name>"."<function_name>"(varchar, varchar, varchar)'
		return fmt.Sprintf("\"%s\".\"%s\".\"%s\"(%s)", db, schema, name, strings.Join(argumentTypes, ", ")), nil
	},
	SetIdentifierArgumentFn: func(base map[string]any, externalName string) {

	},
}
