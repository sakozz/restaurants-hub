package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/restaurants"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/services"
	rest_errors "resturants-hub.com/m/v2/utils"
)

type AdminRestaurantsHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type adminRestaurantsHandler struct {
	service services.UsersService
	dao     restaurants.RestaurantDao
	payload RequestPayload
}

func NewAdminRestaurantsHandler() AdminRestaurantsHandler {
	return &adminRestaurantsHandler{
		service: services.NewUsersService(),
		dao:     restaurants.NewRestaurantDao(),
		payload: NewParamsHandler(),
	}
}

func (ctr *adminRestaurantsHandler) Create(c *gin.Context) {
	newRestaurant := &restaurants.CreateRestaurantPayload{}
	if err := c.ShouldBindJSON(newRestaurant); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := newRestaurant.Validate(); err != nil {
		restErr := rest_errors.NewBadRequestError("Validation error")
		c.JSON(restErr.Status(), restErr)
		return
	}

	restaurant, getErr := ctr.dao.Create(newRestaurant)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	c.JSON(http.StatusOK, restaurant.Serialize(restaurants.AdminList))
}

func (ctr *adminRestaurantsHandler) Get(c *gin.Context) {
	id, idErr := getIdFromUrl(c, false)
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user, getErr := ctr.service.GetUser(id)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	c.JSON(http.StatusOK, user.Serialize(users.Admin))
}

func (ctr *adminRestaurantsHandler) List(c *gin.Context) {
	params := WhitelistQueryParams(c, []string{"profile_id", "name", "email", "phone"})
	result, err := ctr.dao.Search(params)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result.Serialize(restaurants.AdminList))
}
