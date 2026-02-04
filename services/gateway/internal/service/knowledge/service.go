package knowledge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Service struct {
	repo        *repo.KnowledgeRepo
	projectRepo *repo.ProjectRepo
	storagePath string
	maxUpload   int64
	allowedMIME map[string]struct{}
}

func New(repo *repo.KnowledgeRepo, projectRepo *repo.ProjectRepo, storagePath string, maxUploadBytes int64, allowedMIME []string) *Service {
	allowed := make(map[string]struct{}, len(allowedMIME))
	for _, item := range allowedMIME {
		allowed[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}

	return &Service{
		repo:        repo,
		projectRepo: projectRepo,
		storagePath: storagePath,
		maxUpload:   maxUploadBytes,
		allowedMIME: allowed,
	}
}

var (
	ErrInvalidScope      = errors.New("invalid scope")
	ErrInvalidName       = errors.New("invalid name")
	ErrProjectRequired   = errors.New("project required")
	ErrProjectNotFound   = errors.New("project not found")
	ErrKnowledgeNotFound = errors.New("knowledge base not found")
	ErrDocumentNotFound  = errors.New("document not found")
	ErrNoUpdates         = errors.New("no updates")
	ErrFileTooLarge      = errors.New("file too large")
	ErrInvalidMimeType   = errors.New("invalid content type")
	ErrInvalidFile       = errors.New("invalid file")
)

type KnowledgeBaseItem struct {
	ID          int64   `json:"id"`
	OrgID       int64   `json:"org_id"`
	ProjectID   *int64  `json:"project_id"`
	Scope       string  `json:"scope"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	DocCount    int64   `json:"doc_count"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type KnowledgeDocumentItem struct {
	ID              int64   `json:"id"`
	KnowledgeBaseID int64   `json:"knowledge_base_id"`
	OrgID           int64   `json:"org_id"`
	ProjectID       *int64  `json:"project_id"`
	FileName        string  `json:"file_name"`
	ContentType     *string `json:"content_type"`
	SizeBytes       *int64  `json:"size_bytes"`
	StoragePath     string  `json:"storage_path"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (s *Service) ListBases(ctx context.Context, orgID int64, scope string, projectID *int64) ([]KnowledgeBaseItem, error) {
	items, err := s.repo.ListBases(ctx, orgID, scope, projectID, false)
	if err != nil {
		return nil, err
	}

	result := make([]KnowledgeBaseItem, 0, len(items))
	for _, item := range items {
		result = append(result, KnowledgeBaseItem{
			ID:          item.ID,
			OrgID:       item.OrgID,
			ProjectID:   item.ProjectID,
			Scope:       item.Scope,
			Name:        item.Name,
			Description: item.Description,
			DocCount:    item.DocCount,
			CreatedAt:   formatTime(item.CreatedAt),
			UpdatedAt:   formatTime(item.UpdatedAt),
		})
	}
	return result, nil
}

func (s *Service) GetBase(ctx context.Context, orgID, id int64) (*KnowledgeBaseItem, error) {
	base, err := s.repo.GetBase(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if base == nil {
		return nil, nil
	}
	return &KnowledgeBaseItem{
		ID:          base.ID,
		OrgID:       base.OrgID,
		ProjectID:   base.ProjectID,
		Scope:       base.Scope,
		Name:        base.Name,
		Description: base.Description,
		DocCount:    0,
		CreatedAt:   formatTime(base.CreatedAt),
		UpdatedAt:   formatTime(base.UpdatedAt),
	}, nil
}

func (s *Service) CreateBase(ctx context.Context, orgID int64, scope string, name string, description *string, projectID *int64) (*KnowledgeBaseItem, error) {
	scope = strings.ToLower(strings.TrimSpace(scope))
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}

	if scope != "org" && scope != "project" {
		return nil, ErrInvalidScope
	}

	if scope == "project" {
		if projectID == nil {
			return nil, ErrProjectRequired
		}
		if s.projectRepo != nil {
			project, err := s.projectRepo.Get(ctx, orgID, *projectID)
			if err != nil {
				return nil, err
			}
			if project == nil {
				return nil, ErrProjectNotFound
			}
		}
	}

	if description != nil {
		value := strings.TrimSpace(*description)
		description = &value
	}

	base := &model.KnowledgeBase{
		OrgID:       orgID,
		ProjectID:   projectID,
		Scope:       scope,
		Name:        name,
		Description: description,
	}

	if err := s.repo.CreateBase(ctx, base); err != nil {
		return nil, err
	}

	return &KnowledgeBaseItem{
		ID:          base.ID,
		OrgID:       base.OrgID,
		ProjectID:   base.ProjectID,
		Scope:       base.Scope,
		Name:        base.Name,
		Description: base.Description,
		DocCount:    0,
		CreatedAt:   formatTime(base.CreatedAt),
		UpdatedAt:   formatTime(base.UpdatedAt),
	}, nil
}

func (s *Service) UpdateBase(ctx context.Context, orgID, id int64, name *string, description *string) (*KnowledgeBaseItem, error) {
	updates := map[string]any{}

	if name != nil {
		value := strings.TrimSpace(*name)
		if value == "" {
			return nil, ErrInvalidName
		}
		updates["name"] = value
	}
	if description != nil {
		value := strings.TrimSpace(*description)
		updates["description"] = value
	}

	if len(updates) == 0 {
		return nil, ErrNoUpdates
	}

	updates["updated_at"] = time.Now()

	item, err := s.repo.UpdateBase(ctx, orgID, id, updates)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	return &KnowledgeBaseItem{
		ID:          item.ID,
		OrgID:       item.OrgID,
		ProjectID:   item.ProjectID,
		Scope:       item.Scope,
		Name:        item.Name,
		Description: item.Description,
		DocCount:    0,
		CreatedAt:   formatTime(item.CreatedAt),
		UpdatedAt:   formatTime(item.UpdatedAt),
	}, nil
}

func (s *Service) DeleteBase(ctx context.Context, orgID, id int64) (bool, error) {
	docs, err := s.repo.ListDocuments(ctx, orgID, id)
	if err != nil {
		return false, err
	}
	for _, doc := range docs {
		_ = os.Remove(doc.StoragePath)
	}

	return s.repo.DeleteBase(ctx, orgID, id)
}

func (s *Service) ListDocuments(ctx context.Context, orgID, kbID int64) ([]KnowledgeDocumentItem, error) {
	docs, err := s.repo.ListDocuments(ctx, orgID, kbID)
	if err != nil {
		return nil, err
	}

	result := make([]KnowledgeDocumentItem, 0, len(docs))
	for _, doc := range docs {
		result = append(result, mapDocumentItem(&doc))
	}
	return result, nil
}

func (s *Service) CreateDocument(ctx context.Context, orgID, kbID int64, fileName string, contentType *string, sizeBytes *int64, reader io.Reader) (*KnowledgeDocumentItem, error) {
	base, err := s.repo.GetBase(ctx, orgID, kbID)
	if err != nil {
		return nil, err
	}
	if base == nil {
		return nil, ErrKnowledgeNotFound
	}

	if strings.TrimSpace(fileName) == "" {
		return nil, fmt.Errorf("file name required")
	}
	if err := s.ValidateUpload(fileName, contentType, sizeBytes); err != nil {
		return nil, err
	}

	storagePath, err := s.persistFile(orgID, kbID, fileName, reader)
	if err != nil {
		return nil, err
	}

	doc := &model.KnowledgeDocument{
		OrgID:           orgID,
		ProjectID:       base.ProjectID,
		KnowledgeBaseID: kbID,
		FileName:        fileName,
		ContentType:     contentType,
		SizeBytes:       sizeBytes,
		StoragePath:     storagePath,
		Status:          "uploaded",
		Metadata:        datatypes.JSON([]byte(`{}`)),
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		_ = os.Remove(storagePath)
		return nil, err
	}

	item := mapDocumentItem(doc)
	return &item, nil
}

func (s *Service) ValidateUpload(fileName string, contentType *string, sizeBytes *int64) error {
	if strings.TrimSpace(fileName) == "" {
		return ErrInvalidFile
	}
	if sizeBytes == nil || *sizeBytes <= 0 {
		return ErrInvalidFile
	}
	if s.maxUpload > 0 && *sizeBytes > s.maxUpload {
		return ErrFileTooLarge
	}
	if len(s.allowedMIME) > 0 {
		if contentType == nil || strings.TrimSpace(*contentType) == "" {
			return ErrInvalidMimeType
		}
		value := strings.ToLower(strings.TrimSpace(*contentType))
		if _, ok := s.allowedMIME[value]; !ok {
			return ErrInvalidMimeType
		}
	}
	return nil
}

func (s *Service) MaxUploadBytes() int64 {
	return s.maxUpload
}

func (s *Service) DeleteDocument(ctx context.Context, orgID, kbID, docID int64) (bool, error) {
	doc, err := s.repo.DeleteDocument(ctx, orgID, kbID, docID)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}
	_ = os.Remove(doc.StoragePath)
	return true, nil
}

func (s *Service) GetDocument(ctx context.Context, orgID, kbID, docID int64) (*model.KnowledgeDocument, error) {
	return s.repo.GetDocument(ctx, orgID, kbID, docID)
}

func (s *Service) CountDocumentsByOrg(ctx context.Context, orgID int64) (int64, error) {
	return s.repo.CountDocumentsByOrg(ctx, orgID)
}

func (s *Service) persistFile(orgID, kbID int64, fileName string, reader io.Reader) (string, error) {
	baseDir := filepath.Join(s.storagePath, fmt.Sprintf("%d", orgID), fmt.Sprintf("%d", kbID))
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}

	safeName := strings.ReplaceAll(fileName, string(os.PathSeparator), "_")
	fileID := uuid.New().String()
	path := filepath.Join(baseDir, fmt.Sprintf("%s_%s", fileID, safeName))

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		_ = os.Remove(path)
		return "", err
	}

	return path, nil
}

func mapDocumentItem(doc *model.KnowledgeDocument) KnowledgeDocumentItem {
	return KnowledgeDocumentItem{
		ID:              doc.ID,
		KnowledgeBaseID: doc.KnowledgeBaseID,
		OrgID:           doc.OrgID,
		ProjectID:       doc.ProjectID,
		FileName:        doc.FileName,
		ContentType:     doc.ContentType,
		SizeBytes:       doc.SizeBytes,
		StoragePath:     doc.StoragePath,
		Status:          doc.Status,
		CreatedAt:       formatTime(doc.CreatedAt),
		UpdatedAt:       formatTime(doc.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
