package errs

// Failure struct implode runtime & api errors
type Failure struct {
	RuntimeErr error
	ApiErr     string
	ApiErrs    []string
}

// FailureResponse convert json error response into object
type FailureResponse struct {
	ErrorMsg string   `json:"errorMsg,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}
