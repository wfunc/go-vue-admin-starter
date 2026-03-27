package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	entconversation "github.com/Wei-Shaw/sub2api/ent/conversation"
	entconversationmessage "github.com/Wei-Shaw/sub2api/ent/conversationmessage"
	entmenu "github.com/Wei-Shaw/sub2api/ent/menu"
	entrole "github.com/Wei-Shaw/sub2api/ent/role"
	systemconfig "github.com/Wei-Shaw/sub2api/ent/systemconfig"
	entuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/util/password"
)

func seedStarterData(ctx context.Context, client *ent.Client, cfg *config.Config) error {
	superAdmin, err := ensureRole(ctx, client, "Super Admin", "super_admin", []string{"*"}, true)
	if err != nil {
		return err
	}
	_, err = ensureRole(ctx, client, "Operator", "operator", []string{
		"dashboard:view",
		"user:view",
		"role:view",
		"menu:view",
		"system_config:view",
		"audit_log:view",
		"cs_dashboard:view",
		"conversation:view",
		"conversation:reply",
		"conversation:transfer",
		"conversation:resolve",
		"customer:view",
		"ticket:view",
		"ticket_category:view",
		"quick_reply:view",
		"knowledge:view",
	}, true)
	if err != nil {
		return err
	}
	if err := ensureAdmin(ctx, client, cfg, superAdmin.ID); err != nil {
		return err
	}
	if err := ensureMenus(ctx, client); err != nil {
		return err
	}
	if err := ensureConversations(ctx, client); err != nil {
		return err
	}
	if err := ensureConfigs(ctx, client); err != nil {
		return err
	}
	return nil
}

func ensureRole(ctx context.Context, client *ent.Client, name, code string, permissions []string, isSystem bool) (*ent.Role, error) {
	item, err := client.Role.Query().Where(entrole.CodeEQ(code)).Only(ctx)
	if err == nil {
		if item.Name != name || item.Description != name || item.IsSystem != isSystem || !stringSliceEqual(item.Permissions, permissions) {
			return client.Role.UpdateOneID(item.ID).
				SetName(name).
				SetDescription(name).
				SetPermissions(permissions).
				SetIsSystem(isSystem).
				Save(ctx)
		}
		return item, nil
	}
	return client.Role.Create().SetName(name).SetCode(code).SetDescription(name).SetPermissions(permissions).SetIsSystem(isSystem).Save(ctx)
}

func ensureAdmin(ctx context.Context, client *ent.Client, cfg *config.Config, roleID int) error {
	username := strings.TrimSpace(cfg.Seed.AdminUsername)
	if username == "" {
		username = "admin"
	}
	existing, err := client.User.Query().Where(entuser.UsernameEQ(username)).Only(ctx)
	if err == nil {
		if existing.RoleID != roleID {
			return client.User.UpdateOneID(existing.ID).SetRoleID(roleID).Exec(ctx)
		}
		return nil
	}
	hash, err := password.Hash(cfg.Seed.AdminPassword)
	if err != nil {
		return err
	}
	builder := client.User.Create().SetUsername(username).SetPasswordHash(hash).SetNickname(defaultString(cfg.Seed.AdminNickname, "Starter Admin")).SetStatus("active").SetRoleID(roleID)
	if email := strings.TrimSpace(cfg.Seed.AdminEmail); email != "" {
		builder.SetEmail(email)
	}
	_, err = builder.Save(ctx)
	return err
}

func ensureMenus(ctx context.Context, client *ent.Client) error {
	dashboard, err := ensureMenu(ctx, client, menuSeed{Title: "Dashboard", Name: "dashboard", Path: "/dashboard", Component: "dashboard", Icon: "dashboard", Permission: "dashboard:view", Sort: 1})
	if err != nil {
		return err
	}
	customerServiceRoot, err := ensureMenu(ctx, client, menuSeed{Title: "Customer Service", Name: "customer-service", Path: "/customer-service", Component: "customer-service", Icon: "service", Sort: 2})
	if err != nil {
		return err
	}
	systemRoot, err := ensureMenu(ctx, client, menuSeed{Title: "System", Name: "system", Path: "/system", Component: "system", Icon: "system", Sort: 3})
	if err != nil {
		return err
	}
	_, _, _ = dashboard, customerServiceRoot, systemRoot
	customerServiceChildren := []menuSeed{
		{Title: "Workbench", Name: "cs-dashboard", Path: "/customer-service/dashboard", Component: "customer-service/dashboard", Icon: "dashboard", Permission: "cs_dashboard:view", Sort: 10, ParentID: &customerServiceRoot.ID},
		{Title: "Conversations", Name: "conversations", Path: "/customer-service/conversations", Component: "customer-service/conversations", Icon: "chat", Permission: "conversation:view", Sort: 20, ParentID: &customerServiceRoot.ID},
		{Title: "Customers", Name: "customers", Path: "/customer-service/customers", Component: "customer-service/customers", Icon: "users", Permission: "customer:view", Sort: 30, ParentID: &customerServiceRoot.ID},
		{Title: "Tickets", Name: "tickets", Path: "/customer-service/tickets", Component: "customer-service/tickets", Icon: "ticket", Permission: "ticket:view", Sort: 40, ParentID: &customerServiceRoot.ID},
		{Title: "Categories", Name: "ticket-categories", Path: "/customer-service/categories", Component: "customer-service/categories", Icon: "category", Permission: "ticket_category:view", Sort: 50, ParentID: &customerServiceRoot.ID},
		{Title: "Quick Replies", Name: "quick-replies", Path: "/customer-service/quick-replies", Component: "customer-service/quick-replies", Icon: "reply", Permission: "quick_reply:view", Sort: 60, ParentID: &customerServiceRoot.ID},
		{Title: "Knowledge Base", Name: "knowledge-articles", Path: "/customer-service/knowledge", Component: "customer-service/knowledge", Icon: "knowledge", Permission: "knowledge:view", Sort: 70, ParentID: &customerServiceRoot.ID},
	}
	for _, item := range customerServiceChildren {
		if _, err := ensureMenu(ctx, client, item); err != nil {
			return err
		}
	}
	children := []menuSeed{
		{Title: "Users", Name: "users", Path: "/users", Component: "user/users", Icon: "users", Permission: "user:view", Sort: 10, ParentID: &systemRoot.ID},
		{Title: "Roles", Name: "roles", Path: "/roles", Component: "system/roles", Icon: "roles", Permission: "role:view", Sort: 20, ParentID: &systemRoot.ID},
		{Title: "Menus", Name: "menus", Path: "/menus", Component: "system/menus", Icon: "menus", Permission: "menu:view", Sort: 30, ParentID: &systemRoot.ID},
		{Title: "System Config", Name: "system-configs", Path: "/system-configs", Component: "system/configs", Icon: "settings", Permission: "system_config:view", Sort: 40, ParentID: &systemRoot.ID},
		{Title: "Audit Logs", Name: "audit-logs", Path: "/audit-logs", Component: "audit/logs", Icon: "audit", Permission: "audit_log:view", Sort: 50, ParentID: &systemRoot.ID},
	}
	for _, item := range children {
		if _, err := ensureMenu(ctx, client, item); err != nil {
			return err
		}
	}
	return nil
}

func ensureConfigs(ctx context.Context, client *ent.Client) error {
	seeds := []struct {
		Key, Value, Category, Description string
		IsPublic                          bool
	}{
		{"site_name", "Go Vue Admin Starter", "branding", "Display name for the starter admin panel", true},
		{"site_description", "Reusable admin skeleton for future business modules", "branding", "Default marketing description", true},
		{"allow_registration", "false", "security", "Whether public self registration is enabled", false},
	}
	for _, item := range seeds {
		existing, err := client.SystemConfig.Query().Where(systemconfig.KeyEQ(item.Key)).Only(ctx)
		if err == nil {
			_ = existing
			continue
		}
		_, err = client.SystemConfig.Create().SetKey(item.Key).SetValue(item.Value).SetCategory(item.Category).SetDescription(item.Description).SetIsPublic(item.IsPublic).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

type conversationSeed struct {
	TicketNo         string
	Subject          string
	Preview          string
	Channel          string
	Queue            string
	Assignee         string
	Status           string
	Priority         string
	Unread           int
	SLA              string
	CustomerName     string
	CustomerCompany  string
	CustomerContact  string
	CustomerTags     []string
	CustomerPresence string
	CustomerTier     string
	LastOrder        string
	OpenTickets      int
	Satisfaction     string
	LastActiveAt     time.Time
	ClosedAt         *time.Time
	Messages         []conversationMessageSeed
}

type conversationMessageSeed struct {
	Actor       string
	ActorType   string
	MessageType string
	Content     string
	CreatedAt   time.Time
}

func ensureConversations(ctx context.Context, client *ent.Client) error {
	now := time.Now()
	closedAt := now.Add(-2 * time.Hour)
	seeds := []conversationSeed{
		{
			TicketNo:         "CS-2026-042",
			Subject:          "续费后 webhook 没有重新生效",
			Preview:          "客户反馈续费成功，但回调地址没有再次触发。",
			Channel:          "Web Chat",
			Queue:            "Billing Queue",
			Assignee:         "Lena",
			Status:           "processing",
			Priority:         "urgent",
			Unread:           2,
			SLA:              "14 分钟内首响",
			CustomerName:     "Aster Health",
			CustomerCompany:  "Aster Health",
			CustomerContact:  "ops@asterhealth.io · +86 138 0000 1122",
			CustomerTags:     []string{"VIP", "续费客户", "回调异常"},
			CustomerPresence: "online",
			CustomerTier:     "vip",
			LastOrder:        "2026-03-27 14:02",
			OpenTickets:      2,
			Satisfaction:     "94%",
			LastActiveAt:     now.Add(-2 * time.Minute),
			Messages: []conversationMessageSeed{
				{Actor: "客户", ActorType: "customer", MessageType: "message", Content: "我们刚完成续费，但 webhook 还是没有回调。", CreatedAt: now.Add(-9 * time.Minute)},
				{Actor: "Lena", ActorType: "agent", MessageType: "message", Content: "收到，我先帮您核查续费同步和回调日志。", CreatedAt: now.Add(-7 * time.Minute)},
				{Actor: "系统事件", ActorType: "system", MessageType: "event", Content: "Transfer to Billing Queue via Billing Queue. Renewal callback issue linked to payment retry log.", CreatedAt: now.Add(-5 * time.Minute)},
			},
		},
		{
			TicketNo:         "CS-2026-040",
			Subject:          "客户要求导出最近三个月会话记录",
			Preview:          "法务审核需要导出会话归档，等待客服确认导出范围。",
			Channel:          "Email",
			Queue:            "Compliance",
			Assignee:         "Mason",
			Status:           "waiting",
			Priority:         "high",
			Unread:           1,
			SLA:              "32 分钟内回复",
			CustomerName:     "BlueRiver Retail",
			CustomerCompany:  "BlueRiver Retail",
			CustomerContact:  "legal@blueriver.ai · +1 415 010 0222",
			CustomerTags:     []string{"导出请求", "法务审核"},
			CustomerPresence: "offline",
			CustomerTier:     "standard",
			LastOrder:        "2026-03-11 09:40",
			OpenTickets:      1,
			Satisfaction:     "91%",
			LastActiveAt:     now.Add(-8 * time.Minute),
			Messages: []conversationMessageSeed{
				{Actor: "客户", ActorType: "customer", MessageType: "message", Content: "请帮忙导出最近三个月所有客服会话记录。", CreatedAt: now.Add(-42 * time.Minute)},
				{Actor: "系统事件", ActorType: "system", MessageType: "event", Content: "Transfer to Mason via Compliance. Data export request queued for compliance review.", CreatedAt: now.Add(-39 * time.Minute)},
			},
		},
		{
			TicketNo:         "CS-2026-038",
			Subject:          "登录回调地址偶发跳转失败",
			Preview:          "客户提供了两次失败录屏，怀疑和地区网络节点有关。",
			Channel:          "WeCom",
			Queue:            "Technical Support",
			Assignee:         "Mia",
			Status:           "processing",
			Priority:         "medium",
			Unread:           0,
			SLA:              "已在 SLA 内",
			CustomerName:     "Lattice Care",
			CustomerCompany:  "Lattice Care",
			CustomerContact:  "it@latticecare.cn · +86 21 6600 9981",
			CustomerTags:     []string{"登录异常", "企业版"},
			CustomerPresence: "online",
			CustomerTier:     "standard",
			LastOrder:        "2026-02-28 11:20",
			OpenTickets:      3,
			Satisfaction:     "89%",
			LastActiveAt:     now.Add(-18 * time.Minute),
			Messages: []conversationMessageSeed{
				{Actor: "客户", ActorType: "customer", MessageType: "message", Content: "我们这边偶发会跳到失效回调页面。", CreatedAt: now.Add(-80 * time.Minute)},
				{Actor: "Mia", ActorType: "agent", MessageType: "message", Content: "已收到录屏，我这边先核查回调参数与边缘节点日志。", CreatedAt: now.Add(-77 * time.Minute)},
			},
		},
		{
			TicketNo:         "CS-2026-034",
			Subject:          "夜间班次会话已完成回访",
			Preview:          "会话已结束，等待回访问卷",
			Channel:          "Web Chat",
			Queue:            "Night Shift",
			Assignee:         "Iris",
			Status:           "closed",
			Priority:         "low",
			Unread:           0,
			SLA:              "已完成",
			CustomerName:     "Northwind Labs",
			CustomerCompany:  "Northwind Labs",
			CustomerContact:  "support@northwind.dev · +1 212 010 9981",
			CustomerTags:     []string{"回访问卷", "夜间班次"},
			CustomerPresence: "offline",
			CustomerTier:     "standard",
			LastOrder:        "2026-03-01 20:21",
			OpenTickets:      0,
			Satisfaction:     "97%",
			LastActiveAt:     now.Add(-3 * time.Hour),
			ClosedAt:         &closedAt,
			Messages: []conversationMessageSeed{
				{Actor: "客户", ActorType: "customer", MessageType: "message", Content: "问题已经解决，多谢。", CreatedAt: now.Add(-190 * time.Minute)},
				{Actor: "Iris", ActorType: "agent", MessageType: "event", Content: "Transfer to Iris via Night Shift. Night shift agent accepted the conversation.", CreatedAt: now.Add(-188 * time.Minute)},
				{Actor: "Iris", ActorType: "agent", MessageType: "event", Content: "Conversation closed by agent. Follow-up survey queued.", CreatedAt: now.Add(-2 * time.Hour)},
			},
		},
	}
	for _, seed := range seeds {
		if err := ensureConversation(ctx, client, seed); err != nil {
			return err
		}
	}
	return nil
}

func ensureConversation(ctx context.Context, client *ent.Client, seed conversationSeed) error {
	item, err := client.Conversation.Query().Where(entconversation.TicketNoEQ(seed.TicketNo)).Only(ctx)
	if err != nil {
		builder := client.Conversation.Create().
			SetTicketNo(seed.TicketNo).
			SetSubject(seed.Subject).
			SetPreview(seed.Preview).
			SetChannel(seed.Channel).
			SetQueue(seed.Queue).
			SetAssignee(seed.Assignee).
			SetStatus(seed.Status).
			SetPriority(seed.Priority).
			SetUnread(seed.Unread).
			SetSLA(seed.SLA).
			SetCustomerName(seed.CustomerName).
			SetCustomerCompany(seed.CustomerCompany).
			SetCustomerContact(seed.CustomerContact).
			SetCustomerTags(seed.CustomerTags).
			SetCustomerPresence(seed.CustomerPresence).
			SetCustomerTier(seed.CustomerTier).
			SetLastOrder(seed.LastOrder).
			SetOpenTickets(seed.OpenTickets).
			SetSatisfaction(seed.Satisfaction).
			SetLastActiveAt(seed.LastActiveAt)
		if seed.ClosedAt != nil {
			builder.SetClosedAt(*seed.ClosedAt)
		}
		item, err = builder.Save(ctx)
		if err != nil {
			return err
		}
	} else {
		builder := client.Conversation.UpdateOneID(item.ID).
			SetSubject(seed.Subject).
			SetPreview(seed.Preview).
			SetChannel(seed.Channel).
			SetQueue(seed.Queue).
			SetAssignee(seed.Assignee).
			SetStatus(seed.Status).
			SetPriority(seed.Priority).
			SetUnread(seed.Unread).
			SetSLA(seed.SLA).
			SetCustomerName(seed.CustomerName).
			SetCustomerCompany(seed.CustomerCompany).
			SetCustomerContact(seed.CustomerContact).
			SetCustomerTags(seed.CustomerTags).
			SetCustomerPresence(seed.CustomerPresence).
			SetCustomerTier(seed.CustomerTier).
			SetLastOrder(seed.LastOrder).
			SetOpenTickets(seed.OpenTickets).
			SetSatisfaction(seed.Satisfaction).
			SetLastActiveAt(seed.LastActiveAt)
		if seed.ClosedAt != nil {
			builder.SetClosedAt(*seed.ClosedAt)
		} else {
			builder.ClearClosedAt()
		}
		item, err = builder.Save(ctx)
		if err != nil {
			return err
		}
	}
	messageCount, err := client.ConversationMessage.Query().Where(entconversationmessage.ConversationIDEQ(item.ID)).Count(ctx)
	if err != nil {
		return err
	}
	if messageCount > 0 {
		return nil
	}
	for _, message := range seed.Messages {
		if _, err := client.ConversationMessage.Create().
			SetConversationID(item.ID).
			SetActor(message.Actor).
			SetActorType(message.ActorType).
			SetMessageType(message.MessageType).
			SetContent(message.Content).
			SetCreatedAt(message.CreatedAt).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

type menuSeed struct {
	Title, Name, Path, Component, Icon, Permission string
	Sort                                           int
	ParentID                                       *int
}

func ensureMenu(ctx context.Context, client *ent.Client, seed menuSeed) (*ent.Menu, error) {
	item, err := client.Menu.Query().Where(entmenu.NameEQ(seed.Name)).Only(ctx)
	if err == nil {
		updater := client.Menu.UpdateOneID(item.ID).
			SetTitle(seed.Title).
			SetPath(seed.Path).
			SetComponent(seed.Component).
			SetIcon(seed.Icon).
			SetMenuType("menu").
			SetPermission(seed.Permission).
			SetSort(seed.Sort).
			SetHidden(false)
		if seed.ParentID != nil {
			updater.SetParentID(*seed.ParentID)
		} else {
			updater.ClearParentID()
		}
		return updater.Save(ctx)
	}
	builder := client.Menu.Create().SetTitle(seed.Title).SetName(seed.Name).SetPath(seed.Path).SetComponent(seed.Component).SetIcon(seed.Icon).SetMenuType("menu").SetPermission(seed.Permission).SetSort(seed.Sort)
	if seed.ParentID != nil {
		builder.SetParentID(*seed.ParentID)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create menu %s: %w", seed.Name, err)
	}
	return created, nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}
