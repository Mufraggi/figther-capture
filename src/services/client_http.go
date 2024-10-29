package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ClientHttp struct {
	url string
}

func (c *ClientHttp) Send(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("erreur d'accès au fichier : %w", err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier : %w", err)
	}
	defer file.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("video", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("erreur création form file : %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("erreur copie fichier : %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("erreur fermeture writer : %w", err)
	}

	req, err := http.NewRequest("POST", c.url, &requestBody)
	if err != nil {
		return fmt.Errorf("erreur création requête : %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(requestBody.Len())

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur envoi requête : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur serveur : status=%d, body=%s", resp.StatusCode, string(body))
	}
	return nil
}

type IClientHttp interface {
	Send(videoFile string) error
}

func NewClientHttp(url string) IClientHttp {
	return &ClientHttp{url: url}
}
