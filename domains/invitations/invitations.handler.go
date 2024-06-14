package invitations

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/jsonapi"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type InvitationsHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type invitationsHandler struct {
	dao  InvitationsDao
	base jsonapi.BaseHandler
}

func NewInvitationsHandler() InvitationsHandler {
	return &invitationsHandler{
		dao:  NewInvitationDao(),
		base: jsonapi.NewBaseHandler(),
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
	newRecord := &CreateInvitationPayload{}
	mapstructure.Decode(payload.Data, &newRecord)

	/* Authorize request for current user */
	permissions, restErr := ctr.Authorize("create", nil, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	if err := jsonapi.Validate.Struct(newRecord); err != nil {
		restErr := rest_errors.NewValidationError(rest_errors.StructValidationErrors(err))
		c.JSON(restErr.Status(), restErr)
		return
	}

	restaurant, getErr := ctr.dao.Create(newRecord)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	resource := restaurant.MemberFor()
	jsonPayload := jsonapi.NewMemberSerializer[Invitation](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *invitationsHandler) Get(c *gin.Context) {
	id, idErr := jsonapi.GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	invitation, getErr := ctr.dao.Get(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	permissions, restErr := ctr.Authorize("access", invitation, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := invitation.MemberFor()
	jsonapi := jsonapi.NewMemberSerializer[Invitation](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)

}

func (ctr *invitationsHandler) Update(c *gin.Context) {
	id, idErr := jsonapi.GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	/* Check if user exists with given Id */
	invitation, getErr := ctr.dao.Get(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	permissions, restErr := ctr.Authorize("update", invitation, c)
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

	updatedUser, updateErr := ctr.dao.Update(invitation, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := updatedUser.MemberFor()
	jsonapi := jsonapi.NewMemberSerializer[Invitation](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *invitationsHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	_, restErr := ctr.Authorize("accessCollection", nil, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := jsonapi.WhitelistQueryParams(c, []string{"email", "token", "expires_at"})
	result, err := ctr.dao.AuthorizedCollection(params, ctr.base.CurrentUser(c))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor()
	jsonapi := jsonapi.NewCollectionSerializer[Invitation](collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
