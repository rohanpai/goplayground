package main

import (
	&#34;fmt&#34;
	&#34;github.com/clbanning/mxj&#34;
)

var xmlData = []byte(`&lt;?xml version=&#34;1.0&#34; encoding=&#34;UTF-8&#34;?&gt;
&lt;response&gt;
	&lt;lst name=&#34;list1&#34;&gt;
	&lt;/lst&gt;
	&lt;lst name=&#34;list2&#34;&gt;
	&lt;/lst&gt;
	&lt;lst name=&#34;list3&#34;&gt;
		&lt;int name=&#34;docId&#34;&gt;1&lt;/int&gt;
		&lt;lst name=&#34;list3-1&#34;&gt;
			&lt;lst name=&#34;list3-1-1&#34;&gt;
				&lt;lst name=&#34;list3-1-1-1&#34;&gt;
					&lt;int name=&#34;field1&#34;&gt;1&lt;/int&gt;
					&lt;int name=&#34;field2&#34;&gt;2&lt;/int&gt;
					&lt;int name=&#34;field3&#34;&gt;3&lt;/int&gt;
					&lt;int name=&#34;field4&#34;&gt;4&lt;/int&gt;
					&lt;int name=&#34;field5&#34;&gt;5&lt;/int&gt;
				&lt;/lst&gt;
			&lt;/lst&gt;
			&lt;lst name=&#34;list3-1-2&#34;&gt;
				&lt;lst name=&#34;list3-1-2-1&#34;&gt;
					&lt;int name=&#34;field1&#34;&gt;1&lt;/int&gt;
					&lt;int name=&#34;field2&#34;&gt;2&lt;/int&gt;
					&lt;int name=&#34;field3&#34;&gt;3&lt;/int&gt;
					&lt;int name=&#34;field4&#34;&gt;4&lt;/int&gt;
					&lt;int name=&#34;field5&#34;&gt;5&lt;/int&gt;
				&lt;/lst&gt;
			&lt;/lst&gt;
		&lt;/lst&gt;
	&lt;/lst&gt;
&lt;/response&gt;`)

func main() {
	// parse XML into a Map
	m, merr := mxj.NewMapXml(xmlData)
	if merr != nil {
		fmt.Println(&#34;merr:&#34;, merr.Error())
		return
	}

	// extract the &#39;list3-1-1-1&#39; node - there&#39;ll be just 1?
	// NOTE: attribute keys are prepended with &#39;-&#39;
	lstVal, lerr := m.ValuesForPath(&#34;*.*.*.*.*&#34;, &#34;-name:list3-1-1-1&#34;)
	if lerr != nil {
		fmt.Println(&#34;ierr:&#34;, lerr.Error())
		return
	}

	// assuming just one value returned - create a new Map
	mv := mxj.Map(lstVal[0].(map[string]interface{}))

	// extract the &#39;int&#39; values by &#39;name&#39; attribute: &#34;-name&#34;
	// interate over list of &#39;name&#39; values of interest
	var names = []string{&#34;field1&#34;, &#34;field2&#34;, &#34;field3&#34;, &#34;field4&#34;, &#34;field5&#34;}
	for _, n := range names {
		vals, verr := mv.ValuesForKey(&#34;int&#34;, &#34;-name:&#34;&#43;n)
		if verr != nil {
			fmt.Println(&#34;verr:&#34;, verr.Error(), len(vals))
			return
		}
		if len(vals) == 0 { // good to check to avoid PANIC
			continue
		}

		// values for simple elements have key &#39;#text&#39;
		// NOTE: there can be only one value for key &#39;#text&#39;
		fmt.Println(n, &#34;:&#34;, vals[0].(map[string]interface{})[&#34;#text&#34;])
	}
}
