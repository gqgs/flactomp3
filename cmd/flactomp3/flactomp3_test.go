package main

import "testing"

func Test_updateFilename(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		format string
		want   string
	}{
		{
			name:   "path without format",
			path:   "[2020] Shore",
			format: "V0",
			want:   "[2020] Shore [V0]",
		},
		{
			name:   "format between parenthesis",
			path:   "Madvillain - 2004 - Madvillainy (FLAC)",
			format: "V0",
			want:   "Madvillain - 2004 - Madvillainy (V0)",
		},
		{
			name:   "format between brackets",
			path:   "Radiohead - 2001 - Amnesiac (CDP 7243 5 32764 2 3) [FLAC]",
			format: "V0",
			want:   "Radiohead - 2001 - Amnesiac (CDP 7243 5 32764 2 3) [V0]",
		},
		{
			name:   "format as suffix",
			path:   "Gorillaz - Demon Days (2005) FLAC",
			format: "V0",
			want:   "Gorillaz - Demon Days (2005) V0",
		},
		{
			name:   "captalized format",
			path:   "Bonobo - Migration (2017) Flac",
			format: "V0",
			want:   "Bonobo - Migration (2017) V0",
		},
		{
			name:   "format with source",
			path:   "Billie Eilish - dont smile at me (2017) {Bonus Track - 00602567328070} [WEB FLAC]",
			format: "V0",
			want:   "Billie Eilish - dont smile at me (2017) {Bonus Track - 00602567328070} [WEB V0]",
		},
		{
			name:   "lowercase format with source",
			path:   "Milo - Who Told You To Think!!!!! (2017) [WEB flac]",
			format: "V0",
			want:   "Milo - Who Told You To Think!!!!! (2017) [WEB V0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateFilename(tt.path, tt.format); got != tt.want {
				t.Errorf("updateFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}
