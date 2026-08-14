package main

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/signintech/gopdf"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func createPDFCertificateBaptism(certPDF CertificateBaptismPDF) string {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	marginLeftDefault := 20.0
	marginLeftDescription := 30.0
	marginLeftContent := 50.0
	pdf.AddPage()
	err := pdf.AddTTFFont("TTNorms-Bold", "ttf/tt_norms_pro_serif/TT Norms Pro Serif Trial Bold.ttf")
	err = pdf.AddTTFFont("TTNorms-Regular", "ttf/tt_norms_pro_serif/TT Norms Pro Serif Trial Regular.ttf")
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	// Header: Verification and QR Code
	pdf.Br(5)
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "TTNorms-Bold", 10)
	pdf.Cell(nil, "N° de Certificado: ")
	pdf.SetX(108)
	setFontWithSize(&pdf, "TTNorms-Regular", 11)
	pdf.Cell(nil, certPDF.CertID)
	pdf.Image("img/qrcode.jpg", 550, 8, &gopdf.Rect{W: 25, H: 25})
	pdf.SetX(385)
	setFontWithSize(&pdf, "TTNorms-Regular", 8)
	pdf.Cell(nil, certPDF.CertUUID)
	pdf.Br(15)
	pdf.SetLineWidth(0.1)
	pdf.Line(0, 40, 595, 40)
	pdf.Br(20)
	// Organization and Validity
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "TTNorms-Bold", 12)
	pdf.Cell(nil, "Periodo de validez: ")
	pdf.SetX(395)
	pdf.Cell(nil, certPDF.OrgEmiType+" "+certPDF.OrgEmiName)
	pdf.Br(20)
	// Location and Emission and Expiration Date
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "TTNorms-Regular", 12)
	pdf.Cell(nil, certPDF.CertDateEmission+" - "+certPDF.CertDateExpiration)
	pdf.SetX(440)
	setFontWithSize(&pdf, "TTNorms-Bold", 12)
	pdf.Cell(nil, "Diócesis: ")
	setFontWithSize(&pdf, "TTNorms-Regular", 12)
	pdf.Cell(nil, certPDF.OrgEmiDiocese+", Chile")
	pdf.Br(60)
	// Title
	setFontWithSize(&pdf, "TTNorms-Bold", 22)
	pdf.SetX(120)
	pdf.Cell(nil, "Certificado Partida de Bautismo")
	pdf.Br(120)
	// Description
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "Certifico que en el libro ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, "N° "+certPDF.RegBookNumber)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, ", página ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegBookPage)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " se encuentra la siguiente partida:")
	pdf.Br(40)
	// Baptism Info
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "En la ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegOrgBaptism)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " a ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegDateBaptism)
	pdf.Br(20)
	// Baptized Info
	pdf.SetX(marginLeftContent)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, "se bautizó a ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserBaptized)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, ",")
	pdf.Br(20)
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "RUT ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserRUT)
	// Birth Date and Place
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, ", nacido/a el ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserBirthDate)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " en ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserBirthPlace)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, ",")
	pdf.Br(20)
	// Parents
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "hijo/a de ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserFather)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " y de ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserMother)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, ".")
	pdf.Br(20)
	// Godparents
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "Padrinos ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserGodfather)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " y ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegUserGodmother)
	pdf.Br(40)
	// Validator
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, "Doy fe, ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.RegValidator)
	pdf.Br(30)
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, "Notas: ")
	pdf.Rectangle(marginLeftContent, 470, 500, 530, "", 0, 0)
	pdf.Br(110)
	// Signature
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "En constancia sello y firmo en ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.OrgEmiDiocese)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, " el día ")
	setFontWithSize(&pdf, "TTNorms-Bold", 14)
	pdf.Cell(nil, certPDF.CertDateEmission)
	pdf.Br(40)
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "TTNorms-Regular", 14)
	pdf.Cell(nil, "Firma párroco: ")
	pdf.Image("img/firmaparroco.jpg", 250, 590, &gopdf.Rect{W: 100, H: 100})
	pdf.Line(150, 660, 430, 660)
	pdf.Br(100)
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "Sello parroquial: ")
	pdf.Image("img/timbreparroquial.jpg", 160, 690, &gopdf.Rect{W: 80, H: 80})
	pdf.Br(125)
	// Contact Info
	pdf.SetLineWidth(0.1)
	pdf.Line(0, 800, 595, 800)
	pdf.SetX(marginLeftDefault + 90)
	setFontWithSize(&pdf, "TTNorms-Bold", 12)
	pdf.Cell(nil, "Información de contacto: ")
	setFontWithSize(&pdf, "TTNorms-Regular", 12)
	pdf.Cell(nil, certPDF.OrgEmiAddress+", "+certPDF.OrgEmiCommune+", "+certPDF.OrgEmiDiocese)
	pdfFilename := removeAccents(strings.ToLower(strings.ReplaceAll(certPDF.RegUserBaptized, " ", ""))) + "_" + regExpDate(certPDF.CertDateEmission)
	log.Printf("PDF Filename: %s", pdfFilename)
	pdfRoute := "pdf/" + pdfFilename + ".pdf"
	if err := pdf.WritePdf(pdfRoute); err != nil {
		log.Printf("Error writing PDF: %s", err)
		return ""
	}
	return pdfRoute
}

func setFontWithSize(pdf *gopdf.GoPdf, font string, size int) {
	err := pdf.SetFont(font, "", size)
	if err != nil {
		log.Print(err.Error())
	}
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func regExpDate(s string) string {
	redash, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Printf("Error compiling regex: %s", err)
		return ""
	}
	return redash.ReplaceAllString(s, "-")
}
