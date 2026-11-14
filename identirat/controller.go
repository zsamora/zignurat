package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/zsamora/utils"
)

func IdentiratController() http.Handler {
	log.Println("### Identirat v0.1 ###")
	router := gin.Default()
	// Create DB connection (Handles connection pooling and is threadsafe)
	db := connectGORM()
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/readyz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})
	router.SetFuncMap(template.FuncMap{
		"parseDate": utils.FormatDate,
	})
	router.Static("/img", "./img")
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*")
	store := cookie.NewStore([]byte(utils.GetConfig("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(utils.DropInvalidSessionCookie("session", store))
	router.Use(sessions.Sessions("session", store))
	LandingRoutes(router)
	LoginRoutes(router, db)
	AuthRoutes(router, db)
	return router
}
