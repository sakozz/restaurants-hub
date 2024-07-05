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

type RestaurantsHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	MyRestaurant(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}

type restaurantsHandler struct {
	dao      dao.RestaurantDao
	usersDao dao.UsersDao
	base     BaseHandler
}

func NewAdminRestaurantsHandler() RestaurantsHandler {
	return &restaurantsHandler{
		dao:      dao.NewRestaurantDao(),
		usersDao: dao.NewUsersDao(),
		base:     NewBaseHandler(),
	}
}

func (ctr *restaurantsHandler) Create(c *gin.Context) {

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
	newRestaurant := &dto.CreateRestaurantPayload{}
	mapstructure.Decode(payload.Data, &newRestaurant)

	currentUser := ctr.base.CurrentUser(c)
	/* if currentUser is not admin, set managerId to current user */
	if !currentUser.IsAdmin() {
		newRestaurant.ManagerId = currentUser.Id
	}

	/* Authorize request for current user */
	authorizer := authorizer.NewRestaurantsAuthorizer(currentUser, newRestaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("create")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	if err := Validate.Struct(newRestaurant); err != nil {
		restErr := rest_errors.NewValidationError(rest_errors.StructValidationErrors(err))
		c.JSON(restErr.Status(), restErr)
		return
	}

	restaurant, getErr := ctr.dao.CreateRestaurant(newRestaurant)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Set restaurantId to current user */
	if !currentUser.IsAdmin() {
		_, updateErr := ctr.usersDao.UpdateUser(&currentUser.Id, map[string]interface{}{"restaurant_id": restaurant.Id})
		if updateErr != nil {
			c.JSON(updateErr.Status(), updateErr)
			return
		}
	}

	resource := restaurant.MemberFor(currentUser.Role)
	jsonPayload := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *restaurantsHandler) Get(c *gin.Context) {

	id, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	restaurant, getErr := ctr.dao.GetRestaurant(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	currentUser := ctr.base.CurrentUser(c)
	/* Authorize access to resource */
	authorizer := authorizer.NewRestaurantsAuthorizer(currentUser, restaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(currentUser.Role)
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *restaurantsHandler) MyRestaurant(c *gin.Context) {

	restaurant, getErr := ctr.dao.RestaurantByOwnerId(&ctr.base.CurrentUser(c).Id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	currentUser := ctr.base.CurrentUser(c)
	/* Authorize access to resource */
	authorizer := authorizer.NewRestaurantsAuthorizer(currentUser, restaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(currentUser.Role)
	jsonapi := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *restaurantsHandler) Update(c *gin.Context) {
	id, idErr := GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	/* Check if restaurant exists with given Id */
	record, getErr := ctr.dao.GetRestaurant(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	currentUser := ctr.base.CurrentUser(c)
	/* Authorize request for current user */
	authorizer := authorizer.NewRestaurantsAuthorizer(currentUser, record.ManagerId)
	permissions, restErr := authorizer.Authorize("access")
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

	result, updateErr := ctr.dao.UpdateRestaurant(record, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := result.MemberFor(currentUser.Role)
	jsonPayload := serializers.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *restaurantsHandler) List(c *gin.Context) {
	currentUser := ctr.base.CurrentUser(c)
	/* Authorize request for current user */
	authorizer := authorizer.NewRestaurantsAuthorizer(currentUser)
	_, restErr := authorizer.Authorize("accessCollection")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := WhitelistQueryParams(c, []string{"user_id", "name", "email", "phone"})

	// Get authorized collection of restaurants
	result, err := ctr.dao.AuthorizedRestaurantCollection(params, ctr.base.CurrentUser(c))
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
