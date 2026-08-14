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

func NotariatController() http.Handler {
	log.Println("### Notariat v0.1 ###")
	router := gin.Default()
	// Create DB connection (Handles connection pooling and is threadsafe)
	db := connectGORM()
	router.SetFuncMap(template.FuncMap{
		"parseDate":    utils.FormatDate,
		"parseOrgType": ParseOrgType,
	})
	router.LoadHTMLGlob("templates/**/*")
	store := cookie.NewStore([]byte(utils.GetConfig("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("session", store))
	LandingRoutes(router)
	LoginRoutes(router)
	AdminRoutes(router, db)
	InternalRoutes(router, db)
	return router
}
