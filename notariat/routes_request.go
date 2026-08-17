package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

// RequestRoutes registers the public certificate-request flow (no account
// needed to submit a request or look one up by mail) plus the org-side inbox
// for reviewing requests addressed to that organization. A logged-in Usuario
// hitting /getRequests gets their own requests auto-scoped by their account's
// mail instead of an empty search box.
func RequestRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/requestCertificate", utils.AuthMiddleware(2), RequestCertificateForm(db))
	router.POST("/requestCertificate", SaveRequestCertificate(db))
	router.GET("/getRequests", utils.AuthMiddleware(2), SearchRequest(db))
	router.POST("/getRequestsByMail", GetRequestsByMail(db))

	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/requests", Requests(db))
	}
}

func RequestCertificateForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := GetOrganizations(db)
		c.HTML(http.StatusOK, "requestCertificate.html", gin.H{
			"sacrTypes": sacramentTypes,
			"certTypes": certificatePurposes,
			"orgs":      orgs,
		})
	}
}

func SaveRequestCertificate(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sacrTypeInt, _ := utils.ParseInt(c.PostForm("sacr_types"))
		certPurposeInt, _ := utils.ParseInt(c.PostForm("cert_purpose"))
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		mail := c.PostForm("mail")
		request := CertificateRequest{
			SacrType:    uint8(sacrTypeInt),
			Purpose:     uint8(certPurposeInt),
			External:    len(c.PostForm("ext_req")) > 0,
			Names:       c.PostForm("names"),
			Surnames:    c.PostForm("surnames"),
			RUT:         c.PostForm("rut"),
			DateBirth:   utils.ParseDateForm(c.PostForm("date_birth")),
			Mail:        mail,
			Phone:       c.PostForm("phone"),
			OrgUnknown:  len(c.PostForm("org_unknown")) > 0,
			OrgID:       uint(orgIDInt),
			OrgDetails:  c.PostForm("org_details"),
			DateUnknown: len(c.PostForm("date_unknown")) > 0,
			DateEvent:   utils.ParseDateForm(c.PostForm("date_event")),
			DateFrom:    utils.ParseYearForm(c.PostForm("date_from")),
			DateTo:      utils.ParseYearForm(c.PostForm("date_to")),
		}
		createCertificateRequest(db, request)
		reqs := getCertificateRequestsFromMail(db, mail)
		c.HTML(http.StatusOK, "requestsUser.html", gin.H{
			"reqs": reqs,
		})
	}
}

func SearchRequest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqs []CertificateRequest
		if role, ok := c.Get("Role"); ok && role != nil && role.(uint8) == accRoleUsuario {
			if name, ok := c.Get("Name"); ok && name != nil {
				reqs = getCertificateRequestsFromMail(db, name.(string))
			}
		}
		c.HTML(http.StatusOK, "requestsUser.html", gin.H{
			"reqs": reqs,
		})
	}
}

func GetRequestsByMail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mail := c.PostForm("mail")
		reqs := getCertificateRequestsFromMail(db, mail)
		c.HTML(http.StatusOK, "requestTable", gin.H{
			"reqs": reqs,
		})
	}
}

func Requests(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		var reqs []CertificateRequest
		if orgID != 0 {
			reqs = getCertificateRequestsFromOrg(db, orgID)
		}
		c.HTML(http.StatusOK, "requestsOrg.html", gin.H{
			"reqs": reqs,
		})
	}
}
