package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	Username  string
	FlashData string
	Id        int
}

var Data = MetaData{
	Title: "Personal Web",
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var user = User{}

type Project struct {
	Id              int
	ProjectName     string
	StartDate       time.Time
	EndDate         time.Time
	DurationText    string
	Description     string
	Technology      []string
	Image           string
	FormatStartDate string
	FormatEndDate   string
	Author          string
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact-me", ContactMe).Methods("GET")
	route.HandleFunc("/add-project", AddProject).Methods("GET")
	route.HandleFunc("/projects/{id}", ProjectDetails).Methods("GET")
	route.HandleFunc("/add-new-project", middleware.UploadFile(AddNewProject)).Methods("POST")
	route.HandleFunc("/delete-project/{id}", DeleteProject).Methods("GET")
	route.HandleFunc("/update-project-page/{id}", UpdateProjectPage).Methods("GET")
	route.HandleFunc("/update-project/{id}", middleware.UploadFile(UpdateProject)).Methods("POST")
	route.HandleFunc("/register", RegisterPage).Methods("GET")
	route.HandleFunc("/login", LoginPage).Methods("GET")
	route.HandleFunc("/register", Register).Methods("POST")
	route.HandleFunc("/login", Login).Methods("POST")
	route.HandleFunc("/logout", Logout).Methods("GET")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var result []Project
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false

		rows, _ := connection.Conn.Query(context.Background(), "SELECT project.id, title, start_date, end_date, technologies, description, image, users.name as author FROM project LEFT JOIN users ON project.author_id = users.id  ORDER BY end_date DESC")

		for rows.Next() {
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Technology, &each.Description, &each.Image, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			each.FormatStartDate = each.StartDate.Format("Jan 21, 2000")
			each.FormatEndDate = each.EndDate.Format("Jan 21, 2000")

			result = append(result, each)
		}

	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
		Data.Id = session.Values["Id"].(int)

		rows, _ := connection.Conn.Query(context.Background(), "SELECT project.id, title, start_date, end_date, technologies, description, image, users.name as author FROM project LEFT JOIN users ON project.author_id = users.id WHERE project.author_id=$1 ORDER BY id DESC", Data.Id)
		for rows.Next() {
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Technology, &each.Description, &each.Image, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			each.FormatStartDate = each.StartDate.Format("Jan 21, 2000")
			each.FormatEndDate = each.EndDate.Format("Jan 21, 2000")

			result = append(result, each)
		}

	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func ContactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO users(name, email, password) VALUES ($1,$2,$3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user = User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM public.users WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800

	session.AddFlash("Login success", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1

	session.AddFlash("Logout success", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func ProjectDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetails := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM public.project WHERE.id=$1;", id).Scan(&ProjectDetails.Id, &ProjectDetails.ProjectName, &ProjectDetails.StartDate, &ProjectDetails.EndDate, &ProjectDetails.Description, &ProjectDetails.Technology, &ProjectDetails.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":           Data,
		"ProjectDetails": ProjectDetails,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func AddNewProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ProjectName := r.PostForm.Get("project-name")
	StartDate := r.PostForm.Get("start-date")
	EndDate := r.PostForm.Get("end-date")
	Description := r.PostForm.Get("description")
	Technology := r.Form["technology"]
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)
	// var store = sessions.NewCookieStore([]byte("SESSIONID"))
	// session, _ := store.Get(r, "SESSION_ID")
	// author := session.Values["ID"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO project(title, start_date, end_date, description, technologies, author_id, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", ProjectName, StartDate, EndDate, Description, Technology, user.Id, image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-chace, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func UpdateProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/update-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	UpdateData := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM project WHERE id=$1", id).Scan(&UpdateData.Id, &UpdateData.ProjectName, &UpdateData.StartDate, &UpdateData.EndDate, &UpdateData.Description, &UpdateData.Technology, &UpdateData.Image, &UpdateData.Author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var nodeJS, React, Angular, VueJS bool
	for _, A := range UpdateData.Technology {
		if A == "node" {
			nodeJS = true
		}
		if A == "react" {
			React = true
		}
		if A == "angular" {
			Angular = true
		}
		if A == "vuejs" {
			VueJS = true
		}
	}
	StartDateString := UpdateData.StartDate.Format("2006-01-02")
	EndDateString := UpdateData.EndDate.Format("2006-01-02")

	respData := map[string]interface{}{
		"Data":            Data,
		"Id":              id,
		"UpdateData":      UpdateData,
		"StartDateString": StartDateString,
		"EndDateString":   EndDateString,
		"nodeJS":          nodeJS,
		"React":           React,
		"Angular":         Angular,
		"VueJS":           VueJS,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ProjectName := r.PostForm.Get("project-name")
	StartDateString := r.PostForm.Get("start-date")
	EndDateString := r.PostForm.Get("end-date")
	Description := r.PostForm.Get("description")
	Technology := r.Form["technology"]
	StartDate, _ := time.Parse("2006-01-02", StartDateString)
	EndtDate, _ := time.Parse("2006-01-02", EndDateString)
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.project SET title=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7", ProjectName, StartDate, EndtDate, Description, Technology, image, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func CalculateDuration(StartDate string, EndDate string) string {
	StartDateFormated, _ := time.Parse("2006-01-02", StartDate)
	EndtDateFormated, _ := time.Parse("2006-01-02", EndDate)
	Duration := EndtDateFormated.Sub(StartDateFormated)
	DurationHours := Duration.Hours()
	DurationDays := math.Floor(DurationHours / 24)
	DurationWeeks := math.Floor(DurationDays / 7)
	DurationMonths := math.Floor(DurationDays / 30)
	DurationText := ""
	if DurationMonths > 1 {
		DurationText = strconv.FormatFloat(DurationMonths, 'f', 0, 64) + " months"
	} else if DurationMonths > 0 {
		DurationText = strconv.FormatFloat(DurationMonths, 'f', 0, 64) + " month"
	} else {
		if DurationWeeks > 1 {
			DurationText = strconv.FormatFloat(DurationWeeks, 'f', 0, 64) + " weeks"
		} else if DurationWeeks > 0 {
			DurationText = strconv.FormatFloat(DurationWeeks, 'f', 0, 64) + " week"
		} else {
			if DurationDays > 1 {
				DurationText = strconv.FormatFloat(DurationDays, 'f', 0, 64) + " days"
			} else if DurationDays > 0 {
				DurationText = strconv.FormatFloat(DurationDays, 'f', 0, 64) + " day"
			} else {
				DurationText = "less than a day"
			}
		}
	}
	return DurationText
}

func ConvertTechnologyToBoolean(
	Technology1String string,
	Technology2String string,
	Technology3String string,
	Technology4String string) (
	bool,
	bool,
	bool,
	bool) {
	Technology1 := false
	if Technology1String == "on" {
		Technology1 = true
	}
	Technology2 := false
	if Technology2String == "on" {
		Technology2 = true
	}
	Technology3 := false
	if Technology3String == "on" {
		Technology3 = true
	}
	Technology4 := false
	if Technology4String == "on" {
		Technology4 = true
	}
	return Technology1, Technology2, Technology3, Technology4
}
