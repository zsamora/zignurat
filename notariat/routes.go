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
	accRoleUsuario       uint8 = 4
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
		org.GET("/search", IndexSearch(db))
	}
}
func IndexSearch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		booksIndex := getBooksIndexBaptismFromOrg(db, orgID)
		c.HTML(http.StatusOK, "indexSearch.html", gin.H{
			"OrgID":      orgID,
			"booksIndex": booksIndex,
		})
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

// currentValidatorID resolves the logged-in account's own Member ID, used to
// preselect the validator dropdown when a Validador is creating a certificate.
func currentValidatorID(db *gorm.DB, c *gin.Context) uint {
	ownerUUID := currentOwnerUUID(c)
	if ownerUUID == nil {
		return 0
	}
	return GetMemberFromUUID(db, *ownerUUID).ID
}
func parseOrgIDParam(db *gorm.DB, c *gin.Context) uint {
	if idStr := c.PostForm("org_id"); idStr != "" {
		orgIDInt, _ := utils.ParseInt(idStr)
		return uint(orgIDInt)
	}
	if idStr := c.Query("org_id"); idStr != "" {
		orgIDInt, _ := utils.ParseInt(idStr)
		return uint(orgIDInt)
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
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgID := uint(orgIDInt)
		ownerOrgID := currentOwnerOrgID(db, c)
		if ownerOrgID != 0 {
			orgID = ownerOrgID
		}
		roleInt, _ := utils.ParseInt(c.PostForm("role"))
		member := Member{
			Names:         c.PostForm("names"),
			Surnames:      c.PostForm("surnames"),
			Role:          uint8(roleInt),
			OrgID:         orgID,
			SignaturePath: c.PostForm("signature_path"),
		}
		createMember(db, member)
		if ownerOrgID != 0 {
			members := GetMembersFromOrg(db, ownerOrgID)
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
		orgID := parseOrgIDParam(db, c)
		books := getBooksFromOrg(db, orgID)
		c.HTML(http.StatusOK, "books.html", gin.H{
			"OrgID": orgID,
			"books": books,
		})
	}
}
func CreateBookForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := parseOrgIDParam(db, c)
		data := gin.H{"orgs": GetOrganizations(db)}
		if orgID != 0 {
			data["Org"] = GetOrganizationFromID(db, orgID)
		}
		c.HTML(http.StatusOK, "createBook.html", data)
	}
}
func AddBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgID := uint(orgIDInt)
		bookNrInt, _ := utils.ParseInt(c.PostForm("book_nr"))
		book := Book{
			OrgID:       orgID,
			BookType:    2,
			BookNr:      uint8(bookNrInt),
			DateInitial: utils.ParseDateForm(c.PostForm("date_initial")),
			DateFinal:   utils.ParseDateForm(c.PostForm("date_final")),
		}
		createBook(db, book)
		books := getBooksFromOrg(db, orgID)
		c.HTML(http.StatusOK, "books.html", gin.H{
			"OrgID": orgID,
			"books": books,
		})
	}
}
func Certificates(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		certsBap := getCertificatesBaptismFromOrg(db, orgID)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
		})
	}
}
func OrgMembers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		members := GetMembersFromOrg(db, orgID)
		c.HTML(http.StatusOK, "members.html", gin.H{
			"members":     members,
			"memberRoles": memberRoles,
			"createPath":  "/org/createMember",
		})
	}
}
func GetBookRegisters(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookIDStr := c.PostForm("book_id")
		if bookIDStr == "" {
			bookIDStr = c.Query("book_id")
		}
		bookIDInt, _ := utils.ParseInt(bookIDStr)
		bookID := uint(bookIDInt)
		book := getBookFromID(db, bookID)
		regsBap := getRegistersBaptismFromBook(db, bookID)
		c.HTML(http.StatusOK, "registerBaptismAll.html", gin.H{
			"OrgID":    book.OrgID,
			"Book":     book,
			"regs_bap": regsBap,
		})
	}
}
func CreateRegisterBaptismForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgID := uint(orgIDInt)
		bookIDInt, _ := utils.ParseInt(c.PostForm("book_id"))
		bookID := uint(bookIDInt)
		book := getBookFromID(db, bookID)
		indexID := getLastIndexBaptismID(db) + 1
		data := gin.H{
			"Book":    book,
			"IndexID": indexID,
			"books":   getBooksFromOrg(db, orgID),
			"orgs":    GetOrganizations(db),
		}
		if orgID != 0 {
			data["Org"] = GetOrganizationFromID(db, orgID)
		}
		c.HTML(http.StatusOK, "registerBaptismCreate.html", data)
	}
}
func AddRegisterBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookIDInt, _ := utils.ParseInt(c.PostForm("book_id"))
		bookID := uint(bookIDInt)
		indexIDInt, _ := utils.ParseInt(c.PostForm("index_id"))
		indexID := uint(indexIDInt)
		pageNrInt, _ := utils.ParseInt(c.PostForm("page_nr"))
		regNrInt, _ := utils.ParseInt(c.PostForm("reg_nr"))
		orgBaptismInt, _ := utils.ParseInt(c.PostForm("org_baptism"))
		reg := RegisterBaptism{
			BookID:        bookID,
			PageNumber:    uint16(pageNrInt),
			RegNumber:     uint16(regNrInt),
			IndexID:       indexID,
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
		regID := createRegisterBaptism(db, reg)
		regNew := getRegisterBaptismFromID(db, regID)
		updateOrInsertIndexBaptism(db, indexID, regNew)
		book := getBookFromID(db, regNew.BookID)
		org := GetOrganizationFromID(db, regNew.OrgBaptism)
		c.HTML(http.StatusOK, "registerBaptismSingle.html", gin.H{
			"RegBap":  regNew,
			"BookBap": book,
			"OrgBap":  org,
		})
	}
}
func GetRegisterBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIDStr := c.PostForm("reg_id")
		if regIDStr == "" {
			regIDStr = c.Query("reg_id")
		}
		regIDInt, _ := utils.ParseInt(regIDStr)
		regID := uint(regIDInt)
		reg := getRegisterBaptismFromID(db, regID)
		book := getBookFromID(db, reg.BookID)
		org := GetOrganizationFromID(db, reg.OrgBaptism)
		c.HTML(http.StatusOK, "registerBaptismSingle.html", gin.H{
			"RegBap":  reg,
			"BookBap": book,
			"OrgBap":  org,
		})
	}
}
func GetCertificatesFromReg(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regID := uint(regIDInt)
		reg := getRegisterBaptismFromID(db, regID)
		orgID := currentOwnerOrgID(db, c)
		if orgID == 0 {
			orgID = reg.OrgBaptism
		}
		certsBap := getCertificatesBaptismFromOrgAndReg(db, orgID, regID)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
			"RegID":     regID,
		})
	}
}
func CreateCertificateForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regID := uint(regIDInt)
		reg := getRegisterBaptismFromID(db, regID)
		orgID := currentOwnerOrgID(db, c)
		if orgID == 0 {
			orgID = reg.OrgBaptism
		}
		org := GetOrganizationFromID(db, orgID)
		orgBap := GetOrganizationFromID(db, reg.OrgBaptism)
		bookBap := getBookFromID(db, reg.BookID)
		validators := GetMembersFromOrgAndRole(db, org.ID, 1)
		c.HTML(http.StatusOK, "certificateCreate.html", gin.H{
			"RegBap":      reg,
			"Org":         org,
			"OrgBap":      orgBap,
			"BookBap":     bookBap,
			"validators":  validators,
			"ValidatorID": currentValidatorID(db, c),
		})
	}
}
func AddCertificateBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIDInt, _ := utils.ParseInt(c.PostForm("org_id"))
		orgID := uint(orgIDInt)
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		userValInt, _ := utils.ParseInt(c.PostForm("user_val"))
		cert := CertificateBaptism{
			OrgEmisor:      orgID,
			RegID:          uint(regIDInt),
			UserValidator:  uint(userValInt),
			DateEmission:   utils.ParseDateForm(c.PostForm("date_emission")),
			DateExpiration: utils.ParseDateForm(c.PostForm("date_expiration")),
		}
		createCertificateBaptism(db, cert)
		certsBap := getCertificatesBaptismFromOrg(db, orgID)
		c.HTML(http.StatusOK, "certificatesOrg.html", gin.H{
			"certs_bap": certsBap,
			"RegID":     uint(regIDInt),
		})
	}
}
func DownloadPDFCertificateBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		certID := c.PostForm("cert_id")
		certUUID := c.PostForm("cert_uuid")
		dateEmission := c.PostForm("date_emi")
		dateExpiration := c.PostForm("date_exp")
		orgEmisorInt, _ := utils.ParseInt(c.PostForm("org_emisor"))
		orgID := uint(orgEmisorInt)
		org := GetOrganizationFromID(db, orgID)
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regID := uint(regIDInt)
		reg := getRegisterBaptismFromID(db, regID)
		regBook := getBookFromID(db, reg.BookID)
		regOrgBaptism := GetOrganizationFromID(db, reg.OrgBaptism)
		regOrgBaptismTotal := ParseOrgType(regOrgBaptism.OrgType) + " " + regOrgBaptism.Name
		regBaptNameS := ""
		if reg.BaptizedNameS != "" {
			regBaptNameS = " " + reg.BaptizedNameS
		}
		userValInt, _ := utils.ParseInt(c.PostForm("user_val"))
		userValidator := uint(userValInt)
		validator := GetMemberFromID(db, userValidator)
		certPDF := CertificateBaptismPDF{
			CertID:             certID,
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
