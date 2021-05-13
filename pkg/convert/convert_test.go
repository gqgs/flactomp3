package convert

import "testing"

func TestIsConvertible(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"convertible FLAC (lowercase)",
			"14 - Cradling Mother, Cradling Woman.flac",
			true,
		},
		{
			"convertible FLAC (uppercase)",
			"04. Romance.FLAC",
			true,
		},
		{
			"non-convertible",
			"13 - Scar.mp3",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConvertible(tt.path); got != tt.want {
				t.Errorf("IsConvertible() = %v, want %v", got, tt.want)
			}
		})
	}
}
