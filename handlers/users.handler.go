package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/services"
	rest_errors "resturants-hub.com/m/v2/utils"
)

type UsersHandler interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
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

func getIdFromUrl(c *gin.Context, fromQuery bool) (int64, rest_errors.RestErr) {
	paramId := c.Param("id")
	if fromQuery {
		paramId = c.Query("id")
	}
	id, userErr := strconv.ParseInt(paramId, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return id, nil
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

	c.JSON(http.StatusOK, user.Serialize(users.Admin))
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

	result, updateErr := ctr.service.UpdateUser(user, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}
	c.JSON(http.StatusOK, result.Serialize(users.Admin))
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
	c.JSON(http.StatusOK, result.Serialize(users.OwnerUser))
}
