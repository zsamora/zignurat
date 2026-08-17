package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

var (
	memberRoles = map[uint8]string{
		0: "Administrador",
		1: "Validador",
		2: "Miembro",
	}
)

// AccRole values for module 2 (Notariat), mirroring identirat's notariatRoles map.
const (
	accRoleAdministrador uint8 = 1
	accRoleOrganizacion  uint8 = 2
	accRoleValidador     uint8 = 3
)

func AdminRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	admin.Use(utils.RequireRoleMiddleware(2, accRoleAdministrador))
	{
		admin.GET("/organizations", Organizations(db))
		admin.GET("/users", Users(db))
		admin.GET("/createOrg", CreateOrgForm())
		admin.POST("/createOrg", CreateOrg(db))
		admin.GET("/createUser", CreateUserForm())
		admin.POST("/createUser", CreateUser(db))
	}
}
func MemberRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	admin.Use(utils.RequireRoleMiddleware(2, accRoleAdministrador))
	{
		admin.GET("/members", Members(db))
		admin.GET("/createMember", CreateMemberForm(db))
		admin.POST("/createMember", CreateMember(db))
	}
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/members", OrgMembers(db))
		org.GET("/createMember", CreateMemberForm(db))
		org.POST("/createMember", CreateMember(db))
	}
}
func OrgRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/books", Books(db))
		org.GET("/certificates", Certificates(db))
	}
}
func BookRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.POST("/createBook", CreateBookForm(db))
		org.POST("/addBook", AddBook(db))
		org.GET("/getOrgBooks", Books(db))
		org.POST("/getOrgBooks", Books(db))
	}
}
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/getBookRegisters", GetBookRegisters(db))
		org.POST("/getBookRegisters", GetBookRegisters(db))
		org.POST("/createRegisterBaptism", CreateRegisterBaptismForm(db))
		org.POST("/addRegisterBaptism", AddRegisterBaptism(db))
		org.GET("/getRegisterBaptism", GetRegisterBaptism(db))
		org.POST("/getRegisterBaptism", GetRegisterBaptism(db))
	}
}
func CertificateRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRoleMiddleware(2, accRoleOrganizacion, accRoleValidador))
	{
		org.POST("/getCertificatesFromReg", GetCertificatesFromReg(db))
		org.POST("/createCertificate", CreateCertificateForm(db))
		org.POST("/addCertificateBaptism", AddCertificateBaptism(db))
		org.POST("/downloadPDFCertificateBaptism", DownloadPDFCertificateBaptism(db))
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

// currentOwnerOrgID resolves the logged-in account's organization, whether it
// owns an Organization directly (Organización) or a Member that belongs to
// one (Validador, Administrador member role, etc).
func currentOwnerOrgID(db *gorm.DB, c *gin.Context) uint {
	ownerUUID := currentOwnerUUID(c)
	if ownerUUID == nil {
		return 0
	}
	if org := GetOrganizationFromUUID(db, *ownerUUID); org.ID != 0 {
		return org.ID
	}
	return GetMemberFromUUID(db, *ownerUUID).OrgID
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
func CreateOrgForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "orgCreate.html", gin.H{
			"orgTypes": orgTypes,
		})
	}
}
func CreateOrg(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		typeIdInt, _ := utils.ParseInt(c.PostForm("type_id"))
		org := Organization{
			OrgType:  uint8(typeIdInt),
			Name:     c.PostForm("name"),
			Diocese:  c.PostForm("diocese"),
			Commune:  c.PostForm("commune"),
			Address:  c.PostForm("address"),
			SealPath: c.PostForm("seal_path"),
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
func Members(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		members := GetMembers(db)
		c.HTML(http.StatusOK, "members.html", gin.H{
			"members":     members,
			"memberRoles": memberRoles,
			"createPath":  "/admin/createMember",
		})
	}
}
func CreateMemberForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := GetOrganizations(db)
		c.HTML(http.StatusOK, "memberCreate.html", gin.H{
			"orgs":        orgs,
			"memberRoles": memberRoles,
			"submitPath":  c.FullPath(),
		})
	}
}
func CreateMember(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIdInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgId := uint(orgIdInt)
		ownerOrgId := currentOwnerOrgID(db, c)
		if ownerOrgId != 0 {
			orgId = ownerOrgId
		}
		roleInt, _ := utils.ParseInt(c.PostForm("role"))
		member := Member{
			Names:         c.PostForm("names"),
			Surnames:      c.PostForm("surnames"),
			Role:          uint8(roleInt),
			OrgID:         orgId,
			SignaturePath: c.PostForm("signature_path"),
		}
		createMember(db, member)
		if ownerOrgId != 0 {
			members := GetMembersFromOrg(db, ownerOrgId)
			c.HTML(http.StatusOK, "members.html", gin.H{
				"members":     members,
				"memberRoles": memberRoles,
				"createPath":  "/org/createMember",
			})
			return
		}
		members := GetMembers(db)
		c.HTML(http.StatusOK, "members.html", gin.H{
			"members":     members,
			"memberRoles": memberRoles,
			"createPath":  "/admin/createMember",
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
func OrgMembers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := currentOwnerOrgID(db, c)
		members := GetMembersFromOrg(db, orgId)
		c.HTML(http.StatusOK, "members.html", gin.H{
			"members":     members,
			"memberRoles": memberRoles,
			"createPath":  "/org/createMember",
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
		bookBap := getBookFromId(db, reg.BookID)
		validators := GetMembersFromOrgAndRole(db, org.ID, 1)
		c.HTML(http.StatusOK, "certificateCreate.html", gin.H{
			"RegBap":     reg,
			"Org":        org,
			"OrgBap":     orgBap,
			"BookBap":    bookBap,
			"validators": validators,
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
		validator := GetMemberFromId(db, userValidator)
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
			RegValidator:       validator.Names + " " + validator.Surnames,
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
