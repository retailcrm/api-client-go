package retailcrm

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorTest struct {
	suite.Suite
}

func TestError(t *testing.T) {
	suite.Run(t, new(ErrorTest))
}

func (t *ErrorTest) TestFailure_ApiErrorsSlice() {
	b := []byte(`{"success": false,
				"errorMsg": "Failed to activate module",
				"errors": [
					"Your account has insufficient funds to activate integration module",
					"Test error"
				]}`)
	expected := APIErrorsList{
		"0": "Your account has insufficient funds to activate integration module",
		"1": "Test error",
	}

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().ErrorIs(e, ErrGeneric)
	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Assert().Equal(expected, apiErr.Errors())
}

func (t *ErrorTest) TestFailure_ApiErrorsMap() {
	b := []byte(`{"success": false,
				"errorMsg": "Failed to activate module",
				"errors": {"id": "ID must be an integer", "test": "Test error"}}`,
	)
	expected := APIErrorsList{
		"id":   "ID must be an integer",
		"test": "Test error",
	}

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().ErrorIs(e, ErrGeneric)
	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Assert().Equal(expected, apiErr.Errors())
}

func (t *ErrorTest) TestFailure_APIKeyMissing() {
	b := []byte(`{"success": false,
				"errorMsg": "\"apiKey\" is missing."}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrMissingCredentials)
}

func (t *ErrorTest) TestFailure_APIKeyWrong() {
	b := []byte(`{"success": false,
				"errorMsg": "Wrong \"apiKey\" value."}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrInvalidCredentials)
}

func (t *ErrorTest) TestFailure_AccessDenied() {
	b := []byte(`{"success": false,
				"errorMsg": "Access denied."}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrAccessDenied)
}

func (t *ErrorTest) TestFailure_AccountDoesNotExist() {
	b := []byte(`{"success": false,
				"errorMsg": "Account does not exist."}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrAccountDoesNotExist)
}

func (t *ErrorTest) TestFailure_Validation() {
	b := []byte(`{"success": false,
				"errorMsg": "Errors in the entity format",
				"errors": {"name": "name must be provided"}}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrValidation)
	t.Assert().Equal("name must be provided", apiErr.Errors()["name"])
}

func (t *ErrorTest) TestFailure_Validation2() {
	b := []byte(`{"success": false,
				"errorMsg": "Validation error",
				"errors": {"name": "name must be provided"}}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrValidation)
	t.Assert().Equal("name must be provided", apiErr.Errors()["name"])
	t.Assert().Equal("errorMsg: \"Validation error\", errors: [name: \"name must be provided\"]", apiErr.String())
}

func (t *ErrorTest) TestFailure_MissingParameter() {
	b := []byte(`{"success": false,
				"errorMsg": "Parameter 'item' is missing"}`,
	)

	e := CreateAPIError(b)
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Require().ErrorIs(e, ErrMissingParameter)
	t.Assert().Equal("item", apiErr.Errors()["Name"])
}

func (t *ErrorTest) Test_CreateGenericAPIError() {
	e := CreateGenericAPIError("generic error message")
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Assert().ErrorIs(apiErr, ErrGeneric)
	t.Assert().Equal("generic error message", e.Error())
}

func (t *ErrorTest) TestFailure_HTML() {
	e := CreateAPIError([]byte{'<'})
	apiErr, ok := AsAPIError(e)

	t.Require().NotNil(apiErr)
	t.Require().True(ok)
	t.Assert().ErrorIs(apiErr, ErrAccountDoesNotExist)
}
