 import java.io.File;

import java.io.IOException;



import org.apache.avro.Schema;

import org.apache.avro.generic.GenericDatumReader;

import org.apache.avro.generic.GenericRecord;



public class TestSchema {



  public static void main(String[] args) throws IOException {



    Schema schema;

   // schema = new Schema.Parser().parse(new File("/home/vikas/phoenix/logger-schema.avsc"));

   // schema = new Schema.Parser().parse(new File("/home/kartiksura/Downloads/logger-schema.avsc"));
    schema = new Schema.Parser().parse(new File("/home/kartiksura/Downloads/avrophoenix.avsc"));


    GenericDatumReader<GenericRecord> reader = new GenericDatumReader<GenericRecord>(schema);

    System.out.println("[getBinaryAvroRecords] : reader");

  }



}
