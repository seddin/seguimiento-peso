package main

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"path"
)

type Entrada struct {
	Peso *string
	Fecha *string
}

func main () {

	r := httprouter.New()

	r.NotFound = http.FileServer(http.Dir("public"))

	r.GET("/", HomeHandler)
	r.POST("/añadir", AñadirHandler)


	http.ListenAndServe(":10000", r)
}


func HomeHandler (rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, _ := sql.Open("sqlite3", "data.db")
	defer db.Close()

	rows, _ := db.Query("SELECT peso, fecha FROM seguimiento ORDER BY id ASC")
	defer rows.Close()

	entradas := []Entrada{}

	for rows.Next() {
		var peso, fecha string
		rows.Scan(&peso, &fecha)

		entradas = append(entradas, Entrada{ &peso, &fecha})
	}

	RenderHome(rw, &entradas)
}

func RenderHome (rw http.ResponseWriter, entradas *[]Entrada) {
	fp := path.Join("public", "index.html")

	tmpl, _ := template.ParseFiles(fp)
	tmpl.Execute(rw, entradas)
}

func AñadirHandler (rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, _ := sql.Open("sqlite3", "data.db")
	defer db.Close()

	peso := r.FormValue("peso")
	fecha := r.FormValue("fecha")

	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("INSERT INTO seguimiento(peso, fecha) values (?, ?)")
	defer stmt.Close()

	stmt.Exec(peso, fecha)
	tx.Commit()

	http.Redirect(rw, r, "/", 303)
}