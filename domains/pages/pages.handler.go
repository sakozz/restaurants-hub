package pages

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/jsonapi"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type PagesHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}

type pagesHandler struct {
	dao  PagesDao
	base jsonapi.BaseHandler
}

func NewPagesHandler() PagesHandler {
	return &pagesHandler{
		dao:  NewPageDao(),
		base: jsonapi.NewBaseHandler(),
	}
}

func (ctr *pagesHandler) Create(c *gin.Context) {

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
	newRecord := &CreatePagePayload{}
	mapstructure.Decode(payload.Data, &newRecord)

	currentUser := ctr.base.CurrentUser(c)
	/* if currentUser is not admin, set managerId to current user */
	if !currentUser.IsAdmin() {
		newRecord.AuthorId = currentUser.Id
	}

	/* Authorize request for current user */
	authorizer := NewAuthorizer(currentUser, newRecord.AuthorId)
	permissions, restErr := authorizer.Authorize("create")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	/* Generate slug for new record */
	newRecord.Slug = ctr.dao.GenerateSlug(newRecord.Title)

	/* Set authorId to current user */
	newRecord.AuthorId = currentUser.Id

	/* Set restaurantId to current user if current user is manager */
	if currentUser.IsManager() {
		newRecord.RestaurantId = currentUser.RestaurantId
	}

	/* Validate payload data */
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
	resource := restaurant.MemberFor(currentUser.Role)
	jsonPayload := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *pagesHandler) Get(c *gin.Context) {

	slug := jsonapi.GetIdentifierFromUrl(c, "slug", false)
	if slug == "" {
		slugErr := rest_errors.NewBadRequestError("slug is required")
		c.JSON(slugErr.Status(), slugErr)
		return
	}

	restaurant, getErr := ctr.dao.Get(&slug)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	/* Authorize access to resource */
	currentUser := ctr.base.CurrentUser(c)
	authorizer := NewAuthorizer(currentUser, restaurant.AuthorId)
	permissions, restErr := authorizer.Authorize("access")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	meta := map[string]interface{}{
		"permissions": permissions,
	}

	resource := restaurant.MemberFor(currentUser.Role)
	jsonapi := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonapi)
}

func (ctr *pagesHandler) Update(c *gin.Context) {
	slug := jsonapi.GetIdentifierFromUrl(c, "slug", true)
	if slug == "" {
		slugErr := rest_errors.NewBadRequestError("slug is required")
		c.JSON(slugErr.Status(), slugErr)
		return
	}

	/* Check if restaurant exists with given Id */
	record, getErr := ctr.dao.Get(&slug)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	currentUser := ctr.base.CurrentUser(c)
	/* Authorize request for current user */
	authorizer := NewAuthorizer(currentUser, record.AuthorId)
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
	payload.Permit(record.UpdableAttributes(currentUser.Role))

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

	resource := result.MemberFor(currentUser.Role)
	jsonPayload := jsonapi.NewMemberSerializer(resource, nil, nil, meta)
	c.JSON(http.StatusOK, jsonPayload)
}

func (ctr *pagesHandler) List(c *gin.Context) {
	currentUser := ctr.base.CurrentUser(c)
	/* Authorize request for current user */
	authorizer := NewAuthorizer(currentUser)
	_, restErr := authorizer.Authorize("accessCollection")
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	params := jsonapi.WhitelistQueryParams(c, []string{"author_id", "title", "restaurant_id", "visibility"})

	// Get authorized collection of restaurants
	result, err := ctr.dao.AuthorizedCollection(params, ctr.base.CurrentUser(c))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	meta := map[string]interface{}{
		"total": len(result),
	}

	collection := result.CollectionFor(currentUser.Role)
	jsonapi := jsonapi.NewCollectionSerializer(collection, meta)
	c.JSON(http.StatusOK, jsonapi)
}
