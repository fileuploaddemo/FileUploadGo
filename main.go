package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type FileInfo struct {
	Id   int    `json:"id"`
	Size int64  `json:"size"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("www"))
	mux.Handle("/", fs)
	mh := http.HandlerFunc(handleRequest)
	mux.Handle("/files/", mh)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	fmt.Println("Now listening on: http://0.0.0.0:8080")
	fmt.Println("Application started. Press Ctrl+C to shut down.")
	server.ListenAndServe()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("GET: " + r.URL.Path)
	name := path.Base(r.URL.Path)
	fmt.Println("download: " + name)
	file, err := retrieve()
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&file, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("POST: " + r.URL.Path)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("newfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	return
}
func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Println("DELETE: " + r.URL.Path)
	name := path.Base(r.URL.Path)
	fmt.Println("DELETE: " + name)
	if err != nil {
		return
	}
	fileInfo, err := find(name)
	if err != nil {
		return
	}
	err = fileInfo.delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
