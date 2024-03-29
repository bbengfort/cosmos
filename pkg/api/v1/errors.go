package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	unsuccessful = Reply{Success: false}
	notFound     = Reply{Success: false, Error: "resource not found"}
	notAllowed   = Reply{Success: false, Error: "method not allowed"}
)

var (
	ErrMissingID         = errors.New("missing required id")
	ErrMissingField      = errors.New("missing required field")
	ErrInvalidField      = errors.New("invalid or unparsable field")
	ErrRestrictedField   = errors.New("field restricted for request")
	ErrConflictingFields = errors.New("only one field can be set")
	ErrModelIDMismatch   = errors.New("resource id does not match id of endpoint")
	ErrUnparsable        = errors.New("could not parse request")
	ErrUnknownUserRole   = errors.New("unknown user role")
	ErrWeakPassword      = errors.New("password must be at least 8 characters")
)

// Construct a new response for an error or simply return unsuccessful.
func ErrorResponse(err interface{}) Reply {
	if err == nil {
		return unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}

// NotFound returns a JSON 404 response for the API.
// NOTE: we know it's weird to put server-side handlers like NotFound and NotAllowed
// here in the client/api side package but it unifies where we keep our error handling
// mechanisms.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}
