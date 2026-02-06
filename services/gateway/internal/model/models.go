package model

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Status       string    `gorm:"index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
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

type Org struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	Name        string
	OwnerUserID int64     `gorm:"index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type OrgMember struct {
	OrgID     int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"primaryKey"`
	Role      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Project struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	OrgID       int64 `gorm:"index"`
	Name        string
	Type        string `gorm:"index"`
	Description *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type ProjectDocument struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	OrgID     int64 `gorm:"index:idx_project_documents_org_project_updated,priority:1"`
	ProjectID int64 `gorm:"index:idx_project_documents_org_project_updated,priority:2"`
	Title     string
	Content   string         `gorm:"type:text"`
	Tags      datatypes.JSON `gorm:"type:jsonb"`
	Status    string         `gorm:"index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index:idx_project_documents_org_project_updated,priority:3"`
}

type ProjectSkill struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	OrgID       int64 `gorm:"index:idx_project_skills_org_project_updated,priority:1"`
	ProjectID   int64 `gorm:"index:idx_project_skills_org_project_updated,priority:2"`
	Name        string
	Description *string
	Prompt      *string       `gorm:"type:text"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime;index:idx_project_skills_org_project_updated,priority:3"`
}

type ProjectWorkflow struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	OrgID       int64 `gorm:"index:idx_project_workflows_org_project_updated,priority:1"`
	ProjectID   int64 `gorm:"index:idx_project_workflows_org_project_updated,priority:2"`
	Name        string
	Description *string
	Steps       datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;index:idx_project_workflows_org_project_updated,priority:3"`
}

type Wallet struct {
	OrgID         int64     `gorm:"primaryKey"`
	Balance       float64   `gorm:"type:numeric(20,6)"`
	FrozenBalance float64   `gorm:"type:numeric(20,6)"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	OrgID     int64 `gorm:"uniqueIndex:idx_transactions_org_ref,priority:1"`
	Type      string
	Amount    float64        `gorm:"type:numeric(20,6)"`
	RefID     string         `gorm:"uniqueIndex:idx_transactions_org_ref,priority:2"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

type UsageRecord struct {
	ID               int64  `gorm:"primaryKey;autoIncrement"`
	OrgID            int64  `gorm:"index:idx_usage_records_org_created,priority:1;index:idx_usage_records_org_project,priority:1"`
	ProjectID        *int64 `gorm:"index:idx_usage_records_org_project,priority:2"`
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64   `gorm:"type:numeric(20,6)"`
	TraceID          string    `gorm:"index:idx_usage_records_trace"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index:idx_usage_records_org_created,priority:2"`
}

type Conversation struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	OrgID     int64  `gorm:"index:idx_conversations_org_project_updated,priority:1"`
	ProjectID *int64 `gorm:"index:idx_conversations_org_project_updated,priority:2"`
	Title     *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;index:idx_conversations_org_project_updated,priority:3"`
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
	OrgID         *int64 `gorm:"index:idx_audit_logs_org_created,priority:1"`
	UserID        *int64
	TraceID       string `gorm:"index:idx_audit_logs_trace"`
	Action        string
	RequestPath   *string
	RequestMethod *string
	StatusCode    *int
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_audit_logs_org_created,priority:2"`
}

type KnowledgeBase struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	OrgID       int64 `gorm:"index:idx_knowledge_bases_org_scope_project_created,priority:1"`
	ProjectID   *int64
	Scope       string `gorm:"index:idx_knowledge_bases_org_scope_project_created,priority:2"`
	Name        string
	Description *string
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_knowledge_bases_org_scope_project_created,priority:4"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type KnowledgeDocument struct {
	ID              int64 `gorm:"primaryKey;autoIncrement"`
	OrgID           int64 `gorm:"index:idx_knowledge_documents_org_project,priority:1"`
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
