package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tes_dbo/helpers/validate"
	"tes_dbo/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (c *Context) CreateOrder(ctx *gin.Context) {
	t := models.Order{}
	err := ctx.ShouldBindJSON(&t)

	if err != nil {
		Results := &Result{
			Success: false,
			Data:    err.Error(),
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		c.Log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, Results)
		return
	}

	err = c.Validate(&t)

	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]validate.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = validate.ErrorMsg{Field: fe.Field(), Message: validate.GetErrorMsg(fe)}
			}
			Results := &Result{
				Success: false,
				Data:    out,
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			c.Log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, Results)
			return
		} else {
			c.Log.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	err = t.CreateOrder(c.DB)

	if err != nil {
		Results := &Result{
			Success: false,
			Data:    err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Failed Create Order",
		}
		c.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, Results)
		return
	}

	result := gin.H{
		"success": true,
		"data":    t,
		"code":    http.StatusCreated,
		"message": "Success Create Order",
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *Context) DetailOrder(ctx *gin.Context) {
	const param = "id"
	id := (ctx.Param((param)))
	idValue, _ := strconv.Atoi(id)
	u := models.Order{}
	err := u.SelectByOrderId(c.DB, idValue)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusNotFound,
				"message": "Data Not Found",
			}
			c.Log.Warn(err)
			ctx.JSON(http.StatusNotFound, result)
			return
		} else {
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusInternalServerError,
				"message": "Server Error",
			}
			c.Log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, result)
			return
		}
	}
	result := gin.H{
		"success": true,
		"data":    u,
		"code":    http.StatusOK,
		"message": "Success Get Data",
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Context) UpdateOrder(ctx *gin.Context) {
	//get data from user
	m := models.Order{}
	err := ctx.ShouldBindJSON(&m)
	if err != nil {
		Results := &Result{
			Success: false,
			Data:    err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Bad Request !",
		}
		c.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, Results)
		return
	}

	//get data from db
	p := models.Order{}
	err = p.SelectByOrderId(c.DB, int(m.OrderID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warn(err)
			ctx.JSON(http.StatusNotFound, nil)
			return
		} else {
			c.Log.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	//ganti data yang mau diupdate
	p.Product = m.Product
	p.Quantity = m.Quantity

	//lakukan update
	err = p.UpdateOrder(c.DB)

	//cek success or error
	if err != nil {
		result := gin.H{
			"success": false,
			"data":    p,
			"code":    http.StatusInternalServerError,
			"message": "Failed Updating Data",
		}
		c.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result := gin.H{
		"success": true,
		"data":    p,
		"code":    http.StatusOK,
		"message": "Success Updating Data",
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *Context) DeleteOrder(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	r := models.Order{}
	err := r.SelectByOrderId(c.DB, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warn(err)
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusNotFound,
				"message": "Data Not Found",
			}
			ctx.JSON(http.StatusNotFound, result)
			return
		} else {
			c.Log.Error(err)
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error !",
			}
			ctx.JSON(http.StatusInternalServerError, result)
			return
		}
	}

	err = r.DeleteOrder(c.DB)
	if err != nil {
		c.Log.Error(err)
		result := gin.H{
			"success": false,
			"data":    nil,
			"code":    http.StatusInternalServerError,
			"message": "Failed Deleting Data !",
		}
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result := gin.H{
		"success": true,
		"data":    "",
		"code":    http.StatusOK,
		"message": "Success Deleting Data !",
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Context) OrderListByCustomerId(ctx *gin.Context) {
	customerid, _ := strconv.Atoi(ctx.Query("customerId"))
	// fmt.Println("customerId", customerid)
	u := models.Orders{}
	err := u.GetOrderListByCustomerId(c.DB, customerid)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusNotFound,
				"message": "Data Not Found",
			}
			c.Log.Warn(err)
			ctx.JSON(http.StatusNotFound, result)
			return
		} else {
			result := gin.H{
				"success": false,
				"data":    nil,
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
			}
			c.Log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, result)
			return
		}
	}

	if len(u) == 0 {
		result := gin.H{
			"success": true,
			"data":    u,
			"code":    http.StatusNotFound,
			"message": "Data Not Found",
		}

		ctx.JSON(http.StatusNotFound, result)
		return
	}
	result := gin.H{
		"success": true,
		"data":    u,
		"code":    http.StatusOK,
		"message": "Success Get Data",
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Context) AllOrder(ctx *gin.Context) {
	role := c.UserValidation(ctx)
	if role != "admin" {
		Results := &Result{
			Success: false,
			Data:    nil,
			Code:    http.StatusForbidden,
			Message: "Access Denied !",
		}
		ctx.JSON(http.StatusInternalServerError, Results)
		return
	}

	t := models.Orders{}

	err := t.SelectAllOrder(c.DB)

	if err != nil {
		result := gin.H{
			"success": false,
			"data":    nil,
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		}
		c.Log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result := gin.H{
		"success":    true,
		"jumlahData": len(t),
		"data":       t,
		"code":       http.StatusOK,
		"message":    "Success Get All Data",
	}
	ctx.JSON(http.StatusOK, result)
}
