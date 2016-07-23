package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/satori/go.uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/docs", docs)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Service started.")
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, true)
}

// Docs holds URLs for the generated docs.
type Docs struct {
	Consultant string
	Patient    string
}

func docs(w http.ResponseWriter, r *http.Request) {
	// printRequestBody(r)

	//Call to ParseForm makes form fields available.
	err := r.ParseForm()
	if err != nil {

	}

	rq := r.PostFormValue("referralquestion")
	td := r.PostFormValue("testdetails")
	c, p, _ := buildPDF(rq, td)

	t, err := template.ParseFiles("/docs.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, Docs{
		Consultant: c,
		Patient:    p,
	})
}

func buildPDF(rq, td string) (consultant, patient string, err error) {
	font := "Helvetica"

	fmt.Println("Building PDF")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFillColor(230, 230, 230)
	pdf.SetCellMargin(10)
	pdf.AddPage()

	pdf.SetFont(font, "B", 24)
	pdf.Cell(40, 10, "Referral Runway")
	pdf.Ln(15)

	// Referral Question
	pdf.SetFont(font, "B", 16)
	pdf.Cell(40, 10, "Referral Question")
	pdf.Ln(6)
	pdf.SetFont(font, "I", 11)
	// TODO: Different description for patient?
	pdf.Cell(40, 10, "The following referral question has been agreed with the patient.")
	pdf.Ln(10)
	pdf.SetFont(font, "", 11)
	pdf.MultiCell(0, 5, rq, "1", "", true)
	pdf.Ln(15)

	// Test Details
	pdf.SetFont(font, "B", 16)
	pdf.Cell(40, 10, "Test Details")
	pdf.Ln(6)
	pdf.SetFont(font, "I", 11)
	pdf.Cell(40, 10, "The referring GP believes they have completed all reasonable tests in primary care. Details are below.")
	pdf.Ln(10)
	pdf.SetFont(font, "", 11)
	pdf.MultiCell(0, 5, td, "1", "", true)

	b := bytes.NewBuffer([]byte{}) // A Buffer needs no initialization.
	pdf.Output(b)

	auth, err := aws.EnvAuth()
	if err != nil {
		// fmt.Printf("No environment auth: %s", err)
		auth, err = aws.SharedAuth()
	}
	if err != nil {
		fmt.Printf("No auth: %s", err)
	}
	sss := s3.New(auth, aws.EUWest)
	bkt := sss.Bucket("referralrunway")

	// TODO: Random name.
	prefix := uuid.NewV4().String()
	expiry := time.Now().Add(30 * time.Minute)

	// Consultant
	ck := prefix + "-consultant.pdf"
	err = bkt.Put(ck, b.Bytes(), "application/pdf", s3.PublicRead)
	if err != nil {
		fmt.Printf("Failed to write %s to S3: %s", ck, err)
	}
	consultant = bkt.SignedURL(ck, expiry)

	// Patient
	pk := prefix + "-patient.pdf"
	err = bkt.Put(pk, b.Bytes(), "application/pdf", s3.PublicRead)
	if err != nil {
		fmt.Printf("Failed to write %s to S3: %s", pk, err)
	}
	patient = bkt.SignedURL(pk, expiry)

	err = nil
	return
}

func printRequestBody(r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	fmt.Println(buf.String())
}
