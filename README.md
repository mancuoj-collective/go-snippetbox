# Snippetbox

Yet another gist app written in Go from book [Let's Go](https://letsgo.mancuoj.me).

## 2.1

- `go mod init x.x.x`
- `go run .`

## 2.2

- MVC
- handler - controller - app logic, write http response headers and bodies
- router(servemux) - map url patterns and handlers
- web server - listen for requests

- net/http
- `*http.Request` 包含请求信息
- `http.ResponseWriter`
- `http.ListenAndServe` - `ip:port`
- `:http` or `:http-alt` - 没有设置端口就从 `/etc/services` 获取

- 启动服务
- 侦听端口
- 接受新请求
- 将请求传给 servemux
- 检查 URL 路径
- 派发给匹配的 handler

- `/` 包罗万象，相当于通配符 `/**`
- `go run` - compile and execute binary in /tmp
  - `go run .`
  - `go run main.go`
  - `go run snippetbox.hh`
- `.` 代表当前文件夹

## 2.3

- / home
- /snippet/view snippetView
- /snippet/create snippetCreate

- url patterns - fixed paths and subtree paths (end with a trailing slash / )
  - subtree paths - `/foo` will automatically redirected `/foo/`
- 更长的路径优先级更高

- hostname matching
- `mux.HandleFunc("foo.example.org/", fooHandler)`
- `mux.HandleFunc("bar.example.org/", barHandler)`

- servemux is pretty lightweight
  - no request method
  - no variables
  - no regex

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
