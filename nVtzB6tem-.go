package main

import (
	"fmt"
	"regexp"
	"strings"
)

var pattern = `^TID:`
var negate = true

var content = `TID: [0] [BAM] [2015-11-27 23:51:19,549] ERROR {org.wso2.carbon.hadoop.hive.jdbc.storage.db.DBOperation} -  Failed to write data to database {org.wso2.carbon.hadoop.hive.jdbc.storage.db.DBOperation}
org.h2.jdbc.JdbcSQLException: NULL not allowed for column "CONSUMERKEY"; SQL statement:
INSERT INTO API_RESPONSE_SUMMARY_DAY (time,resourcepath,context,servicetime,total_response_count,version,tzoffset,consumerkey,epoch,userid,apipublisher,api) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) [90006-140]
        at org.h2.message.DbException.getJdbcSQLException(DbException.java:327)
        at org.h2.message.DbException.get(DbException.java:167)
        at org.h2.message.DbException.get(DbException.java:144)
        at org.h2.table.Column.validateConvertUpdateSequence(Column.java:294)
        at org.h2.table.Table.validateConvertUpdateSequence(Table.java:621)
        at org.h2.command.dml.Insert.insertRows(Insert.java:116)
        at org.h2.command.dml.Insert.update(Insert.java:82)
        at org.h2.command.CommandContainer.update(CommandContainer.java:70)
        at org.h2.command.Command.executeUpdate(Command.java:199)
        at org.h2.jdbc.JdbcPreparedStatement.executeUpdateInternal(JdbcPreparedStatement.java:141)`


func main() {
	regex, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		fmt.Println("Failed to compile pattern: ", err)
		return
	}

	lines := strings.Split(content, "\n")
	fmt.Printf("matches\tline\n")
	for _, line := range lines {
		matches := regex.MatchString(line)
		if negate {
			matches = !matches
		}
		fmt.Printf("%v\t%v\n", matches, line)
	}
}