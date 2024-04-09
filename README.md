
# 自用go工具类


1. http请求工具类

Usage:

GET:
```go
httpResult :=  httputil.Get[string]("http://www.example.com")
```

POST:
```go
form:= url.Values{}
form.Add("name", "binbug")
httpResult :=  httputil.PostForm[string]("http://www.example.com", form)
```

2. json工具类
