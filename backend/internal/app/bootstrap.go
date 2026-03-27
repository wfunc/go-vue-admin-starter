package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/ent"
	entmenu "github.com/Wei-Shaw/sub2api/ent/menu"
	entrole "github.com/Wei-Shaw/sub2api/ent/role"
	systemconfig "github.com/Wei-Shaw/sub2api/ent/systemconfig"
	entuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/util/password"
)

func seedStarterData(ctx context.Context, client *ent.Client, cfg *config.Config) error {
	superAdmin, err := ensureRole(ctx, client, "Super Admin", "super_admin", []string{"*"}, true)
	if err != nil { return err }
	_, err = ensureRole(ctx, client, "Operator", "operator", []string{"dashboard:view", "user:view", "role:view", "menu:view", "system_config:view", "audit_log:view"}, true)
	if err != nil { return err }
	if err := ensureAdmin(ctx, client, cfg, superAdmin.ID); err != nil { return err }
	if err := ensureMenus(ctx, client); err != nil { return err }
	if err := ensureConfigs(ctx, client); err != nil { return err }
	return nil
}

func ensureRole(ctx context.Context, client *ent.Client, name, code string, permissions []string, isSystem bool) (*ent.Role, error) {
	item, err := client.Role.Query().Where(entrole.CodeEQ(code)).Only(ctx)
	if err == nil {
		return item, nil
	}
	return client.Role.Create().SetName(name).SetCode(code).SetDescription(name).SetPermissions(permissions).SetIsSystem(isSystem).Save(ctx)
}

func ensureAdmin(ctx context.Context, client *ent.Client, cfg *config.Config, roleID int) error {
	username := strings.TrimSpace(cfg.Seed.AdminUsername)
	if username == "" { username = "admin" }
	existing, err := client.User.Query().Where(entuser.UsernameEQ(username)).Only(ctx)
	if err == nil {
		if existing.RoleID != roleID {
			return client.User.UpdateOneID(existing.ID).SetRoleID(roleID).Exec(ctx)
		}
		return nil
	}
	hash, err := password.Hash(cfg.Seed.AdminPassword)
	if err != nil { return err }
	builder := client.User.Create().SetUsername(username).SetPasswordHash(hash).SetNickname(defaultString(cfg.Seed.AdminNickname, "Starter Admin")).SetStatus("active").SetRoleID(roleID)
	if email := strings.TrimSpace(cfg.Seed.AdminEmail); email != "" { builder.SetEmail(email) }
	_, err = builder.Save(ctx)
	return err
}

func ensureMenus(ctx context.Context, client *ent.Client) error {
	dashboard, err := ensureMenu(ctx, client, menuSeed{Title: "Dashboard", Name: "dashboard", Path: "/dashboard", Component: "dashboard", Icon: "dashboard", Permission: "dashboard:view", Sort: 1})
	if err != nil { return err }
	systemRoot, err := ensureMenu(ctx, client, menuSeed{Title: "System", Name: "system", Path: "/system", Component: "system", Icon: "system", Sort: 2})
	if err != nil { return err }
	_, _ = dashboard, systemRoot
	children := []menuSeed{
		{Title: "Users", Name: "users", Path: "/users", Component: "user/users", Icon: "users", Permission: "user:view", Sort: 10, ParentID: &systemRoot.ID},
		{Title: "Roles", Name: "roles", Path: "/roles", Component: "system/roles", Icon: "roles", Permission: "role:view", Sort: 20, ParentID: &systemRoot.ID},
		{Title: "Menus", Name: "menus", Path: "/menus", Component: "system/menus", Icon: "menus", Permission: "menu:view", Sort: 30, ParentID: &systemRoot.ID},
		{Title: "System Config", Name: "system-configs", Path: "/system-configs", Component: "system/configs", Icon: "settings", Permission: "system_config:view", Sort: 40, ParentID: &systemRoot.ID},
		{Title: "Audit Logs", Name: "audit-logs", Path: "/audit-logs", Component: "audit/logs", Icon: "audit", Permission: "audit_log:view", Sort: 50, ParentID: &systemRoot.ID},
	}
	for _, item := range children {
		if _, err := ensureMenu(ctx, client, item); err != nil { return err }
	}
	return nil
}

func ensureConfigs(ctx context.Context, client *ent.Client) error {
	seeds := []struct { Key, Value, Category, Description string; IsPublic bool }{
		{"site_name", "Go Vue Admin Starter", "branding", "Display name for the starter admin panel", true},
		{"site_description", "Reusable admin skeleton for future business modules", "branding", "Default marketing description", true},
		{"allow_registration", "false", "security", "Whether public self registration is enabled", false},
	}
	for _, item := range seeds {
		existing, err := client.SystemConfig.Query().Where(systemconfig.KeyEQ(item.Key)).Only(ctx)
		if err == nil { _ = existing; continue }
		_, err = client.SystemConfig.Create().SetKey(item.Key).SetValue(item.Value).SetCategory(item.Category).SetDescription(item.Description).SetIsPublic(item.IsPublic).Save(ctx)
		if err != nil { return err }
	}
	return nil
}

type menuSeed struct {
	Title, Name, Path, Component, Icon, Permission string
	Sort int
	ParentID *int
}

func ensureMenu(ctx context.Context, client *ent.Client, seed menuSeed) (*ent.Menu, error) {
	item, err := client.Menu.Query().Where(entmenu.NameEQ(seed.Name)).Only(ctx)
	if err == nil { return item, nil }
	builder := client.Menu.Create().SetTitle(seed.Title).SetName(seed.Name).SetPath(seed.Path).SetComponent(seed.Component).SetIcon(seed.Icon).SetMenuType("menu").SetPermission(seed.Permission).SetSort(seed.Sort)
	if seed.ParentID != nil { builder.SetParentID(*seed.ParentID) }
	created, err := builder.Save(ctx)
	if err != nil { return nil, fmt.Errorf("create menu %s: %w", seed.Name, err) }
	return created, nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" { return fallback }
	return strings.TrimSpace(value)
}
