package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

type notifier struct {
	title string
	path  string
}

func NewNotifier(title, path string) *notifier {
	return &notifier{
		title: title,
		path:  path,
	}
}

func (n *notifier) Start() {
	beeep.Notify(n.title, fmt.Sprintf("Processing: %s", n.path), "/usr/share/icons/gnome/32x32/emblems/emblem-documents.png")

}

func (n *notifier) Error(err error) {
	beeep.Notify(n.title, fmt.Sprintf("Error: %s (%q)", n.path, err), "/usr/share/icons/gnome/32x32/emblems/emblem-important.png")

}

func (n *notifier) Success() {
	beeep.Notify(n.title, fmt.Sprintf("Done (%s)", n.path), "/usr/share/icons/gnome/32x32/emblems/emblem-default.png")
}
