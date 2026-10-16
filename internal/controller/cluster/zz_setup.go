// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	clusterdatabase "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/database/clusterdatabase"
	clusterprocedurepython "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/procedure/clusterprocedurepython"
	providerconfig "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/providerconfig"
	clusteraccountrole "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/clusteraccountrole"
	clusteraccountrolegrant "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/clusteraccountrolegrant"
	clusteraccountroleprivilegegrant "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/clusteraccountroleprivilegegrant"
	clusterserviceuser "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/clusterserviceuser"
	clusteruser "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/clusteruser"
	clusterschema "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/schema/clusterschema"
	clusterwarehouse "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/warehouse/clusterwarehouse"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		clusterdatabase.Setup,
		clusterprocedurepython.Setup,
		providerconfig.Setup,
		clusteraccountrole.Setup,
		clusteraccountrolegrant.Setup,
		clusteraccountroleprivilegegrant.Setup,
		clusterserviceuser.Setup,
		clusteruser.Setup,
		clusterschema.Setup,
		clusterwarehouse.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		clusterdatabase.SetupGated,
		clusterprocedurepython.SetupGated,
		providerconfig.SetupGated,
		clusteraccountrole.SetupGated,
		clusteraccountrolegrant.SetupGated,
		clusteraccountroleprivilegegrant.SetupGated,
		clusterserviceuser.SetupGated,
		clusteruser.SetupGated,
		clusterschema.SetupGated,
		clusterwarehouse.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
