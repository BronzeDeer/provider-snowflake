package serviceuser

import (
	clusterresource "github.com/BronzeDeer/provider-snowflake/config/cluster/service_user"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	// Reuse cluster configure method
	clusterresource.Configure(p)
}
