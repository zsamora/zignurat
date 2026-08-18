package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

func InscriptionRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/inscriptions", Inscriptions(db))
		org.POST("/createInscription", CreateInscription(db))
	}
}

func Inscriptions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		inscriptions := getInscriptionsFromOrg(db, orgID)
		c.HTML(http.StatusOK, "orgInscriptions.html", gin.H{
			"OrgID":        orgID,
			"Org":          GetOrganizationFromID(db, orgID),
			"inscriptions": inscriptions,
		})
	}
}

func CreateInscription(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgID := uint(orgIDInt)
		inscription := Inscription{
			OrgID: orgID,
		}
		createInscription(db, inscription)
		inscriptions := getInscriptionsFromOrg(db, orgID)
		c.HTML(http.StatusOK, "orgInscriptions.html", gin.H{
			"OrgID":        orgID,
			"Org":          GetOrganizationFromID(db, orgID),
			"inscriptions": inscriptions,
		})
	}
}
