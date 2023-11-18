# 2

1. project structure
2. web server
3. routes and handlers
4. HTTP responses, headers and status code
5. URL query parameters
6. render HTML pages
7. serve static files

## 2.1

1. `go mod init x.x.x`
2. `go run .`
3. `go run` - compile and execute binary in `/tmp`
4. `go run .` - `go run main.go` - `go run snippetbox.hh` - `.` 代表当前文件夹


## 2.2

1. MVC
2. handler - controller - app logic, write http response headers and bodies
3. router(servemux) - map url patterns and handlers
4. web server - listen for requests
5. net/http
6. `*http.Request` 包含请求信息
7. `http.ResponseWriter`
8. `http.ListenAndServe` - `ip:port`
9. `:http` or `:http-alt` - 没有设置端口就从 `/etc/services` 获取
10. 启动服务 -> 侦听端口 -> 接受新请求 -> 将请求传给 servemux -> 检查 URL 路径 -> 派发给匹配的 handler
11. `/` 包罗万象，相当于通配符 `/**`

## 2.3

1. url patterns - fixed paths and subtree paths (end with a trailing slash / )
2. subtree paths - `/foo` will automatically redirected `/foo/`
3. 更长的路径优先级更高
4. servemux is pretty lightweight, no request method, no variables, no regex

## 2.4

1. 405 Method not allowed，尽量用 http 提供的常量
2. 一个请求只能调用一次 `w.WriteHeader()`
3. 写 header 要提前，调用 `w.Write()` 就相当于默认返回 200 成功了，`w.Header().Set()` 也要提前
4. `w.WriteHeader(405) and w.Write(...)` 简写为 `http.Error(w, 405, "Method Not Allowed")`
5. 每个请求，Go 会自动生成 `Date, Content-Length, Content-Type` 三个 header
6. Go 会通过 `http.DetectContentType()` 猜测 Content-Type，猜不出来就用 `Content-Type: application/octet-stream`
7. 区分不了纯文本和 JSON，所以要手动设置 `w.Header().Set("Content-Type", "application/json")`
8. `Header(), Set(), Add(), Del(), Get(), Values()`
9. header 会通过` textproto.CanonicalMIMEHeaderKey()` 自动规范化，大写开头，`-` 连字符
10. 使用 `w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}` 跳过规范
11. `Del()` 不能删除系统自动生成的 header，使用`Nil，w.Header()["Date"] = nil` 删除

## 2.5

1. `/snippet/view?id=1` query string
2. `r.URL.Query().Get()` return string value or empty `""`
3. check positive integer `strconv.Atoi()`
4. `io.Writer` type is an interface, `http.ResponseWriter` has a `w.Write()` method satisfies the interface
5. in practice, see `io.Writer` interface, pass `http.ResponseWriter`, will be sent as the body of the HTTP response

## 2.6

1. project structure
2. `cmd/` 特定于此 app 的代码
3. `internal/` 不特定于此 app 的代码，潜在的可重用部分代码
4. `ui/` user interface

## 2.7

1. use `html/template` to parse and render
2. `template.ParseFiles("./ui/html/pages/home.tmpl.html")` get template set
3. `ts.Execute(w, nil)` write template content as the response body, nil represents no dynamic data
4. `{{define "base"}}  {{end}}`
5. `{{template "title" .}}` and `{{template "main" .}}`
6. `.` represents any dynamic data
7. `ts.ExecuteTemplate(w, "base", nil)` use the base template you defined
8. 用 `{{template xx .}}` 调用其他模板，使用 `{{block xx .}} {{end}}` 在调用模板不存在时显示默认内容
9. `embed` 包可以将文件嵌入 Go 程序中，而不是从磁盘中读取，后面将说明

## 2.8

1. `http.FileServer` 通过 HTTP 从特定目录提供文件，处理 `/static/` 路由
2. `fileServer := http.FileServer(http.Dir("./ui/static/"))` 删去前导路径
3. `mux.Handle("/static/", http.StripPrefix("/static", fileServer))` 删去 static 前缀
4. 搜索文件前，自动通过 `path.Clean()` 清除 `.` `..`
5. 支持 Range Requests，大文件续传友好
6. 支持 `Last-Modified` 和 `If-Modified-Since` 标头，如果文件自用户上次请求以来没有更改，则发送 304 Not Modified 状态码而不是文件本身
7. Content-Type 是使用 `mime.TypeByExtension()` 函数从文件扩展名自动设置的，可以使用 `mime.AddExtensionType()` 函数添加自己的自定义扩展和内容类型

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    // 不会自动清理路径，需手动 filepath.Clean()
    http.ServeFile(w, r, "./ui/static/file.zip")
}
```

8. 通过增加空白 index.html 文件可以禁用文件夹目录列表，最简单也是最啥比的方法，最好通过自定义 `http.FileSystem` 来实现
9. https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
10. https://gist.github.com/alexedwards/3b40775846535d0014ab1ff477e4a568

## 2.9

1. handler 应该是满足 http.Handler interface 的，需要实现 `ServeHTTP(http.ResponseWriter, *http.Request)`
2. `handleFunc()` 是一个语法糖，transform it into a handler, 给我们的普通函数自动添加了该函数，是我们可以直接传入
3. `func ListenAndServe(addr string, handler Handler) error` 我们传递了 servemux 进去，本身不提供响应，将其传给第二个 handler
4. 收到 http 请求，调用 servemux 的 ServeHTTP 方法，该方法通过 URL 路径找到对应 handler，然后再调用 handler 中的 ServeHTTP
5. Go Web 就是一系列的调用 ServeHTTP 的过程
6. 所有请求都会有一个单独的 goroutine 提供服务，服务繁忙时，可能会同时调用代码，访问共享资源时需要防止竞争

# 3

1. command-line flags config
2. application log messages
3. dependencies
4. 集中错误处理

## 3.1

1. `addr := flag.String("addr", ":4000", "HTTP network address")` return a pointer, need reference `*`
2. `-addr=":3333"` command-line flag, it's optional
3. 自动转换为 String，不能转换就 log an error and exit
4. flag.Int() Bool() Float64()
5. `flag.Parse()` 读完必须 parse 一下才能用
6. `-help` automatically help with flag
7. `os.GetEnv()` 没有默认值，返回值只有字符串，没有自动的帮助信息
8. `-flag=true` == `-flag`
9. `-flag=false` 只有一种写法
10. `flag.StringVar(&cfg.addr, .....)` parse command-line flags into the memory address of pre-existing variables

## 3.2

1. go standard logger 会默认信息前加上日期时间，然后写入 standard error stream，会显示在终端上
2. log.Fatal 写入后还会调用 os.Exit(1) 退出
3. INFO - information messages - stdout + prefix "INFO" - 文件描述符1 可省略
4. ERROR - error messages - stderr + prefix "ERROR" + file name + line number - 文件描述符2
5. 自定义 logger `errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)`
6. log.Llongfile 包含完整文件路径，log.LUTC 输出 UTC 时间
7. | bitwise or to join flags
8. `go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log`，`>>` 追加模式，文件不存在时创建一个
9. 创建一个带有 addr handler errorLog 的自定义 http.Server 结构体来调用 ListenAndServe


## 3.3

1. put dependencies into a custom `application` struct
2. then define your handler as methods against `application`
3. https://www.alexedwards.net/blog/organising-database-access

## 3.4

1. debug.Stack() get a stack trace for the current goroutine and append it to the log message
2. 将错误处理全部交由 app 自己的 methods，做一层抽象

## 3.5

1. 抽象出路由
2. main 函数只留处理 flags，创建依赖，启动服务器

# 4

1. 用 MySQL 存储数据
2. database driver 充当 MySQL 和 app 之间的中介
3. 独立的 models 包，数据库逻辑与 app 逻辑分离
4. database/sql 包
5. 防止 SQL 注入
6. transactions 事务

## 4.1

1. 安装
2. 新建数据库
3. 新建表
4. 插入数据
5. 增加用户，赋予权限


```sql
-- sudo apt install mysql-server
-- sudo mysql

-- 字符集和排序规则
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE snippetbox;

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- 在 created 字段创建索引，提高查询性能
CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

CREATE USER 'web'@'localhost';
-- snippetbox.* 表示该数据库的所有表
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

-- mysql -D snippetbox -u web -p
-- -D 指定数据库
-- -u 指定用户
-- -p 输入密码
SELECT id, title, expires FROM snippets;
DROP TABLE snippets;
```

## 4.2

1. `go get github.com/go-sql-driver/mysql@v1`
2. `@version`

## 4.3

1. `go get -u` to update
2. `@none` to remove
3. delete all reference in code, `go mod tidy -v` to remove

## 4.4

1. 连接数据库 `sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")`
2. `parseTime` 将 SQL 中的 TIME 和 DATE 字段转换为 Go 中的 `time.Time`
3. 返回的是一个 `sql.DB` object, it's a pool of many connections
4. 实际连接后是在第一次使用时，lazily establish
5. 连接池长期存在，在 main 中定义一次，传递给其他 handlers
6. `_ "github.com/go-sql-driver/mysql"` 使用空白标识符导入后运行 driver's init() 向 `database/sql` 注册自己
7. `defer db.Close()` 关闭

## 4.5