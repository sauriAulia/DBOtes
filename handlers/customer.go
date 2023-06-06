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

func (c *Context) InsertCustomer(ctx *gin.Context) {
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

	t := models.Customer{}

	err := ctx.ShouldBindJSON(&t)

	if err != nil {
		Results := &Result{
			Success: false,
			Data:    err.Error(),
			Code:    http.StatusInternalServerError,
			Message: "Gagal mengirimkan data",
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
				Message: "Input yang dimasukkan salah",
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

	err = t.InsertCustomer(c.DB)

	if err != nil {
		Results := &Result{
			Success: false,
			Data:    err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Gagal membuat role baru",
		}
		c.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, Results)
		return
	}

	result := gin.H{
		"success": true,
		"data":    t,
		"code":    http.StatusCreated,
		"message": "Success Insert Customer",
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *Context) DetailCustomer(ctx *gin.Context) {
	const param = "id"
	id := (ctx.Param((param)))
	idValue, _ := strconv.Atoi(id)
	u := models.Customer{}
	err := u.SelectById(c.DB, idValue)

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

func (c *Context) UpdateCustomer(ctx *gin.Context) {
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

	//get data from user
	m := models.Customer{}
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
	p := models.Customer{}
	err = p.SelectById(c.DB, int(m.CustomerID))
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
	p.FullName = m.FullName
	p.Email = m.Email
	p.PhoneNumber = m.PhoneNumber

	//lakukan update
	err = p.UpdateCustomer(c.DB)

	//cek success or error
	if err != nil {
		result := gin.H{
			"success": false,
			"data":    p,
			"code":    http.StatusInternalServerError,
			"message": "Failed Editing Data",
		}
		c.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result := gin.H{
		"success": true,
		"data":    p,
		"code":    http.StatusOK,
		"message": "Success Editing Data",
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *Context) DeleteCustomer(ctx *gin.Context) {
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

	id, _ := strconv.Atoi(ctx.Param("id"))

	r := models.Customer{}

	err := r.SelectById(c.DB, id)

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

	err = r.DeleteCustomer(c.DB)
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

func (c *Context) CustomerList(ctx *gin.Context) {
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

	p := models.Customers{}
	err := p.SelectAll(c.DB)

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
		"jumlahData": len(p),
		"data":       p,
		"code":       http.StatusOK,
		"message":    "Success Get All Data",
	}
	ctx.JSON(http.StatusOK, result)
}
