package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"snowflake_database": config.NameAsIdentifier,
	// since we might have multiple schemas with the same name from different databases in the same namespace or cluster we need to decouple metadata.name and forProvider.name
	// name in spec determines externalname and database + "." + "externalName" is the id
	"snowflake_schema": {
		GetExternalNameFn: func(tfstate map[string]any) (string, error) {
			v, ok := tfstate["id"]

			if !ok {
				return "", errors.New("Could not set external name based on id")
			}

			vStr, ok := v.(string)

			if !ok {
				return "", errors.Errorf("provider returned non string value '%v' for property 'id'", v)
			}

			idParts := strings.Split(vStr, ".")
			schemaName := idParts[len(idParts)-1]
			schemaName = strings.Trim(schemaName, "\"")

			return schemaName, nil
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

			// Need to build the id from name directly if external name is empty
			if externalName == "" {
				v, ok = parameters["name"]

				if !ok {
					return "", errors.New("Could not set external name based on spec.forProvider.name")
				}

				externalName, ok = v.(string)

				if !ok {
					return "", errors.Errorf("provider returned non string value '%v' for property 'name'", v)
				}

			}
			return fmt.Sprintf("%s.%s", db, externalName), nil
		},
		// External Name value is the same as value of the terraform "name" field
		SetIdentifierArgumentFn: func(base map[string]any, externalName string) {

		},
		// Do not take metadata.name as name value, otherwise we will get loads of collisions unless sticking to a "1 DB per NS" rule
		DisableNameInitializer: true,
	},
}

func idWithStub() config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		en, _ := config.IDAsExternalName(tfstate)
		return en, nil
	}
	return e
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
