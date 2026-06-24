package utils

func CPFCheckDigitsValid(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}

	d1 := (sum * 10) % 11
	if d1 == 10 {
		d1 = 0
	}

	if d1 != int(cpf[9]-'0') {
		return false
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}

	d2 := (sum * 10) % 11
	if d2 == 10 {
		d2 = 0
	}

	return d2 == int(cpf[10]-'0')
}
