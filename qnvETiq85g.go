package main

import (
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;regexp&#34;
	&#34;strings&#34;

	//	&#34;golang.org/x/crypto/ssh/terminal&#34;
)

// Blacklist type is a map of Nodes with string keys
type Blacklist map[string]*Node

// Source type is a map of Srcs with string keys
type Source map[string]*Src

// Node configuration record
type Node struct {
	Disable          bool
	IP               string
	Include, Exclude []string
	Source
}

// Src record, struct for Source map
type Src struct {
	Disable bool
	Desc    string
	Prfx    string
	URL     string
}

// String returns pretty print for the Blacklist struct
func (b Blacklist) String() (result string) {
	//	cols, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	cols := 20
	for pkey := range b {
		result &#43;= fmt.Sprintf(&#34;Node: %v\n\tDisabled: %v\n&#34;, pkey, b[pkey].Disable)
		result &#43;= fmt.Sprintf(&#34;\tRedirect IP: %v\n\tExclude(s):\n&#34;, b[pkey].IP)
		for _, exclude := range b[pkey].Exclude {
			result &#43;= fmt.Sprintf(&#34;\t\t%v\n&#34;, exclude)
		}
		result &#43;= fmt.Sprintf(&#34;\tInclude(s):\n&#34;)
		for _, include := range b[pkey].Include {
			result &#43;= fmt.Sprintf(&#34;\t\t%v\n&#34;, include)
		}
		for skey, src := range b[pkey].Source {
			result &#43;= fmt.Sprintf(&#34;\tSource: %v\n\t\tDisabled: %v\n&#34;, skey, src.Disable)
			result &#43;= fmt.Sprintf(&#34;\t\tDescription: %v\n&#34;, b[pkey].Source[skey].Desc)
			result &#43;= fmt.Sprintf(&#34;\t\tPrefix: %v\n\t\tURL: %v\n&#34;, b[pkey].Source[skey].Prfx, b[pkey].Source[skey].URL)
		}
		result &#43;= fmt.Sprintln(strings.Repeat(&#34;-&#34;, cols/2))
	}
	return result
}

// ToBool converts a string (&#34;true&#34; or &#34;false&#34;) to it&#39;s boolean equivalent
func ToBool(s string) (b bool) {
	if len(s) == 0 {
		log.Fatal(&#34;ERROR: variable empty, cannot convert to boolean&#34;)
	}
	switch s {
	case &#34;false&#34;:
		b = false
	case &#34;true&#34;:
		b = true
	}
	return b
}

// Get extracts nodes from a EdgeOS/VyOS configuration structure
func Get(cfg string) {
	type re struct {
		brkt, cmnt, desc, dsbl, leaf, misc, mlti, mpty, name, node *regexp.Regexp
	}

	rx := &amp;re{}
	rx.brkt = regexp.MustCompile(`[}]`)
	rx.cmnt = regexp.MustCompile(`^([\/*]&#43;).*([*\/]&#43;)$`)
	rx.desc = regexp.MustCompile(`^(?:description)&#43;\s&#34;?([^&#34;]&#43;)?&#34;?$`)
	rx.dsbl = regexp.MustCompile(`^(disabled)&#43;\s([\S]&#43;)$`)
	rx.leaf = regexp.MustCompile(`^(source)&#43;\s([\S]&#43;)\s[{]{1}$`)
	rx.misc = regexp.MustCompile(`^([\w-]&#43;)$`)
	rx.mlti = regexp.MustCompile(`^((?:include|exclude)&#43;)\s([\S]&#43;)$`)
	rx.mpty = regexp.MustCompile(`^$`)
	rx.name = regexp.MustCompile(`^([\w-]&#43;)\s([\S]&#43;)$`)
	rx.node = regexp.MustCompile(`^([\w-]&#43;)\s[{]{1}$`)

	cfgtree := make(map[string]*Node)

	var tnode string
	var leaf string

	for _, line := range strings.Split(cfg, &#34;\n&#34;) {
		line = strings.TrimSpace(line)
		switch {
		case rx.mlti.MatchString(line):
			{
				IncExc := rx.mlti.FindStringSubmatch(line)
				switch IncExc[1] {
				case &#34;exclude&#34;:
					cfgtree[tnode].Exclude = append(cfgtree[tnode].Exclude, IncExc[2])
				case &#34;include&#34;:
					cfgtree[tnode].Include = append(cfgtree[tnode].Include, IncExc[2])
				}
			}
		case rx.node.MatchString(line):
			{
				node := rx.node.FindStringSubmatch(line)
				tnode = node[1]
				cfgtree[tnode] = &amp;Node{}
				cfgtree[tnode].Source = make(map[string]*Src)
			}
		case rx.leaf.MatchString(line):
			src := rx.leaf.FindStringSubmatch(line)
			leaf = src[2]

			if src[1] == &#34;source&#34; {
				cfgtree[tnode].Source[leaf] = &amp;Src{}
			}
		case rx.dsbl.MatchString(line):
			{
				disabled := rx.dsbl.FindStringSubmatch(line)
				cfgtree[tnode].Disable = ToBool(disabled[1])
			}
		case rx.name.MatchString(line):
			{
				name := rx.name.FindStringSubmatch(line)
				switch name[1] {
				case &#34;prefix&#34;:
					cfgtree[tnode].Source[leaf].Prfx = name[2]
				case &#34;url&#34;:
					cfgtree[tnode].Source[leaf].URL = name[2]
				case &#34;description&#34;:
					cfgtree[tnode].Source[leaf].Desc = name[2]
				case &#34;dns-redirect-ip&#34;:
					cfgtree[tnode].IP = name[2]
				}
			}
		case rx.desc.MatchString(line) || rx.cmnt.MatchString(line) || rx.misc.MatchString(line):
			break
		}
		// fmt.Printf(&#34;%s\n&#34;, line)
	}
	fmt.Println(cfgtree)
}

func main() {
	cfgtree := make(Blacklist)
	for _, k := range []string{&#34;root&#34;, &#34;hosts&#34;, &#34;domains&#34;} {
		cfgtree[k] = &amp;Node{}
		cfgtree[k].Source = make(Source)
	}
	cfgtree[&#34;hosts&#34;].Exclude = append(cfgtree[&#34;hosts&#34;].Exclude, &#34;rackcdn.com&#34;, &#34;schema.org&#34;)
	cfgtree[&#34;hosts&#34;].Include = append(cfgtree[&#34;hosts&#34;].Include, &#34;msdn.com&#34;, &#34;badgits.org&#34;)
	cfgtree[&#34;hosts&#34;].IP = &#34;192.168.168.1&#34;
	cfgtree[&#34;hosts&#34;].Source[&#34;hpHosts&#34;] = &amp;Src{URL: &#34;http://www.bonzon.com&#34;, Prfx: &#34;127.0.0.0&#34;}
	fmt.Println(cfgtree)
	fmt.Println(cfgtree[&#34;hosts&#34;])
	Get(testdata)
}

var testdata = `blacklist {
			disabled false
			dns-redirect-ip 0.0.0.0
			domains {
					include adsrvr.org
					include adtechus.net
					include advertising.com
					include centade.com
					include doubleclick.net
					include free-counter.co.uk
					include intellitxt.com
					include kiosked.com
					source malc0de {
							description &#34;List of zones serving malicious executables observed by malc0de.com/database/&#34;
							prefix &#34;zone &#34;
							url http://malc0de.com/bl/ZONES
					}
			}
			exclude 122.2o7.net
			exclude 1e100.net
			exclude adobedtm.com
			exclude akamai.net
			exclude amazon.com
			exclude amazonaws.com
			exclude apple.com
			exclude ask.com
			exclude avast.com
			exclude bitdefender.com
			exclude cdn.visiblemeasures.com
			exclude cloudfront.net
			exclude coremetrics.com
			exclude edgesuite.net
			exclude freedns.afraid.org
			exclude github.com
			exclude githubusercontent.com
			exclude google.com
			exclude googleadservices.com
			exclude googleapis.com
			exclude googleusercontent.com
			exclude gstatic.com
			exclude gvt1.com
			exclude gvt1.net
			exclude hb.disney.go.com
			exclude hp.com
			exclude hulu.com
			exclude images-amazon.com
			exclude msdn.com
			exclude paypal.com
			exclude rackcdn.com
			exclude schema.org
			exclude skype.com
			exclude smacargo.com
			exclude sourceforge.net
			exclude ssl-on9.com
			exclude ssl-on9.net
			exclude static.chartbeat.com
			exclude storage.googleapis.com
			exclude windows.net
			exclude yimg.com
			exclude ytimg.com
			hosts {
					include beap.gemini.yahoo.com
					source adaway {
							description &#34;Blocking mobile ad providers and some analytics providers&#34;
							prefix &#34;127.0.0.1 &#34;
							url http://adaway.org/hosts.txt
					}
					source malwaredomainlist {
							description &#34;127.0.0.1 based host and domain list&#34;
							prefix &#34;127.0.0.1 &#34;
							url http://www.malwaredomainlist.com/hostslist/hosts.txt
					}
					source openphish {
							description &#34;OpenPhish automatic phishing detection&#34;
							prefix http
							url https://openphish.com/feed.txt
					}
					source someonewhocares {
							description &#34;Zero based host and domain list&#34;
							prefix 0.0.0.0
							url http://someonewhocares.org/hosts/zero/
					}
					source volkerschatz {
							description &#34;Ad server blacklists&#34;
							prefix http
							url http://www.volkerschatz.com/net/adpaths
					}
					source winhelp2002 {
							description &#34;Zero based host and domain list&#34;
							prefix &#34;0.0.0.0 &#34;
							url http://winhelp2002.mvps.org/hosts.txt
					}
					source yoyo {
							description &#34;Fully Qualified Domain Names only - no prefix to strip&#34;
							prefix &#34;&#34;
							url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&amp;showintro=1&amp;mimetype=plaintext
					}
			}
	}`
