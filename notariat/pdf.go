package main

import (
	"log"

	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

func createPDFCertificateBaptism(certPDF CertificateBaptismPDF) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	marginLeftDefault := 20.0
	marginLeftDescription := 30.0
	marginLeftContent := 50.0
	pdf.AddPage()
	if err := pdf.AddTTFFontData("GoBold", gobold.TTF); err != nil {
		return nil, err
	}
	if err := pdf.AddTTFFontData("GoRegular", goregular.TTF); err != nil {
		return nil, err
	}
	// Header: Verification and QR Code
	pdf.Br(5)
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "GoBold", 10)
	pdf.Cell(nil, "N° de Certificado: ")
	pdf.SetX(108)
	setFontWithSize(&pdf, "GoRegular", 11)
	pdf.Cell(nil, certPDF.CertID)
	pdf.SetX(385)
	setFontWithSize(&pdf, "GoRegular", 8)
	pdf.Cell(nil, certPDF.CertUUID)
	pdf.Br(15)
	pdf.SetLineWidth(0.1)
	pdf.Line(0, 40, 595, 40)
	pdf.Br(20)
	// Organization and Validity
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "GoBold", 12)
	pdf.Cell(nil, "Periodo de validez: ")
	pdf.SetX(395)
	pdf.Cell(nil, certPDF.OrgEmiType+" "+certPDF.OrgEmiName)
	pdf.Br(20)
	// Location and Emission and Expiration Date
	pdf.SetX(marginLeftDefault)
	setFontWithSize(&pdf, "GoRegular", 12)
	pdf.Cell(nil, certPDF.CertDateEmission+" - "+certPDF.CertDateExpiration)
	pdf.SetX(440)
	setFontWithSize(&pdf, "GoBold", 12)
	pdf.Cell(nil, "Diócesis: ")
	setFontWithSize(&pdf, "GoRegular", 12)
	pdf.Cell(nil, certPDF.OrgEmiDiocese+", Chile")
	pdf.Br(60)
	// Title
	setFontWithSize(&pdf, "GoBold", 22)
	pdf.SetX(120)
	pdf.Cell(nil, "Certificado Partida de Bautismo")
	pdf.Br(120)
	// Description
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "Certifico que en el libro ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, "N° "+certPDF.RegBookNumber)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, ", página ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegBookPage)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " se encuentra la siguiente partida:")
	pdf.Br(40)
	// Baptism Info
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "En la ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegOrgBaptism)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " a ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegDateBaptism)
	pdf.Br(20)
	// Baptized Info
	pdf.SetX(marginLeftContent)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, "se bautizó a ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserBaptized)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, ",")
	pdf.Br(20)
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "RUT ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserRUT)
	// Birth Date and Place
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, ", nacido/a el ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserBirthDate)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " en ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserBirthPlace)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, ",")
	pdf.Br(20)
	// Parents
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "hijo/a de ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserFather)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " y de ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserMother)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, ".")
	pdf.Br(20)
	// Godparents
	pdf.SetX(marginLeftContent)
	pdf.Cell(nil, "Padrinos ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserGodfather)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " y ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegUserGodmother)
	pdf.Br(40)
	// Validator
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, "Doy fe, ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.RegValidator)
	pdf.Br(30)
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, "Notas: ")
	pdf.Rectangle(marginLeftContent, 470, 500, 530, "", 0, 0)
	pdf.Br(110)
	// Signature
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "En constancia sello y firmo en ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.OrgEmiDiocese)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, " el día ")
	setFontWithSize(&pdf, "GoBold", 14)
	pdf.Cell(nil, certPDF.CertDateEmission)
	pdf.Br(40)
	pdf.SetX(marginLeftDescription)
	setFontWithSize(&pdf, "GoRegular", 14)
	pdf.Cell(nil, "Firma párroco: ")
	if err := pdf.Image("img/firmaparroco.jpg", 250, 590, &gopdf.Rect{W: 100, H: 100}); err != nil {
		return nil, err
	}
	pdf.Line(150, 660, 430, 660)
	pdf.Br(100)
	pdf.SetX(marginLeftDescription)
	pdf.Cell(nil, "Sello parroquial: ")
	if err := pdf.Image("img/timbreparroquial.jpg", 160, 690, &gopdf.Rect{W: 80, H: 80}); err != nil {
		return nil, err
	}
	pdf.Br(125)
	// Contact Info
	pdf.SetLineWidth(0.1)
	pdf.Line(0, 800, 595, 800)
	pdf.SetX(marginLeftDefault + 90)
	setFontWithSize(&pdf, "GoBold", 12)
	pdf.Cell(nil, "Información de contacto: ")
	setFontWithSize(&pdf, "GoRegular", 12)
	pdf.Cell(nil, certPDF.OrgEmiAddress+", "+certPDF.OrgEmiCommune+", "+certPDF.OrgEmiDiocese)
	return pdf.GetBytesPdfReturnErr()
}

func setFontWithSize(pdf *gopdf.GoPdf, font string, size int) {
	err := pdf.SetFont(font, "", size)
	if err != nil {
		log.Print(err.Error())
	}
}
