package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

var (
	adminNavbarItems = map[string]string{
		"/admin/organizations": "Organizaciones",
		"/admin/users":         "Usuarios",
	}
)

func LandingRoutes(router *gin.Engine) {
	router.GET("/", utils.AuthMiddleware(2), Home())
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "landing.html", gin.H{})
	})
}
func LoginRoutes(router *gin.Engine) {
	router.GET("/loginForm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	router.POST("/loginForm", LoginRequest())
	router.GET("/logout", Logout())
}
func AdminRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	admin.Use(utils.AuthMiddleware(2))
	{
		admin.GET("/organizations", Organizations(db))
		admin.GET("/users", Users(db))
	}
	router.GET("/createOrg", CreateOrgForm(db))
	router.POST("/createOrg", CreateOrg(db))
	router.GET("/createUser", CreateUserForm())
	router.POST("/createUser", CreateUser(db))
}
func InternalRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/internal/organizations", GetOrganizationsJSON(db))
}

type ownerSummary struct {
	UUID        uuid.UUID `json:"uuid"`
	DisplayName string    `json:"display_name"`
}

func GetOrganizationsJSON(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := GetOrganizations(db)
		summaries := make([]ownerSummary, 0, len(orgs))
		for _, org := range orgs {
			summaries = append(summaries, ownerSummary{
				UUID:        org.UUID,
				DisplayName: ParseOrgType(org.OrgType) + " - " + org.Name,
			})
		}
		c.JSON(http.StatusOK, summaries)
	}
}

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("Role")
		name := c.MustGet("Name")
		if role != nil && name != nil {
			getRoleHome(c, role.(uint8), name.(string))
		} else {
			c.HTML(http.StatusUnauthorized, "index.html", gin.H{
				"LoggedIn":    false,
				"navbarItems": nil,
				"Name":        nil,
			})
		}
	}
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func LoginRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form Login
		if err := c.ShouldBind(&form); err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"LoggedIn":    false,
				"navbarItems": nil,
				"Name":        nil,
			})
		} else {
			identiratUrl := utils.GetConfig("IDENTIRAT_URL")
			identiratPort := utils.GetConfig("IDENTIRAT_PORT")
			identiratLogin := utils.GetConfig("IDENTIRAT_LOGIN")
			requestURL := fmt.Sprintf("http://%s:%s/%s", identiratUrl, identiratPort, identiratLogin)
			requestJSON := []byte(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, form.User, form.Password))
			bodyReader := bytes.NewReader(requestJSON)
			req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
			if err != nil {
				log.Printf("# Notariat: Failed creating POST request: %s\n", err)
			} else {
				req.Header.Set("Content-Type", "application/json")
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("# Notariat: Failed executing POST request: %s\n", err)
				} else {
					defer res.Body.Close()
					log.Printf("client: status code: %d\n", res.StatusCode)
					if res.StatusCode == 200 {
						resBody, err := io.ReadAll(res.Body)
						if err != nil {
							log.Printf("# Notariat: Failed reading POST request response: %s\n", err)
						} else {
							var jsonResponse utils.JWTResponse
							json.Unmarshal(resBody, &jsonResponse)
							if jsonResponse.Token != "" {
								claims, err := utils.CheckJWTToken(jsonResponse.Token, c)
								if err == nil {
									log.Printf("- Logged in (ID: %d, Name: %s, Module: %d, Role: %d, Owner UUID: %v, Exp. Time: %s)",
										claims.ID, claims.Name, claims.Module, claims.AccRole, claims.OwnerUUID, claims.RegisteredClaims.ExpiresAt)
									utils.SetJWTTokenFromSession(c, &jsonResponse.Token, &jsonResponse.Refresh)
									role := claims.AccRole
									name := claims.Name
									getRoleHome(c, role, name)
									return
								}
							} else {
								log.Printf("!! No JWT Token received")
							}
						}
					} else {
						log.Printf("!! Error status code: %d", res.StatusCode)
					}
				}
			}
			utils.SetJWTTokenFromSession(c, nil, nil)
			c.HTML(http.StatusOK, "index.html", gin.H{
				"LoggedIn":    false,
				"navbarItems": nil,
				"Name":        nil,
			})
		}
	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Logged out")
		utils.SetJWTTokenFromSession(c, nil, nil)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"LoggedIn":    false,
			"navbarItems": nil,
			"Name":        nil,
		})
	}
}

func Users(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := GetUsers(db)
		c.HTML(http.StatusOK, "users.html", gin.H{
			"users": users,
		})
	}
}
func Organizations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := GetOrganizations(db)
		c.HTML(http.StatusOK, "organizations.html", gin.H{
			"orgs": orgs,
		})
	}
}
func CreateOrgForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := GetUsers(db)
		c.HTML(http.StatusOK, "orgCreate.html", gin.H{
			"orgTypes": orgTypes,
			"users":    users,
		})
	}
}
func CreateOrg(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		typeIdInt, _ := utils.ParseInt(c.PostForm("type_id"))
		userIdInt, _ := utils.ParseInt(c.PostForm("user_id"))
		org := Organization{
			OrgType: uint8(typeIdInt),
			Name:    c.PostForm("name"),
			Diocese: c.PostForm("diocese"),
			Commune: c.PostForm("commune"),
			Address: c.PostForm("address"),
			AdminID: uint(userIdInt),
		}
		createOrganization(db, org)
		orgs := GetOrganizations(db)
		c.HTML(http.StatusOK, "organizations.html", gin.H{
			"orgs": orgs,
		})
	}
}
func CreateUserForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "userCreate.html", gin.H{})
	}
}
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User{
			Names:     c.PostForm("names"),
			Surnames:  c.PostForm("surnames"),
			DateBirth: utils.ParseDateForm(c.PostForm("date_birth")),
		}
		createUser(db, user)
		users := GetUsers(db)
		c.HTML(http.StatusOK, "users.html", gin.H{
			"users": users,
		})
	}
}

func getRoleHome(c *gin.Context, role uint8, name string) {
	switch role {
	case 1:
		log.Printf("## Administrator")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"LoggedIn":    true,
			"navbarItems": adminNavbarItems,
			"Name":        name,
		})
	case 2:
		log.Printf("## Organización")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"LoggedIn":    true,
			"navbarItems": nil,
			"Name":        name,
		})
	default:
		log.Printf("!! No role defined for this value")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"LoggedIn":    false,
			"navbarItems": nil,
			"Name":        name,
		})
	}
}
