package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type RequestPayload interface {
	Require([]string) *payloadHandler
	Permit([]string) *payloadHandler
	SetData(map[string]interface{}) *payloadHandler
}

type payloadHandler struct {
	Data   map[string]interface{}
	Errors []rest_errors.RestErr
}

func NewParamsHandler() RequestPayload {
	return &payloadHandler{}
}

func WhitelistQueryParams(c *gin.Context, allowdKeys []string) url.Values {
	allowdKeys = append(allowdKeys, []string{"page", "sort"}...)
	params := c.Request.URL.Query()
	for key, _ := range params {
		attr := strings.Split(key, "__")[0]
		if !slices.Contains(allowdKeys, attr) {
			delete(params, key)
		}
	}

	return params
}

func getIdFromUrl(c *gin.Context, fromQuery bool) (int64, rest_errors.RestErr) {
	paramId := c.Param("id")
	if fromQuery {
		paramId = c.Query("id")
	}
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return 0, rest_errors.NewBadRequestError("id should be a number")
	}
	return id, nil
}

func (p *payloadHandler) Require(attrs []string) *payloadHandler {
	p.Errors = []rest_errors.RestErr{}
	for _, attr := range attrs {
		if p.Data[attr] == nil {
			p.Errors = append(p.Errors, rest_errors.InvalidError(fmt.Sprintf("Required Field not found: %s ", attr)))
		}
	}

	return p
}

func (p *payloadHandler) Permit(attrs []string) *payloadHandler {
	for key, _ := range p.Data {
		if !slices.Contains(attrs, key) {
			delete(p.Data, key)
		}
	}
	return p
}

func (p *payloadHandler) ClearEmpty() {
	for key, value := range p.Data {
		if value == nil || value == "" {
			delete(p.Data, key)
		}
	}
}

func (p *payloadHandler) SetData(payload map[string]interface{}) *payloadHandler {
	p.Errors = []rest_errors.RestErr{}
	data := payload["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	if attributes == nil {
		p.Errors = append(p.Errors, rest_errors.InvalidError("Attributes not found"))
	}

	p.Data = attributes
	return p
}

var (
	Validate = validator.New(validator.WithRequiredStructEnabled())
)

/* func buildUpdatePayload(isPartialUpdat bool, payload interface{}) interface{} {

	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Serialize(authType)
	}
	return result
} */
