package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BuscaCep Buscador de informações
// por meio de CEP e receber dados
// como Bairro, Cidade, Endereço, Complemento e UF
type BuscaCep interface {
	// Search pesquisar cep e obter Result
	Search(cep string) (Result, error)
}

// Cep Buscador de cep via Soap
type Cep struct {
	//Erro para depuração
	Error           error
	SoapCorreiosUrl string
}

// Payload dados para enviar a api
type Payload struct {
	//Tag soapenv:Envelope
	XMLName xml.Name `xml:"soapenv:Envelope"`
	//Atributo xmlns:soapenv
	Env string `xml:"xmlns:soapenv,attr"`
	//Atributo xmlns:cli
	Cli string `xml:"xmlns:cli,attr"`
	//Tag de consultar para obter os dados consultado
	Cep string `xml:"soapenv:Body>cli:consultaCEP>cep"`
}

// CepResponse dados da consulta
type CepResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Body    struct {
		Text    string `xml:",chardata"`
		Consult struct {
			Text   string `xml:",chardata"`
			Ns2    string `xml:"ns2,attr"`
			Result Result `xml:"return"`
		} `xml:"consultaCEPResponse"`
	} `xml:"Body"`
}

// Result Informações do cep consultado
// Cuidado com dados inexistente como
// Complemento ou até mesmo um Endereco
type Result struct {
	UF          string `xml:"uf"`
	Cep         string `xml:"cep"`
	Cidade      string `xml:"cidade"`
	Bairro      string `xml:"bairro"`
	Endereco    string `xml:"end"`
	Complemento string `xml:"complemento2"`
}

type CepFault struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Body    struct {
		Text  string `xml:",chardata"`
		Fault Fault  `xml:"Fault"`
	} `xml:"Body"`
}

type Fault struct {
	Text        string `xml:",chardata"`
	Faultstring string `xml:"faultstring"`
}

// NewBuscaCep Nova instancia
func NewBuscaCep() BuscaCep {
	return &Cep{
		SoapCorreiosUrl: "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente",
	}
}

// Search pesquisar por meio do cep e receber um Result
func (c *Cep) Search(cep string) (Result, error) {
	var result Result
	if t := len(cep); t < 8 {
		return result, errors.New(fmt.Sprintf("o cep deve contér 8 digitos, atual: %d digitos", t))
	}
	payload := bytes.NewBuffer(GetBody(cep))
	res, err := http.Post(c.SoapCorreiosUrl, "application/xml", payload)
	if err != nil {
		c.Error = err
		return result, errors.New("não foi possível consultar o cep, tente novamente")
	}
	out, _ := io.ReadAll(res.Body)
	if strings.Contains(string(out), "faultstring") {
		var cepFault CepFault
		_ = xml.Unmarshal(ToUTF8(out), &cepFault)
		return result, errors.New(cepFault.Body.Fault.Faultstring)
	} else {
		var cepResponse CepResponse
		err = xml.Unmarshal(ToUTF8(out), &cepResponse)
		if err != nil {
			return result, errors.New("não foi possível obter o resultado da cep")
		}
		return cepResponse.Body.Consult.Result, nil
	}
}

// GetBody gear xml com os dados para consultar e retonar []byte
func GetBody(cep string) []byte {
	payload := Payload{
		Env: "http://schemas.xmlsoap.org/soap/envelope/",
		Cli: "http://cliente.bean.master.sigep.bsb.correios.com.br/",
		Cep: cep,
	}
	out, _ := xml.Marshal(payload)
	return out
}
