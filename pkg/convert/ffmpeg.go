package convert

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ffmpeg struct{}

func NewFFmpegConverter() *ffmpeg {
	return &ffmpeg{}
}

func (f *ffmpeg) Convert(relativePath, baseFolder, outFolder string) error {
	mp3FileName := strings.TrimSuffix(relativePath, ".flac") + ".mp3"
	newPath := filepath.Join(outFolder, mp3FileName)
	originalPath := filepath.Join(baseFolder, relativePath)
	log.Printf("encoding: %q", mp3FileName)
	cmd := exec.Command("ffmpeg", "-i", originalPath, "-q:a", "0", "-threads", "0", "-loglevel", "error", newPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
