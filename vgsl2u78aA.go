package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;math&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
)

const MaxSerialNo = math.MaxUint64 - math.MaxUint32

func SerialNoInc(sn string) (string, error) {
	n, err := strconv.ParseUint(sn, 10, 64)
	if err == nil &amp;&amp; n &gt;= MaxSerialNo {
		err = errors.New(&#34;serial number invalid after increment&#34;)
	}
	if err != nil {
		return &#34;&#34;, fmt.Errorf(&#34;SerialNoInc(%#v): %s&#34;, sn, err.Error())
	}
	s := strconv.FormatUint(n&#43;1, 10)
	z := len(sn) - len(s)
	if z &lt; 0 {
		z = 0
	}
	return sn[:z] &#43; s, nil
}

func pubSerialStrInc(ss, sep string) (string, error) {
	i := strings.Index(ss, sep)
	if i &lt; 0 {
		return &#34;&#34;, fmt.Errorf(&#34;pubSerialStrInc(%#v, %#v): separator %#v not found&#34;, ss, sep, sep)
	}
	i &#43;= len(sep)
	if i &gt;= len(ss) {
		return &#34;&#34;, fmt.Errorf(&#34;pubSerialStrInc(%#v, %#v): serial number not found&#34;, ss, sep)
	}
	sn, err := SerialNoInc(ss[i:])
	if err != nil {
		return &#34;&#34;, fmt.Errorf(&#34;pubSerialStrInc(%#v, %#v): %s&#34;, ss, sep, err.Error())
	}
	return ss[:i] &#43; sn, nil
}

func main() {
	serialStrs := []struct {
		ss, sep string
	}{
		{&#34;xs-00009&#34;, &#34;-&#34;},
		{&#34;yxs&lt;&gt;000099&#34;, &#34;&lt;&gt;&#34;},
		{&#34;yxzs-99&#34;, &#34;-&#34;},
		{&#34;yxzs-&#34; &#43; strconv.FormatUint(MaxSerialNo-1, 10), &#34;-&#34;},
		{&#34;yxzs-0000&#34; &#43; strconv.FormatUint(MaxSerialNo-1, 10), &#34;-&#34;},
		{&#34;#01&#34;, &#34;#&#34;},
		{&#34;01&#34;, &#34;&#34;},
		{&#34;err-09&#34;, &#34;=&#34;},
		{&#34;err-&#34;, &#34;-&#34;},
		{&#34;err-09x&#34;, &#34;-&#34;},
		{&#34;err-&#34; &#43; strconv.FormatUint(MaxSerialNo, 10), &#34;-&#34;},
		{&#34;err-&#34; &#43; strconv.FormatUint(math.MaxUint64, 10), &#34;-&#34;},
	}
	for _, serialStr := range serialStrs {
		si, err := pubSerialStrInc(serialStr.ss, serialStr.sep)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(serialStr.ss, &#34;&#43;&#43;  ==&gt;&#34;, si)
	}
}
