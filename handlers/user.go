package handlers

import (
	"net/http"
	"tes_dbo/models"
	"tes_dbo/utils/token"

	"github.com/gin-gonic/gin"
)

func (c *Context) UserRegistration(ctx *gin.Context) {
	type RegisterInput struct {
		FullName    string `json:"fullName" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Role        string `json:"role" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
	}

	var input RegisterInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.FullName = input.FullName
	u.Username = input.Username
	u.Password = input.Password
	u.Email = input.Email
	u.PhoneNumber = input.PhoneNumber
	u.Role = input.Role

	err := u.InsertUser(c.DB)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	result := gin.H{
		"success": true,
		"data":    u,
		"code":    http.StatusOK,
		"message": "Success Create Customer",
	}

	ctx.JSON(http.StatusOK, result)

}

func (c *Context) Login(ctx *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheck(c.DB, u.Username, u.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	// ctx.JSON(http.StatusOK, gin.H{"token": token})

	result := gin.H{
		"success": true,
		"token":   token,
		"code":    http.StatusOK,
		"message": "Login Success",
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Context) LoginInfo(ctx *gin.Context) {

	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(c.DB, user_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{
		"success": true,
		"data":    u,
		"code":    http.StatusOK,
		"message": "Success Get Login Info",
	}

	ctx.JSON(http.StatusOK, result)
}
