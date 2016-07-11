package main

import (
	&#34;fmt&#34;
	&#34;strings&#34;
	&#34;time&#34;
)

func main() {
	for i, layout := range strToTimeFormats {
		for j, strval := range strToTimeFormats {
			strval = strings.Replace(strval, &#34;_2&#34;, &#34; 2&#34;, 1) // _2 标记表示用用空格替代前导0, 值里面是没有_2前面的_

			if j != 0 { // 对于值 ISO8601 最后是 Z 表示使用 UTC 时间
				strval = strings.Replace(strval, &#34;Z&#34;, &#34;-&#34;, 1) // Z 标记在值里面被 -,&#43; 号替换
			}

			t, err := time.Parse(layout, strval)
			out := t.Format(layout)

			// 只有当i==j的时候输出才会和处理后的 strval 相等
			if err == nil &amp;&amp; i != j &amp;&amp; out != strval || err != nil &amp;&amp; i == j {
				fmt.Println(&#34;Failed layout&#34;, err, &#34;\n&#34;, i, layout, &#34;\n&#34;, j, strval)
				return
			}
			if i != j {
				continue
			}

			// 测试 Nanosecond 部分, 通过给 layout 加上 &#34; .000000000&#34; 保留所有尾部的 0 来判断,
			// .9 是省略尾部的 0
			t = t.Add(time.Duration(123456780))
			out = t.Format(layout &#43; &#34; .000000000&#34;)

			if out != strval&#43;&#34; .123456780&#34; {
				fmt.Println(&#34;Failed Nanosecond&#34;, i, &#34;\n&#34;, strval&#43;&#34; .123456780&#34;, &#34;\n&#34;, out)
				return
			}
		}
	}
	fmt.Println(&#34;Okay&#34;)
}

// strToTimeFormats 定义了可预见的格式
// 最初的代码来源于 https://github.com/gosexy/to
// 注意含有 Nanosecond 的格式部分省略也是可以兼容的
// 要注意对于值 ISO8601 最后是 Z 表示使用 UTC 时间
// Go 在没有标明时区是使用的是 GMT/UTC 时区
var strToTimeFormats = []string{
	&#34;2006-01-02T15:04:05Z&#34;, // ISO8601 TimeZone is UTC
	&#34;2006-01-02&#34;,
	&#34;2006-01-02 15:04&#34;,
	&#34;2006-01-02 15:04:05&#34;, // &#34;2006-01-02 15:04:05.000000000&#34;,
	&#34;2006-01-02T15:04&#34;,
	&#34;2006-01-02T15:04:05&#34;, // &#34;2006-01-02T15:04:05.000000000&#34;,
	&#34;01/02/2006&#34;,
	&#34;01/02/2006 15:04&#34;,
	&#34;01/02/2006 15:04:05&#34;, // &#34;01/02/2006 15:04:05.000000000&#34;,
	&#34;01/02/06&#34;,
	&#34;01/02/06 15:04&#34;,
	&#34;01/02/06 15:04:05&#34;, // &#34;01/02/06 15:04:05.000000000&#34;,
	&#34;_2/Jan/2006 15:04:05&#34;,
	&#34;Jan _2, 2006&#34;,
	time.ANSIC,
	time.Kitchen,
	time.Stamp,                 //time.StampMilli, time.StampMicro, time.StampNano,
	&#34;2006-01-02 15:04:05-0700&#34;, // follow with timezone
	&#34;2006-01-02 15:04:05 -0700&#34;,
	&#34;2006-01-02 15:04:05Z07:00&#34;,
	&#34;2006-01-02 15:04:05 Z07:00&#34;,
	&#34;Mon Jan _2 15:04:05 -0700 MST 2006&#34;, //&#34;Mon Jan _2 15:04:05.000000000 -0700 MST 2006&#34;,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339, //time.RFC3339Nano,
}
