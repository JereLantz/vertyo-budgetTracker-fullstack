package main

import (
	"html/template"
	"log"
	"net/http"
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

func addNewTrans(w http.ResponseWriter, r *http.Request){
    testTa := transaction{
        Id: 1,
        Desc: "sdfjlkas",
        Amount: 54.32,
    }

    w.WriteHeader(200)
    apiTempl.ExecuteTemplate(w, "TransDispl", testTa)
}

type transaction struct{
    Id int
    Desc string
    Amount float64
}

func main(){
    handler := http.NewServeMux()
    server := http.Server{
        Addr: ":42069",
        Handler: handler,
    }

    // pages
    handler.HandleFunc("GET /", displayHome)

    // Api
    handler.HandleFunc("POST /api/addNewTa", addNewTrans)

    log.Printf("http server started on port %s\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
