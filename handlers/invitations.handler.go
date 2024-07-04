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
)

type InvitationsHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type invitationsHandler struct {
	dao  dao.InvitationsDao
	base BaseHandler
}

func NewInvitationsHandler() InvitationsHandler {
	return &invitationsHandler{
		dao:  dao.NewInvitationDao(),
		base: NewBaseHandler(),
	}
}

func (ctr *invitationsHandler) Create(c *gin.Context) {

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
	newRecord := &dto.CreateInvitationPayload{}
	mapstructure.Decode(payload.Data, &newRecord)

	/* Authorize request for current user */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewInvitationAuthorizer(currentUser, newRecord.Email)
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

	restaurant, getErr := ctr.dao.CreateInvitation(newRecord)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	resource := restaurant.MemberFor()
	jsonPayload := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *invitationsHandler) Get(c *gin.Context) {
	id, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	invitation, getErr := ctr.dao.GetInvitation(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewInvitationAuthorizer(currentUser, invitation.Email)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := invitation.MemberFor()
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)

}

func (ctr *invitationsHandler) Update(c *gin.Context) {
	id, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	/* Check if user exists with given Id */
	invitation, getErr := ctr.dao.GetInvitation(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewInvitationAuthorizer(currentUser, invitation.Email)
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
	payload.Require([]string{"id"}).Permit(invitation.UpdableAttributes())

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

	updatedUser, updateErr := ctr.dao.UpdateInvitation(invitation, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := updatedUser.MemberFor()
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *invitationsHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := authorizer.NewInvitationAuthorizer(currentUser)
	_, restErr := authorizer.Authorize("accessCollection")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := WhitelistQueryParams(c, []string{"email", "token", "expires_at"})
	result, err := ctr.dao.AuthorizedInvitationsCollection(params, ctr.base.CurrentUser(c))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor()
	jsonapi := serializers.NewCollectionSerializer(collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
