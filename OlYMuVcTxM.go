package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;io/ioutil&#34;
)

type Nvd struct {
	Entry
}

type Entry struct {
	VulnerableConfiguration
	VulnerableSoftwareList
	VulnCveId                string `xml:&#34;entry&gt;cve-id&#34;`
	VulnPublishedDatTime     string `xml:&#34;entry&gt;published-datetime&#34;`
	VulnLastModifiedDateTime string `xml:&#34;entry&gt;last-modified-datetime&#34;`
	Cvss
	VulnCweId string `xml:&#34;entry&gt;cwe&gt;id,attr` // BUG won&#39;t parse
	VulnReferences
	VulnSummary string `xml:&#34;entry&gt;summary&#34;`
}

type VulnerableConfiguration struct {
	CpeLangLogicalTest
}

// BUG this isn&#39;t being parsed 
type CpeLangLogicalTest struct {
	CpeLangFactRef []string `xml:entry&gt;vulnerable-configuration&gt;fact-ref,attr`
}

type VulnerableSoftwareList struct {
	Product []string `xml:&#34;entry&gt;vulnerable-software-list&gt;product&#34;`
}

type Cvss struct {
	CvssBaseMetrics
}

type CvssBaseMetrics struct {
	Score                 string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;score&#34;`
	AccessVector          string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;access-vector&#34;`
	AccessComplexity      string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;access-complexity&#34;`
	Authentication        string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;authentication&#34;`
	ConfidentialityImpact string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;confidentiality-impact&#34;`
	IntegrityImpact       string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;integrity-impact&#34;`
	AvailabilityImpact    string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;availability-impact&#34;`
	Source                string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;source&#34;`
	GeneratedOnDateTime   string `xml:&#34;entry&gt;cvss&gt;base_metrics&gt;generated-on-datetime&#34;`
}

type VulnReferences struct {
	Source    string `xml:&#34;entry&gt;references&gt;source&#34;`
	Reference string `xml:&#34;entry&gt;references&gt;reference&#34;`
}

func main() {
		nvdXml, err := ioutil.ReadFile(&#34;nvdtruncated.xml&#34;)
		//nvdXml, err := ioutil.ReadFile(&#34;nvdcve-2.0-modified.xml&#34;)

		if err != nil {
			fmt.Println(&#34;Error opening file:&#34;, err)
			return
		}
		var v Nvd
		err = xml.Unmarshal([]byte(nvdXml), &amp;v)
		if err != nil {
			fmt.Printf(&#34;error: %v&#34;, err)
			return
		}

		fmt.Println(v.Entry)
}


/* SAMPLE DATA HERE  &#34;nvdtruncated.xml&#34;

&lt;?xml version=&#39;1.0&#39; encoding=&#39;UTF-8&#39;?&gt;
&lt;nvd xmlns:cvss=&#34;http://scap.nist.gov/schema/cvss-v2/0.2&#34; xmlns=&#34;http://scap.nist.gov/schema/feed/vulnerability/2.0&#34; xmlns:vuln=&#34;http://scap.nist.gov/schema/vulnerability/0.4&#34; xmlns:scap-core=&#34;http://scap.nist.gov/schema/scap-core/0.1&#34; xmlns:patch=&#34;http://scap.nist.gov/schema/patch/0.1&#34; xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xmlns:cpe-lang=&#34;http://cpe.mitre.org/language/2.0&#34; nvd_xml_version=&#34;2.0&#34; pub_date=&#34;2012-08-15T16:01:01&#34; xsi:schemaLocation=&#34;http://scap.nist.gov/schema/patch/0.1 http://nvd.nist.gov/schema/patch_0.1.xsd http://scap.nist.gov/schema/scap-core/0.1 http://nvd.nist.gov/schema/scap-core_0.1.xsd http://scap.nist.gov/schema/feed/vulnerability/2.0 http://nvd.nist.gov/schema/nvd-cve-feed_2.0.xsd&#34;&gt;
  &lt;entry id=&#34;CVE-2005-4895&#34;&gt;
    &lt;vuln:vulnerable-configuration id=&#34;http://nvd.nist.gov/&#34;&gt;
      &lt;cpe-lang:logical-test negate=&#34;false&#34; operator=&#34;OR&#34;&gt;
        &lt;cpe-lang:fact-ref name=&#34;cpe:/a:csilvers:gperftools:0.3&#34; /&gt;
        &lt;cpe-lang:fact-ref name=&#34;cpe:/a:csilvers:gperftools:0.2&#34; /&gt;
        &lt;cpe-lang:fact-ref name=&#34;cpe:/a:csilvers:gperftools:0.1&#34; /&gt;
      &lt;/cpe-lang:logical-test&gt;
    &lt;/vuln:vulnerable-configuration&gt;
    &lt;vuln:vulnerable-software-list&gt;
      &lt;vuln:product&gt;cpe:/a:csilvers:gperftools:0.3&lt;/vuln:product&gt;
      &lt;vuln:product&gt;cpe:/a:csilvers:gperftools:0.1&lt;/vuln:product&gt;
      &lt;vuln:product&gt;cpe:/a:csilvers:gperftools:0.2&lt;/vuln:product&gt;
    &lt;/vuln:vulnerable-software-list&gt;
    &lt;vuln:cve-id&gt;CVE-2005-4895&lt;/vuln:cve-id&gt;
    &lt;vuln:published-datetime&gt;2012-07-25T15:55:01.273-04:00&lt;/vuln:published-datetime&gt;
    &lt;vuln:last-modified-datetime&gt;2012-08-09T00:00:00.000-04:00&lt;/vuln:last-modified-datetime&gt;
    &lt;vuln:cvss&gt;
      &lt;cvss:base_metrics&gt;
        &lt;cvss:score&gt;5.0&lt;/cvss:score&gt;
        &lt;cvss:access-vector&gt;NETWORK&lt;/cvss:access-vector&gt;
        &lt;cvss:access-complexity&gt;LOW&lt;/cvss:access-complexity&gt;
        &lt;cvss:authentication&gt;NONE&lt;/cvss:authentication&gt;
        &lt;cvss:confidentiality-impact&gt;NONE&lt;/cvss:confidentiality-impact&gt;
        &lt;cvss:integrity-impact&gt;NONE&lt;/cvss:integrity-impact&gt;
        &lt;cvss:availability-impact&gt;PARTIAL&lt;/cvss:availability-impact&gt;
        &lt;cvss:source&gt;http://nvd.nist.gov&lt;/cvss:source&gt;
        &lt;cvss:generated-on-datetime&gt;2012-07-26T08:38:00.000-04:00&lt;/cvss:generated-on-datetime&gt;
      &lt;/cvss:base_metrics&gt;
    &lt;/vuln:cvss&gt;
    &lt;vuln:cwe id=&#34;CWE-189&#34; /&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;MISC&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://kqueue.org/blog/2012/03/05/memory-allocator-security-revisited/&#34; xml:lang=&#34;en&#34;&gt;http://kqueue.org/blog/2012/03/05/memory-allocator-security-revisited/&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;CONFIRM&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://code.google.com/p/gperftools/source/browse/tags/perftools-0.4/ChangeLog&#34; xml:lang=&#34;en&#34;&gt;http://code.google.com/p/gperftools/source/browse/tags/perftools-0.4/ChangeLog&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:summary&gt;Multiple integer overflows in TCMalloc (tcmalloc.cc) in gperftools before 0.4 make it easier for context-dependent attackers to perform memory-related attacks such as buffer overflows via a large size value, which causes less memory to be allocated than expected.&lt;/vuln:summary&gt;
  &lt;/entry&gt;
  &lt;entry id=&#34;CVE-2006-4330&#34;&gt;
    &lt;vuln:vulnerable-configuration id=&#34;http://nvd.nist.gov/&#34;&gt;
      &lt;cpe-lang:logical-test negate=&#34;false&#34; operator=&#34;OR&#34;&gt;
        &lt;cpe-lang:fact-ref name=&#34;cpe:/a:wireshark:wireshark:0.99.2&#34; /&gt;
      &lt;/cpe-lang:logical-test&gt;
    &lt;/vuln:vulnerable-configuration&gt;
    &lt;vuln:vulnerable-software-list&gt;
      &lt;vuln:product&gt;cpe:/a:wireshark:wireshark:0.99.2&lt;/vuln:product&gt;
    &lt;/vuln:vulnerable-software-list&gt;
    &lt;vuln:cve-id&gt;CVE-2006-4330&lt;/vuln:cve-id&gt;
    &lt;vuln:published-datetime&gt;2006-08-24T16:04:00.000-04:00&lt;/vuln:published-datetime&gt;
    &lt;vuln:last-modified-datetime&gt;2012-08-13T21:59:57.393-04:00&lt;/vuln:last-modified-datetime&gt;
    &lt;vuln:cvss&gt;
      &lt;cvss:base_metrics&gt;
        &lt;cvss:score&gt;4.3&lt;/cvss:score&gt;
        &lt;cvss:access-vector approximated=&#34;true&#34;&gt;NETWORK&lt;/cvss:access-vector&gt;
        &lt;cvss:access-complexity&gt;MEDIUM&lt;/cvss:access-complexity&gt;
        &lt;cvss:authentication&gt;NONE&lt;/cvss:authentication&gt;
        &lt;cvss:confidentiality-impact&gt;NONE&lt;/cvss:confidentiality-impact&gt;
        &lt;cvss:integrity-impact&gt;NONE&lt;/cvss:integrity-impact&gt;
        &lt;cvss:availability-impact&gt;PARTIAL&lt;/cvss:availability-impact&gt;
        &lt;cvss:source&gt;http://nvd.nist.gov&lt;/cvss:source&gt;
        &lt;cvss:generated-on-datetime&gt;2006-08-26T17:11:00.000-04:00&lt;/cvss:generated-on-datetime&gt;
      &lt;/cvss:base_metrics&gt;
    &lt;/vuln:cvss&gt;
    &lt;vuln:assessment_check name=&#34;oval:org.mitre.oval:def:9869&#34; href=&#34;http://oval.mitre.org/repository/data/getDef?id=oval:org.mitre.oval:def:9869&#34; system=&#34;http://oval.mitre.org/XMLSchema/oval-definitions-5&#34; /&gt;
    &lt;vuln:assessment_check name=&#34;oval:org.mitre.oval:def:14684&#34; href=&#34;http://oval.mitre.org/repository/data/getDef?id=oval:org.mitre.oval:def:14684&#34; system=&#34;http://oval.mitre.org/XMLSchema/oval-definitions-5&#34; /&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;CERT-VN&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.kb.cert.org/vuls/id/808832&#34; xml:lang=&#34;en&#34;&gt;VU#808832&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;CONFIRM&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.wireshark.org/security/wnpa-sec-2006-02.html&#34; xml:lang=&#34;en&#34;&gt;http://www.wireshark.org/security/wnpa-sec-2006-02.html&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;PATCH&#34;&gt;
      &lt;vuln:source&gt;BID&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.securityfocus.com/bid/19690&#34; xml:lang=&#34;en&#34;&gt;19690&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;PATCH&#34;&gt;
      &lt;vuln:source&gt;SECTRACK&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://securitytracker.com/id?1016736&#34; xml:lang=&#34;en&#34;&gt;1016736&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/21597&#34; xml:lang=&#34;en&#34;&gt;21597&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;CONFIRM&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;https://issues.rpath.com/browse/RPL-597&#34; xml:lang=&#34;en&#34;&gt;https://issues.rpath.com/browse/RPL-597&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;XF&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://xforce.iss.net/xforce/xfdb/28553&#34; xml:lang=&#34;en&#34;&gt;wireshark-esp-offbyone(28553)&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;XF&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://xforce.iss.net/xforce/xfdb/28550&#34; xml:lang=&#34;en&#34;&gt;wireshark-scsi-dos(28550)&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;VUPEN&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.vupen.com/english/advisories/2006/3370&#34; xml:lang=&#34;en&#34;&gt;ADV-2006-3370&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;BUGTRAQ&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.securityfocus.com/archive/1/archive/1/444323/100/0/threaded&#34; xml:lang=&#34;en&#34;&gt;20060825 rPSA-2006-0158-1 tshark wireshark&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;REDHAT&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.redhat.com/support/errata/RHSA-2006-0658.html&#34; xml:lang=&#34;en&#34;&gt;RHSA-2006:0658&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;MANDRIVA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://www.mandriva.com/security/advisories?name=MDKSA-2006:152&#34; xml:lang=&#34;en&#34;&gt;MDKSA-2006:152&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;CONFIRM&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://support.avaya.com/elmodocs2/security/ASA-2006-227.htm&#34; xml:lang=&#34;en&#34;&gt;http://support.avaya.com/elmodocs2/security/ASA-2006-227.htm&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;GENTOO&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://security.gentoo.org/glsa/glsa-200608-26.xml&#34; xml:lang=&#34;en&#34;&gt;GLSA-200608-26&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/22378&#34; xml:lang=&#34;en&#34;&gt;22378&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/21885&#34; xml:lang=&#34;en&#34;&gt;21885&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/21682&#34; xml:lang=&#34;en&#34;&gt;21682&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;UNKNOWN&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/21649&#34; xml:lang=&#34;en&#34;&gt;21649&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:references xml:lang=&#34;en&#34; reference_type=&#34;VENDOR_ADVISORY&#34;&gt;
      &lt;vuln:source&gt;SECUNIA&lt;/vuln:source&gt;
      &lt;vuln:reference href=&#34;http://secunia.com/advisories/21619&#34; xml:lang=&#34;en&#34;&gt;21619&lt;/vuln:reference&gt;
    &lt;/vuln:references&gt;
    &lt;vuln:scanner&gt;
      &lt;vuln:definition name=&#34;oval:org.mitre.oval:def:9869&#34; href=&#34;http://oval.mitre.org/repository/data/DownloadDefinition?id=oval:org.mitre.oval:def:9869&#34; system=&#34;http://oval.mitre.org/XMLSchema/oval-definitions-5&#34; /&gt;
    &lt;/vuln:scanner&gt;
    &lt;vuln:summary&gt;Unspecified vulnerability in the SCSI dissector in Wireshark (formerly Ethereal) 0.99.2 allows remote attackers to cause a denial of service (crash) via unspecified vectors.&lt;/vuln:summary&gt;
  &lt;/entry&gt;

&lt;/nvd&gt;



*/