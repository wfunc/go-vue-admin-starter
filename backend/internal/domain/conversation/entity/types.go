package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type Customer struct {
	Name         string   `json:"name"`
	Company      string   `json:"company"`
	Contact      string   `json:"contact"`
	Tags         []string `json:"tags"`
	LastOrder    string   `json:"last_order"`
	OpenTickets  int      `json:"open_tickets"`
	Satisfaction string   `json:"satisfaction"`
	Presence     string   `json:"presence"`
	Tier         string   `json:"tier"`
}

type Message struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Sender  string    `json:"sender"`
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
}

type History struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
}

type Item struct {
	ID         int       `json:"id"`
	TicketNo   string    `json:"ticket_no"`
	Subject    string    `json:"subject"`
	Preview    string    `json:"preview"`
	Channel    string    `json:"channel"`
	Queue      string    `json:"queue"`
	Assignee   string    `json:"assignee"`
	Status     string    `json:"status"`
	Priority   string    `json:"priority"`
	Unread     int       `json:"unread"`
	LastActive time.Time `json:"last_active"`
	SLA        string    `json:"sla"`
	Customer   Customer  `json:"customer"`
}

type Detail struct {
	Item
	Messages []Message `json:"messages"`
	History  []History `json:"history"`
}

type QueueItem struct {
	ID         int       `json:"id"`
	TicketNo   string    `json:"ticket_no"`
	Subject    string    `json:"subject"`
	Customer   string    `json:"customer"`
	Priority   string    `json:"priority"`
	Channel    string    `json:"channel"`
	Status     string    `json:"status"`
	LastActive time.Time `json:"last_active"`
}

type AgentStatus struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Queue      string `json:"queue"`
	ActiveLoad int    `json:"active_load"`
}

type TransferRecord struct {
	ID      int       `json:"id"`
	Actor   string    `json:"actor"`
	Time    time.Time `json:"time"`
	Summary string    `json:"summary"`
}

type Summary struct {
	WaitingCount     int64            `json:"waiting_count"`
	ProcessingCount  int64            `json:"processing_count"`
	SLARiskCount     int64            `json:"sla_risk_count"`
	ResolvedToday    int64            `json:"resolved_today"`
	OnlineAgentCount int64            `json:"online_agent_count"`
	TransferCount    int64            `json:"transfer_count"`
	QueueItems       []QueueItem      `json:"queue_items"`
	AgentStatuses    []AgentStatus    `json:"agent_statuses"`
	RecentTransfers  []TransferRecord `json:"recent_transfers"`
}

type ReplyRequest struct {
	Content      string `json:"content" binding:"required"`
	InternalNote bool   `json:"internal_note"`
}

type TransferRequest struct {
	Assignee string `json:"assignee"`
	Queue    string `json:"queue"`
	Note     string `json:"note"`
}

type ResolveRequest struct {
	Note string `json:"note"`
}

func FromEnt(model *ent.Conversation) Item {
	return Item{
		ID:         model.ID,
		TicketNo:   model.TicketNo,
		Subject:    model.Subject,
		Preview:    model.Preview,
		Channel:    model.Channel,
		Queue:      model.Queue,
		Assignee:   model.Assignee,
		Status:     model.Status,
		Priority:   model.Priority,
		Unread:     model.Unread,
		LastActive: model.LastActiveAt,
		SLA:        model.SLA,
		Customer: Customer{
			Name:         model.CustomerName,
			Company:      model.CustomerCompany,
			Contact:      model.CustomerContact,
			Tags:         model.CustomerTags,
			LastOrder:    model.LastOrder,
			OpenTickets:  model.OpenTickets,
			Satisfaction: model.Satisfaction,
			Presence:     model.CustomerPresence,
			Tier:         model.CustomerTier,
		},
	}
}

func DetailFromEnt(model *ent.Conversation) Detail {
	item := Detail{
		Item:     FromEnt(model),
		Messages: make([]Message, 0, len(model.Edges.Messages)),
		History:  make([]History, 0, len(model.Edges.Messages)),
	}
	for _, message := range model.Edges.Messages {
		messageItem := Message{
			ID:      message.ID,
			Author:  message.Actor,
			Sender:  message.ActorType,
			Type:    message.MessageType,
			Time:    message.CreatedAt,
			Content: message.Content,
		}
		item.Messages = append(item.Messages, messageItem)
		if message.MessageType == "event" {
			item.History = append(item.History, History{
				ID:          message.ID,
				Title:       message.Actor,
				Time:        message.CreatedAt,
				Description: message.Content,
			})
		}
	}
	return item
}
