package main

import (
	&#34;crypto/aes&#34;
	&#34;crypto/cipher&#34;
	&#34;encoding/hex&#34;
	&#34;fmt&#34;
	&#34;strings&#34;
)

func min(a, b int) int {
	if a &lt; b {
		return a
	}
	return b
}

//forward takes a multi-line string of ACSII art and returns
//a slice of hex-encoded strings which when decoded in a certain way prints out
//said ASCII art. You can then post it on the internet to show how funny you are.
func forward(input, key string) (secret []string) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(c, make([]byte, aes.BlockSize))
	lines := strings.Split(input, &#34;\n&#34;)
	secret = make([]string, 0)

	for i := range lines {
		output := make([]byte, len(lines[i]))
		stream.XORKeyStream(output, []byte(lines[i]))
		secret = append(secret, fmt.Sprint(hex.EncodeToString(output)))
	}
	return
}

var data = []string{&#34;92ebfd9371a19a7e8eb7c824af1da09c426cc689173020dc6f6eaad99203b72881339d6637fdf4d3e3&#34;,
	&#34;b68e4644674f85af27e38f52a6a6554e11e3c304d3aa0a39fffa24514eee274f74f41ad6210ea9ff56&#34;,
	&#34;80aa490e173c61c099e04ed9e7f18165791e4840e35c313eac0de04fa21f97d78ecdea03f5862489ce&#34;,
	&#34;e91f7a9ccab0ff90c336dbdc38660621681bfb2d09fc7ca51796f29a648f0908a2d4ed427453887080&#34;,
	&#34;ee1825f272ffdea73a8496693d24e9bd12fa4605c314a60ffc45564a8afb7536dfcbdbd01ded7948ef&#34;,
	&#34;07f589bb76828371d59d5ee8948905a5807a72d464033f0a1ce4ea71fb8bc72ea83f3f095cce2248a0&#34;,
	&#34;b69008c4456e8c5508d531876ca9a44a9e70b6401c9eb32bb44a71e1377438a4ead041ee154598caebf1d39c&#34;,
	&#34;964c3d3ec640ffceb20009d59ac0c33b7f873baadb8bcd722c31af00d4f4a12664e2a18f5c3eb5e6095c0e5c&#34;,
	&#34;8e10c74e7f51e578a2a7f3791eea1f6347236d78cdbf03807aae581c129f883435ea43e26b1aab700f2bf632&#34;,
	&#34;e632a89374a9656fc9d7ee96c940615f1b69f561e676a3190fabb3a57cff1d4d4cae311519562305b01093ad&#34;,
	&#34;4eeb79381dfe0de69c5946042f756e0d35ae3516fee7c51894daf605a54d4909b1e9ca6ece9897fc237391f9&#34;,
	&#34;6d0d88e8ffbe2b25f08201294e16b468e2e55f9e91656591aacc10075c7455d86796098ccf0a4d4f9168435f&#34;,
	&#34;e1e8d5697f1cc39f22ed6fc871ca3ab69e2fa722103b5f24c8a0a5e711439d1e18f64d9ffd808fcc8f4b74468b81069b81ef290207e1521678340b5e4ffe40b6a24f147f4c&#34;,
	&#34;38801d2345805bdad452e9fc4e939cd8a168c14e0e7b966c1a52d53768a1e0e39158b8245386ad2f287e0a0ec642e3dbe6e938d3d25820a9182aff65073970629a3266902b&#34;,
	&#34;70e4b2eddcf8e197b8676ba3ba7596312cee61d0acf25eed4b9077989d60cefd95f6e8ab63429d7d73cd76d9a172cb04917af596dd9cbc32521e345441e09725c77b684759&#34;,
	&#34;9627c02c3e7d21c90f7554f7de679d49e605a88db56c5874735cc98ade8b4b1b7d2f76212623ae23aea38b55597b30cd11737dc62d104e8848d1bce62f8226968492f23536&#34;,
	&#34;8a0aea66d01a24ff1071efeec96db2f3e572f0165c7ca2dfa0bffa30f4a85fa359376da74ff87119b68e9a3a3dea854bf94e8bc589ba228fe636477f4af1928153fa17dfd4&#34;,
	&#34;ff99c4ac3aa40ff06af55b7aef27b0fc12afbaba428894deea84804377f924514acb3aea4b73fa5e3414719ca3c056587af9481b46528a616940eef2f6b38a9b4420f80188&#34;,
	&#34;37141bb08cd28c3de2be0c0f851cf0131e0974c7c0e902690c3d892297b0f10ba627cc604c44e0229435877eb6defa683570f03f28c50c&#34;,
	&#34;7df3dac053f3a0578766f3132483bbc74bd3a99f782d28e1a60cd19fdd4daa5ac2681b87499281b2dc52a6393c405d4c3c78bbd46bcca2&#34;,
	&#34;2677f13da1a8598573a524f701e1e5d6e0f6a50fb182d3faaef33cb43a0970bb92e14f1499d069d85b3e46de94c177ddff18c2f49387d7&#34;,
	&#34;0b20d18955b3f8c5db838fa309fe44d6ed95b8e1816e53d98ce0bd38031d87ebed5e5c71a8e2ddf2d620e0b75c66d99ad3bf5177bdc80c&#34;,
	&#34;7e8deb44e4e19db2e12fb3974c44fb3eb25ae22a07ce27b622152f36de3c996147765f7938fc78364fdfca2c8ec78578585dfab2ea50c2&#34;,
	&#34;9270b682dd883f0f6422b5bbe72280ba815c2501575f850f1ab9f4be943d9922466e2e63d124c89dc86abcf6937e8c416022fc83d9b367&#34;,
}

func main() {
	c, err := aes.NewCipher([]byte(&#34;thats hot lololo&#34;))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(c, make([]byte, aes.BlockSize))
	for _, line := range data {
		data, err := hex.DecodeString(line)
		if err != nil {
			panic(err)
		}
		stream.XORKeyStream(data, data)
		fmt.Println(string(data))
	}
}
