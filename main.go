package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bmizerany/pat"
	"gopkg.in/gomail.v2"
)

type Name struct {
	Name string
}

type Message struct {
	Name    string
	Email   string
	Content string
}

type Page struct {
	Key string
}

func main() {
	port := os.Getenv("PORT")

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(index))
	mux.Post("/", http.HandlerFunc(send))
	mux.Get("/confirmation", http.HandlerFunc(confirmation))
	mux.Get("/fail", http.HandlerFunc(fail))

	log.Println("Listening...")
	http.ListenAndServe(":"+port, mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := &Page{
		Key: os.Getenv("DATA_SITEKEY"),
	}
	render(w, "templates/main.html", data)
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/confirmation.html", nil)
}

func fail(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/fail.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	temp, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := temp.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	m := &Message{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Content: r.FormValue("content"),
	}

	if m.Name == "" || m.Email == "" || m.Content == "" {
		http.Redirect(w, r, "/fail", http.StatusSeeOther)
	} else {
		pwd := os.Getenv("MAIL_PASSWORD")
		email := os.Getenv("EMAIL")

		msg := gomail.NewMessage()
		msg.SetHeader("From", fmt.Sprintf("Jay <%s>", email))
		msg.SetHeader("To", fmt.Sprintf("Jay <%s>", email))
		msg.SetAddressHeader("reply-to", m.Email, "Contactee")
		msg.SetHeader("Subject", "Contact")
		msg.SetBody("text/html", fmt.Sprintf("<p>From %s,</p><p>%s</p>", m.Name, m.Content))
		d := gomail.NewDialer("smtp.gmail.com", 587, email, pwd)
		if err := d.DialAndSend(msg); err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
	}
}
