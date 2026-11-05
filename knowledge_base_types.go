package retailcrm

import "encoding/json"

// KnowledgeBaseSection is a section in the company document tree.
type KnowledgeBaseSection struct {
	ID             int    `json:"id"`
	ParentID       *int   `json:"parentId"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	DocumentsCount int    `json:"documentsCount"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// KnowledgeBaseSectionInput changes only non-nil fields. ParentID=0 moves a section to the root.
type KnowledgeBaseSectionInput struct {
	Name     *string `json:"name,omitempty"`
	ParentID *int    `json:"parentId,omitempty"`
}

// KnowledgeBaseDocumentInput changes only non-nil fields; pointers preserve false, zero and empty text.
// SectionID=0 removes the section. Content pointing to an empty string clears the text.
type KnowledgeBaseDocumentInput struct {
	Title               *string `json:"title,omitempty"`
	Content             *string `json:"content,omitempty"`
	SectionID           *int    `json:"sectionId,omitempty"`
	AvailableForCopilot *bool   `json:"availableForCopilot,omitempty"`
	ExternalID          *string `json:"externalId,omitempty"`
	// ClearExternalID sends externalId=null, taking precedence over ExternalID.
	ClearExternalID bool `json:"-"`
}

// KnowledgeBaseDocument is a company document with Markdown content and an optional original file.
type KnowledgeBaseDocument struct {
	ID                  int                        `json:"id"`
	SectionID           *int                       `json:"sectionId"`
	Title               string                     `json:"title"`
	Content             *string                    `json:"content"`
	AvailableForCopilot bool                       `json:"availableForCopilot"`
	ExternalID          *string                    `json:"externalId"`
	ContentHash         *string                    `json:"contentHash"`
	File                *KnowledgeBaseDocumentFile `json:"file"`
	CreatedAt           string                     `json:"createdAt"`
	UpdatedAt           string                     `json:"updatedAt"`
}

// KnowledgeBaseDocumentFile identifies the original file and its download URL.
type KnowledgeBaseDocumentFile struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Type *string `json:"type"`
	Size int     `json:"size"`
	URL  string  `json:"url"`
}

// KnowledgeBaseDocumentReference identifies a document in history or a file attachment.
type KnowledgeBaseDocumentReference struct {
	ID         int    `json:"id,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}

// KnowledgeBaseDocumentHistoryRecord records a change; content values are hashes, not text.
type KnowledgeBaseDocumentHistoryRecord struct {
	ID        int                             `json:"id"`
	CreatedAt string                          `json:"createdAt"`
	Created   bool                            `json:"created,omitempty"`
	Deleted   bool                            `json:"deleted,omitempty"`
	Source    string                          `json:"source"`
	Field     string                          `json:"field"`
	OldValue  any                             `json:"oldValue"`
	NewValue  any                             `json:"newValue"`
	User      *User                           `json:"user,omitempty"`
	APIKey    *APIKey                         `json:"apiKey,omitempty"`
	Document  *KnowledgeBaseDocumentReference `json:"document"`
	File      *KnowledgeBaseDocumentFile      `json:"file,omitempty"`
}

// KnowledgeBaseDocumentsFilter selects documents. DirectOnly without SectionID selects root documents.
type KnowledgeBaseDocumentsFilter struct {
	IDs                 []int    `url:"ids,omitempty,brackets"`
	ExternalIDs         []string `url:"externalIds,omitempty,brackets"`
	SectionID           int      `url:"sectionId,omitempty"`
	DirectOnly          *bool    `url:"directOnly,omitempty,int"`
	AvailableForCopilot *bool    `url:"availableForCopilot,omitempty,int"`
	HasFile             *bool    `url:"hasFile,omitempty,int"`
	Search              string   `url:"search,omitempty"`
}

// KnowledgeBaseDocumentsRequest controls filtering and pagination.
type KnowledgeBaseDocumentsRequest struct {
	Filter KnowledgeBaseDocumentsFilter `url:"filter,omitempty"`
	Limit  int                          `url:"limit,omitempty"`
	Cursor string                       `url:"cursor,omitempty"`
}

// KnowledgeBaseDocumentsHistoryFilter retrieves changes after SinceID or within a date range.
// Dates use the format "2006-01-02 15:04:05".
type KnowledgeBaseDocumentsHistoryFilter struct {
	DocumentID int    `url:"documentId,omitempty"`
	SinceID    int    `url:"sinceId,omitempty"`
	StartDate  string `url:"startDate,omitempty"`
	EndDate    string `url:"endDate,omitempty"`
}

// KnowledgeBaseDocumentsHistoryRequest controls history pagination.
type KnowledgeBaseDocumentsHistoryRequest struct {
	Filter KnowledgeBaseDocumentsHistoryFilter `url:"filter,omitempty"`
	Limit  int                                 `url:"limit,omitempty"`
	Page   int                                 `url:"page,omitempty"`
}

// KnowledgeBaseSectionsResponse contains all sections without pagination.
type KnowledgeBaseSectionsResponse struct {
	SuccessfulResponse
	Sections []KnowledgeBaseSection `json:"sections"`
}

// KnowledgeBaseDocumentsResponse contains a page of documents.
type KnowledgeBaseDocumentsResponse struct {
	SuccessfulResponse
	Pagination *CursorPagination       `json:"pagination"`
	Documents  []KnowledgeBaseDocument `json:"documents"`
}

// KnowledgeBaseDocumentResponse contains one document.
type KnowledgeBaseDocumentResponse struct {
	SuccessfulResponse
	Document KnowledgeBaseDocument `json:"document"`
}

// KnowledgeBaseDocumentsHistoryResponse contains a page of document changes.
type KnowledgeBaseDocumentsHistoryResponse struct {
	SuccessfulResponse
	GeneratedAt string                               `json:"generatedAt"`
	Pagination  *Pagination                          `json:"pagination"`
	History     []KnowledgeBaseDocumentHistoryRecord `json:"history"`
}

// MarshalJSON preserves omitted fields and allows explicitly clearing the external ID.
func (input KnowledgeBaseDocumentInput) MarshalJSON() ([]byte, error) {
	type documentInput KnowledgeBaseDocumentInput
	if input.ClearExternalID {
		return json.Marshal(struct {
			documentInput
			ExternalID *string `json:"externalId"`
		}{documentInput: documentInput(input)})
	}
	return json.Marshal(documentInput(input))
}
