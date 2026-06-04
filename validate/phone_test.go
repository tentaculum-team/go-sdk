package validate

import "testing"

func TestPhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		cfg     []PhoneConfig
		wantErr bool
	}{
		{"valid digits", "11999998888", nil, false},
		{"valid with plus", "+5511999998888", nil, false},
		{"empty", "", nil, true},
		{"too short", "123", nil, true},
		{"letters", "11abc998888", nil, true},
		{"formatted rejected when digits only", "(11) 99999-8888", nil, true},
		{"formatted allowed", "(11) 99999-8888", []PhoneConfig{{MaxChars: 20, MinChars: 8, DigitsOnly: false, AllErrors: true}}, false},
		{"require plus missing", "5511999998888", []PhoneConfig{{MaxChars: 20, MinChars: 8, RequirePlus: true, DigitsOnly: true, AllErrors: true}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Phone(tt.phone, tt.cfg...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Phone(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}
