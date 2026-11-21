package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"

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
	admin.Use(utils.RequireRole(2, accRoleAdministrador))
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
	admin.Use(utils.RequireRole(2, accRoleAdministrador))
	{
		admin.GET("/members", Members(db))
		admin.GET("/createMember", CreateMemberForm(db))
		admin.POST("/createMember", CreateMember(db))
	}
	org := router.Group("/org")
	org.Use(utils.RequireRole(2, accRoleOrganizacion, accRoleValidador))
	{
		org.GET("/members", OrgMembers(db))
		org.GET("/createMember", CreateMemberForm(db))
		org.POST("/createMember", CreateMember(db))
	}
}
func OrgRoutes(router *gin.Engine, db *gorm.DB) {
	org := router.Group("/org")
	org.Use(utils.RequireRole(2, accRoleOrganizacion, accRoleValidador))
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
	read := router.Group("/org")
	read.Use(utils.RequireRole(2, accRoleAdministrador, accRoleOrganizacion, accRoleValidador))
	{
		read.GET("/getOrgBooks", Books(db))
		read.POST("/getOrgBooks", Books(db))
	}
	write := router.Group("/org")
	write.Use(utils.RequireRole(2, accRoleOrganizacion, accRoleValidador))
	{
		write.POST("/createBook", CreateBookForm(db))
		write.POST("/addBook", AddBook(db))
	}
}
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	read := router.Group("/org")
	read.Use(utils.RequireRole(2, accRoleAdministrador, accRoleOrganizacion, accRoleValidador))
	{
		read.GET("/getBookRegisters", GetBookRegisters(db))
		read.POST("/getBookRegisters", GetBookRegisters(db))
		read.GET("/getRegisterBaptism", GetRegisterBaptism(db))
		read.POST("/getRegisterBaptism", GetRegisterBaptism(db))
	}
	write := router.Group("/org")
	write.Use(utils.RequireRole(2, accRoleOrganizacion, accRoleValidador))
	{
		write.POST("/createRegisterBaptism", CreateRegisterBaptismForm(db))
		write.POST("/addRegisterBaptism", AddRegisterBaptism(db))
	}
}
func CertificateRoutes(router *gin.Engine, db *gorm.DB) {
	read := router.Group("/org")
	read.Use(utils.RequireRole(2, accRoleAdministrador, accRoleOrganizacion, accRoleValidador))
	{
		read.POST("/getCertificatesFromReg", GetCertificatesFromReg(db))
		read.POST("/downloadPDFCertificateBaptism", DownloadPDFCertificateBaptism(db))
	}
	issue := router.Group("/org")
	issue.Use(utils.RequireRole(2, accRoleValidador))
	{
		issue.POST("/createCertificate", CreateCertificateForm(db))
		issue.POST("/addCertificateBaptism", AddCertificateBaptism(db))
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
	if role, _ := c.Get("Role"); role == accRoleAdministrador {
		if idStr := c.PostForm("org_id"); idStr != "" {
			orgIDInt, _ := utils.ParseInt(idStr)
			return uint(orgIDInt)
		}
		if idStr := c.Query("org_id"); idStr != "" {
			orgIDInt, _ := utils.ParseInt(idStr)
			return uint(orgIDInt)
		}
	}
	return currentOwnerOrgID(db, c)
}
func hasRole(c *gin.Context, roles ...uint8) bool {
	role, _ := c.Get("Role")
	r, ok := role.(uint8)
	return ok && slices.Contains(roles, r)
}
func canAccessOrg(db *gorm.DB, c *gin.Context, orgID uint) bool {
	if orgID == 0 {
		return false
	}
	return hasRole(c, accRoleAdministrador) || orgID == currentOwnerOrgID(db, c)
}
func notFound(c *gin.Context) {
	c.String(http.StatusNotFound, "No encontrado")
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
			"OrgID":   orgID,
			"books":   books,
			"CanEdit": hasRole(c, accRoleOrganizacion, accRoleValidador),
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
		orgID := currentOwnerOrgID(db, c)
		if orgID == 0 {
			notFound(c)
			return
		}
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
			"OrgID":   orgID,
			"books":   books,
			"CanEdit": hasRole(c, accRoleOrganizacion, accRoleValidador),
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
		if !canAccessOrg(db, c, book.OrgID) {
			notFound(c)
			return
		}
		regsBap := getRegistersBaptismFromBook(db, bookID)
		c.HTML(http.StatusOK, "registerBaptismAll.html", gin.H{
			"OrgID":    book.OrgID,
			"Book":     book,
			"regs_bap": regsBap,
			"CanEdit":  hasRole(c, accRoleOrganizacion, accRoleValidador),
		})
	}
}
func CreateRegisterBaptismForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		bookIDInt, _ := utils.ParseInt(c.PostForm("book_id"))
		bookID := uint(bookIDInt)
		book := getBookFromID(db, bookID)
		if !canAccessOrg(db, c, book.OrgID) {
			notFound(c)
			return
		}
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
		if !canAccessOrg(db, c, getBookFromID(db, bookID).OrgID) {
			notFound(c)
			return
		}
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
			"RegBap":   regNew,
			"BookBap":  book,
			"OrgBap":   org,
			"CanIssue": hasRole(c, accRoleValidador),
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
		if !canAccessOrg(db, c, book.OrgID) {
			notFound(c)
			return
		}
		org := GetOrganizationFromID(db, reg.OrgBaptism)
		c.HTML(http.StatusOK, "registerBaptismSingle.html", gin.H{
			"RegBap":   reg,
			"BookBap":  book,
			"OrgBap":   org,
			"CanIssue": hasRole(c, accRoleValidador),
		})
	}
}
func GetCertificatesFromReg(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		regID := uint(regIDInt)
		reg := getRegisterBaptismFromID(db, regID)
		book := getBookFromID(db, reg.BookID)
		if !canAccessOrg(db, c, book.OrgID) {
			notFound(c)
			return
		}
		certsBap := getCertificatesBaptismFromOrgAndReg(db, book.OrgID, regID)
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
		bookBap := getBookFromID(db, reg.BookID)
		if !canAccessOrg(db, c, bookBap.OrgID) {
			notFound(c)
			return
		}
		org := GetOrganizationFromID(db, bookBap.OrgID)
		orgBap := GetOrganizationFromID(db, reg.OrgBaptism)
		validator := GetMemberFromID(db, currentValidatorID(db, c))
		c.HTML(http.StatusOK, "certificateCreate.html", gin.H{
			"RegBap":    reg,
			"Org":       org,
			"OrgBap":    orgBap,
			"BookBap":   bookBap,
			"Validator": validator,
		})
	}
}
func AddCertificateBaptism(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := currentOwnerOrgID(db, c)
		regIDInt, _ := utils.ParseInt(c.PostForm("reg_id"))
		reg := getRegisterBaptismFromID(db, uint(regIDInt))
		validatorID := currentValidatorID(db, c)
		if validatorID == 0 || !canAccessOrg(db, c, getBookFromID(db, reg.BookID).OrgID) {
			notFound(c)
			return
		}
		cert := CertificateBaptism{
			OrgEmisor:      orgID,
			RegID:          reg.ID,
			UserValidator:  validatorID,
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
		certIDInt, _ := utils.ParseInt(c.PostForm("cert_id"))
		cert := getCertificateBaptismFromID(db, uint(certIDInt))
		if !canAccessOrg(db, c, cert.OrgEmisor) {
			notFound(c)
			return
		}
		org := GetOrganizationFromID(db, cert.OrgEmisor)
		reg := getRegisterBaptismFromID(db, cert.RegID)
		regBook := getBookFromID(db, reg.BookID)
		regOrgBaptism := GetOrganizationFromID(db, reg.OrgBaptism)
		regOrgBaptismTotal := ParseOrgType(regOrgBaptism.OrgType) + " " + regOrgBaptism.Name
		regBaptNameS := ""
		if reg.BaptizedNameS != "" {
			regBaptNameS = " " + reg.BaptizedNameS
		}
		validator := GetMemberFromID(db, cert.UserValidator)
		certPDF := CertificateBaptismPDF{
			CertID:             utils.FormatUint(uint64(cert.ID)),
			CertUUID:           cert.UUID.String(),
			CertDateEmission:   utils.FormatDate(cert.DateEmission),
			CertDateExpiration: utils.FormatDate(cert.DateExpiration),
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
		pdfBytes, err := createPDFCertificateBaptism(certPDF)
		if err != nil {
			log.Printf("Error generating PDF for certificate %d: %s", cert.ID, err)
			c.String(http.StatusInternalServerError, "Error generando el PDF")
			return
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=certificado-%d.pdf", cert.ID))
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}
