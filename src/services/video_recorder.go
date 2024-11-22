package services

import (
	"fmt"
	"os/exec"
	"time"
)

type IVideoRecorder interface {
	Rec(filename string) (*string, error)
}
type VideoRecorder struct {
	duration   time.Duration
	cameraPort int
}

func (v *VideoRecorder) Rec(filename string) (*string, error) {
	/*cmd := exec.Command("ffmpeg",
		"-f", "v4l2", // Périphérique vidéo Linux
		"-input_format", "mp4v", // Format d'entrée natif
		"-s", "1280x720", // Full HD
		"-i", "/dev/video0", // Périphérique webcam
		"-t", "00:02:00", // Durée personnalisable
		filename, // Fichier de sortie
	)*/
	/*	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation", // Utilisation du format AVFoundation sur macOS
		"-framerate", "30", // Fréquence d'images souhaitée (par exemple 30 fps)
		"-video_device_index", "0", // Index de votre webcam, généralement "0" pour la première caméra
		"-s", "1280x720", // Résolution souhaitée
		"-t", "00:02:00", // Durée personnalisée
		"-i", "0", // Le périphérique d'entrée pour la webcam sur macOS
		"-c:v", "libx264", // Utilisation du codec vidéo H.264
		"-preset", "fast", // Préréglage d'encodage rapide
		"-crf", "23", // Qualité d'encodage (plus bas = meilleure qualité)
		filename, // Fichier de sortie
	)*/
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation", // Utilisation du format AVFoundation sur macOS
		"-framerate", "24", // Fréquence d'images souhaitée (par exemple 30 fps)
		"-video_device_index", "0", // Index de votre webcam
		"-s", "1280x720", // Résolution réduite pour alléger la vidéo
		"-t", "00:02:00", // Durée personnalisée
		"-i", "0", // Le périphérique d'entrée pour la webcam
		"-c:v", "libx264", // Utilisation du codec vidéo H.264
		"-preset", "ultrafast", // Préréglage d'encodage plus rapide pour réduire la taille
		"-crf", "23",
		"-an",    // Désactiver l'audio pour économiser de l'espace (optionnel)
		filename, // Fichier de sortie
	)
	fmt.Printf("Running command: %v\n", cmd)

	// Exécute la commande et capture les erreurs
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'exécution de ffmpeg: %v", err)
	}
	return &filename, nil
}

func NewVideoRecorder(duration time.Duration, cameraPort int) IVideoRecorder {
	return &VideoRecorder{duration: duration, cameraPort: cameraPort}
}
