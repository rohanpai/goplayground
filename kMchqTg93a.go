package main

import(
    &#34;log&#34;
    &#34;database/sql&#34;
    // 引入 mymysql 的驱动，驱动将会自动注册。注意，不要丢失前面的下划线（_）。
    _ &#34;github.com/mikespook/mymysql/godrv&#34;
)

func main() {
    // 使用 mymysql 驱动打开一个 sql.DB。
    // 对于 mymysql 中，dsn 有多种写法。
    // 使用 tcp 协议：[tcp://addr/]dbname/user/password[?params]
    // 使用 unix sock：[unix://sockpath/]dbname/user/password[?params]
    // 括号中是可选的内容，当协议信息未指定时，默认使用 tcp://127.0.0.1:3306/。 
    // params 部分可设置两个参数
    //   charset - 用于 &#39;set names&#39; 设置连接编码。
    //   keepalive - 每 keepalive 秒向服务器发送 PING。
    //
    // 如果密码含有斜线（/），则用星号（*）代替。
    // 如果密码含有星号（*），用两个星号（**）代替。
    //   pass/wd =&gt; pass*wd
    //   pass*wd =&gt; pass**wd

    db, err := sql.Open(&#34;mymysql&#34;, &#34;tcp://127.0.0.1:3306/test/root/xxiyy?charset=utf8&amp;keepalive=1200&#34;)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 使用 Prepare 进行数据插入
    stmt, err := db.Prepare(&#34;insert into `test` (`key`, `value`) values (?, ?)&#34;)
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    rslt, err := stmt.Exec(&#34;name&#34;, &#34;foobar&#34;)
    if err != nil {
        log.Fatal(err)
    }
    if a, err := rslt.RowsAffected(); err != nil {
        log.Print(err)
    } else {
        log.Printf(&#34;[INS]Affected rows=%d&#34;, a)
    }
    if id, err := rslt.LastInsertId(); err != nil {
        log.Print(err)
    } else {
        log.Printf(&#34;[INS]Last insert id=%d&#34;, id)
    }

    // 查询多行数据
    rows, err := db.Query(&#34;select * from `test`&#34;)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var k, v string
        rows.Scan(&amp;k, &amp;v)
        log.Printf(&#34;[ROWS]key=%s, value=%s&#34;, k, v)
    }

    // 查询单行数据
    row := db.QueryRow(&#34;select * from `test` where `key` = ?&#34;, &#34;name&#34;)
    var k, v string
    row.Scan(&amp;k, &amp;v)
    log.Printf(&#34;[ROW]key=%s, value=%s&#34;, k, v)

    // 删除
    rslt, err = db.Exec(&#34;delete from `test`&#34;)
    if a, err := rslt.RowsAffected(); err != nil {
        log.Print(err)
    } else {
        log.Printf(&#34;[DEL]Affected rows=%d&#34;, a)
    }
    if id, err := rslt.LastInsertId(); err != nil {
        log.Print(err)
    } else {
        log.Printf(&#34;[DEL]Last insert id=%d&#34;, id)
    }
}