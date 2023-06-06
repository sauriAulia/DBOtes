package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"tes_dbo/handlers"
	"tes_dbo/helpers/env"
	"tes_dbo/helpers/log"
	"tes_dbo/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	l := log.NewLog("Log", "")
	defer l.Close()
	l.Info("Set Up Log Started")

	cs := env.Get().ConnectionString()
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})

	if err != nil {
		l.Fatal("Connection To DB Failed")
		return
	}
	l.Info("Connection to DB Successfully")

	err = models.MigrateModel(db)
	if err != nil {
		l.Error(err)
		return
	}
	l.Info("Auto Migration Successfully")

	gin.SetMode(gin.ReleaseMode) //enable when you want to production
	g := gin.Default()

	//Setting Config CORS
	c := cors.DefaultConfig()
	c.AllowWildcard = true
	c.AllowCredentials = true
	c.AllowOrigins = []string{"*"}
	//c.AllowOrigins = []string{"https://*.nozyra.com"}
	c.AddAllowHeaders("Authorization", "Content-Type")
	c.AddExposeHeaders("Authorization", "Content-Type")

	//Use CORS
	g.Use(cors.New(c))

	//use Validation
	v := validator.New()

	//Inisialisasi Handler / Endpoint
	h := handlers.Context{Gin: g, DB: db, Log: l, Validator: v}

	h.API("api")

	Host := env.Get().AppHost
	Port := env.Get().AppPort

	Address := fmt.Sprintf("%s:%d", Host, Port)
	s := &http.Server{Addr: Address, Handler: g}

	l.Infof("start listen and serve at %s: %v", Host, Port)

	err = s.ListenAndServe()
	if err != nil {
		l.Fatal("failed to connect to server")
		return
	}

}
