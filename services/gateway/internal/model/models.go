package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string     `gorm:"default:user;index"` // admin, user, developer
	Status       string     `gorm:"index"`
	LastLoginAt  *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`

	Profile *UserProfile `gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	UserID      int64 `gorm:"primaryKey"`
	DisplayName *string
	FullName    *string
	Title       *string
	AvatarURL   *string
	Bio         *string
	Phone       *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type UserSettings struct {
	UserID    int64     `gorm:"primaryKey"`
	Theme     string    `gorm:"default:system"`
	Locale    string    `gorm:"default:zh-CN"`
	Timezone  string    `gorm:"default:Asia/Shanghai"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Project struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	UserID      int64 `gorm:"index"`
	Name        string
	Type        string `gorm:"index"`
	Description *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type ProjectDocument struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"index:idx_project_documents_user_project_updated,priority:1"`
	ProjectID int64 `gorm:"index:idx_project_documents_user_project_updated,priority:2"`
	Title     string
	Content   string         `gorm:"type:text"`
	Tags      datatypes.JSON `gorm:"type:jsonb"`
	Status    string         `gorm:"index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index:idx_project_documents_user_project_updated,priority:3"`
}

type ProjectSkill struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	UserID      int64 `gorm:"index:idx_project_skills_user_project_updated,priority:1"`
	ProjectID   int64 `gorm:"index:idx_project_skills_user_project_updated,priority:2"`
	Name        string
	Description *string
	Prompt      *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;index:idx_project_skills_user_project_updated,priority:3"`
}

type ProjectWorkflow struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	UserID      int64 `gorm:"index:idx_project_workflows_user_project_updated,priority:1"`
	ProjectID   int64 `gorm:"index:idx_project_workflows_user_project_updated,priority:2"`
	Name        string
	Description *string
	Steps       datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;index:idx_project_workflows_user_project_updated,priority:3"`
}

type Wallet struct {
	UserID        int64     `gorm:"primaryKey"`
	Balance       float64   `gorm:"type:numeric(20,6)"`
	FrozenBalance float64   `gorm:"type:numeric(20,6)"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"uniqueIndex:idx_transactions_user_ref,priority:1"`
	Type      string
	Amount    float64        `gorm:"type:numeric(20,6)"`
	RefID     string         `gorm:"uniqueIndex:idx_transactions_user_ref,priority:2"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

type UsageRecord struct {
	ID               int64  `gorm:"primaryKey;autoIncrement"`
	UserID           int64  `gorm:"index:idx_usage_records_user_created,priority:1;index:idx_usage_records_user_project,priority:1"`
	ProjectID        *int64 `gorm:"index:idx_usage_records_user_project,priority:2"`
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64   `gorm:"type:numeric(20,6)"`
	TraceID          string    `gorm:"index:idx_usage_records_trace"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index:idx_usage_records_user_created,priority:2"`
}

type Conversation struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	UserID    int64  `gorm:"index:idx_conversations_user_project_updated,priority:1"`
	ProjectID *int64 `gorm:"index:idx_conversations_user_project_updated,priority:2"`
	Title     *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;index:idx_conversations_user_project_updated,priority:3"`
}

type Message struct {
	ID             int64 `gorm:"primaryKey;autoIncrement"`
	ConversationID int64 `gorm:"index:idx_messages_conversation_created,priority:1"`
	Role           string
	Content        string
	Model          *string
	TraceID        *string
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_messages_conversation_created,priority:2"`
}

type AuditLog struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	UserID        *int64 `gorm:"index:idx_audit_logs_user_created,priority:1"`
	TraceID       string `gorm:"index:idx_audit_logs_trace"`
	Action        string
	RequestPath   *string
	RequestMethod *string
	StatusCode    *int
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_audit_logs_user_created,priority:2"`
}

type KnowledgeBase struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	UserID      int64 `gorm:"index:idx_knowledge_bases_user_scope_project_created,priority:1"`
	ProjectID   *int64
	Scope       string `gorm:"index:idx_knowledge_bases_user_scope_project_created,priority:2"`
	Name        string
	Description *string
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_knowledge_bases_user_scope_project_created,priority:4"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type KnowledgeDocument struct {
	ID              int64 `gorm:"primaryKey;autoIncrement"`
	UserID          int64 `gorm:"index:idx_knowledge_documents_user_project,priority:1"`
	ProjectID       *int64
	KnowledgeBaseID int64 `gorm:"index:idx_knowledge_documents_kb_created,priority:1"`
	FileName        string
	ContentType     *string
	SizeBytes       *int64
	StoragePath     string
	Status          string
	Metadata        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index:idx_knowledge_documents_kb_created,priority:2"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
}

type Model struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         `gorm:"uniqueIndex:idx_models_name_provider,priority:1"`
	Provider     string         `gorm:"uniqueIndex:idx_models_name_provider,priority:2;index"`
	PriceInput   float64        `gorm:"type:numeric(20,6)"`
	PriceOutput  float64        `gorm:"type:numeric(20,6)"`
	Currency     string         `gorm:"default:USD"`
	Capabilities datatypes.JSON `gorm:"type:jsonb"`
	Status       string         `gorm:"default:active;index"`
	Metadata     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
}

type Plan struct {
	ID                int64 `gorm:"primaryKey;autoIncrement"`
	Name              string
	Status            string    `gorm:"default:active;index"`
	IncludedTokens    int64     `gorm:"default:0"`
	IncludedRequests  int64     `gorm:"default:0"`
	ResetIntervalDays int       `gorm:"default:30"`
	Price             float64   `gorm:"type:numeric(20,6)"`
	Currency          string    `gorm:"default:CNY"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

type PlanSubscription struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	UserID    int64      `gorm:"index:idx_plan_subscriptions_user_status,priority:1;index:idx_plan_subscriptions_user_start_end,priority:1"`
	PlanID    int64      `gorm:"index"`
	Status    string     `gorm:"default:active;index:idx_plan_subscriptions_user_status,priority:2"`
	StartAt   time.Time  `gorm:"index:idx_plan_subscriptions_user_start_end,priority:2"`
	EndAt     *time.Time `gorm:"index:idx_plan_subscriptions_user_start_end,priority:3"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

type PlanUsage struct {
	ID             int64      `gorm:"primaryKey;autoIncrement"`
	SubscriptionID int64      `gorm:"index:idx_plan_usages_subscription_period,priority:1"`
	UserID         int64      `gorm:"index:idx_plan_usages_user_period,priority:1"`
	PeriodStart    time.Time  `gorm:"index:idx_plan_usages_subscription_period,priority:2;index:idx_plan_usages_user_period,priority:2"`
	PeriodEnd      *time.Time `gorm:"index:idx_plan_usages_subscription_period,priority:3;index:idx_plan_usages_user_period,priority:3"`
	UsedTokens     int64      `gorm:"default:0"`
	UsedRequests   int64      `gorm:"default:0"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}

type RiskPolicy struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	Name      string
	Scope     string    `gorm:"index"`
	UserID    *int64    `gorm:"index"`
	ProjectID *int64    `gorm:"index"`
	Status    string    `gorm:"default:active;index"`
	Priority  int       `gorm:"default:0;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type RateLimit struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	PolicyID      int64     `gorm:"index"`
	WindowSeconds int       `gorm:"index"`
	MaxRequests   int       `gorm:"default:0"`
	MaxTokens     int       `gorm:"default:0"`
	Status        string    `gorm:"default:active;index"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type IPRule struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	PolicyID  int64  `gorm:"index"`
	Type      string `gorm:"index"`
	IP        *string
	CIDR      *string
	Status    string    `gorm:"default:active;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type BudgetCap struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	PolicyID  int64     `gorm:"index"`
	Cycle     string    `gorm:"index"`
	MaxCost   float64   `gorm:"type:numeric(20,6)"`
	Currency  string    `gorm:"default:CNY"`
	Status    string    `gorm:"default:active;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
