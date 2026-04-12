package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	accountRoleCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/account_role"
	databaseCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/database"
	accountRoleGrantCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/grant_account_role"
	schemaCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/schema"
	userCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/user"
	warehouseCluster "github.com/BronzeDeer/provider-snowflake/config/cluster/warehouse"
	accountRoleNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/account_role"
	databaseNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/database"
	accountRoleGrantNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/grant_account_role"
	schemaNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/schema"
	userNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/user"
	warehouseNamespaced "github.com/BronzeDeer/provider-snowflake/config/namespaced/warehouse"
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
		accountRoleGrantCluster.Configure,
		userCluster.Configure,
		warehouseCluster.Configure,
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
		accountRoleGrantNamespaced.Configure,
		userNamespaced.Configure,
		warehouseNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
