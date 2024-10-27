package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Mufraggi/figther-capture/src/services"
	"gocv.io/x/gocv"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func testVideoService() {
	inputVideo := "video.mp4"
	logo := "bl-logo.png" // Assurez-vous que le logo est au format PNG pour la transparence
	outputVideo := "video_avec_logo.mp4"

	// Ajoute le logo à la vidéo
	videoService := services.NewVideoService(logo)
	err := videoService.AddLogoToVideo(inputVideo, outputVideo)
	if err != nil {
		fmt.Println("Erreur:", err)
	} else {
		fmt.Println("Logo ajouté avec succès à la vidéo !")
	}
}
func createFileName() string {
	timestamp := time.Now()
	dateStr := timestamp.Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s.mp4", dateStr)
}

func initRecord() string {
	webcam, _ := gocv.OpenVideoCapture(0)
	defer webcam.Close()

	// Nom du fichier vidéo à enregistrer
	videoFile := createFileName()

	// Capture de vidéo
	services.RecordVideo(webcam, videoFile)
	return videoFile
}

const videoDuration = 6 * time.Minute // Durée de 6 minutes

func UploadFileTest() {
	credentialsPath := "service-account.json"
	folderID := "12C0lsnLpMHJW5olW5rLrYmxag3kqdbep"
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Impossible de se connecter au service Drive: %v", err)
	}

	// Créer les métadonnées du fichier
	f := &drive.File{
		Name:    "mon_fichier.mp4",
		Parents: []string{folderID}, // Spécifiez le dossier où le fichier sera stocké
	}

	// Ouvrir le fichier local
	file, err := os.Open("video_output.mp4")
	if err != nil {
		log.Fatalf("Impossible d'ouvrir le fichier: %v", err)
	}
	defer file.Close()

	// Upload du fichier
	fileDrive, err := srv.Files.Create(f).Media(file).Do()
	if err != nil {
		log.Fatalf("Impossible de créer le fichier dans Google Drive: %v", err)
	}

	// Confirmation de l'upload
	fmt.Printf("Fichier uploadé avec succès. ID : %s\n", fileDrive.Id)
}

func uplaod_video(filePath string) {

	// Spécifiez le chemin du fichier vidéo
	url := "http://localhost:8080/video"

	// Ouvrir le fichier
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	// Créer un buffer pour stocker les données du formulaire
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Ajouter le fichier au formulaire
	part, err := writer.CreateFormFile("video", filePath)
	if err != nil {
		fmt.Println("Erreur lors de la création du formulaire :", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Erreur lors de la copie du fichier dans le formulaire :", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Erreur lors de la fermeture du writer :", err)
		return
	}

	// Créer la requête HTTP POST
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return
	}
	fmt.Println("Réponse du serveur :", string(body))
}

func main() {
	path := initRecord()
	uplaod_video(path)
}
