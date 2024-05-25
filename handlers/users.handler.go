package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/jsonapi"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
	"resturants-hub.com/m/v2/services"
)

type UsersHandler interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
	Profile(c *gin.Context)
	List(c *gin.Context)
}

type usersHandler struct {
	service services.UsersService
	payload RequestPayload
}

func NewUsersHandler() UsersHandler {
	return &usersHandler{
		service: services.NewUsersService(),
		payload: NewParamsHandler(),
	}
}

func (ctr *usersHandler) Get(c *gin.Context) {
	userId, idErr := getIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user, getErr := ctr.service.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	resource := user.MemberFor(users.AdminDetails)
	jsonapi := jsonapi.NewMemberSerializer[users.AdminDetailItem](resource, nil, nil)
	c.JSON(http.StatusOK, jsonapi)

}

func (ctr *usersHandler) Profile(c *gin.Context) {
	session, ok := c.Get("currentSession")
	if !ok {
		restError := rest_errors.InvalidError("unauthorized")
		c.JSON(restError.Status(), restError)
		return
	}

	user, getErr := ctr.service.GetUser(session.(*users.Session).ProfileId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	resource := user.MemberFor(users.OwnerDetails)
	jsonapi := jsonapi.NewMemberSerializer[users.OwnerDetailItem](resource, nil, nil)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *usersHandler) Update(c *gin.Context) {
	userId, idErr := getIdFromUrl(c, false)
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

	resource := updatedUser.MemberFor(users.OwnerDetails)
	jsonapi := jsonapi.NewMemberSerializer[users.AdminDetailItem](resource, nil, nil)
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
	params := WhitelistQueryParams(c, []string{"first_name", "email", "id", "last_name"})
	result, err := ctr.service.SearchUser(params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(users.AdminList)
	jsonapi := jsonapi.NewCollectionSerializer[users.AdminListItem](collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
