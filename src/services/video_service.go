package services

import (
	"fmt"
	"os"
	"os/exec"
)

type IVideoService interface {
	AddLogoToVideo(inputVideo, outputVideo string) error
}

type videoService struct {
	logoPath string
}

func (v *videoService) AddLogoToVideo(inputVideo, outputVideo string) error {
	if _, err := os.Stat(inputVideo); os.IsNotExist(err) {
		return fmt.Errorf("le fichier vidéo d'entrée n'existe pas : %s", inputVideo)
	}
	if _, err := os.Stat(v.logoPath); os.IsNotExist(err) {
		return fmt.Errorf("le fichier logo n'existe pas : %s", v.logoPath)
	}

	cmd := exec.Command("ffmpeg", "-i", inputVideo, "-i", v.logoPath, "-filter_complex",
		"[1:v] scale=200:200 [logo]; [0:v][logo] overlay=10:H-h-10", "-c:a", "copy", outputVideo)

	// Affiche la commande pour le debug
	fmt.Printf("Running command: %v\n", cmd)

	// Exécute la commande et capture les erreurs
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution de ffmpeg: %v", err)
	}
	return nil
}

func NewVideoService(logoPath string) IVideoService {
	return &videoService{logoPath: logoPath}
}
