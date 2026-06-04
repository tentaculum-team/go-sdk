package validate

import "testing"

func TestPostal(t *testing.T) {
	tests := []struct {
		name    string
		postal  string
		wantErr bool
	}{
		{"br cep bare", "01310100", false},
		{"br cep hyphen", "01310-100", false},
		{"us zip+4", "90210-1234", false},
		{"uk with space", "SW1A 1AA", false},
		{"empty", "", true},
		{"too short", "abc", true},
		{"invalid char", "0131$100", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Postal(tt.postal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Postal(%q) error = %v, wantErr %v", tt.postal, err, tt.wantErr)
			}
		})
	}
}
