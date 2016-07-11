package main
import(
 "code.google.com/p/leveldb-go/leveldb/table"
 "code.google.com/p/leveldb-go/leveldb/db"
)

func p(r []byte,e error){
  println(string(r))
}

func main(){
 dbname:="my.db"
 dbfs:=db.DefaultFileSystem
 
 if true {
	 f0 ,_:=dbfs.Create(dbname)
	 w:=table.NewWriter(f0 ,nil)  	  
	 w.Set([]byte("google"),[]byte("oh yes"),nil)
	 
	 //try uncomment below line ...
	 //w.Set([]byte("google"),[]byte("uncomment this line,then nothing display.BUG?"),nil)
	 
	 
	 w.Close() //must call to write to file
	 
  }
  
  f1,_:=dbfs.Open(dbname)  
  r:=table.NewReader(f1 ,nil)   
  p( r.Get([]byte("google") ,nil) )     
 
}