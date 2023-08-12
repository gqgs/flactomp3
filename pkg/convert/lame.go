package convert

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type lame struct{}

func NewLameConverter() *lame {
	return &lame{}
}

func (c *lame) Convert(relativePath, baseFolder, outFolder string) error {
	newFilename := strings.TrimSuffix(relativePath, ".flac") + ".mp3"
	newPath := filepath.Join(outFolder, newFilename)
	originalPath := filepath.Join(baseFolder, relativePath)
	log.Printf("encoding: %q", newFilename)
	cmd := exec.Command("ffmpeg", "-i", originalPath, "-c:v", "copy", "-c:a", "libmp3lame", "-q:a", "0", "-threads", "0", "-loglevel", "error", newPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
