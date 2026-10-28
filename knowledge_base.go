package retailcrm

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

func (c *Client) knowledgeBaseGet(path string, parameters any, result any) (int, error) {
	if parameters != nil {
		values, err := query.Values(parameters)
		if err != nil {
			return 0, err
		}
		if encoded := values.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	data, status, err := c.GetRequest(path)
	if err != nil {
		return status, err
	}
	return status, json.Unmarshal(data, result)
}

func (c *Client) knowledgeBasePost(path, field string, input any, by string, result any) (int, error) {
	values := url.Values{}
	if field != "" {
		encoded, err := marshalToString(input)
		if err != nil {
			return 0, err
		}
		values.Set(field, encoded)
	}
	if by != "" {
		values.Set("by", by)
	}
	data, status, err := c.PostRequest(path, values)
	if err != nil {
		return status, err
	}
	return status, json.Unmarshal(data, result)
}

// KnowledgeBaseSections returns all sections.
func (c *Client) KnowledgeBaseSections() (KnowledgeBaseSectionsResponse, int, error) {
	var response KnowledgeBaseSectionsResponse
	status, err := c.knowledgeBaseGet("/knowledge-base/sections", nil, &response)
	return response, status, err
}

// KnowledgeBaseSectionCreate creates a section.
func (c *Client) KnowledgeBaseSectionCreate(section KnowledgeBaseSectionInput) (CreateResponse, int, error) {
	var response CreateResponse
	status, err := c.knowledgeBasePost("/knowledge-base/sections/create", "section", section, "", &response)
	return response, status, err
}

// KnowledgeBaseSectionEdit updates only supplied section fields.
func (c *Client) KnowledgeBaseSectionEdit(id int, section KnowledgeBaseSectionInput) (CreateResponse, int, error) {
	var response CreateResponse
	status, err := c.knowledgeBasePost(fmt.Sprintf("/knowledge-base/sections/%d/edit", id), "section", section, "", &response)
	return response, status, err
}

// KnowledgeBaseSectionDelete deletes a section, all its subsections and their documents.
func (c *Client) KnowledgeBaseSectionDelete(id int) (SuccessfulResponse, int, error) {
	var response SuccessfulResponse
	status, err := c.knowledgeBasePost(fmt.Sprintf("/knowledge-base/sections/%d/delete", id), "", nil, "", &response)
	return response, status, err
}

// KnowledgeBaseDocuments returns a filtered page of documents.
func (c *Client) KnowledgeBaseDocuments(parameters KnowledgeBaseDocumentsRequest) (KnowledgeBaseDocumentsResponse, int, error) {
	var response KnowledgeBaseDocumentsResponse
	status, err := c.knowledgeBaseGet("/knowledge-base/documents", parameters, &response)
	return response, status, err
}

// KnowledgeBaseDocument gets a document; by is "id" or "externalId" (empty uses externalId).
func (c *Client) KnowledgeBaseDocument(id, by string) (KnowledgeBaseDocumentResponse, int, error) {
	var response KnowledgeBaseDocumentResponse
	status, err := c.knowledgeBaseGet("/knowledge-base/documents/"+url.PathEscape(id), struct {
		By string `url:"by,omitempty"`
	}{by}, &response)
	return response, status, err
}

// KnowledgeBaseDocumentCreate creates a document.
func (c *Client) KnowledgeBaseDocumentCreate(document KnowledgeBaseDocumentInput) (CreateResponse, int, error) {
	var response CreateResponse
	status, err := c.knowledgeBasePost("/knowledge-base/documents/create", "document", document, "", &response)
	return response, status, err
}

// KnowledgeBaseDocumentEdit updates supplied fields; by is "id" or "externalId".
func (c *Client) KnowledgeBaseDocumentEdit(id string, document KnowledgeBaseDocumentInput, by string) (CreateResponse, int, error) {
	var response CreateResponse
	status, err := c.knowledgeBasePost("/knowledge-base/documents/"+url.PathEscape(id)+"/edit", "document", document, by, &response)
	return response, status, err
}

// KnowledgeBaseDocumentDelete deletes a document; by is "id" or "externalId".
func (c *Client) KnowledgeBaseDocumentDelete(id, by string) (SuccessfulResponse, int, error) {
	var response SuccessfulResponse
	status, err := c.knowledgeBasePost("/knowledge-base/documents/"+url.PathEscape(id)+"/delete", "", nil, by, &response)
	return response, status, err
}

// KnowledgeBaseDocumentsHistory returns document changes; sections have no history.
func (c *Client) KnowledgeBaseDocumentsHistory(parameters KnowledgeBaseDocumentsHistoryRequest) (KnowledgeBaseDocumentsHistoryResponse, int, error) {
	var response KnowledgeBaseDocumentsHistoryResponse
	status, err := c.knowledgeBaseGet("/knowledge-base/documents/history", parameters, &response)
	return response, status, err
}
