package validate

import "testing"

func TestCountry(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		cfg     []CountryConfig
		wantErr bool
	}{
		{"BR", "BR", nil, false},
		{"US", "US", nil, false},
		{"lowercase rejected by default", "br", nil, true},
		{"lowercase allowed", "br", []CountryConfig{{Uppercase: false, AllErrors: true}}, false},
		{"empty", "", nil, true},
		{"too long", "BRA", nil, true},
		{"digits", "B1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Country(tt.code, tt.cfg...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Country(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			}
		})
	}
}
