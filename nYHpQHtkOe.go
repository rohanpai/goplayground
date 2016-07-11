package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
)

type SoapEnvelope struct {
	XMLName xml.Name `xml:&#34;http://schemas.xmlsoap.org/soap/envelope/ Envelope&#34;`
	Xsi     string   `xml:&#34;xmlns xsi,attr&#34;`
	Xsd     string   `xml:&#34;xmlns xsd,attr&#34;`
	Soap    string   `xml:&#34;xmlns soap,attr&#34;`
	Body    SoapBody
}

type SoapBody struct {
	XMLName xml.Name `xml:&#34;http://schemas.xmlsoap.org/soap/envelope/ Body&#34;`
	ObterUsuarios ObterUsuarios
}

type ObterUsuarios struct {
	XMLName      xml.Name   `xml:&#34;obterUsuarios&#34;`
	Param        UsuarioDto `xml:&#34;param&#34;`
}

type UsuarioDto struct {
	Agencia      string `xml:&#34;agencia&#34;`
	Conta        string `xml:&#34;conta&#34;`
	Banco        string `xml:&#34;banco&#34;`
	Imei         string `xml:&#34;IMEI&#34;`
	Linha        string `xml:&#34;linha&#34;`
}	set.Xsi = &#34;http://www.w3.org/2001/XMLSchema-i	set.Xsi = &#34;http://www.w3.org/2001/XMLSchema-instance&#34;
	set.Xsd = &#34;http://www.w3.org/2001/XMLSchema&#34;
	set.Soap = &#34;http://schemas.xmlsoap.org/soap/envelope/&#34;
	set.Body.ObterUsuarios.XMLNamespace = &#34;http:/	set.Xsi = &#34;http://www.w3.org/2001/XMLSchema-instance&#34;
	set.Xsd = &#34;http://www.w3.org/2001/XMLSchema&#34;
	set.Soap = &#34;http://schemas.xmlsoap.org/soap/envelope/&#34;
	set.Body.ObterUsuarios.XMLNamespace = &#34;http://webservice.auth.app.bsbr.altec.com/&#34;
/webservice.auth.app.bsbr.altec.com/&#34;
	set.Xsi = &#34;http://www.w3.org/2001/XMLSchema-instance&#34;
	set.Xsd = &#34;http://www.w3.org/2001/XMLSchema&#34;
	set.Soap = &#34;http://schemas.xmlsoap.org/soap/envelope/&#34;
	set.Body.ObterUsuarios.XMLNamespace = &#34;http://webservice.auth.app.bsbr.altec.com/&#34;
nstance&#34;
	set.Xsd = &#34;http://www.w3.org/2001/XMLSchema&#34;
	set.Soap = &#34;http://schemas.xmlsoap.org/soap/envelope/&#34;
	set.Body.ObterUsuarios.XMLNamespace = &#34;http://webservice.auth.app.bsbr.altec.com/&#34;


var (
	xmlTest = `&lt;?xml version=&#34;1.0&#34; encoding=&#34;utf-8&#34;?&gt;
&lt;soap:Envelope xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xmlns:xsd=&#34;http://www.w3.org/2001/XMLSchema&#34; xmlns:soap=&#34;http://schemas.xmlsoap.org/soap/envelope/&#34;&gt;
  &lt;soap:Body&gt;
    &lt;obterUsuarios xmlns=&#34;http://webservice.auth.app.bsbr.altec.com/&#34;&gt;
      &lt;param xmlns=&#34;&#34;&gt;
        &lt;agencia&gt;1111&lt;/agencia&gt;
        &lt;banco&gt;0033&lt;/banco&gt;
        &lt;conta&gt;111111111&lt;/conta&gt;
        &lt;IMEI&gt;0102030405&lt;/IMEI&gt;
        &lt;linha&gt;983157131&lt;/linha&gt;
      &lt;/param&gt;
    &lt;/obterUsuarios&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;
`
)

func main() {
	fmt.Println(&#34;lendo soap...&#34;)
	
	se := SoapEnvelope{}
	se.Body = SoapBody{}

	err := xml.Unmarshal([]byte(xmlTest), &amp;se)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(&#34;%s\n&#34;, se.Body.ObterUsuarios.Param.Agencia)

	fmt.Println(&#34;gerando soap...&#34;)
	
	set := SoapEnvelope{}
	set.Body.ObterUsuarios.Param.Agencia = &#34;1111&#34;
	set.Body.ObterUsuarios.Param.Banco = &#34;0033&#34;
	set.Body.ObterUsuarios.Param.Conta = &#34;000111111111&#34;
	set.Body.ObterUsuarios.Param.Imei = &#34;01020304&#34;
	set.Body.ObterUsuarios.Param.Linha = &#34;983157131&#34;

	b, err := xml.MarshalIndent(&amp;set, &#34;&#34;, &#34;   &#34;)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(`&lt;?xml version=&#34;1.0&#34; encoding=&#34;utf-8&#34;?&gt;` &#43; string(b))

}