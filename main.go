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
    if r.URL.Path != "/"{
        w.WriteHeader(404)
        return
    }
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

func displayAllTasFromDb(db *sql.DB, w http.ResponseWriter){
    transactions, err := getTAsFromDB(db)
    if err != nil {
        w.WriteHeader(500)
        log.Printf("error fetching all transactions from the db %s\n", err)
        return
    }
    apiTempl.ExecuteTemplate(w, "DispAllTas", transactions)
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

func deleteTransaction(db *sql.DB, w http.ResponseWriter, r *http.Request) error{
    deleteQuery := `DELETE FROM transactions where id = ?;`
    log.Println(r)
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        w.WriteHeader(500)
        log.Printf("error parsing the id from the request %s\n", err)
        return err
    }

    _, err = db.Exec(deleteQuery, id)
    if err != nil {
        log.Printf("error sending the delete request to the db %s\n", err)
        w.WriteHeader(500)
        return err
    }

    w.WriteHeader(200)
    return nil
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

    // pages
    handler.HandleFunc("GET /", displayHome)

    // Api
    handler.HandleFunc("POST /api/addNewTa", func(w http.ResponseWriter, r *http.Request) {
        addNewTrans(db, w, r)
    })
    handler.HandleFunc("POST /api/fetchAllTas", func(w http.ResponseWriter, r *http.Request) {
        displayAllTasFromDb(db, w)
    })
    handler.HandleFunc("DELETE /api/delTrans/{id}", func(w http.ResponseWriter, r *http.Request) {
        deleteTransaction(db, w, r)
    })

    // Files
    handler.Handle("GET /files/js/main.js", http.StripPrefix("/files/js", http.FileServer(http.Dir("./"))))

    /*TODO:
    Add the possibility to edit already added inputs or outputs.
    Add a chart (e.g. Chart.js) to visualize the distribution of income and expenditure.
    Add categories for income and expenditure (e.g. food, transport, entertainment).
    Add the possibility to filter income and expenditure by category.
    Tyylit
    */

    log.Printf("http server started on port %s\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
