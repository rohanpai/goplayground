package main

import (
	&#34;fmt&#34;
	&#34;github.com/alyu/logger&#34;
	&#34;log/syslog&#34;
)

var lg *logger.Logger4go

func main() {
	// get a new logger instance named &#34;example&#34; and with prefix example
	lg = logger.Get(&#34;example&#34;)
	lg.Info(&#34;This is not written out, we need to add a handler first&#34;)

	// log to console/stdout
	lg.AddConsoleHandler()
	lg.Info(&#34;This will be written out to stdout&#34;)

	// log to file. as default the log will be rotated 5 times with a
	// max filesize of 1MB starting with sequence no 1, daily rotate and compression disabled
	_, err := lg.AddStdFileHandler(&#34;/tmp/logger.log&#34;)
	if err != nil {
		fmt.Errorf(&#34;%v&#34;, err)
	}
	lg.Alert(&#34;This is an alert message written to the console and log file&#34;)

	// log to syslog
	protocol := &#34;&#34; // tcp|udp
	ipaddr := &#34;&#34;
	sh, err := lg.AddSyslogHandler(protocol, ipaddr, syslog.LOG_INFO|syslog.LOG_LOCAL0, &#34;example&#34;)
	if err != nil {
		fmt.Errorf(&#34;%v&#34;, err)
	}
	lg.Notice(&#34;This is a critical message written to the console, log file and syslog&#34;)
	lg.Notice(&#34;The format written to syslog is the same as for the console and log file&#34;)
	err = sh.Out.Err(&#34;This is a message to syslog without any preformatted header, it just contains this message&#34;)
	if err != nil {
		fmt.Errorf(&#34;%v&#34;, err)
	}

	// filter logs
	lg.SetFilter(logger.DEBUG | logger.INFO)
	lg.Alert(&#34;This message should not be shown&#34;)
	lg.Debug(&#34;This debug message is filtered through&#34;)
	lg.Info(&#34;As well as this info message&#34;)

	lg = logger.GetWithFlags(&#34;micro&#34;, logger.Ldate|logger.Ltime|logger.Lmicroseconds)
	lg.Info(&#34;This is written out with micrseconds precision&#34;)

	// get standard logger
	log := logger.Std()
	log.Info(&#34;Standard logger always has a console handler&#34;)

	// add a file handler which rotates 5 files with a maximum size of 5MB starting with sequence no 1, daily midnight rotation disabled
	// and with compress logs enabled
	log.AddFileHandler(&#34;/tmp/logger2.log&#34;, uint(5*logger.MB), 5, 1, true, false)

	// add a file handler which keeps logs for 5 days with no filesize limit starting with sequence no 1, daily midnight rotation
	// and  compress logs enabled
	log.AddFileHandler(&#34;/tmp/logger3.log&#34;, 0, 5, 1, true, true)

	// add a file handler with only one daily midnight rotation and compress logs enabled
	log.AddFileHandler(&#34;/tmp/logger3.log&#34;, 0, 1, 1, true, true)

	// Same as above
	fh, _ := log.AddStdFileHandler(&#34;/tmp/logger4.log&#34;)
	fh.SetSize(0)
	fh.SetRotate(1)
	fh.SetCompress(true)
	fh.SetDaily(true)
}
