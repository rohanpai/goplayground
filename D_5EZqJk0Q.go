package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;strings&#34;
	&#34;time&#34;
)

var queensNote string = `결국은 앞으로 신규는 분명히 순환 출자를 막겠다
이렇게 결정을 했습니다.  긴데 기존 것을 머 1년후고 2년후고
다 해소를 하라할적에 그 지분을 계속 유지하기 위해서 더 건설적인데
써야될 것을 같다가 그 고리 끊는데에 쓰게 될거라는게 일케 훤히 보이는
거고 또 미래성장 동력에 투자하는 거보다는 오히려 그것을 막는데 급급하고
또 그렇게 경영권이 이제 약해지게 되며는 외국 자본에도 그게 넘어갈 수가
있습니다. 사실은요 그래서 그런 여러가지를 생각할 때 더군다나 이케 경제가
어려운 상황에서 지금 기존에 이케 허용됐던 것을 인제 지금부터 아닌 걸루다가
딱 끊는 이런 경제정책은 더군다나 이렇게 어려운 시기에는 국민한테도
별로 도움이 안된다.
국민들에게는 실제 어떤 일자리가 만들어지고 가뜩이나 우리 성장 잠재력두 자꾸
떨어지구 있는데 그런쪽에 많이 투자를 하도록 하는 것이 저는 더 실질적인
도움이 된다 이르케 생각하고 있습니다.`

func main() {
	pgh := NewPghRize(nil, 50, 30)
	pgh.SetString(strings.Replace(queensNote, &#34;\n&#34;, &#34; &#34;, -1))

	pause := func(ms int) { time.Sleep(time.Millisecond * time.Duration(ms)) }

	col := 0
	for w := range pgh.Iter() {
		for _, c := range w {
			fmt.Printf(&#34;%c&#34;, c)
			delay := rand.Intn(250)
			if c == rune(&#39;…&#39;) {
				delay &#43;= 250 &#43; rand.Intn(1000)
			}
			pause(delay)
			col &#43;= 2
		}
		fmt.Printf(&#34;%c&#34;, &#39; &#39;)
		col &#43;= 1

		if col &gt;= 80-1 {
			fmt.Println()
			col = 0
		}
	}
	fmt.Println()
}

// -----------------------------------------------------------------------------

type PghRize struct {
	ijs              []string // 간투어들
	words            []string // 소스
	sIdx, cIdx, sCnt int      // 공백 위치 / 전체 공백 개수
	skipIj           bool
	sRatio, eRatio   int // 간투어 삽입 비율 %
}

func NewPghRize(ijs []string, startRatio, endRatio int) *PghRize {
	b := new(PghRize)

	if ijs == nil {
		b.ijs = []string{&#34;어&#34;, &#34;이&#34;, &#34;또&#34;, &#34;그&#34;, &#34;그러니까&#34;, &#34;그게&#34;, &#34;거&#34;}
	} else {
		b.ijs = ijs
	}

	b.sRatio = startRatio
	b.eRatio = endRatio

	rand.Seed(time.Now().UTC().UnixNano())

	return b

}

func (b *PghRize) Reset() {
	b.sIdx = 0
	b.cIdx = 0
}

func (b *PghRize) SetString(s string) {
	b.words = strings.Split(s, &#34; &#34;)
	b.sCnt = len(b.words) - 1
	b.Reset()
}

func (b *PghRize) pickInterjection(repeat int) (ret string) {
	if repeat &lt;= 0 {
		repeat = 1
	}

	lastIdx := -1
	for n := 0; n &lt; repeat; n&#43;&#43; {
		if n != 0 {
			ret &#43;= &#34; &#34;
		}

		i := rand.Intn(len(b.ijs))
		// TODO: ijs가 1개면 안되는데..
		for i == lastIdx {
			i = rand.Intn(len(b.ijs))
		}
		ret &#43;= b.ijs[i] &#43; &#34;…&#34;
		lastIdx = i
	}
	/* log.Println(repeat) */

	return
}

// 공백의 위치에 따라 다른 비율로 간투어 삽입
func (b *PghRize) chooseInterjectionOf(cIdx int) string {
	if b.skipIj {
		b.skipIj = false
		return &#34;&#34;
	}

	currRatio := b.sRatio &#43; ((b.eRatio-b.sRatio)*cIdx)/(b.sCnt-1)
	/* log.Println(&#34;currRatio = &#34;, currRatio) */

	var s string
	if rand.Intn(100) &lt; currRatio {
		s = b.pickInterjection(rand.Intn(9) - 5)
		b.skipIj = true //rand.Intn(100) &lt; 70 /* 간투어 중복 허용? true */
	} else {
		s = &#34;&#34;
		b.skipIj = true
	}
	return s
}

func (b *PghRize) Iter() &lt;-chan string {
	c := make(chan string)
	go func() {
		for b.cIdx &lt;= b.sCnt {
			s := b.chooseInterjectionOf(b.cIdx)
			if strings.EqualFold(s, &#34;&#34;) {
				c &lt;- b.words[b.cIdx]
				b.cIdx &#43;= 1
			} else {
				c &lt;- s
			}
		}
		close(c)
	}()
	return c
}
