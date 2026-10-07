package tag

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"

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
	DisableNameInitializer:  true,
	SetIdentifierArgumentFn: config.NameAsIdentifier.SetIdentifierArgumentFn,
	GetExternalNameFn: func(tfstate map[string]any) (string, error) {
		v, ok := tfstate["id"]

		if !ok {
			return "", errors.New("Missing required field 'id' in terraform state")
		}

		id, ok := v.(string)

		if !ok {
			return "", errors.Errorf("provider returned non string value '%v' for property 'id'", v)
		}

		parts := strings.Split(id, ".")

		if len(parts) != 3 {
			return "", errors.Errorf("Expected id to be made of 3 parts separated by dots, but got %d parts", len(parts))
		}

		return strings.Trim(parts[2], "\""), nil
	},
	GetIDFn: func(ctx context.Context, externalName string, parameters, terraformProviderConfig map[string]any) (string, error) {
		v, ok := parameters["database"]

		if !ok {
			return "", errors.New("Missing required field 'database' in terraform state")
		}

		db, ok := v.(string)

		if !ok {
			return "", errors.Errorf("provider returned non string value '%v' for property 'database'", v)
		}

		v, ok = parameters["schema"]

		if !ok {
			return "", errors.New("Missing required field 'schema' in terraform state")
		}

		schema, ok := v.(string)

		if !ok {
			return "", errors.Errorf("provider returned non string value '%v' for property 'schema'", v)
		}

		// Need to build the id from name directly if external name is empty
		if externalName == "" {
			v, ok = parameters["name"]

			if !ok {
				return "", errors.New("Could not build id based on spec.forProvider.name")
			}

			externalName, ok = v.(string)

			if !ok {
				return "", errors.Errorf("provider returned non string value '%v' for property 'name'", v)
			}

		}
		fmt.Println("======================================= HI! ===============================")
		fmt.Printf("\"%s\".\"%s\".\"%s\"\n", db, schema, externalName)
		fmt.Println("======================================= !IH ===============================")
		return fmt.Sprintf("\"%s\".\"%s\".\"%s\"", db, schema, externalName), nil
	},
}
