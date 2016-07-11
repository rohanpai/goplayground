package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
)

func main() {
	inputLine := `START,Gooseberry-SA,0x0001FFC100DB7646,3068140336,GMT&#43;05:30-Calcutta,11/20/2015,11:56:51.0,0,0,0,VoIP,IP-TO-IP,DEFAULT,,,9977554432,,0,,0,,0,,RL_BELUR_EXT_EMS,1,GOOSEBERRY-SA:BELUR_EXT_EMS,10.54.26.102,10.54.80.143,BELUR_INT_EMS,,10.54.24.101:38262/10.54.80.140:6000,,10.54.26.101:38134/10.54.80.143:6000,0,,,0x0080023A,,,,2,&#34;SIP,1-32415@10.54.80.140,%22sipp%22;tag=32415SIPpTag001,%22sut%22;tag=gK0cfd65f5,0,,,,sip:9977554432@10.54.24.102:5060,,,,sip:sipp@10.54.80.140:5555,sip:9977554432@10.54.24.102:5060,,,,,,,0,0,,0,0,,,,,,,,1,0,0,0,,,,,,,,0,,,,,,,,,0,,,,,,,,,,,&#34;,12,12,0,5,,,0x0a,9977554432,1,1,,BELUR_EXT_EMS,&#34;SIP,474775912_58694532@10.54.26.102,%22sipp%22;tag=gK0c7d662e,;tag=4829SIPpTag0114382283,0,,,,sip:&#43;19977554432@10.54.80.143:6781;user=phone,,,,sip:sipp@10.54.26.102:5060,sip:10.54.80.143:6781;transport=UDP,,,,,,,0,0,,0,0,,,,,,,,1,0,0,0,,,,,,,,0,,,,,,,,,0,,,,,,,,,,,&#34;,,110,,,1,1,,,2,0x1C4C8168,0,,,,,,0,,,,1,,,,,,,6,,,,&#34;sipp&#34;,2,1,1,1,1,,0,,,1,7,1,,,10.54.24.102,10.54.80.140,60121,16,8,,,,,,,,,0,,,TANDEM,,,,,,,13,1,,,,,,,,,,,,,,,,0,,,,,,,,0,,,&#34;2,1,0,3&#34;,0,,,,,,,,,,,,,,,,,,,,,,,,,,` &#43; &#34;\n&#34;

	fmt.Println(inputLine)

	encoded, _ := json.Marshal(inputLine)
	fmt.Println(string(encoded))
}
