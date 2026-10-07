package tag

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
			IgnoredFields: []string{},
		}
	})
}

var ExternalName config.ExternalName = config.ExternalName{
	// Needed to avoid clash on k8s resource name if tags in differetn schemas are managed in the same namespace
	DisableNameInitializer: true,
	IdentifierFields: []string{
		"name",
		"database",
		"schema",
	},
	GetExternalNameFn: config.IDAsExternalName,
	GetIDFn: func(ctx context.Context, externalName string, parameters, terraformProviderConfig map[string]any) (string, error) {
		if externalName != "" {
			return externalName, nil
		}

		v, ok := parameters["database"]
		if !ok {
			return "", errors.New("Resource missing required field 'database'")
		}

		db, ok := v.(string)
		if !ok {
			return "", errors.New("Field 'database' was not expected type string")
		}

		v, ok = parameters["schema"]
		if !ok {
			return "", errors.New("resource missing required field 'schema'")
		}

		schema, ok := v.(string)
		if !ok {
			return "", errors.New("field 'schema' was not expected type string")
		}

		v, ok = parameters["name"]
		if !ok {
			return "", errors.New("resource missing required field 'name'")
		}

		name, ok := v.(string)
		if !ok {
			return "", errors.New("field 'name' was not expected type string")
		}

		return fmt.Sprintf("\"%s\".\"%s\".\"%s\"", db, schema, name), nil
	},
	SetIdentifierArgumentFn: func(base map[string]any, externalName string) {
		// Check if externalName is set and is made from 3 dot-delmited parts
		// If not set or externalName is of an invalid format, do not mutate the resource,
		//  instead see if the already set fields in spec are sufficient to build an id which will then be set on the resource as external-name
		if externalName != "" {
			parts := strings.Split(externalName, ".")
			if len(parts) == 3 {
				base["database"] = strings.Trim(parts[0], "\"")
				base["schema"] = strings.Trim(parts[1], "\"")
				base["name"] = strings.Trim(parts[2], "\"")
			}
		}

	},
}
