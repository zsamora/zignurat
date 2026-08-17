package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zsamora/utils"
)

var (
	adminNavbarItems = map[string]string{
		"/admin/organizations": "Organizaciones",
		"/admin/users":         "Usuarios",
		"/admin/members":       "Miembros",
	}
	orgNavbarItems = map[string]string{
		"/org/books":        "Registro",
		"/org/certificates": "Verificación",
		"/org/members":      "Miembros",
	}
	accRoleNames = map[uint8]string{
		accRoleAdministrador: "Administrador",
		accRoleOrganizacion:  "Organización",
		accRoleValidador:     "Validador",
	}
)

func LandingRoutes(router *gin.Engine) {
	router.GET("/", utils.AuthMiddleware(2), Home())
	router.GET("/home", utils.AuthMiddleware(2), func(c *gin.Context) {
		c.HTML(http.StatusOK, "help.html", gin.H{})
	})
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
func getRoleHome(c *gin.Context, role uint8, name string) {
	switch role {
	case 1:
		log.Printf("## Administrator")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"LoggedIn":    true,
			"navbarItems": adminNavbarItems,
			"Name":        name,
			"Role":        accRoleNames[role],
		})
	case 2, 3:
		log.Printf("## Organización / Validador")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"LoggedIn":    true,
			"navbarItems": orgNavbarItems,
			"Name":        name,
			"Role":        accRoleNames[role],
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
