package discovery

import "testing"

func TestValidateIssuer(t *testing.T) {
	tests := []struct {
		name    string
		issuer  string
		wantErr bool
	}{
		{"empty", "", true},
		{"invalid URL", "sdfgsdfgsdfgsdfg", true},
		{"HTTP", "http://example.com", true},
		{"Empty Query", "https://example.com/path?", true},
		{"Query", "https://example.com/path?q=1", true},
		{"Empty Fragment", "https://example.com/path#", true},
		{"Fragment", "https://example.com/path#frag", true},
		{"Good no path", "https://example.com", false},
		{"Good empty path", "https://example.com/", false},
		{"Good with path", "https://example.com/some/path", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateIssuer(tt.issuer); (err != nil) != tt.wantErr {
				t.Errorf("ValidateIssuer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
