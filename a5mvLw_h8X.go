package main

import (
	&#34;encoding/xml&#34;
	&#34;fmt&#34;
	&#34;html/template&#34;
)

var xmlContent = `
&lt;?xml version=&#39;1.0&#39; encoding=&#39;UTF-8&#39;?&gt;&lt;?xml-stylesheet type=&#39;text/xsl&#39; href=&#39;http://solidot.org.feedsportal.com/xsl/eng/rss.xsl&#39;?&gt;
&lt;rss xmlns:rdf=&#34;http://www.w3.org/1999/02/22-rdf-syntax-ns#&#34; xmlns:taxo=&#34;http://purl.org/rss/1.0/modules/taxonomy/&#34; xmlns:sy=&#34;http://purl.org/rss/1.0/modules/syndication/&#34; xmlns:dc=&#34;http://purl.org/dc/elements/1.1/&#34; xmlns:itunes=&#34;http://www.itunes.com/dtds/podcast-1.0.dtd&#34; xmlns:slash=&#34;http://purl.org/rss/1.0/modules/slash/&#34; version=&#34;2.0&#34;&gt;
&lt;channel&gt;
  &lt;title&gt;Solidot&lt;/title&gt;
  &lt;link&gt;http://www.solidot.org&lt;/link&gt;
  &lt;description&gt;奇客的资讯，重要的东西&lt;/description&gt;
  &lt;language&gt;en-us&lt;/language&gt;
  &lt;copyright&gt;Copyright © 2006, Solidot.org&lt;/copyright&gt;
  &lt;pubDate&gt;Thu, 26 Feb 2015 14:15:27 GMT&lt;/pubDate&gt;
  &lt;lastBuildDate&gt;Thu, 26 Feb 2015 14:15:27 GMT&lt;/lastBuildDate&gt;
  &lt;ttl&gt;2&lt;/ttl&gt;
  &lt;sy:updatePeriod&gt;hourly&lt;/sy:updatePeriod&gt;
  &lt;sy:updateFrequency&gt;1&lt;/sy:updateFrequency&gt;
  &lt;sy:updateBase&gt;1970-01-01T00:00:00Z&lt;/sy:updateBase&gt;
  &lt;dc:creator&gt;admin@solidot.org&lt;/dc:creator&gt;
  &lt;dc:subject&gt;Technology&lt;/dc:subject&gt;
  &lt;dc:publisher&gt;Solidot.org&lt;/dc:publisher&gt;
  &lt;dc:language&gt;en-us&lt;/dc:language&gt;
  &lt;dc:rights&gt;Copyright © 2006, Solidot.org&lt;/dc:rights&gt;
  &lt;image&gt;&lt;title&gt;Solidot&lt;/title&gt;&lt;url&gt;http://solidot.org/images/topics/topicslash.gif&lt;/url&gt;&lt;link&gt;http://www.solidot.org&lt;/link&gt;&lt;/image&gt;
  &lt;item&gt;
    &lt;title&gt;Reddit禁止未经当事人同意的裸体图片或视频&lt;/title&gt;
    &lt;link&gt;http://solidot.org.feedsportal.com/c/33236/f/556826/s/43d1e9a4/sc/21/l/0L0Ssolidot0Borg0Cstory0Dsid0F43138/story01.htm&lt;/link&gt;
    &lt;description&gt;Reddit引入了“非自愿色情”的隐私条款，从3月10日起禁止用户在网站上张贴未经当事人同意的裸体或性活动照片、视频或数字图像的链接。Reddit认为非自愿色情是一种性骚扰行为，表示决不会容忍这种行为。Google也几乎在同一时间宣布，其博客服务平台从3月23日起，不再能公开发表含有色情内容的图像和视频。&amp;lt;img width=&#39;1&#39; height=&#39;1&#39; src=&#39;http://solidot.org.feedsportal.com/c/33236/f/556826/s/43d1e9a4/sc/21/mf.gif&#39; border=&#39;0&#39;/&amp;gt;&amp;lt;br clear=&#39;all&#39;/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/1/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/1/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/2/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/2/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/3/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/rc/3/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/a2.htm&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/a2.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;img width=&#34;1&#34; height=&#34;1&#34; src=&#34;http://pi.feedsportal.com/r/222166538220/u/0/f/556826/c/33236/s/43d1e9a4/sc/21/a2t.img&#34; border=&#34;0&#34;/&amp;gt;&lt;/description&gt;
    &lt;pubDate&gt;Thu, 26 Feb 2015 14:14:59 GMT&lt;/pubDate&gt;
    &lt;guid isPermaLink=&#34;false&#34;&gt;http://www.solidot.org/story?sid=43138&lt;/guid&gt;
    &lt;slash:department&gt;/r/nsfw要被关了&lt;/slash:department&gt;
    &lt;dc:creator&gt;WinterIsComing&lt;/dc:creator&gt;
  &lt;/item&gt;
  &lt;item&gt;
    &lt;title&gt;Facebook删除一位台湾教授的账号&lt;/title&gt;
    &lt;link&gt;http://solidot.org.feedsportal.com/c/33236/f/556826/s/43d104a2/sc/21/l/0L0Ssolidot0Borg0Cstory0Dsid0F43137/story01.htm&lt;/link&gt;
    &lt;description&gt;Facebook删除了台湾中山大学生物学教授顏聖紘（已不存在，新建账号）的账号，理由是他违反了用户条款（如图所示），但没有给出详细解释。账户关闭后，顏聖紘在Facebook上建立的课程等活动页面也一并消失。Facebook删除账号被认为是因为顏教授的言论招致了一些人的反感，导致对方向Facebook投诉封杀其言论。社交巨人有专人负责检查投诉，删除账号的决定权最终在它手中。此事引发了言论管制的争论，Facebook之前为避免被土耳其封锁而屏蔽批评穆罕默德的帖子，并冻结一中国作家的账号。&amp;lt;img width=&#39;1&#39; height=&#39;1&#39; src=&#39;http://solidot.org.feedsportal.com/c/33236/f/556826/s/43d104a2/sc/21/mf.gif&#39; border=&#39;0&#39;/&amp;gt;&amp;lt;br clear=&#39;all&#39;/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/1/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/1/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/2/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/2/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/3/rc.htm&#34; rel=&#34;nofollow&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/rc/3/rc.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;br/&amp;gt;&amp;lt;a href=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/a2.htm&#34;&amp;gt;&amp;lt;img src=&#34;http://da.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/a2.img&#34; border=&#34;0&#34;/&amp;gt;&amp;lt;/a&amp;gt;&amp;lt;img width=&#34;1&#34; height=&#34;1&#34; src=&#34;http://pi.feedsportal.com/r/222166529049/u/0/f/556826/c/33236/s/43d104a2/sc/21/a2t.img&#34; border=&#34;0&#34;/&amp;gt;&lt;/description&gt;
    &lt;pubDate&gt;Thu, 26 Feb 2015 13:34:14 GMT&lt;/pubDate&gt;
    &lt;guid isPermaLink=&#34;false&#34;&gt;http://www.solidot.org/story?sid=43137&lt;/guid&gt;
    &lt;slash:department&gt;政治要正确&lt;/slash:department&gt;
    &lt;dc:creator&gt;WinterIsComing&lt;/dc:creator&gt;
  &lt;/item&gt;
&lt;/channel&gt;
&lt;/rss&gt;
`

type Rss2 struct {
	XMLName xml.Name `xml:&#34;rss&#34;`
	Version string   `xml:&#34;version,attr&#34;`
	// Required
	Title       string `xml:&#34;channel&gt;title&#34;`
	Link        string `xml:&#34;channel&gt;link&#34;`
	Description string `xml:&#34;channel&gt;description&#34;`
	// Optional
	PubDate  string `xml:&#34;channel&gt;pubDate&#34;`
	ItemList []Item `xml:&#34;channel&gt;item&#34;`
}

type Item struct {
	// Required
	Title       string        `xml:&#34;title&#34;`
	Link        string        `xml:&#34;link&#34;`
	Description template.HTML `xml:&#34;description&#34;`
	// Optional
	Content  template.HTML `xml:&#34;encoded&#34;`
	PubDate  string        `xml:&#34;pubDate&#34;`
	Comments string        `xml:&#34;comments&#34;`
}

func main() {
	r := Rss2{}
	err := xml.Unmarshal([]byte(xmlContent), &amp;r)
	if err != nil {
		panic(err)
	}
	for _, item := range r.ItemList {
		fmt.Println(item)
	}
}
