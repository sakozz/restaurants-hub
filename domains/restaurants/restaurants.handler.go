package restaurants

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/jsonapi"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type RestaurantsHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	MyRestaurant(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}

type restaurantsHandler struct {
	dao      RestaurantDao
	usersDao users.UsersDao
	base     jsonapi.BaseHandler
}

func NewAdminRestaurantsHandler() RestaurantsHandler {
	return &restaurantsHandler{
		dao:      NewRestaurantDao(),
		usersDao: users.NewUsersDao(),
		base:     jsonapi.NewBaseHandler(),
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
	newRestaurant := &CreateRestaurantPayload{}
	mapstructure.Decode(payload.Data, &newRestaurant)

	currentUser := ctr.base.CurrentUser(c)
	/* if currentUser is not admin, set managerId to current user */
	if !currentUser.IsAdmin() {
		newRestaurant.ManagerId = currentUser.Id
	}

	/* Authorize request for current user */
	authorizer := NewAuthorizer(currentUser, newRestaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("create")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	if err := jsonapi.Validate.Struct(newRestaurant); err != nil {
		restErr := rest_errors.NewValidationError(rest_errors.StructValidationErrors(err))
		c.JSON(restErr.Status(), restErr)
		return
	}

	restaurant, getErr := ctr.dao.Create(newRestaurant)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Set restaurantId to current user */
	if !currentUser.IsAdmin() {
		_, updateErr := ctr.usersDao.Update(&currentUser.Id, map[string]interface{}{"restaurant_id": restaurant.Id})
		if updateErr != nil {
			c.JSON(updateErr.Status(), updateErr)
			return
		}
	}

	resource := restaurant.MemberFor(AdminDetails)
	jsonPayload := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *restaurantsHandler) Get(c *gin.Context) {

	id, idErr := jsonapi.GetIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	restaurant, getErr := ctr.dao.Get(&id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	authorizer := NewAuthorizer(ctr.base.CurrentUser(c), restaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(AdminDetails)
	jsonapi := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *restaurantsHandler) MyRestaurant(c *gin.Context) {

	restaurant, getErr := ctr.dao.RestaurantByOwnerId(&ctr.base.CurrentUser(c).Id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	authorizer := NewAuthorizer(ctr.base.CurrentUser(c), restaurant.ManagerId)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(OwnerDetails)
	jsonapi := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *restaurantsHandler) Update(c *gin.Context) {
	id, idErr := jsonapi.GetIdFromUrl(c, false)
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

	/* Authorize request for current user */
	authorizer := NewAuthorizer(ctr.base.CurrentUser(c), record.ManagerId)
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

	result, updateErr := ctr.dao.Update(record, payload.Data)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}

	resource := result.MemberFor(AdminDetails)
	jsonPayload := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *restaurantsHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	authorizer := NewAuthorizer(ctr.base.CurrentUser(c))
	_, restErr := authorizer.Authorize("accessCollection")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := jsonapi.WhitelistQueryParams(c, []string{"user_id", "name", "email", "phone"})

	// Get authorized collection of restaurants
	result, err := ctr.dao.AuthorizedCollection(params, ctr.base.CurrentUser(c))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(AdminList)
	jsonapi := jsonapi.NewCollectionSerializer(collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
