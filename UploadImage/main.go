package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var uploadTemplate, _ = template.ParseFiles("upload.html")

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		uploadTemplate.Execute(w, nil)
		return
	}
	f, _, err := r.FormFile("image")
	check(err)
	defer f.Close()
	t, err := ioutil.TempFile(".", "image-")
	check(err)
	defer t.Close()
	_, err = io.Copy(t, f)
	check(err)
	http.Redirect(w, r, "/view?id="+t.Name()[6:], 302)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

var errorTemplate, _ = template.ParseFiles("error.html")

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				w.WriteHeader(500)
				errorTemplate.Execute(w, e)
			}
		}()
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/view", errorHandler(view))
	log.Println("Init app on :8080")
	http.ListenAndServe(":8080", nil)
}
