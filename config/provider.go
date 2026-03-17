package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	accountRoleCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/account_role"
	databaseCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/database"
	schemaCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/schema"
	accountRoleNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/account_role"
	databaseNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/database"
	schemaNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/schema"
)

const (
	resourcePrefix = "snowflake"
	modulePath     = "github.com/BronzeDeer/provider-snowflake"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("snowflake.bronze-deer.de"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		databaseCluster.Configure,
		schemaCluster.Configure,
		accountRoleCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("snowflake.m.bronze-deer.de"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		databaseNamespaced.Configure,
		schemaNamespaced.Configure,
		accountRoleNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
