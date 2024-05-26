package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/domains/restaurants"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/jsonapi"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
	"resturants-hub.com/m/v2/services"
)

type AdminRestaurantsHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}

type adminRestaurantsHandler struct {
	service     services.UsersService
	dao         restaurants.RestaurantDao
	payload     RequestPayload
	currentUser *users.User
}

func NewAdminRestaurantsHandler() AdminRestaurantsHandler {
	return &adminRestaurantsHandler{
		service: services.NewUsersService(),
		dao:     restaurants.NewRestaurantDao(),
		payload: NewParamsHandler(),
	}
}

func (ctr *adminRestaurantsHandler) Create(c *gin.Context) {
	/* Authorize request for current user */
	if ok, restErr := ctr.authorize("create", c); !ok {
		c.JSON(restErr.Status(), restErr)
		return
	}

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
	payload := ctr.payload.SetData(mapBody)
	newRestaurant := &restaurants.CreateRestaurantPayload{}
	mapstructure.Decode(payload.Data, &newRestaurant)

	if err := Validate.Struct(newRestaurant); err != nil {
		restErr := rest_errors.NewValidationError(rest_errors.StructValidationErrors(err))
		c.JSON(restErr.Status(), restErr)
		return
	}

	restaurant, getErr := ctr.dao.Create(newRestaurant)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	resource := restaurant.MemberFor(restaurants.AdminDetails)
	jsonPayload := jsonapi.NewMemberSerializer[restaurants.AdminDetailItem](resource, nil, nil)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *adminRestaurantsHandler) Get(c *gin.Context) {
	if ok, restErr := ctr.authorize("access", c); !ok {
		c.JSON(restErr.Status(), restErr)
		return
	}
	id, idErr := getIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	restaurant, getErr := ctr.dao.Get(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	resource := restaurant.MemberFor(restaurants.AdminDetails)
	jsonapi := jsonapi.NewMemberSerializer[restaurants.AdminDetailItem](resource, nil, nil)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *adminRestaurantsHandler) Update(c *gin.Context) {
	/* Authorize request for current user */
	if ok, restErr := ctr.authorize("update", c); !ok {
		c.JSON(restErr.Status(), restErr)
		return
	}

	id, idErr := getIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	/* Check if restaurant exists with given Id */
	record, getErr := ctr.dao.Get(&id)
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
	payload.Permit(record.AdminUpdableAttributes())

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

	result, updateErr := ctr.dao.Update(record, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := result.MemberFor(restaurants.AdminDetails)
	jsonPayload := jsonapi.NewMemberSerializer[restaurants.AdminDetailItem](resource, nil, nil)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *adminRestaurantsHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	if ok, restErr := ctr.authorize("accessCollection", c); !ok {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := WhitelistQueryParams(c, []string{"profile_id", "name", "email", "phone"})

	// Get authorized collection of restaurants
	result, err := ctr.dao.AuthorizedCollection(params, ctr.currentUser)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(restaurants.AdminList)
	jsonapi := jsonapi.NewCollectionSerializer[restaurants.AdminListItem](collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *adminRestaurantsHandler) authorize(action string, c *gin.Context) (bool, rest_errors.RestErr) {

	// Get current user from context
	userData, ok := c.Get("currentUser")
	if !ok {
		return false, rest_errors.NewUnauthorizedError("unauthorized")
	}

	ctr.currentUser = userData.(*users.User)

	// Check if user is authorized to access the resource
	return ctr.currentUser.Can(action, consts.Restaurants)
}
