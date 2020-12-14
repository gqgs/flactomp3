package convert

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type opus struct{}

func NewOpusConverter() *opus {
	return &opus{}
}

func (c *opus) Convert(relativePath, baseFolder, outFolder string) error {
	newFilename := strings.TrimSuffix(relativePath, ".flac") + ".opus"
	newPath := filepath.Join(outFolder, newFilename)
	originalPath := filepath.Join(baseFolder, relativePath)
	log.Printf("encoding: %q", newFilename)
	cmd := exec.Command("ffmpeg", "-i", originalPath, "-c:a", "libopus", "-vbr", "1", "-application", "audio", "-b:a", "128k", "-threads", "0", "-loglevel", "error", newPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
