package retailcrm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnowledgeBaseWrites(t *testing.T) {
	name, text, externalID, zero, disabled := "Shipping", "", "doc-7", 0, false
	cases := []struct {
		name, path, field, payload, by string
		call                           func(*Client) (bool, int, error)
	}{
		{name: "section create", path: "sections/create", field: "section", payload: `{"name":"Shipping"}`, by: "", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseSectionCreate(KnowledgeBaseSectionInput{Name: &name})
			return resp.Success, status, err
		}},
		{name: "section edit", path: "sections/7/edit", field: "section", payload: `{"parentId":0}`, by: "", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseSectionEdit(7, KnowledgeBaseSectionInput{ParentID: &zero})
			return resp.Success, status, err
		}},
		{name: "section delete", path: "sections/7/delete", field: "", payload: ``, by: "", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseSectionDelete(7)
			return resp.Success, status, err
		}},
		{name: "document create", path: "documents/create", field: "document", payload: `{"title":"Shipping","externalId":"doc-7"}`, by: "", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseDocumentCreate(KnowledgeBaseDocumentInput{Title: &name, ExternalID: &externalID})
			return resp.Success, status, err
		}},
		{name: "document edit", path: "documents/7/edit", field: "document", payload: `{"content":"","sectionId":0,"availableForCopilot":false}`, by: "id", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseDocumentEdit("7", KnowledgeBaseDocumentInput{Content: &text, SectionID: &zero, AvailableForCopilot: &disabled}, "id")
			return resp.Success, status, err
		}},
		{name: "document delete", path: "documents/doc-7/delete", field: "", payload: ``, by: "externalId", call: func(c *Client) (bool, int, error) {
			resp, status, err := c.KnowledgeBaseDocumentDelete("doc-7", "externalId")
			return resp.Success, status, err
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v5/knowledge-base/"+testCase.path, req.URL.Path)
				assert.Equal(t, "key", req.Header.Get("X-API-KEY"))
				assert.NoError(t, req.ParseForm())
				expected := url.Values{}
				if testCase.field != "" {
					assert.JSONEq(t, testCase.payload, req.PostForm.Get(testCase.field))
					expected.Set(testCase.field, req.PostForm.Get(testCase.field))
				}
				if testCase.by != "" {
					expected.Set("by", testCase.by)
				}
				assert.Equal(t, expected, req.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, `{"success":true,"id":7}`)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			success, status, err := testCase.call(New(server.URL, "key").WithHTTPClient(server.Client()))
			require.NoError(t, err)
			assert.True(t, success)
			assert.Equal(t, http.StatusOK, status)
		})
	}
}

func TestKnowledgeBaseReads(t *testing.T) {
	disabled, enabled := false, true
	cases := []struct {
		name, path, response string
		query                url.Values
		call                 func(*Client) (any, int, error)
	}{
		{name: "sections", path: "sections", query: url.Values{}, response: `{"success":true,"sections":[{"id":1,"parentId":null,"name":"Shipping","documentsCount":2}]}`,
			call: func(c *Client) (any, int, error) { return c.KnowledgeBaseSections() }},
		{name: "document", path: "documents/7", query: url.Values{"by": {"id"}}, response: `{"success":true,"document":{"id":7,"content":null,"sectionId":null,"file":{"id":9,"name":"a.pdf","type":null,"size":10,"url":"https://example.com/a"}}}`,
			call: func(c *Client) (any, int, error) { return c.KnowledgeBaseDocument("7", "id") }},
		{name: "external ID", path: "documents/doc-7", query: url.Values{}, response: `{"success":true,"document":{"id":7,"externalId":"doc-7"}}`,
			call: func(c *Client) (any, int, error) { return c.KnowledgeBaseDocument("doc-7", "") }},
		{name: "documents", path: "documents", query: url.Values{
			"limit": {"20"}, "cursor": {"opaque+/="}, "filter[ids][]": {"7", "8"}, "filter[externalIds][]": {"doc-7"},
			"filter[sectionId]": {"1"}, "filter[directOnly]": {"1"}, "filter[availableForCopilot]": {"0"},
			"filter[hasFile]": {"0"}, "filter[search]": {"Shipping"},
		}, response: `{"success":true,"pagination":{"nextCursor":"next+/="},"documents":[{"id":7,"title":"Shipping","content":"# Shipping","availableForCopilot":false}]}`,
			call: func(c *Client) (any, int, error) {
				return c.KnowledgeBaseDocuments(KnowledgeBaseDocumentsRequest{Limit: 20, Cursor: "opaque+/=", Filter: KnowledgeBaseDocumentsFilter{
					IDs: []int{7, 8}, ExternalIDs: []string{"doc-7"}, SectionID: 1, DirectOnly: &enabled, AvailableForCopilot: &disabled, HasFile: &disabled, Search: "Shipping",
				}})
			}},
		{name: "history", path: "documents/history", query: url.Values{"limit": {"20"}, "filter[sinceId]": {"12"}, "filter[documentId]": {"7"}, "filter[startDate]": {"2026-10-26 10:00:00"}, "filter[endDate]": {"2026-10-26 11:00:00"}},
			response: `{"success":true,"generatedAt":"2026-10-26 12:00:00","history":[{"id":13,"field":"section","oldValue":null,"newValue":{"id":2},"document":{"id":7,"externalId":"doc-7"}}]}`,
			call: func(c *Client) (any, int, error) {
				return c.KnowledgeBaseDocumentsHistory(KnowledgeBaseDocumentsHistoryRequest{Limit: 20, Filter: KnowledgeBaseDocumentsHistoryFilter{SinceID: 12, DocumentID: 7, StartDate: "2026-10-26 10:00:00", EndDate: "2026-10-26 11:00:00"}})
			}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v5/knowledge-base/"+testCase.path, req.URL.Path)
				assert.Equal(t, testCase.query, req.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, testCase.response)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			response, status, err := testCase.call(New(server.URL, "key").WithHTTPClient(server.Client()))
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			encoded, err := json.Marshal(response)
			require.NoError(t, err)
			// Compare the supplied fixture fields; omitted response fields have their Go zero values.
			var expected, actual map[string]any
			require.NoError(t, json.Unmarshal([]byte(testCase.response), &expected))
			require.NoError(t, json.Unmarshal(encoded, &actual))
			assert.Equal(t, expected["success"], actual["success"])
			switch result := response.(type) {
			case KnowledgeBaseSectionsResponse:
				assert.Equal(t, 1, result.Sections[0].ID)
				assert.Nil(t, result.Sections[0].ParentID)
			case KnowledgeBaseDocumentResponse:
				assert.Equal(t, 7, result.Document.ID)
				if testCase.name == "document" {
					require.NotNil(t, result.Document.File)
					assert.Equal(t, 9, result.Document.File.ID)
					assert.Nil(t, result.Document.Content)
				}
			case KnowledgeBaseDocumentsResponse:
				assert.Equal(t, "# Shipping", *result.Documents[0].Content)
				assert.Equal(t, "next+/=", result.Pagination.NextCursor)
			case KnowledgeBaseDocumentsHistoryResponse:
				assert.Equal(t, 13, result.History[0].ID)
				assert.Equal(t, map[string]any{"id": float64(2)}, result.History[0].NewValue)
			}
		})
	}
}

func TestKnowledgeBaseFileAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/api/v5/files/9/edit", req.URL.Path)
		assert.NoError(t, req.ParseForm())
		assert.JSONEq(t, `{"filename":"a.pdf","attachment":[{"knowledgeBaseDocument":{"id":7}}]}`, req.PostForm.Get("file"))
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, `{"success":true,"file":{"id":9,"attachment":[{"knowledgeBaseDocument":{"id":7}}]}}`)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	response, _, err := New(server.URL, "key").WithHTTPClient(server.Client()).FileEdit(9, File{Filename: "a.pdf", Attachment: []Attachment{{KnowledgeBaseDocument: &KnowledgeBaseDocumentReference{ID: 7}}}})
	require.NoError(t, err)
	require.Len(t, response.File.Attachment, 1)
	assert.Equal(t, 7, response.File.Attachment[0].KnowledgeBaseDocument.ID)
}

func TestKnowledgeBaseErrors(t *testing.T) {
	for _, status := range []int{400, 403, 404, 500} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, err := io.WriteString(w, `{"success":false,"errorMsg":"Failure"}`)
			assert.NoError(t, err)
		}))
		_, actual, err := New(server.URL, "key").WithHTTPClient(server.Client()).KnowledgeBaseDocument("missing", "")
		server.Close()
		require.Error(t, err)
		assert.Equal(t, status, actual)
	}
}

func TestKnowledgeBaseClearExternalID(t *testing.T) {
	value := "old-id"
	for _, input := range []KnowledgeBaseDocumentInput{
		{ClearExternalID: true},
		{ExternalID: &value, ClearExternalID: true},
	} {
		encoded, err := json.Marshal(input)
		require.NoError(t, err)
		assert.JSONEq(t, `{"externalId":null}`, string(encoded))
	}
	encoded, err := json.Marshal(KnowledgeBaseDocumentInput{})
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(encoded))
}

func TestKnowledgeBaseDocumentsCursorPagination(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls++
		expected := url.Values{"limit": {"20"}}
		body := `{"success":true,"documents":[{"id":7}],"pagination":{"nextCursor":"next+/="}}`
		if calls == 2 {
			expected.Set("cursor", "next+/=")
			body = `{"success":true,"documents":[{"id":8}],"pagination":{"nextCursor":""}}`
		}
		assert.Equal(t, expected, req.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, body)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := New(server.URL, "key").WithHTTPClient(server.Client())
	request := KnowledgeBaseDocumentsRequest{Limit: 20}
	var ids []int
	for {
		response, _, err := client.KnowledgeBaseDocuments(request)
		require.NoError(t, err)
		require.NotNil(t, response.Pagination)
		for _, document := range response.Documents {
			ids = append(ids, document.ID)
		}
		request.Cursor = response.Pagination.NextCursor
		if request.Cursor == "" || calls > 2 {
			break
		}
	}
	assert.Equal(t, []int{7, 8}, ids)
	assert.Equal(t, 2, calls)
	assert.Empty(t, request.Cursor)
}
