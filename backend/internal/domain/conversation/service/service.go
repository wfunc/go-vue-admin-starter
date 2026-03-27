package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	conversationentity "github.com/Wei-Shaw/sub2api/internal/domain/conversation/entity"
	conversationrepository "github.com/Wei-Shaw/sub2api/internal/domain/conversation/repository"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
)

type Service struct {
	repo *conversationrepository.Repository
}

func New(repo *conversationrepository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, params pagination.Params, status, tier, keyword string) ([]conversationentity.Item, int64, error) {
	items, total, err := s.repo.List(ctx, params, status, tier, keyword)
	if err != nil {
		return nil, 0, err
	}
	result := make([]conversationentity.Item, 0, len(items))
	for _, item := range items {
		result = append(result, conversationentity.FromEnt(item))
	}
	return result, total, nil
}

func (s *Service) Get(ctx context.Context, id int) (conversationentity.Detail, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return conversationentity.Detail{}, errorx.NotFound("conversation not found")
	}
	detail := conversationentity.DetailFromEnt(item)
	for index, history := range detail.History {
		if strings.TrimSpace(history.Title) == "" {
			detail.History[index].Title = "Event"
		}
	}
	sort.Slice(detail.History, func(i, j int) bool {
		return detail.History[i].Time.After(detail.History[j].Time)
	})
	return detail, nil
}

func (s *Service) Summary(ctx context.Context) (conversationentity.Summary, error) {
	items, err := s.repo.All(ctx)
	if err != nil {
		return conversationentity.Summary{}, err
	}
	now := time.Now()
	summary := conversationentity.Summary{
		QueueItems:      make([]conversationentity.QueueItem, 0, 3),
		AgentStatuses:   make([]conversationentity.AgentStatus, 0),
		RecentTransfers: make([]conversationentity.TransferRecord, 0, 5),
	}
	agentLoads := make(map[string]int)
	agentQueues := make(map[string]string)
	agentStatus := make(map[string]string)
	transferRecords := make([]conversationentity.TransferRecord, 0)

	for _, item := range items {
		switch item.Status {
		case "waiting":
			summary.WaitingCount++
		case "processing":
			summary.ProcessingCount++
		case "closed":
			if item.ClosedAt != nil && sameDay(*item.ClosedAt, now) {
				summary.ResolvedToday++
			}
		}
		if item.Priority == "urgent" || (item.Status == "waiting" && time.Since(item.LastActiveAt) > 15*time.Minute) {
			summary.SLARiskCount++
		}
		if item.Status != "closed" && len(summary.QueueItems) < 3 {
			summary.QueueItems = append(summary.QueueItems, conversationentity.QueueItem{
				ID:         item.ID,
				TicketNo:   item.TicketNo,
				Subject:    item.Subject,
				Customer:   item.CustomerName,
				Priority:   item.Priority,
				Channel:    item.Channel,
				Status:     item.Status,
				LastActive: item.LastActiveAt,
			})
		}
		if assignee := strings.TrimSpace(item.Assignee); assignee != "" {
			if item.Status == "processing" {
				agentLoads[assignee]++
				agentStatus[assignee] = "online"
			} else if _, exists := agentStatus[assignee]; !exists {
				agentStatus[assignee] = "offline"
			}
			agentQueues[assignee] = item.Queue
		}
		for _, message := range item.Edges.Messages {
			if message.MessageType == "event" && strings.Contains(strings.ToLower(message.Content), "transfer") {
				summary.TransferCount++
				transferRecords = append(transferRecords, conversationentity.TransferRecord{
					ID:      message.ID,
					Actor:   defaultActor(message.Actor),
					Time:    message.CreatedAt,
					Summary: message.Content,
				})
			}
		}
	}

	agentNames := make([]string, 0, len(agentStatus))
	for name := range agentStatus {
		agentNames = append(agentNames, name)
		if agentStatus[name] == "online" {
			summary.OnlineAgentCount++
		}
	}
	sort.Strings(agentNames)
	for _, name := range agentNames {
		summary.AgentStatuses = append(summary.AgentStatuses, conversationentity.AgentStatus{
			Name:       name,
			Status:     agentStatus[name],
			Queue:      agentQueues[name],
			ActiveLoad: agentLoads[name],
		})
	}
	sort.Slice(transferRecords, func(i, j int) bool {
		return transferRecords[i].Time.After(transferRecords[j].Time)
	})
	if len(transferRecords) > 5 {
		transferRecords = transferRecords[:5]
	}
	summary.RecentTransfers = transferRecords
	return summary, nil
}

func (s *Service) Reply(ctx context.Context, id int, current appauth.CurrentUser, input conversationentity.ReplyRequest) error {
	trimmedContent := strings.TrimSpace(input.Content)
	if trimmedContent == "" {
		return errorx.BadRequest("content is required")
	}
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	actorName := displayName(current)
	actorType := "agent"
	messageType := "message"
	if input.InternalNote {
		actorName = "Internal Note"
		actorType = "system"
		messageType = "note"
	}
	return s.repo.Reply(ctx, id, actorName, actorType, messageType, trimmedContent)
}

func (s *Service) Transfer(ctx context.Context, id int, current appauth.CurrentUser, input conversationentity.TransferRequest) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	targetAssignee := strings.TrimSpace(input.Assignee)
	if targetAssignee == "" {
		targetAssignee = "Tier-2 Queue"
	}
	targetQueue := strings.TrimSpace(input.Queue)
	if targetQueue == "" {
		targetQueue = "Escalation Queue"
	}
	note := strings.TrimSpace(input.Note)
	content := fmt.Sprintf("Transfer to %s via %s", targetAssignee, targetQueue)
	if note != "" {
		content = fmt.Sprintf("%s. %s", content, note)
	}
	return s.repo.Transfer(ctx, id, displayName(current), targetAssignee, targetQueue, content)
}

func (s *Service) Resolve(ctx context.Context, id int, current appauth.CurrentUser, input conversationentity.ResolveRequest) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	content := "Conversation closed by agent"
	if note := strings.TrimSpace(input.Note); note != "" {
		content = fmt.Sprintf("%s. %s", content, note)
	}
	return s.repo.Resolve(ctx, id, displayName(current), "Conversation closed and waiting for follow-up survey", content)
}

func (s *Service) Count(ctx context.Context) (int64, error) { return s.repo.Count(ctx) }

func displayName(current appauth.CurrentUser) string {
	if strings.TrimSpace(current.Nickname) != "" {
		return strings.TrimSpace(current.Nickname)
	}
	return strings.TrimSpace(current.Username)
}

func defaultActor(value string) string {
	if strings.TrimSpace(value) == "" {
		return "System"
	}
	return value
}

func sameDay(left, right time.Time) bool {
	ly, lm, ld := left.Date()
	ry, rm, rd := right.Date()
	return ly == ry && lm == rm && ld == rd
}
