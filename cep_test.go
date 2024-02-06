package main

import (
	"testing"
)

// CEP de São Paulo
const cep = "01007070"

// Pesquisar cep valido
func TestCep_Search(t *testing.T) {
	var cepTests = []struct {
		cep string
	}{
		{cep},                    //Correct
		{cep[:len(cep)-1]},       //Incorrect
		{cep[:len(cep)-1] + "1"}, //Not found
		{cep + "00"},             //Invalid
	}

	for _, tt := range cepTests {
		t.Run(tt.cep, func(t *testing.T) {
			b := NewBuscaCep()
			result, err := b.Search(tt.cep)
			if err != nil {
				t.Log(b.Error)
			} else {
				t.Log(result)
			}
		})
	}
}

// Pesquisar e gerar erro durante a consulta
func TestCep_SearchError(t *testing.T) {
	b := BuscaCep{}
	_, err := b.Search(cep)
	t.Log(err)
}

// Pesquisar e gerar erro durante a consulta
func TestCep_SearchResponseError(t *testing.T) {
	b := BuscaCep{SoapCorreiosUrl: "https://github.com/felipelyp/"}
	_, err := b.Search(cep)
	if err == nil {
		t.Fatal("Não era pra ser nulo")
	}
	t.Log(err)
}

// Gerar XML com o cep
func TestCep_GetBody(t *testing.T) {
	out := string(GetBody(cep))
	t.Log("XML Gerado: " + out)
}

func TestToUTF8(t *testing.T) {
	cht := []byte("\126\117\103\312")
	t.Log(string(ToUTF8(cht)))
}
