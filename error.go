package retailcrm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var missingParameterMatcher = regexp.MustCompile(`^Parameter \'([\w\]\[\_\-]+)\' is missing$`)
var (
	// ErrMissingCredentials will be returned if no API key was provided to the API.
	ErrMissingCredentials = NewAPIError(`apiKey is missing`)
	// ErrInvalidCredentials will be returned if provided API key is invalid.
	ErrInvalidCredentials = NewAPIError(`wrong "apiKey" value`)
	// ErrAccessDenied will be returned in case of "Access denied" error.
	ErrAccessDenied = NewAPIError("access denied")
	// ErrAccountDoesNotExist will be returned if target system does not exist.
	ErrAccountDoesNotExist = NewAPIError("account does not exist")
	// ErrValidation will be returned in case of validation errors.
	ErrValidation = NewAPIError("validation error")
	// ErrMissingParameter will be returned if parameter is missing.
	// Underlying error messages list will contain parameter name in the "Name" key.
	ErrMissingParameter = NewAPIError("missing parameter")
	// ErrGeneric will be returned if error cannot be classified as one of the errors above.
	ErrGeneric = NewAPIError("API error")
)

// APIErrorsList struct.
type APIErrorsList map[string]string

// APIError returns when an API error was occurred.
type APIError interface {
	error
	fmt.Stringer
	withWrapped(error) APIError
	withErrors(APIErrorsList) APIError
	Unwrap() error
	Errors() APIErrorsList
}

type apiError struct {
	ErrorMsg   string        `json:"errorMsg,omitempty"`
	ErrorsList APIErrorsList `json:"errors,omitempty"`
	wrapped    error
}

// CreateAPIError from the provided response data. Different error types will be returned depending on the response,
// all of them can be matched using errors.Is. APi errors will always implement APIError interface.
func CreateAPIError(dataResponse []byte) error {
	a := &apiError{}

	if len(dataResponse) > 0 && dataResponse[0] == '<' {
		return ErrAccountDoesNotExist
	}

	if err := json.Unmarshal(dataResponse, &a); err != nil {
		return err
	}

	var found APIError
	switch a.ErrorMsg {
	case `"apiKey" is missing.`:
		found = ErrMissingCredentials
	case `Wrong "apiKey" value.`:
		found = ErrInvalidCredentials
	case "Access denied.":
		found = ErrAccessDenied
	case "Account does not exist.":
		found = ErrAccountDoesNotExist
	case "Errors in the entity format":
		fallthrough
	case "Validation error":
		found = ErrValidation
	default:
		if param, ok := asMissingParameterErr(a.ErrorMsg); ok {
			return a.withWrapped(ErrMissingParameter).withErrors(APIErrorsList{"Name": param})
		}
		found = ErrGeneric
	}

	result := NewAPIError(a.ErrorMsg).withWrapped(found)
	if len(a.ErrorsList) > 0 {
		return result.withErrors(a.ErrorsList)
	}

	return result
}

// CreateGenericAPIError for the situations when API response cannot be processed, but response was actually received.
func CreateGenericAPIError(message string) APIError {
	return NewAPIError(message).withWrapped(ErrGeneric)
}

// NewAPIError returns API error with the provided message.
func NewAPIError(message string) APIError {
	return &apiError{ErrorMsg: message}
}

// AsAPIError returns APIError and true if provided error is an APIError or contains wrapped APIError.
// Returns (nil, false) otherwise.
func AsAPIError(err error) (APIError, bool) {
	apiErr := unwrapAPIError(err)
	return apiErr, apiErr != nil
}

func unwrapAPIError(err error) APIError {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(APIError); ok { // nolint:errorlint
		return apiErr
	}

	wrapper, ok := err.(interface { // nolint:errorlint
		Unwrap() error
	})
	if ok {
		return unwrapAPIError(wrapper.Unwrap())
	}

	return nil
}

// asMissingParameterErr returns true if "Parameter 'name' is missing" error message is provided.
func asMissingParameterErr(message string) (string, bool) {
	matches := missingParameterMatcher.FindAllStringSubmatch(message, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1], true
	}
	return "", false
}

// Error returns errorMsg field from the response.
func (e *apiError) Error() string {
	return e.ErrorMsg
}

// Unwrap returns wrapped error. It is usually one of the predefined types like ErrGeneric or ErrValidation.
// It can be used directly, but it's main purpose is to make errors matchable via errors.Is call.
func (e *apiError) Unwrap() error {
	return e.wrapped
}

// Errors returns errors field from the response.
func (e *apiError) Errors() APIErrorsList {
	return e.ErrorsList
}

// String returns string representation of an APIError.
func (e *apiError) String() string {
	var sb strings.Builder
	sb.Grow(256) // nolint:gomnd
	sb.WriteString(fmt.Sprintf(`errorMsg: "%s"`, e.Error()))

	if len(e.Errors()) > 0 {
		i := 0
		useIndex := true
		errorList := make([]string, len(e.Errors()))

		for index, errText := range e.Errors() {
			if i == 0 && index == "0" {
				useIndex = false
			}

			if useIndex {
				errorList[i] = fmt.Sprintf(`%s: "%s"`, index, errText)
			} else {
				errorList[i] = errText
			}

			i++
		}

		sb.WriteString(", errors: [" + strings.Join(errorList, ", ") + "]")
	}

	return sb.String()
}

// withWrapped is a wrapped setter.
func (e *apiError) withWrapped(err error) APIError {
	e.wrapped = err
	return e
}

// withErrors is an ErrorsList setter.
func (e *apiError) withErrors(m APIErrorsList) APIError {
	e.ErrorsList = m
	return e
}
