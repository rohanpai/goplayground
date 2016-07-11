package main

import &#34;fmt&#34;
import &#34;encoding/xml&#34;

type MyRespEnvelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName     xml.Name
	GetResponse completeResponse `xml:&#34;activationPack_completeResponse&#34;`
}

type completeResponse struct {
	XMLName xml.Name `xml:&#34;activationPack_completeResponse&#34;`
	Id      string   `xml:&#34;Id,attr&#34;`
	MyVar   string   `xml:&#34;activationPack_completeResult&#34;`
}

func main() {

	Soap := []byte(`&lt;?xml version=&#34;1.0&#34; encoding=&#34;UTF-8&#34;?&gt;
&lt;soap:Envelope SOAP-ENV:encodingStyle=&#34;http://schemas.xmlsoap.org/soap/encoding/&#34; xmlns:SOAP-ENV=&#34;http://schemas.xmlsoap.org/soap/envelope/&#34; xmlns:xsd=&#34;http://www.w3.org/2001/XMLSchema&#34; xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xmlns:SOAP-ENC=&#34;http://schemas.xmlsoap.org/soap/encoding/&#34;&gt;
&lt;soap:Body&gt;
&lt;activationPack_completeResponse Id=&#34;http://tempuri.org/&#34;&gt;
&lt;activationPack_completeResult xsi:type=&#34;xsd:string&#34;&gt;Active&lt;/activationPack_completeResult&gt;
&lt;/activationPack_completeResponse&gt;
&lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;`)

	res := &amp;MyRespEnvelope{}
	err := xml.Unmarshal(Soap, res)

	fmt.Println(res.Body, err)
}
