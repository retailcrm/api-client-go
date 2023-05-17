package retailcrm

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// TemplateItemTypeText is a type for text chunk in template.
	TemplateItemTypeText uint8 = iota
	// TemplateItemTypeVar is a type for variable in template.
	TemplateItemTypeVar
	QuickReply  ButtonType = "QUICK_REPLY"
	PhoneNumber ButtonType = "PHONE_NUMBER"
	URL         ButtonType = "URL"
)

const (
	// TemplateVarCustom is a custom variable type.
	TemplateVarCustom = "custom"
	// TemplateVarName is a name variable type.
	TemplateVarName = "name"
	// TemplateVarFirstName is a first name variable type.
	TemplateVarFirstName = "first_name"
	// TemplateVarLastName is a last name variable type.
	TemplateVarLastName = "last_name"
)

// templateVarAssoc for checking variable validity, only for internal use.
var templateVarAssoc = map[string]interface{}{
	TemplateVarCustom:    nil,
	TemplateVarName:      nil,
	TemplateVarFirstName: nil,
	TemplateVarLastName:  nil,
}

type Text struct {
	Parts   []string `json:"parts"`
	Example []string `json:"example,omitempty"`
}

type Media struct {
	Example string `json:"example,omitempty"`
}

type Header struct {
	Text     *Text  `json:"text,omitempty"`
	Document *Media `json:"document,omitempty"`
	Image    *Media `json:"image,omitempty"`
	Video    *Media `json:"video,omitempty"`
}

type TemplateItemList []BodyTemplateItem

// BodyTemplateItem is a part of template.
type BodyTemplateItem struct {
	Text    string
	VarType string
	Type    uint8
}

// MarshalJSON controls how BodyTemplateItem will be marshaled into JSON.
func (t BodyTemplateItem) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case TemplateItemTypeText:
		return json.Marshal(t.Text)
	case TemplateItemTypeVar:
		return json.Marshal(map[string]interface{}{
			"var": t.VarType,
		})
	}

	return nil, errors.New("unknown BodyTemplateItem type")
}

// UnmarshalJSON will correctly unmarshal BodyTemplateItem.
func (t *BodyTemplateItem) UnmarshalJSON(b []byte) error {
	var obj interface{}
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}

	switch bodyPart := obj.(type) {
	case string:
		t.Type = TemplateItemTypeText
		t.Text = bodyPart
	case map[string]interface{}:
		// {} case
		if len(bodyPart) == 0 {
			t.Type = TemplateItemTypeVar
			t.VarType = TemplateVarCustom
			return nil
		}

		if varTypeCurr, ok := bodyPart["var"].(string); ok {
			if _, ok := templateVarAssoc[varTypeCurr]; !ok {
				return fmt.Errorf("invalid placeholder var '%s'", varTypeCurr)
			}

			t.Type = TemplateItemTypeVar
			t.VarType = varTypeCurr
		} else {
			return errors.New("invalid BodyTemplateItem")
		}
	default:
		return errors.New("invalid BodyTemplateItem")
	}

	return nil
}

type ButtonType string

type Button struct {
	Type        ButtonType `json:"type"`
	URL         string     `json:"url,omitempty"`
	Text        string     `json:"text,omitempty"`
	PhoneNumber string     `json:"phoneNumber,omitempty"`
	Example     []string   `json:"example,omitempty"`
}
