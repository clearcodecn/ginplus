#### Ginplus 
> 扩展 gin 的 html 渲染形式.


#### Install 
```shell
go get github.com/clearcodecn/ginplus
```

#### Usage
```go
func main(){
    // 开启session. 
    ginplus.SessionKey = viper.GetString("session.key")
	
	// 创建默认 gin 
    g := ginplus.Default()
	// 设置 模板渲染. 
    g.HTMLRender = ginplus.NewTemplateManager("templates", ".gohtml", ginplus.IsDebugging())
    // 静态文件
    g.Static("/static", "static")

    // 动态设置主题
	g.Use(function(ctx *ginplus.Context){
        ctx.SetTheme("admin")
        
        // 动态设置模板变量
		ctx.Assign("key","value")
        
        ctx.Next()
    })
    
	g.Get("/index",func(ctx *ginplus.Context){
		// 渲染模板： templates/admin/index.gohtml
		ctx.Html(200,"index",ginplus.H{})
    })

    
	// session 管理
	ctx.SessionGet("key") // 获取session
	ctx.Session("key","val") // 设置session
}

```


#### flash session
设置flash
```go
	// flash Session
    ctx.AddFlash("username", req.Username)
    ctx.AddFlash("password", req.Password)
    ctx.AddFlash("message", "账号或密码错误")
```
在模板中使用
```html

<form class="layui-form" action="/admin/login" method="post">
  <div class="layui-form-item logo-title">
      <h1>用户登录</h1>
      {{ if hasSession "message" }}
          <p style="color: red; ">{{ flash "message" }}</p>
      {{ end }}
  </div>
  <div class="layui-form-item">
      <label class="layui-icon layui-icon-username" for="UserName"></label>
      <input type="text" name="username" id="username" placeholder="用户名或者邮箱" autocomplete="off"
             class="layui-input" value="{{ flash "username" }}">
  </div>
  <div class="layui-form-item">
      <label class="layui-icon layui-icon-password" for="Password"></label>
      <input type="password" name="password" id="password"
             placeholder="密码" autocomplete="off" class="layui-input" value="{{ flash "password" }}">
  </div>
  <div class="layui-form-item">
      <button class="layui-btn layui-btn-fluid" lay-submit lay-filter="form-login">登 入</button>
  </div>
</form>
```