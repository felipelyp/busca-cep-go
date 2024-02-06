package main

// ToUTF8 Converter charset=ISO-8859-1 para UTF8
// O resultado da api chega como charset=ISO-8859-1
// eu não consigo usar xml.NewDecoder diretamente
func ToUTF8(buffer []byte) []byte {
	buf := make([]rune, len(buffer))
	for i, b := range buffer {
		buf[i] = rune(b)
	}
	return []byte(string(buf))
}
