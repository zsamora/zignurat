package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	orgNavbarItems = map[string]string{
		"/org/books":        "Registro",
		"/org/certificates": "Verificación",
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
func OrgRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.AuthMiddleware(2))
	{
		org.GET("/books", Books(db))
		org.GET("/certificates", Certificates(db))
	}
}
func BookRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/createBook", CreateBookForm(db))
	router.POST("/addBook", AddBook(db))
	router.GET("/getOrgBooks", Books(db))
	router.POST("/getOrgBooks", Books(db))
}
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/getBookRegisters", GetBookRegisters(db))
	router.POST("/getBookRegisters", GetBookRegisters(db))
	router.POST("/createRegisterBaptism", CreateRegisterBaptismForm(db))
	router.POST("/addRegisterBaptism", AddRegisterBaptism(db))
	router.GET("/getRegisterBaptism", GetRegisterBaptism(db))
	router.POST("/getRegisterBaptism", GetRegisterBaptism(db))
}
func CertificateRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/getCertificatesFromReg", GetCertificatesFromReg(db))
	router.POST("/createCertificate", CreateCertificateForm(db))
	router.POST("/addCertificateBaptism", AddCertificateBaptism(db))
	router.POST("/downloadPDFCertificateBaptism", DownloadPDFCertificateBaptism(db))
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

func currentOwnerUUID(c *gin.Context) *uuid.UUID {
	if v, exists := c.Get("OwnerUUID"); exists && v != nil {
		if ptr, ok := v.(*uuid.UUID); ok && ptr != nil {
			return ptr
		}
	}
	return nil
}
func currentOwnerOrgID(db *gorm.DB, c *gin.Context) uint {
	ownerUUID := currentOwnerUUID(c)
	if ownerUUID == nil {
		return 0
	}
	return GetOrganizationFromUUID(db, *ownerUUID).ID
}
func parseOrgIDParam(db *gorm.DB, c *gin.Context) uint {
	if idStr := c.PostForm("org_id"); idStr != "" {
		orgIdInt, _ := utils.ParseInt(idStr)
		return uint(orgIdInt)
	}
	if idStr := c.Query("org_id"); idStr != "" {
		orgIdInt, _ := utils.ParseInt(idStr)
		return uint(orgIdInt)
	}
	return currentOwnerOrgID(db, c)
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

func Books(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := parseOrgIDParam(db, c)
		books := getBooksFromOrg(db, orgId)
		c.HTML(http.StatusOK, "books.html", gin.H{
			"OrgID": orgId,
			"books": books,
		})
	}
}
func CreateBookForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := parseOrgIDParam(db, c)
		data := gin.H{"orgs": GetOrganizations(db)}
		if orgId != 0 {
			data["Org"] = GetOrganizationFromId(db, orgId)
		}
		c.HTML(http.StatusOK, "createBook.html", data)
	}
}
func AddBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIdInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgId := uint(orgIdInt)
		bookNrInt, _ := utils.ParseInt(c.PostForm("book_nr"))
		book := Book{
			OrgID:       orgId,
			BookType:    2,
			BookNr:      uint8(bookNrInt),
			DateInitial: utils.ParseDateForm(c.PostForm("date_initial")),
			DateFinal:   utils.ParseDateForm(c.PostForm("date_final")),
		}
		createBook(db, book)
		books := getBooksFromOrg(db, orgId)
		c.HTML(http.StatusOK, "books.html", gin.H{
			"OrgID": orgId,
			"books": books,
		})
	}
}
func Certificates(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := currentOwnerOrgID(db, c)
		certsBap := getCertificatesBaptismFromOrg(db, orgId)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
		})
	}
}

func GetBookRegisters(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookIdStr := c.PostForm("book_id")
		if bookIdStr == "" {
			bookIdStr = c.Query("book_id")
		}
		bookIdInt, _ := utils.ParseInt(bookIdStr)
		bookId := uint(bookIdInt)
		book := getBookFromId(db, bookId)
		regsBap := getRegistersBaptismFromBook(db, bookId)
		c.HTML(http.StatusOK, "registerBaptismAll.html", gin.H{
			"OrgID":    book.OrgID,
			"Book":     book,
			"regs_bap": regsBap,
		})
	}
}
func CreateRegisterBaptismForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIdInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgId := uint(orgIdInt)
		bookIdInt, _ := utils.ParseInt(c.PostForm("book_id"))
		bookId := uint(bookIdInt)
		book := getBookFromId(db, bookId)
		indexId := getLastIndexBaptismID(db) + 1
		data := gin.H{
			"Book":    book,
			"IndexID": indexId,
			"books":   getBooksFromOrg(db, orgId),
			"orgs":    GetOrganizations(db),
		}
		if orgId != 0 {
			data["Org"] = GetOrganizationFromId(db, orgId)
		}
		c.HTML(http.StatusOK, "registerBaptismCreate.html", data)
	}
}
func AddRegisterBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookIdInt, _ := utils.ParseInt(c.PostForm("book_id"))
		bookId := uint(bookIdInt)
		indexIdInt, _ := utils.ParseInt(c.PostForm("index_id"))
		indexId := uint(indexIdInt)
		pageNrInt, _ := utils.ParseInt(c.PostForm("page_nr"))
		regNrInt, _ := utils.ParseInt(c.PostForm("reg_nr"))
		orgBaptismInt, _ := utils.ParseInt(c.PostForm("org_baptism"))
		reg := RegisterBaptism{
			BookID:        bookId,
			PageNumber:    uint16(pageNrInt),
			RegNumber:     uint16(regNrInt),
			IndexID:       indexId,
			OrgBaptism:    uint(orgBaptismInt),
			Baptizer:      c.PostForm("baptizer"),
			DateBaptism:   utils.ParseDateForm(c.PostForm("date_baptism")),
			BaptizedNameF: c.PostForm("baptized_nf"),
			BaptizedNameS: c.PostForm("baptized_ns"),
			RUT:           c.PostForm("rut"),
			DateBirth:     utils.ParseDateForm(c.PostForm("date_birth")),
			PlaceBirth:    c.PostForm("place_birth"),
			FatherName:    c.PostForm("father_name"),
			FatherSurname: c.PostForm("father_surname"),
			MotherName:    c.PostForm("mother_name"),
			MotherSurname: c.PostForm("mother_surname"),
			Godfather:     c.PostForm("godfather"),
			Godmother:     c.PostForm("godmother"),
		}
		regId := createRegisterBaptism(db, reg)
		regNew := getRegisterBaptismFromId(db, regId)
		updateOrInsertIndexBaptism(db, indexId, regNew)
		book := getBookFromId(db, regNew.BookID)
		org := GetOrganizationFromId(db, regNew.OrgBaptism)
		c.HTML(http.StatusOK, "registerBaptismSingle.html", gin.H{
			"RegBap":  regNew,
			"BookBap": book,
			"OrgBap":  org,
		})
	}
}
func GetRegisterBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIdStr := c.PostForm("reg_id")
		if regIdStr == "" {
			regIdStr = c.Query("reg_id")
		}
		regIdInt, _ := utils.ParseInt(regIdStr)
		regId := uint(regIdInt)
		reg := getRegisterBaptismFromId(db, regId)
		book := getBookFromId(db, reg.BookID)
		org := GetOrganizationFromId(db, reg.OrgBaptism)
		c.HTML(http.StatusOK, "registerBaptismSingle.html", gin.H{
			"RegBap":  reg,
			"BookBap": book,
			"OrgBap":  org,
		})
	}
}

func GetCertificatesFromReg(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIdInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regId := uint(regIdInt)
		reg := getRegisterBaptismFromId(db, regId)
		orgId := currentOwnerOrgID(db, c)
		if orgId == 0 {
			orgId = reg.OrgBaptism
		}
		certsBap := getCertificatesBaptismFromOrgAndReg(db, orgId, regId)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
			"RegID":     regId,
		})
	}
}
func CreateCertificateForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIdInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regId := uint(regIdInt)
		reg := getRegisterBaptismFromId(db, regId)
		orgId := currentOwnerOrgID(db, c)
		if orgId == 0 {
			orgId = reg.OrgBaptism
		}
		org := GetOrganizationFromId(db, orgId)
		orgBap := GetOrganizationFromId(db, reg.OrgBaptism)
		adminBap := GetUserFromId(db, orgBap.AdminID)
		bookBap := getBookFromId(db, reg.BookID)
		c.HTML(http.StatusOK, "certificateCreate.html", gin.H{
			"RegBap":   reg,
			"Org":      org,
			"OrgBap":   orgBap,
			"BookBap":  bookBap,
			"AdminBap": adminBap,
		})
	}
}
func AddCertificateBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIdInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgId := uint(orgIdInt)
		regIdInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		userValInt, _ := utils.ParseInt(c.PostForm("user_val"))
		cert := CertificateBaptism{
			OrgEmisor:      orgId,
			RegID:          uint(regIdInt),
			UserValidator:  uint(userValInt),
			DateEmission:   utils.ParseDateForm(c.PostForm("date_emission")),
			DateExpiration: utils.ParseDateForm(c.PostForm("date_expiration")),
		}
		createCertificateBaptism(db, cert)
		certsBap := getCertificatesBaptismFromOrg(db, orgId)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
			"RegID":     uint(regIdInt),
		})
	}
}
func DownloadPDFCertificateBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		certId := c.PostForm("cert_id")
		certUUID := c.PostForm("cert_uuid")
		dateEmission := c.PostForm("date_emi")
		dateExpiration := c.PostForm("date_exp")
		orgEmisorInt, _ := utils.ParseInt(c.PostForm("org_emisor"))
		orgId := uint(orgEmisorInt)
		org := GetOrganizationFromId(db, orgId)
		regIdInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regId := uint(regIdInt)
		reg := getRegisterBaptismFromId(db, regId)
		regBook := getBookFromId(db, reg.BookID)
		regOrgBaptism := GetOrganizationFromId(db, reg.OrgBaptism)
		regOrgBaptismTotal := ParseOrgType(regOrgBaptism.OrgType) + " " + regOrgBaptism.Name
		regBaptNameS := ""
		if reg.BaptizedNameS != "" {
			regBaptNameS = " " + reg.BaptizedNameS
		}
		userValInt, _ := utils.ParseInt(c.PostForm("user_val"))
		userValidator := uint(userValInt)
		user := GetUserFromId(db, userValidator)
		certPDF := CertificateBaptismPDF{
			CertID:             certId,
			CertUUID:           certUUID,
			CertDateEmission:   dateEmission,
			CertDateExpiration: dateExpiration,
			OrgEmiDiocese:      org.Diocese,
			OrgEmiType:         ParseOrgType(org.OrgType),
			OrgEmiName:         org.Name,
			OrgEmiCommune:      org.Commune,
			OrgEmiAddress:      org.Address,
			RegBookNumber:      utils.FormatUint(uint64(regBook.BookNr)),
			RegBookPage:        utils.FormatUint(uint64(reg.PageNumber)),
			RegOrgBaptism:      regOrgBaptismTotal,
			RegDateBaptism:     utils.FormatDate(reg.DateBaptism),
			RegUserBaptized:    reg.BaptizedNameF + regBaptNameS + " " + reg.FatherSurname + " " + reg.MotherSurname,
			RegUserRUT:         reg.RUT,
			RegUserBirthDate:   utils.FormatDate(reg.DateBirth),
			RegUserBirthPlace:  reg.PlaceBirth,
			RegUserFather:      reg.FatherName + " " + reg.FatherSurname,
			RegUserMother:      reg.MotherName + " " + reg.MotherSurname,
			RegUserGodfather:   reg.Godfather,
			RegUserGodmother:   reg.Godmother,
			RegValidator:       user.Names + " " + user.Surnames,
		}
		pdfFilepath := createPDFCertificateBaptism(certPDF)
		if pdfFilepath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando el PDF"})
			return
		}
		if _, err := os.Stat(pdfFilepath); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(pdfFilepath))
		c.Header("Content-Type", "application/pdf")
		c.File(pdfFilepath)
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
			"navbarItems": orgNavbarItems,
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
