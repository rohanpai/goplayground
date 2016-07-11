 import java.io.File;

import java.io.IOException;



import org.apache.avro.Schema;

import org.apache.avro.generic.GenericDatumReader;

import org.apache.avro.generic.GenericRecord;



public class TestSchema {



  public static void main(String[] args) throws IOException {



    Schema schema;

   // schema = new Schema.Parser().parse(new File(&#34;/home/vikas/phoenix/logger-schema.avsc&#34;));

   // schema = new Schema.Parser().parse(new File(&#34;/home/kartiksura/Downloads/logger-schema.avsc&#34;));
    schema = new Schema.Parser().parse(new File(&#34;/home/kartiksura/Downloads/avrophoenix.avsc&#34;));


    GenericDatumReader&lt;GenericRecord&gt; reader = new GenericDatumReader&lt;GenericRecord&gt;(schema);

    System.out.println(&#34;[getBinaryAvroRecords] : reader&#34;);

  }



}
