package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InternalRoutes is Notariat API to expose its organizations, members and users to other modules.
func InternalRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/internal/organizations", GetOrganizationsJSON(db))
	router.GET("/internal/members", GetMembersJSON(db))
	router.GET("/internal/users", GetUsersJSON(db))
}

type ownerStruct struct {
	UUID        uuid.UUID `json:"uuid"`
	DisplayName string    `json:"display_name"`
}

func GetOrganizationsJSON(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := GetOrganizations(db)
		summaries := make([]ownerStruct, 0, len(orgs))
		for _, org := range orgs {
			summaries = append(summaries, ownerStruct{
				UUID:        org.UUID,
				DisplayName: ParseOrgType(org.OrgType) + " - " + org.Name,
			})
		}
		c.JSON(http.StatusOK, summaries)
	}
}
func GetMembersJSON(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		members := GetMembers(db)
		summaries := make([]ownerStruct, 0, len(members))
		for _, member := range members {
			summaries = append(summaries, ownerStruct{
				UUID:        member.UUID,
				DisplayName: member.Names + " " + member.Surnames + " (" + memberRoles[member.Role] + " - " + member.Organization.Name + ")",
			})
		}
		c.JSON(http.StatusOK, summaries)
	}
}
func GetUsersJSON(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := GetUsers(db)
		summaries := make([]ownerStruct, 0, len(users))
		for _, user := range users {
			summaries = append(summaries, ownerStruct{
				UUID:        user.UUID,
				DisplayName: user.Names + " " + user.Surnames,
			})
		}
		c.JSON(http.StatusOK, summaries)
	}
}
