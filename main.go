package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var pages *template.Template
var apiTempl *template.Template

func init(){
    var err error

    pages, err = template.ParseGlob("./pages/*.html")
    if err != nil {
        log.Fatalf("Error parsing page templates: %s\n", err)
    }

    apiTempl, err = template.ParseGlob("./apiTemplates/*.html")
    if err != nil {
        log.Fatalf("Error parsing api templates: %s\n", err)
    }
}

func displayHome(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(200)
    pages.ExecuteTemplate(w, "Home", nil)
}

func getTAsFromDB(db *sql.DB) ([]transaction, error){
    transActions := []transaction{}

    fetchAllQuery := `SELECT id, amount, description FROM transactions;`

    row, err := db.Query(fetchAllQuery)
    defer row.Close()
    if err != nil {
        return nil, err
    }

    for row.Next(){
        newTa := transaction{}
        err = row.Scan(&newTa.Id, &newTa.Amount, &newTa.Desc)
        if err != nil {
            return nil, err
        }
        transActions = append(transActions, newTa)
    }

    return transActions, nil
}

func addNewTrans(db *sql.DB, w http.ResponseWriter, r *http.Request){
    r.ParseForm()
    transAm, err := strconv.ParseFloat(r.FormValue("transAm"), 64)
    if err != nil {
        w.WriteHeader(400)
        return
    }

    newTrans := transaction{
        Desc: r.FormValue("transDesc"),
        Amount: transAm,
    }

    addNewTransDB(db, newTrans)

    w.WriteHeader(200)
    apiTempl.ExecuteTemplate(w, "TransDispl", newTrans)
}

func addNewTransDB(db *sql.DB, ta transaction) error{
    query := `INSERT INTO transactions (description, amount)
    VALUES (?,?);`

    _, err := db.Exec(query, ta.Desc, ta.Amount)
    if err != nil {
        return err
    }

    return nil
}

type transaction struct{
    Id int
    Desc string
    Amount float64
}

func dbInit() (*sql.DB, error){
    query := ` CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT,
        amount INTEGER NOT NULL);
    `

    db, err := sql.Open("sqlite3", "data.db")
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(query)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func main(){
    handler := http.NewServeMux()
    server := http.Server{
        Addr: ":42069",
        Handler: handler,
    }

    db, err := dbInit()
    if err != nil {
        log.Fatalf("error initializing the database: %s\n", err)
    }
    defer db.Close()

    log.Print(getTAsFromDB(db))
    // pages
    handler.HandleFunc("GET /", displayHome)

    // Api
    handler.HandleFunc("POST /api/addNewTa", func(w http.ResponseWriter, r *http.Request) {
        addNewTrans(db, w, r)
    })

    log.Printf("http server started on port %s\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
