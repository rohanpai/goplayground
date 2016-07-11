/*
func New(out io.Writer, prefix string, flag int) *Logger

out:    The out variable sets the destination to which log data will be written.
prefix: The prefix appears at the beginning of each generated log line.
flags:  The flag argument defines the logging properties.
*/

// Sample program to show how to extend the log package
// from the standard library.
package main

import (
	&#34;io&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;os&#34;
)

var (
	// Trace is for full detailed messages.
	Trace *log.Logger

	// Info is for important messages.
	Info *log.Logger

	// Warning is for need to know issue messages.
	Warning *log.Logger

	// Error is for error messages.
	Error *log.Logger
)

// main is the entry point for the application.
func main() {
	// Open a file for warnings.
	warnings, err := os.OpenFile(&#34;warnings.log&#34;, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(&#34;Failed to open warning log file&#34;)
	}
	defer warnings.Close()

	// Open a file for errors.
	errors, err := os.OpenFile(&#34;errors.log&#34;, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(&#34;Failed to open errors log file&#34;)
	}
	defer errors.Close()

	// Create a multi writer for errors.
	multi := io.MultiWriter(errors, os.Stderr)

	// Init the log package for each message type.
	initLog(ioutil.Discard, os.Stdout, warnings, multi)

	// Test each log type.
	Trace.Println(&#34;I have something standard to say.&#34;)
	Info.Println(&#34;Important Information.&#34;)
	Warning.Println(&#34;There is something you need to know about.&#34;)
	Error.Println(&#34;Something has failed.&#34;)
}

// initLog sets the devices for each log type.
func initLog(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle,
		&#34;TRACE: &#34;,
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		&#34;INFO: &#34;,
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		&#34;WARNING: &#34;,
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		&#34;ERROR: &#34;,
		log.Ldate|log.Ltime|log.Lshortfile)
}
