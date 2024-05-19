package controller

import (
	"fmt"

	status "github.com/Go_CleanArch/common/const"
	userService "github.com/Go_CleanArch/usecase/service/user"
	"github.com/gin-gonic/gin"
)

// UserController is user controlller
type UserController struct {
	userService userService.UserService
}

// NewUserController is the constructor for UserController
func NewUserController(userService userService.UserService) *UserController {
	return &UserController{userService: userService}
}

// Create action: POST /users
func (uc *UserController) UserController(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := uc.userService.CreateUserService(ctx, c)
	if err != nil {
		fmt.Println(err)
	} else {
		c.JSON(
			status.SuccessStatusMap["CREATED"].StatusCode,
			result,
		)
	}
}

// Show action: POST /users/login
func (uc *UserController) LoginControler(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := uc.userService.LoginService(ctx, c)
	if err != nil {
		fmt.Println(err)
	} else {
		c.JSON(
			status.SuccessStatusMap["OK"].StatusCode,
			result,
		)
	}
}
