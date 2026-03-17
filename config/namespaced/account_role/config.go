package accountrole

import (
	accountrole "github.com/BronzeDeer/provider-snowflake/config/cluster/account_role"
	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	// Reuse cluster configure method
	accountrole.Configure(p)
}
