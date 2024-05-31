package users

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/jsonapi"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type UsersHandler interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
	Profile(c *gin.Context)
	List(c *gin.Context)
}

type usersHandler struct {
	service     UsersService
	dao         UsersDao
	payload     jsonapi.RequestPayload
	currentUser *User
}

func NewUsersHandler() UsersHandler {
	return &usersHandler{
		service: NewUsersService(),
		dao:     NewUserDao(),
		payload: jsonapi.NewParamsHandler(),
	}
}

func (ctr *usersHandler) Get(c *gin.Context) {
	userId, idErr := jsonapi.GetIdFromUrl(c, false)
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
	permissions, restErr := ctr.Authorize("access", user, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := user.MemberFor(AdminDetails)
	jsonapi := jsonapi.NewMemberSerializer[AdminDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)

}

func (ctr *usersHandler) Profile(c *gin.Context) {
	session, ok := c.Get("currentSession")
	if !ok {
		restError := rest_errors.InvalidError("unauthorized")
		c.JSON(restError.Status(), restError)
		return
	}

	user, getErr := ctr.service.GetUser(session.(*Session).UserId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	permissions, restErr := ctr.Authorize("access", user, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions":    permissions,
		"appPermissions": user.Permissions(),
	}

	resource := user.MemberFor(OwnerDetails)
	jsonapi := jsonapi.NewMemberSerializer[OwnerDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *usersHandler) Update(c *gin.Context) {
	userId, idErr := jsonapi.GetIdFromUrl(c, false)
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
	permissions, restErr := ctr.Authorize("update", user, c)
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
	payload := ctr.payload.SetData(mapBody)
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

	resource := updatedUser.MemberFor(OwnerDetails)
	jsonapi := jsonapi.NewMemberSerializer[AdminDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *usersHandler) Delete(c *gin.Context) {
	// userId, idErr := getUserId(c.Param("user_id"))
	// if idErr != nil {
	// 	c.JSON(idErr.Status(), idErr)
	// 	return
	// }

	// if err := services.UsersService.DeleteUser(userId); err != nil {
	// 	c.JSON(err.Status(), err)
	// 	return
	// }
	// c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (ctr *usersHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	_, restErr := ctr.Authorize("accessCollection", nil, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := jsonapi.WhitelistQueryParams(c, []string{"first_name", "email", "id", "last_name"})
	result, err := ctr.dao.AuthorizedCollection(params, ctr.currentUser)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(AdminList)
	jsonapi := jsonapi.NewCollectionSerializer[AdminListItem](collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
