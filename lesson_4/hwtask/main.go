package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	filePath := h.UploadDir + "/" + header.Filename
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fileLink := h.HostAddr + "/" + header.Filename

	fmt.Fprintln(w, fileLink)
	fmt.Fprintf(w, "File %s has been successfully uploaded", header.Filename)
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	// Добавить в пример с файловым сервером возможность
	// получить список всех файлов на сервере (имя, расширение, размер в байтах)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		files, err := ioutil.ReadDir(uploadHandler.UploadDir)
		if err != nil {
			log.Fatal(err)
		}
		// С помощью query-параметра, реализовать
		// фильтрацию выводимого списка по расширению
		// (то есть, выводить только .png файлы, или только .jpeg)
		// example: http://localhost/?ext=.txt
		fileExt := r.FormValue("ext")
		for _, file := range files {
			if filepath.Ext(file.Name()) == fileExt || fileExt == "" {
				fileName, _ := json.Marshal(file.Name())
				fileSize, _ := json.Marshal(file.Size())
				fmt.Fprintf(w, "%s size: %s bites\n", fileName, string(fileSize))
			}
		}
	})
	go http.ListenAndServe(":80", nil)

	dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         "localhost:8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()

}

// добавить новый метод
// который лазает по деректории смотрит на файлы и показывает список файло
// в установленом формате (имя расширение размер в байтах) пусть это будето JSON
// будет запрос GET на такой то URL
