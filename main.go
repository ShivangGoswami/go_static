package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type staticHandler struct {
}

func (sh *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("Path requested:", path)

	data, err := ioutil.ReadFile(string("./html/" + path + ".html"))
	if err == nil {
		w.Write(data)
	} else {
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func templateWriter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("Path requested:", path)
	r.ParseForm()
	type format struct {
		Name template.HTML
	}
	data := format{
		Name: template.HTML(r.Form.Get("aname")),
	}
	if data.Name == "" {
		w.Write([]byte("<h1>Please provide a name!</h1>"))
		return
	}
	const temp = `<h1>Hi! My name is {{.Name}}</h1>`
	t := template.Must(template.New("dynamic").Parse(temp))
	if err := t.Execute(w, data); err != nil {
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./html")))
	mux.Handle("/static", new(staticHandler))
	mux.HandleFunc("/dynamic", templateWriter)
	log.Fatal(http.ListenAndServe(":1234", mux))
}
