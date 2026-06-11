package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kavya/content-engine/internal/logger"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	logger.InfoLog.Printf("Uploading File")
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		logger.ErrorLog.Println(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	logger.InfoLog.Printf("Uploaded File: %s", handler.Filename)
	logger.InfoLog.Printf("File Size: %d", handler.Size)
	logger.InfoLog.Printf("MIME Header: %v", handler.Header)

	dst, err := CreateFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}
}

func CreateFile(filename string) (*os.File, error) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return nil, err
	}
	return dst, nil

}

func IsValidFileType(file []byte) bool {
	fileType := http.DetectContentType(file)

	return fileType == "application/pdf"
}

func SetupRoutes(port int) {
	http.HandleFunc("/upload", UploadFile)

	addr := fmt.Sprintf(":%d", port)

	logger.InfoLog.Printf("Server starting on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}