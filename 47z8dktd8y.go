package main
import(
 &#34;code.google.com/p/leveldb-go/leveldb/table&#34;
 &#34;code.google.com/p/leveldb-go/leveldb/db&#34;
)

func p(r []byte,e error){
  println(string(r))
}

func main(){
 dbname:=&#34;my.db&#34;
 dbfs:=db.DefaultFileSystem
 
 if true {
	 f0 ,_:=dbfs.Create(dbname)
	 w:=table.NewWriter(f0 ,nil)  	  
	 w.Set([]byte(&#34;google&#34;),[]byte(&#34;oh yes&#34;),nil)
	 
	 //try uncomment below line ...
	 //w.Set([]byte(&#34;google&#34;),[]byte(&#34;uncomment this line,then nothing display.BUG?&#34;),nil)
	 
	 
	 w.Close() //must call to write to file
	 
  }
  
  f1,_:=dbfs.Open(dbname)  
  r:=table.NewReader(f1 ,nil)   
  p( r.Get([]byte(&#34;google&#34;) ,nil) )     
 
}