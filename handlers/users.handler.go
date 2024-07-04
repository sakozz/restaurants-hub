package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/authorizer"
	"resturants-hub.com/m/v2/dao"
	"resturants-hub.com/m/v2/dto"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
	"resturants-hub.com/m/v2/serializers"
	"resturants-hub.com/m/v2/services"
)

type UsersHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	Profile(c *gin.Context)
	List(c *gin.Context)
}

type usersHandler struct {
	service services.UsersService
	dao     dao.UsersDao
	base    BaseHandler
}

func NewUsersHandler() UsersHandler {
	return &usersHandler{
		service: services.NewUsersService(),
		dao:     dao.NewUsersDao(),
		base:    NewBaseHandler(),
	}
}

func (ctr *usersHandler) Create(c *gin.Context) {

	/* Extract request body as map */
	var mapBody map[string]interface{}
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	/* extract data as json/map  */
	json.Unmarshal(data, &mapBody)

	/* Parse jsonapi payload and set attributes to data*/
	payload := ctr.base.SetData(mapBody)
	newRecord := &dto.CreateUserPayload{}
	mapstructure.Decode(payload.Data, &newRecord)

	/* Authorize request for current user */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewUsersAuthorizer(currentUser, newRecord.Id)
	permissions, restErr := authorizer.Authorize("create")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	if err := Validate.Struct(newRecord); err != nil {
		restErr := rest_errors.NewValidationError(rest_errors.StructValidationErrors(err))
		c.JSON(restErr.Status(), restErr)
		return
	}

	user, getErr := ctr.dao.CreateUser(newRecord)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	resource := user.MemberFor(currentUser.Role)
	jsonPayload := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *usersHandler) Get(c *gin.Context) {
	userId, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user, getErr := ctr.service.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewUsersAuthorizer(currentUser, user.Id)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := user.MemberFor(currentUser.Role)
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)

}

func (ctr *usersHandler) Profile(c *gin.Context) {
	session, ok := c.Get("currentSession")
	if !ok {
		restError := rest_errors.InvalidError("unauthorized")
		c.JSON(restError.Status(), restError)
		return
	}

	user, getErr := ctr.service.GetUser(session.(*dto.Session).UserId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewUsersAuthorizer(currentUser, user.Id)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions":    permissions,
		"appPermissions": user.Permissions(),
	}

	resource := user.MemberFor(currentUser.Role)
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *usersHandler) Update(c *gin.Context) {
	userId, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	/* Check if user exists with given Id */
	user, getErr := ctr.service.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewUsersAuthorizer(currentUser, user.Id)
	permissions, restErr := authorizer.Authorize("update")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	/* Extract request body as map */
	var mapBody map[string]interface{}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	/* Validate required params and whitelisted payload data */
	json.Unmarshal(jsonData, &mapBody)
	payload := ctr.base.SetData(mapBody)
	payload.Require([]string{"id"}).Permit(user.UpdableAttributes())

	/* Skip empty data and patch with only new data if the update is partial(PATCH) */
	isPartial := c.Request.Method == http.MethodPatch
	if isPartial {
		payload.ClearEmpty()
	}

	/* Return error if payload has eroor for require/permit */
	if len(payload.Errors) > 0 {
		c.JSON(payload.Errors[0].Status(), payload.Errors)
		return
	}

	updatedUser, updateErr := ctr.service.UpdateUser(user, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := updatedUser.MemberFor(currentUser.Role)
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *usersHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewUsersAuthorizer(currentUser)
	_, restErr := authorizer.Authorize("accessCollection")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := WhitelistQueryParams(c, []string{"first_name", "email", "id", "last_name"})
	result, err := ctr.dao.AuthorizedUsersCollection(params, currentUser)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(currentUser.Role)
	jsonapi := serializers.NewCollectionSerializer(collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
