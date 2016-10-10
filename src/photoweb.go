package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	uploadDir = "./uploads"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, "<html><body><form method=\"post\" action=\"/upload\" enctype=\"multipart/form-data\">"+
			"Choose an image to upload:<input name=\"image\" type=\"file\" />"+
			"<input type=\"submit\" value=\"upload\" />"+
			"</form></body></html>")

		return
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		t, err := os.Create(uploadDir + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view?id="+filename,
			http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageID := r.FormValue("id")
	imagePath := uploadDir + "/" + imageID
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}

	var listHTML string
	for _, fileInfo := range fileInfoArr {
		imgid := fileInfo.Name()
		listHTML += "<li><a href=\"/view?id=" + imgid + "\">" + "<img src=\"/view?id=" + imgid + "\" width=\"100\" height=\"100\" />" + "</a></li>"
	}
	io.WriteString(w, "<html><body><ol>"+listHTML+"</ol></body></html>")
}

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
