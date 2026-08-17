package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

var (
	adminNavbarItems = map[string]string{
		"/admin/accounts": "Cuentas",
	}
	modules = map[uint8]string{
		1: "Identirat",
		2: "Notariat",
	}
	identiratRoles = map[uint8]string{
		1: "Administrador",
	}
	notariatRoles = map[uint8]string{
		1: "Administrador",
		2: "Organización",
		3: "Validador",
		4: "Usuario",
	}
)

func LandingRoutes(router *gin.Engine) {
	router.GET("/", utils.AuthMiddleware(uint8(1)), Home())
	router.GET("/home", utils.AuthMiddleware(uint8(1)), func(c *gin.Context) {
		c.HTML(http.StatusOK, "help.html", gin.H{})
	})
}
func LoginRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/loginForm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	router.POST("/loginForm", LoginForm(db))
	router.POST("/loginJSON", LoginJSON(db))
	router.GET("/logout", Logout())
}
func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	admin.Use(utils.AuthMiddleware(uint8(1)))
	{
		admin.GET("/accounts", Accounts(db))
		admin.GET("/createAccount", CreateAccount())
		admin.GET("/getAccRoles", CreateAccountGetAccountRoles())
		admin.GET("/getOwners", CreateAccountGetOwners())
		admin.POST("/createAccount", SaveAccount(db))
		admin.POST("/setOwner", SetOwner(db))
		admin.POST("/changePassword", ChangePassword(db))
	}
}

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("Role")
		name := c.MustGet("Name")
		if role != nil && name != nil {
			getRoleView(c, role.(uint8), name.(string))
		} else {
			log.Printf("!! No role or name")
			c.HTML(http.StatusOK, "index.html", gin.H{
				"LoggedIn":    false,
				"navbarItems": nil,
				"Name":        name,
			})
		}
	}
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func LoginForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SetJWTTokenFromSession(c, nil, nil)
		var form Login
		if err := c.ShouldBind(&form); err == nil {
			username := form.User
			foundAccount := getAccountByUsername(db, username)
			if foundAccount != (Account{}) {
				log.Printf("Account found: ", foundAccount)
				password := form.Password
				accountPW := foundAccount.Password
				if VerifyPassword(accountPW, password) {
					log.Printf("Password match")
					name := foundAccount.Username
					accountID := foundAccount.ID
					module := foundAccount.Module
					accountRole := foundAccount.AccRole
					ownerUUID := foundAccount.OwnerUUID
					switch module {
					case 1:
						log.Printf("## Identirat Account")
						jwt_token, refresh_token := utils.GenerateJWTTokens(accountID, name, module, accountRole, ownerUUID)
						utils.SetJWTTokenFromSession(c, &jwt_token, &refresh_token)
						getRoleView(c, accountRole, name)
						return
					default:
						log.Printf("!! External module account, cannot login to Identirat")
					}
				} else {
					log.Printf("!! Password doesn't match")
				}
			} else {
				log.Printf("!! Account not found")
			}
		} else {
			log.Printf("!! Not enough arguments in form")
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"LoggedIn":    false,
			"navbarItems": nil,
			"Name":        nil,
		})
	}
}

func LoginJSON(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			log.Printf("Not enough arguments in form")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			username := json.User
			foundAccount := getAccountByUsername(db, username)
			if foundAccount == (Account{}) {
				log.Printf("Account not found")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Account or password doesn't match"})
			} else {
				log.Println("Account found:", foundAccount)
				password := json.Password
				accountPW := foundAccount.Password
				if VerifyPassword(accountPW, password) {
					log.Printf("Password match")
					accountID := foundAccount.ID
					module := foundAccount.Module
					accountRole := foundAccount.AccRole
					ownerUUID := foundAccount.OwnerUUID
					name := foundAccount.Username
					switch module {
					case 1:
						log.Printf("## Identirat Account")
						name = name + " (" + identiratRoles[accountRole] + ")"
						jwt_token, refresh_token := utils.GenerateJWTTokens(accountID, name, module, accountRole, ownerUUID)
						c.JSON(http.StatusOK, gin.H{"token": jwt_token, "refresh": refresh_token})
					case 2:
						log.Printf("## Notariat Account")
						name = name + " (" + notariatRoles[accountRole] + ")"
						jwt_token, refresh_token := utils.GenerateJWTTokens(accountID, name, module, accountRole, ownerUUID)
						c.JSON(http.StatusOK, gin.H{"token": jwt_token, "refresh": refresh_token})
					default:
						log.Printf("!! Account doesn't belong to modules available")
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Account invalid"})
					}
				} else {
					log.Printf("!! Password incorrect")
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Account or password doesn't match"})
				}
			}
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

func Accounts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accounts := buildAccountRows(getAccounts(db))
		c.HTML(http.StatusOK, "accounts.html", gin.H{
			"accounts": accounts,
		})
	}
}

func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "createAccount.html", gin.H{
			"modules": modules,
		})
	}
}

func CreateAccountGetAccountRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		module, _ := utils.ParseInt(c.Query("module"))
		switch module {
		case 1:
			c.HTML(http.StatusOK, "accountRoles.html", gin.H{
				"accountRoles": identiratRoles,
			})
		case 2:
			c.HTML(http.StatusOK, "accountRoles.html", gin.H{
				"accountRoles": notariatRoles,
			})
		default:
			log.Println("!! No module defined for this value")
			c.HTML(http.StatusOK, "accountRoles.html", gin.H{
				"accountRoles": nil,
			})
		}
	}
}

func CreateAccountGetOwners() gin.HandlerFunc {
	return func(c *gin.Context) {
		module, _ := utils.ParseInt(c.Query("module"))
		accountRole, _ := utils.ParseInt(c.Query("acc_role"))
		owners := ownersForAccount(uint8(module), uint8(accountRole))
		c.HTML(http.StatusOK, "accountOwners.html", gin.H{
			"Module":  module,
			"AccRole": accountRole,
			"owners":  owners,
		})
	}
}

func SaveAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("pw")
		hashedPassword := HashPassword(password)
		moduleInt, _ := utils.ParseInt(c.PostForm("module"))
		module := uint8(moduleInt)
		accountRoleInt, _ := utils.ParseInt(c.PostForm("acc_role"))
		accountRole := uint8(accountRoleInt)
		var ownerUUID *uuid.UUID
		if raw := c.PostForm("owner_uuid"); raw != "" {
			parsed := utils.ParseUUID(raw)
			ownerUUID = &parsed
		}
		account := Account{
			Username:  username,
			Password:  hashedPassword,
			Module:    module,
			AccRole:   accountRole,
			OwnerUUID: ownerUUID,
		}
		createAccount(db, account)
		accounts := buildAccountRows(getAccounts(db))
		c.HTML(http.StatusOK, "accounts.html", gin.H{
			"accounts": accounts,
		})
	}
}

func SetOwner(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accIdInt, _ := utils.ParseInt(c.PostForm("acc_id"))
		account := getAccountById(db, int64(accIdInt))
		ownerUUIDRaw := c.PostForm("owner_uuid")
		ownerUUID := utils.ParseUUID(ownerUUIDRaw)
		account.OwnerUUID = &ownerUUID
		updateAccount(db, account)
		accounts := buildAccountRows(getAccounts(db))
		c.HTML(http.StatusOK, "accounts.html", gin.H{
			"accounts": accounts,
		})
	}
}

func ChangePassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accIdInt, _ := utils.ParseInt(c.PostForm("acc_id"))
		account := getAccountById(db, int64(accIdInt))
		pw := c.PostForm("pw")
		if pw == "" {
			c.HTML(http.StatusOK, "changePassword.html", gin.H{
				"AccountID": account.ID,
				"Username":  account.Username,
			})
			return
		}
		account.Password = HashPassword(pw)
		updateAccount(db, account)
		accounts := buildAccountRows(getAccounts(db))
		c.HTML(http.StatusOK, "accounts.html", gin.H{
			"accounts": accounts,
		})
	}
}

func getRoleView(c *gin.Context, role uint8, name string) {
	switch role {
	case 1:
		log.Printf("## Administrator")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"LoggedIn":    true,
			"navbarItems": adminNavbarItems,
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
