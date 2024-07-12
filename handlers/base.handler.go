package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
	"resturants-hub.com/m/v2/dto"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type BaseHandler interface {
	Require([]string) *baseHandler
	Permit([]string) *baseHandler
	SetData(map[string]interface{}) *baseHandler
	CurrentUser(*gin.Context) *dto.BaseUser
}

type baseHandler struct {
	Data   map[string]interface{}
	Errors []rest_errors.RestErr
}

func NewBaseHandler() BaseHandler {
	return &baseHandler{}
}

func WhitelistQueryParams(c *gin.Context, allowedKeys []string) url.Values {
	allowedKeys = append(allowedKeys, []string{"page", "sort", "size"}...)
	params := c.Request.URL.Query()
	for key, _ := range params {
		attr := strings.Split(key, "__")[0]
		if !slices.Contains(allowedKeys, attr) {
			delete(params, key)
		}
	}

	return params
}

func GetIdFromUrl(c *gin.Context, fromQuery bool) (int64, rest_errors.RestErr) {
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

func GetIdentifierFromUrl(c *gin.Context, keyName string, fromQuery bool) string {
	identifier := c.Param(keyName)
	if fromQuery {
		identifier = c.Query(keyName)
	}
	return identifier
}

func (p *baseHandler) Require(attrs []string) *baseHandler {
	p.Errors = []rest_errors.RestErr{}
	for _, attr := range attrs {
		if p.Data[attr] == nil {
			p.Errors = append(p.Errors, rest_errors.InvalidError(fmt.Sprintf("Required Field not found: %s ", attr)))
		}
	}

	return p
}

func (p *baseHandler) Permit(attrs []string) *baseHandler {
	for key, _ := range p.Data {
		if !slices.Contains(attrs, key) {
			delete(p.Data, key)
		}
	}
	return p
}

func (p *baseHandler) ClearEmpty() {
	for key, value := range p.Data {
		if value == nil || value == "" {
			delete(p.Data, key)
		}
	}
}

func (p *baseHandler) SetData(payload map[string]interface{}) *baseHandler {
	p.Errors = []rest_errors.RestErr{}
	data := payload["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	if attributes == nil {
		p.Errors = append(p.Errors, rest_errors.InvalidError("Attributes not found"))
	}

	p.Data = attributes
	return p
}

func (p *baseHandler) CurrentUser(c *gin.Context) *dto.BaseUser {
	// Get current user from context
	userData, ok := c.Get("currentUser")
	if !ok {
		return nil
	}
	return userData.(*dto.BaseUser)
}

var (
	Validate = validator.New(validator.WithRequiredStructEnabled())
)
