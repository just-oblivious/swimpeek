package picker

import (
	"errors"
	"fmt"
	"slices"
	"swimpeek/pkg/laneclient"
	"time"

	"github.com/charmbracelet/huh"
)

// PickTenant shows a tenant picker dialog and returns the selected tenant.
func PickTenant(tenants []laneclient.Tenant) (laneclient.Tenant, error) {
	var tenant laneclient.Tenant

	opts := make([]huh.Option[laneclient.Tenant], len(tenants))
	for i, t := range tenants {
		opts[i] = huh.NewOption(t.Name, t)
	}

	// Sort the tenant by user count in descending order, this puts the most active tenant at the top.
	slices.SortFunc(opts, func(a, b huh.Option[laneclient.Tenant]) int {
		return b.Value.UserCount - a.Value.UserCount
	})

	tenantForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[laneclient.Tenant]().
				Title("Available tenants").
				Value(&tenant).
				Options(opts...).
				Description("Select a tenant to dump the configuration from."),
			huh.NewNote().DescriptionFunc(func() string {
				return fmt.Sprintf(" Tenant: %s\n     ID: %s\n  Users: %d\nCreated: %s", tenant.Name, tenant.Id, tenant.UserCount, tenant.CreatedDateTime.Format(time.DateTime))
			}, &tenant),
		),
	).WithTheme(huh.ThemeDracula()).WithLayout(huh.LayoutStack)

	if err := tenantForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return tenant, fmt.Errorf("user aborted tenant selection")
		}
		return tenant, err
	}
	return tenant, nil
}
