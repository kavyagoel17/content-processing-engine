package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/kavya/content-engine/internal/logger"
	"github.com/kavya/content-engine/internal/services"
)

type FileMetadata struct {
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

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

	metadata := FileMetadata{
		Filename:    handler.Filename,
		Size:        handler.Size,
		ContentType: handler.Header.Get("Content-Type"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
	
	dst, filePath, err := CreateFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	logger.InfoLog.Printf("Reading file: %s", filePath)
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}

	services.Extractor(filePath)
}

func CreateFile(filename string) (*os.File, string, error) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	path := filepath.Join("uploads", filename)
	dst, err := os.Create(path)
	if err != nil {
		return nil, "", err
	}
	return dst, path, nil

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