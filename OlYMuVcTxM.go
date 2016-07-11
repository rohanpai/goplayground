package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Nvd struct {
	Entry
}

type Entry struct {
	VulnerableConfiguration
	VulnerableSoftwareList
	VulnCveId                string `xml:"entry>cve-id"`
	VulnPublishedDatTime     string `xml:"entry>published-datetime"`
	VulnLastModifiedDateTime string `xml:"entry>last-modified-datetime"`
	Cvss
	VulnCweId string `xml:"entry>cwe>id,attr` // BUG won't parse
	VulnReferences
	VulnSummary string `xml:"entry>summary"`
}

type VulnerableConfiguration struct {
	CpeLangLogicalTest
}

// BUG this isn't being parsed 
type CpeLangLogicalTest struct {
	CpeLangFactRef []string `xml:entry>vulnerable-configuration>fact-ref,attr`
}

type VulnerableSoftwareList struct {
	Product []string `xml:"entry>vulnerable-software-list>product"`
}

type Cvss struct {
	CvssBaseMetrics
}

type CvssBaseMetrics struct {
	Score                 string `xml:"entry>cvss>base_metrics>score"`
	AccessVector          string `xml:"entry>cvss>base_metrics>access-vector"`
	AccessComplexity      string `xml:"entry>cvss>base_metrics>access-complexity"`
	Authentication        string `xml:"entry>cvss>base_metrics>authentication"`
	ConfidentialityImpact string `xml:"entry>cvss>base_metrics>confidentiality-impact"`
	IntegrityImpact       string `xml:"entry>cvss>base_metrics>integrity-impact"`
	AvailabilityImpact    string `xml:"entry>cvss>base_metrics>availability-impact"`
	Source                string `xml:"entry>cvss>base_metrics>source"`
	GeneratedOnDateTime   string `xml:"entry>cvss>base_metrics>generated-on-datetime"`
}

type VulnReferences struct {
	Source    string `xml:"entry>references>source"`
	Reference string `xml:"entry>references>reference"`
}

func main() {
		nvdXml, err := ioutil.ReadFile("nvdtruncated.xml")
		//nvdXml, err := ioutil.ReadFile("nvdcve-2.0-modified.xml")

		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		var v Nvd
		err = xml.Unmarshal([]byte(nvdXml), &v)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		fmt.Println(v.Entry)
}


/* SAMPLE DATA HERE  "nvdtruncated.xml"

<?xml version='1.0' encoding='UTF-8'?>
<nvd xmlns:cvss="http://scap.nist.gov/schema/cvss-v2/0.2" xmlns="http://scap.nist.gov/schema/feed/vulnerability/2.0" xmlns:vuln="http://scap.nist.gov/schema/vulnerability/0.4" xmlns:scap-core="http://scap.nist.gov/schema/scap-core/0.1" xmlns:patch="http://scap.nist.gov/schema/patch/0.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cpe-lang="http://cpe.mitre.org/language/2.0" nvd_xml_version="2.0" pub_date="2012-08-15T16:01:01" xsi:schemaLocation="http://scap.nist.gov/schema/patch/0.1 http://nvd.nist.gov/schema/patch_0.1.xsd http://scap.nist.gov/schema/scap-core/0.1 http://nvd.nist.gov/schema/scap-core_0.1.xsd http://scap.nist.gov/schema/feed/vulnerability/2.0 http://nvd.nist.gov/schema/nvd-cve-feed_2.0.xsd">
  <entry id="CVE-2005-4895">
    <vuln:vulnerable-configuration id="http://nvd.nist.gov/">
      <cpe-lang:logical-test negate="false" operator="OR">
        <cpe-lang:fact-ref name="cpe:/a:csilvers:gperftools:0.3" />
        <cpe-lang:fact-ref name="cpe:/a:csilvers:gperftools:0.2" />
        <cpe-lang:fact-ref name="cpe:/a:csilvers:gperftools:0.1" />
      </cpe-lang:logical-test>
    </vuln:vulnerable-configuration>
    <vuln:vulnerable-software-list>
      <vuln:product>cpe:/a:csilvers:gperftools:0.3</vuln:product>
      <vuln:product>cpe:/a:csilvers:gperftools:0.1</vuln:product>
      <vuln:product>cpe:/a:csilvers:gperftools:0.2</vuln:product>
    </vuln:vulnerable-software-list>
    <vuln:cve-id>CVE-2005-4895</vuln:cve-id>
    <vuln:published-datetime>2012-07-25T15:55:01.273-04:00</vuln:published-datetime>
    <vuln:last-modified-datetime>2012-08-09T00:00:00.000-04:00</vuln:last-modified-datetime>
    <vuln:cvss>
      <cvss:base_metrics>
        <cvss:score>5.0</cvss:score>
        <cvss:access-vector>NETWORK</cvss:access-vector>
        <cvss:access-complexity>LOW</cvss:access-complexity>
        <cvss:authentication>NONE</cvss:authentication>
        <cvss:confidentiality-impact>NONE</cvss:confidentiality-impact>
        <cvss:integrity-impact>NONE</cvss:integrity-impact>
        <cvss:availability-impact>PARTIAL</cvss:availability-impact>
        <cvss:source>http://nvd.nist.gov</cvss:source>
        <cvss:generated-on-datetime>2012-07-26T08:38:00.000-04:00</cvss:generated-on-datetime>
      </cvss:base_metrics>
    </vuln:cvss>
    <vuln:cwe id="CWE-189" />
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>MISC</vuln:source>
      <vuln:reference href="http://kqueue.org/blog/2012/03/05/memory-allocator-security-revisited/" xml:lang="en">http://kqueue.org/blog/2012/03/05/memory-allocator-security-revisited/</vuln:reference>
    </vuln:references>    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>CONFIRM</vuln:source>
      <vuln:reference href="http://code.google.com/p/gperftools/source/browse/tags/perftools-0.4/ChangeLog" xml:lang="en">http://code.google.com/p/gperftools/source/browse/tags/perftools-0.4/ChangeLog</vuln:reference>
    </vuln:references>
    <vuln:summary>Multiple integer overflows in TCMalloc (tcmalloc.cc) in gperftools before 0.4 make it easier for context-dependent attackers to perform memory-related attacks such as buffer overflows via a large size value, which causes less memory to be allocated than expected.</vuln:summary>
  </entry>
  <entry id="CVE-2006-4330">
    <vuln:vulnerable-configuration id="http://nvd.nist.gov/">
      <cpe-lang:logical-test negate="false" operator="OR">
        <cpe-lang:fact-ref name="cpe:/a:wireshark:wireshark:0.99.2" />
      </cpe-lang:logical-test>
    </vuln:vulnerable-configuration>
    <vuln:vulnerable-software-list>
      <vuln:product>cpe:/a:wireshark:wireshark:0.99.2</vuln:product>
    </vuln:vulnerable-software-list>
    <vuln:cve-id>CVE-2006-4330</vuln:cve-id>
    <vuln:published-datetime>2006-08-24T16:04:00.000-04:00</vuln:published-datetime>
    <vuln:last-modified-datetime>2012-08-13T21:59:57.393-04:00</vuln:last-modified-datetime>
    <vuln:cvss>
      <cvss:base_metrics>
        <cvss:score>4.3</cvss:score>
        <cvss:access-vector approximated="true">NETWORK</cvss:access-vector>
        <cvss:access-complexity>MEDIUM</cvss:access-complexity>
        <cvss:authentication>NONE</cvss:authentication>
        <cvss:confidentiality-impact>NONE</cvss:confidentiality-impact>
        <cvss:integrity-impact>NONE</cvss:integrity-impact>
        <cvss:availability-impact>PARTIAL</cvss:availability-impact>
        <cvss:source>http://nvd.nist.gov</cvss:source>
        <cvss:generated-on-datetime>2006-08-26T17:11:00.000-04:00</cvss:generated-on-datetime>
      </cvss:base_metrics>
    </vuln:cvss>
    <vuln:assessment_check name="oval:org.mitre.oval:def:9869" href="http://oval.mitre.org/repository/data/getDef?id=oval:org.mitre.oval:def:9869" system="http://oval.mitre.org/XMLSchema/oval-definitions-5" />
    <vuln:assessment_check name="oval:org.mitre.oval:def:14684" href="http://oval.mitre.org/repository/data/getDef?id=oval:org.mitre.oval:def:14684" system="http://oval.mitre.org/XMLSchema/oval-definitions-5" />
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>CERT-VN</vuln:source>
      <vuln:reference href="http://www.kb.cert.org/vuls/id/808832" xml:lang="en">VU#808832</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>CONFIRM</vuln:source>
      <vuln:reference href="http://www.wireshark.org/security/wnpa-sec-2006-02.html" xml:lang="en">http://www.wireshark.org/security/wnpa-sec-2006-02.html</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="PATCH">
      <vuln:source>BID</vuln:source>
      <vuln:reference href="http://www.securityfocus.com/bid/19690" xml:lang="en">19690</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="PATCH">
      <vuln:source>SECTRACK</vuln:source>
      <vuln:reference href="http://securitytracker.com/id?1016736" xml:lang="en">1016736</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/21597" xml:lang="en">21597</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>CONFIRM</vuln:source>
      <vuln:reference href="https://issues.rpath.com/browse/RPL-597" xml:lang="en">https://issues.rpath.com/browse/RPL-597</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>XF</vuln:source>
      <vuln:reference href="http://xforce.iss.net/xforce/xfdb/28553" xml:lang="en">wireshark-esp-offbyone(28553)</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>XF</vuln:source>
      <vuln:reference href="http://xforce.iss.net/xforce/xfdb/28550" xml:lang="en">wireshark-scsi-dos(28550)</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>VUPEN</vuln:source>
      <vuln:reference href="http://www.vupen.com/english/advisories/2006/3370" xml:lang="en">ADV-2006-3370</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>BUGTRAQ</vuln:source>
      <vuln:reference href="http://www.securityfocus.com/archive/1/archive/1/444323/100/0/threaded" xml:lang="en">20060825 rPSA-2006-0158-1 tshark wireshark</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>REDHAT</vuln:source>
      <vuln:reference href="http://www.redhat.com/support/errata/RHSA-2006-0658.html" xml:lang="en">RHSA-2006:0658</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>MANDRIVA</vuln:source>
      <vuln:reference href="http://www.mandriva.com/security/advisories?name=MDKSA-2006:152" xml:lang="en">MDKSA-2006:152</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>CONFIRM</vuln:source>
      <vuln:reference href="http://support.avaya.com/elmodocs2/security/ASA-2006-227.htm" xml:lang="en">http://support.avaya.com/elmodocs2/security/ASA-2006-227.htm</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>GENTOO</vuln:source>
      <vuln:reference href="http://security.gentoo.org/glsa/glsa-200608-26.xml" xml:lang="en">GLSA-200608-26</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/22378" xml:lang="en">22378</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/21885" xml:lang="en">21885</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/21682" xml:lang="en">21682</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="UNKNOWN">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/21649" xml:lang="en">21649</vuln:reference>
    </vuln:references>
    <vuln:references xml:lang="en" reference_type="VENDOR_ADVISORY">
      <vuln:source>SECUNIA</vuln:source>
      <vuln:reference href="http://secunia.com/advisories/21619" xml:lang="en">21619</vuln:reference>
    </vuln:references>
    <vuln:scanner>
      <vuln:definition name="oval:org.mitre.oval:def:9869" href="http://oval.mitre.org/repository/data/DownloadDefinition?id=oval:org.mitre.oval:def:9869" system="http://oval.mitre.org/XMLSchema/oval-definitions-5" />
    </vuln:scanner>
    <vuln:summary>Unspecified vulnerability in the SCSI dissector in Wireshark (formerly Ethereal) 0.99.2 allows remote attackers to cause a denial of service (crash) via unspecified vectors.</vuln:summary>
  </entry>

</nvd>



*/