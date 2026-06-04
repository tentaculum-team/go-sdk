package validate

import "testing"

func TestCNPJ(t *testing.T) {
	tests := []struct {
		name    string
		cnpj    string
		cfg     []CNPJConfig
		wantErr bool
	}{
		{"valid bare", "11222333000181", nil, false},
		{"valid formatted", "11.222.333/0001-81", nil, false},
		{"empty", "", nil, true},
		{"wrong length", "112223330001", nil, true},
		{"all same digits", "11111111111111", nil, true},
		{"bad check digits", "11222333000182", nil, true},
		{"bad check digits but skipped", "11222333000182", []CNPJConfig{{CheckDigits: false, AllErrors: true}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CNPJ(tt.cnpj, tt.cfg...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CNPJ(%q) error = %v, wantErr %v", tt.cnpj, err, tt.wantErr)
			}
		})
	}
}
