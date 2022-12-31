package discovery

import (
	"testing"
)

func TestParseMetaData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"good json", "", false},
	},
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMetaData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
