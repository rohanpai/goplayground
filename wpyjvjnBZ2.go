package main

import &#34;fmt&#34;
import &#34;errors&#34;

func main() {
	var pool, err = NewSimpleMemPool(1024, 1024)
	
	if err != nil {
		fmt.Println(&#34;Create Memory Pool Failed!&#34;)
	}
	
	var mem = pool.Alloc(512)
	
	if mem == nil {
		fmt.Println(&#34;Allocate Memory Failed!&#34;)
	}
	
	fmt.Println(&#34;Hello, playground&#34;)
}

//
// 简单的内存池实现，用于避免频繁的零散内存申请
//
type SimpleMemPool struct {
	memPool     []byte
	memPoolSize int
	maxPackSize int
}

//
// 创建一个简单内存池，预先申请&#39;memPoolSize&#39;大小的内存，每次分配内存时从中切割出来，直到剩余空间不够分配，再重新申请一块。
// 参数&#39;maxPackSize&#39;用于限制外部申请内存允许的最大长度，所以这个值必须小于等于&#39;memPoolSize&#39;。
//
func NewSimpleMemPool(memPoolSize, maxPackSize int) (*SimpleMemPool, error) {
	if maxPackSize &gt; memPoolSize {
		return nil, errors.New(&#34;maxPackSize &gt; memPoolSize&#34;)
	}

	return &amp;SimpleMemPool{
		memPool:     make([]byte, memPoolSize),
		memPoolSize: memPoolSize,
		maxPackSize: maxPackSize,
	}, nil
}

//
// 申请一块内存，如果&#39;size&#39;超过&#39;maxPackSize&#39;设置将返回nil
//
func (this *SimpleMemPool) Alloc(size int) (result []byte) {
	if size &gt; this.maxPackSize {
		return nil
	}

	if len(this.memPool) &lt; size {
		this.memPool = make([]byte, this.memPoolSize)
	}

	result = this.memPool[0:size]
	this.memPool = this.memPool[size:]

	return result
}