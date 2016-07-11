package main

import (
	"encoding/xml"
	"fmt"
)

type SoapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi     string   `xml:"xmlns xsi,attr"`
	Xsd     string   `xml:"xmlns xsd,attr"`
	Soap    string   `xml:"xmlns soap,attr"`
	Body    SoapBody
}

type SoapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	ObterUsuarios ObterUsuarios
}

type ObterUsuarios struct {
	XMLName      xml.Name   `xml:"obterUsuarios"`
	Param        UsuarioDto `xml:"param"`
}

type UsuarioDto struct {
	Agencia      string `xml:"agencia"`
	Conta        string `xml:"conta"`
	Banco        string `xml:"banco"`
	Imei         string `xml:"IMEI"`
	Linha        string `xml:"linha"`
}	set.Xsi = "http://www.w3.org/2001/XMLSchema-i	set.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	set.Xsd = "http://www.w3.org/2001/XMLSchema"
	set.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	set.Body.ObterUsuarios.XMLNamespace = "http:/	set.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	set.Xsd = "http://www.w3.org/2001/XMLSchema"
	set.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	set.Body.ObterUsuarios.XMLNamespace = "http://webservice.auth.app.bsbr.altec.com/"
/webservice.auth.app.bsbr.altec.com/"
	set.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	set.Xsd = "http://www.w3.org/2001/XMLSchema"
	set.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	set.Body.ObterUsuarios.XMLNamespace = "http://webservice.auth.app.bsbr.altec.com/"
nstance"
	set.Xsd = "http://www.w3.org/2001/XMLSchema"
	set.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	set.Body.ObterUsuarios.XMLNamespace = "http://webservice.auth.app.bsbr.altec.com/"


var (
	xmlTest = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <obterUsuarios xmlns="http://webservice.auth.app.bsbr.altec.com/">
      <param xmlns="">
        <agencia>1111</agencia>
        <banco>0033</banco>
        <conta>111111111</conta>
        <IMEI>0102030405</IMEI>
        <linha>983157131</linha>
      </param>
    </obterUsuarios>
  </soap:Body>
</soap:Envelope>
`
)

func main() {
	fmt.Println("lendo soap...")
	
	se := SoapEnvelope{}
	se.Body = SoapBody{}

	err := xml.Unmarshal([]byte(xmlTest), &se)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", se.Body.ObterUsuarios.Param.Agencia)

	fmt.Println("gerando soap...")
	
	set := SoapEnvelope{}
	set.Body.ObterUsuarios.Param.Agencia = "1111"
	set.Body.ObterUsuarios.Param.Banco = "0033"
	set.Body.ObterUsuarios.Param.Conta = "000111111111"
	set.Body.ObterUsuarios.Param.Imei = "01020304"
	set.Body.ObterUsuarios.Param.Linha = "983157131"

	b, err := xml.MarshalIndent(&set, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(`<?xml version="1.0" encoding="utf-8"?>` + string(b))

}