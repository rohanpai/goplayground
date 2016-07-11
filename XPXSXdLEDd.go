// 获取当天的美剧时间表HTML
func get_day_html(html string, day string) []string {
	re_str := fmt.Sprintf(&#34;&lt;td class=\&#34;(cur|ihbg)\&#34;&gt;.&#43;?&lt;dt&gt;%s&lt;/dt&gt;(.&#43;?)&lt;/dl&gt;&#34;, day) // day 为“1号”
	fmt.Println(re_str)
	re, _ := regexp.Compile(re_str)
	src := re.FindAllString(html, -1)
	return src
}


这里是HTML代码。我想拿到第一个dd到最后一个dd之间的内容。

&lt;td class=&#34;cur&#34;&gt;
                &lt;dl&gt;
                    &lt;dt&gt;1号&lt;/dt&gt;
                    &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/28516&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34; style=&#34;display: none;&#34;&gt;&lt;span&gt; 杀手信徒 The Following  S02E11&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;杀手信徒&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S02E11&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/30032&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34; style=&#34;display: none;&#34;&gt;&lt;span&gt; 未来青年 The Tomorrow People  S01E18&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;未来青年&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E18&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/27750&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 皇家律师 Silk  S03E06&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;皇家律师&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S03E06&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/27621&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 家族风云 Dallas  S03E06&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;家族风云&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S03E06&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/29965&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 智能缉凶 Intelligence  S01E13&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;智能缉凶&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E13&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/26154&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 老爸老妈的浪漫史 How I Met Your Mother  S09E23&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;老爸老妈的浪漫史&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S09E23&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/26154&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 老爸老妈的浪漫史 How I Met Your Mother  S09E24&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;老爸老妈的浪漫史&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S09E24&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/11171&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 我欲为人 Being Human US  S04E12&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;我欲为人&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S04E12&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/10452&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 射手 Archer  S05E10&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;射手&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S05E10&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/31346&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 瑞克与莫蒂 Rick and Morty  S01E10&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;瑞克与莫蒂&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E10&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/30659&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 极品老妈 Mom  S01E21&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;极品老妈&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E21&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/31875&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 独夫 The widower  S101E03&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;独夫&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;E03&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/29964&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 罪恶黑名单 The Blacklist  S01E18&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;罪恶黑名单&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E18&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd class=&#34;even&#34;&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/31940&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 损友的美好时代 Friends With Better Lives  S01E01&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;损友的美好时代&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S01E01&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                  &lt;dd&gt;
                        &lt;a onmouseout=&#34;hidee(this)&#34; onmouseover=&#34;show(this)&#34; href=&#34;http://www.yyets.com/resource/11097&#34; target=&#34;_blank&#34;&gt;
                        &lt;div class=&#34;floatSpan&#34;&gt;&lt;span&gt; 识骨寻踪 Bones  S09E19&lt;span&gt;&lt;/span&gt;&lt;/span&gt;&lt;/div&gt;
                        &lt;font class=&#34;fa1&#34;&gt;识骨寻踪&lt;/font&gt;
                        &lt;font class=&#34;fa1&#34; style=&#34;color:white&#34;&gt;S09E19&lt;/font&gt;
                        &lt;/a&gt;

                    &lt;/dd&gt;
                                  &lt;/dl&gt;
            &lt;/td&gt;