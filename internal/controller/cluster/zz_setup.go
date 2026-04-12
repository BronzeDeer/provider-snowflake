// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	database "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/database/database"
	providerconfig "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/providerconfig"
	accountrole "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/accountrole"
	accountrolegrant "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/accountrolegrant"
	user "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/rbac/user"
	schema "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/schema/schema"
	warehouse "github.com/BronzeDeer/provider-snowflake/internal/controller/cluster/warehouse/warehouse"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		database.Setup,
		providerconfig.Setup,
		accountrole.Setup,
		accountrolegrant.Setup,
		user.Setup,
		schema.Setup,
		warehouse.Setup,
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
		database.SetupGated,
		providerconfig.SetupGated,
		accountrole.SetupGated,
		accountrolegrant.SetupGated,
		user.SetupGated,
		schema.SetupGated,
		warehouse.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
