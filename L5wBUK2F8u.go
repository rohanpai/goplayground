package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
)

type RequestStruct struct {
	XMLName xml.Name `xml:&#34;request&#34;`
}

func main() {
	var xmlstring string = `&lt;request signature=&#34;1c7f34333d4d-6904c8647b5f-4ec7aee89801-02bc&#34;&gt;
&lt;xs:schema xmlns:xs=&#34;http://www.w3.org/2001/XMLSchema&#34;&gt;
	&lt;xs:element name=&#34;methodCall&#34;&gt;
		&lt;xs:complexType&gt;
			&lt;xs:sequence&gt;
				&lt;xs:element name=&#34;SessionState&#34; type=&#34;xs:string&#34; minOccurs=&#34;0&#34;/&gt;
				&lt;xs:element name=&#34;trace&#34; type=&#34;xs:boolean&#34; minOccurs=&#34;0&#34;/&gt;
				&lt;xs:element name=&#34;params&#34;&gt;
					&lt;xs:complexType&gt;
						&lt;xs:sequence&gt;
							&lt;xs:element name=&#34;actorID&#34; type=&#34;xs:string&#34; minOccurs=&#34;1&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;verb&#34; type=&#34;xs:string&#34; minOccurs=&#34;1&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[Activity type verb that describes the activity such as post.]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;content&#34; type=&#34;xs:string&#34; minOccurs=&#34;1&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;extra&#34; type=&#34;xs:string&#34; minOccurs=&#34;0&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;objectURI&#34; type=&#34;xs:string&#34; minOccurs=&#34;0&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;activityStreamID&#34; type=&#34;xs:string&#34; minOccurs=&#34;1&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;title&#34; type=&#34;xs:string&#34; minOccurs=&#34;0&#34; maxOccurs=&#34;1&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;visibility&#34; type=&#34;activityVisibilityType&#34; minOccurs=&#34;0&#34; default=&#34;friend&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;location&#34; type=&#34;geoLocationType&#34; minOccurs=&#34;0&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;type&#34; type=&#34;xs:string&#34; minOccurs=&#34;0&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;objectRefUrn&#34; type=&#34;xs:anyURI&#34; minOccurs=&#34;0&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
							&lt;xs:element name=&#34;activityType&#34; type=&#34;activityType&#34; minOccurs=&#34;0&#34; &gt;
								&lt;xs:annotation&gt;
									&lt;xs:documentation&gt;
										&lt;![CDATA[XX]]&gt;
									&lt;/xs:documentation&gt;
								&lt;/xs:annotation&gt;
							&lt;/xs:element&gt;
						&lt;/xs:sequence&gt;
					&lt;/xs:complexType&gt;
				&lt;/xs:element&gt;
			&lt;/xs:sequence&gt;
			&lt;xs:attribute name=&#34;service&#34; type=&#34;xs:string&#34; use=&#34;required&#34;/&gt;
			&lt;xs:attribute name=&#34;method&#34; type=&#34;xs:string&#34; use=&#34;required&#34;/&gt;
		&lt;/xs:complexType&gt;
	&lt;/xs:element&gt;
&lt;/xs:schema&gt;	
&lt;/request&gt;`

	response := RequestStruct{}
	xml.Unmarshal([]byte(xmlstring), &amp;response)
	fmt.Println(response)

}
