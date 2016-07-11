package main

import (
	&#34;bytes&#34;
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;os&#34;
)

const s = `&lt;?xml version=&#34;1.0&#34; encoding=&#34;utf-8&#34;?&gt;
&lt;manifest xmlns:android=&#34;http://schemas.android.com/apk/res/android&#34; package=&#34;org.golang.todo.{{.NAME}}&#34; android:versionCode=&#34;1&#34; android:versionName=&#34;1.0&#34;&gt;
  &lt;!-- comment --&gt;
  &lt;uses-sdk android:minSdkVersion=&#34;9&#34; /&gt;
  &lt;application android:label=&#34;{{.LABEL}}&#34; android:debuggable=&#34;true&#34; android:icon=&#34;@drawable/ic_launcher&#34;&gt;
    &lt;activity android:name=&#34;org.golang.app.GoNativeActivity&#34; android:label=&#34;{{.LABEL}}&#34; android:configChanges=&#34;orientation|keyboardHidden&#34;&gt;
      &lt;meta-data android:name=&#34;android.app.lib_name&#34; android:value=&#34;{{.NAME}}&#34; /&gt;
      &lt;intent-filter&gt;
        &lt;action android:name=&#34;android.intent.action.MAIN&#34; /&gt;
        &lt;category android:name=&#34;android.intent.category.LAUNCHER&#34; /&gt;
      &lt;/intent-filter&gt;
    &lt;/activity&gt;
  &lt;/application&gt;
&lt;/manifest&gt;
`

type Tag struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []interface{}
}

func (t *Tag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = t.Name
	start.Attr = t.Attr
	e.EncodeToken(start)
	for _, v := range t.Children {
		switch v.(type) {
		case *Tag:
			child := v.(*Tag)
			if err := e.Encode(child); err != nil {
				return err
			}
		case xml.CharData:
			e.EncodeToken(v.(xml.CharData))
		case xml.Comment:
			e.EncodeToken(v.(xml.Comment))
		}
	}
	e.EncodeToken(start.End())
	return nil
}

func (t *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t.Name = start.Name
	t.Attr = start.Attr
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			var data *Tag
			if err := d.DecodeElement(&amp;data, &amp;tok); err != nil {
				return err
			}
			t.Children = append(t.Children, data)
		case xml.CharData:
			t.Children = append(t.Children, token.(xml.CharData).Copy())
		case xml.Comment:
			t.Children = append(t.Children, token.(xml.Comment).Copy())
		}
	}
	return nil
}

func main() {
	v := &amp;Tag{}
	xml.NewDecoder(bytes.NewBuffer([]byte(s))).Decode(&amp;v)
	fmt.Println(v.Name)
	fmt.Println(v.Attr)
	fmt.Println(v.Children)
	xml.NewEncoder(os.Stdout).Encode(v)
}
