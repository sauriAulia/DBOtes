package handlers

import (
	"fmt"
	"strings"

	"tes_dbo/helpers/log"

	"tes_dbo/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	jwt "github.com/dgrijalva/jwt-go"
)

type Context struct {
	Gin       *gin.Engine
	DB        *gorm.DB
	Log       *log.AppLog
	Validator *validator.Validate
}

func (c *Context) Validate(i interface{}) error {
	return c.Validator.Struct(i)
}

func (c *Context) UserValidation(ctx *gin.Context) (role interface{}) {
	infoCustomer := ctx.Request.Header.Get("Authorization")

	splitToken := ""
	if len(strings.Split(infoCustomer, " ")) == 2 {
		splitToken = strings.Split(infoCustomer, " ")[1]
	}

	token, _ := jwt.Parse(splitToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token error")
		}
		return []byte("yoursecretstring"), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["role"]
}

type Result struct {
	Success bool
	Data    interface{}
	Code    int
	Message string
}

func (c *Context) API(group string) {
	protected := c.Gin.Group(group)
	{
		protected.Use(middlewares.JwtAuthMiddleware())

		//Customer
		protected.POST("/insertCustomer", c.InsertCustomer)
		protected.GET("/detailCustomer/:id", c.DetailCustomer)
		protected.PUT("/updateCustomer", c.UpdateCustomer)
		protected.DELETE("/deleteCustomer/:id", c.DeleteCustomer)
		protected.GET("/listCustomer", c.CustomerList)

		// Order
		protected.POST("/createOrder", c.CreateOrder)
		protected.GET("/detailOrder/:id", c.DetailOrder)
		protected.PUT("/updateOrder", c.UpdateOrder)
		protected.DELETE("/deleteOrder/:id", c.DeleteOrder)
		protected.GET("/listOrder", c.OrderListByCustomerId) //list order by idCustomer
		protected.GET("/allOrder", c.AllOrder)

		protected.GET("/loginInfo", c.LoginInfo)

	}

	public := c.Gin.Group(group)
	{
		public.POST("/register/user", c.UserRegistration)
		public.POST("/login", c.Login)
	}
}
