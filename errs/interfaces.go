package errs

// Error implements generic error interface
type Error interface {
	error
	ApiError() string
	ApiErrors() map[string]string
}
