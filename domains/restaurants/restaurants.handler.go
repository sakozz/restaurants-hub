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

type AdminRestaurantsHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}

type adminRestaurantsHandler struct {
	dao         RestaurantDao
	payload     jsonapi.RequestPayload
	currentUser *users.User
}

func NewAdminRestaurantsHandler() AdminRestaurantsHandler {
	return &adminRestaurantsHandler{
		dao:     NewRestaurantDao(),
		payload: jsonapi.NewParamsHandler(),
	}
}

func (ctr *adminRestaurantsHandler) Create(c *gin.Context) {

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
	newRestaurant := &CreateRestaurantPayload{}
	mapstructure.Decode(payload.Data, &newRestaurant)

	/* Authorize request for current user */
	permissions, restErr := ctr.Authorize("create", nil, c)
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
	resource := restaurant.MemberFor(AdminDetails)
	jsonPayload := jsonapi.NewMemberSerializer[AdminDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *adminRestaurantsHandler) Get(c *gin.Context) {

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

	permissions, restErr := ctr.Authorize("access", restaurant, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(AdminDetails)
	jsonapi := jsonapi.NewMemberSerializer[AdminDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *adminRestaurantsHandler) Update(c *gin.Context) {
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
	permissions, restErr := ctr.Authorize("access", record, c)
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
	jsonPayload := jsonapi.NewMemberSerializer[AdminDetailItem](resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *adminRestaurantsHandler) List(c *gin.Context) {
	/* Authorize request for current user */
	_, restErr := ctr.Authorize("access", nil, c)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := jsonapi.WhitelistQueryParams(c, []string{"profile_id", "name", "email", "phone"})

	// Get authorized collection of restaurants
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
