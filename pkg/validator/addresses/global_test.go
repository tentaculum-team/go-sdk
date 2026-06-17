package validator

import "testing"

func TestGlobalAddressCases(t *testing.T) {
	ok := func(name string, err error) {
		if err != nil {
			t.Errorf("%s: expected valid, got %v", name, err)
		}
	}
	bad := func(name string, err error) {
		if err == nil {
			t.Errorf("%s: expected error, got nil", name)
		}
	}

	ok("city São Paulo", City("São Paulo"))
	ok("city Saint-Étienne", City("Saint-Étienne"))
	ok("city 100 Mile House", City("100 Mile House"))
	ok("city O'Fallon", City("O'Fallon"))
	ok("city Washington, D.C.", City("Washington, D.C."))
	bad("city empty", City("   "))
	bad("city injection", City("<script>"))

	ok("postal UK", Postal("SW1A 1AA"))
	ok("postal US+4", Postal("12345-6789"))
	ok("postal JP", Postal("100-0001"))
	ok("postal NL", Postal("1234 AB"))
	bad("postal symbols", Postal("12$45"))

	ok("number 221B", AddressNumber("221B"))
	ok("number s/n", AddressNumber("s/n"))
	ok("number 12/3", AddressNumber("12/3"))
	bad("number onlydigits", AddressNumber("221B", HouseNumberConfig{OnlyNumbers: true, AllErrors: true}))

	ok("street intl", Street("Av. Paulista, 1000 (Bela Vista)"))
	ok("complement", Complement("Apt 4B / 2nd floor"))
	ok("state name", StateRegion("New South Wales"))
	ok("state code", StateRegion("SP"))
	ok("district", District("CT-Étienne"))
	ok("label", Label("Casa de praia #2"))
	ok("country", Country("BR"))
	bad("country bad", Country("Brazil"))
}
