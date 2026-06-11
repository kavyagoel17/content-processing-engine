package services

import (
	"bytes"

	"github.com/kavya/content-engine/internal/logger"
	"github.com/ledongthuc/pdf"
)

func Extractor(filePath string) {
	pdf, reader, err := pdf.Open(filePath)
	if err != nil {
		logger.ErrorLog.Println(err)
	}

	var buf bytes.Buffer
	b, err := reader.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := buf.String()
	logger.InfoLog.Printf("Extracted text:\n", content)

	defer pdf.Close()
}