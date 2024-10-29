package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ClientHttp struct {
	url string
}

func (c *ClientHttp) Send(filePath string) error {

	// Spécifiez le chemin du fichier vidéo
	//url := "http://localhost:8080/video"

	// Ouvrir le fichier
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return err
	}
	defer file.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("video", filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

type IClientHttp interface {
	Send(videoFile string) error
}

func NewClientHttp(url string) IClientHttp {
	return &ClientHttp{url: url}
}
