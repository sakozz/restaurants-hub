package restaurants

import (
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/users"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

func (ctr *adminRestaurantsHandler) Authorize(action string, c *gin.Context) (bool, rest_errors.RestErr) {

	// Get current user from context
	userData, ok := c.Get("currentUser")
	if !ok {
		return false, rest_errors.NewUnauthorizedError("unauthorized")
	}

	ctr.currentUser = userData.(*users.User)

	// Check if user is authorized to access the resource
	return ctr.currentUser.Can(action, consts.Restaurants)
}
