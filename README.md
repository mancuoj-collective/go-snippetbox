# Snippetbox

Yet another gist app written in Go from book [Let's Go](https://letsgo.mancuoj.me).

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
