package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
)

type resp struct {
	Code int
	Msg string `json:"msg"`
	Data interface{} `json:"data"`

}

func (r resp) response(c *gin.Context)  {
	c.JSON(r.Code, resp{
		Msg: r.Msg,
		Data: r.Data,
	})
	c.Abort()
}




func Index(c *gin.Context) {
	var res []*proxyItem
	showProxyInfo(&res)
	c.HTML(200, "index.html", map[string]interface{}{
		"data": res,
	})
}

func Delete(c *gin.Context) {
	request := &proxyItem{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		resp{
			Code: 400,
			Msg:  "参数解析失败",
			Data: nil,
		}.response(c)
		return
	}

	err = deleteProxy(request.ListenAddr, request.ListenPort)
	if err != nil {
		resp{
			Code: 400,
			Msg:  "删除失败: " + err.Error(),
			Data: nil,
		}.response(c)
		return
	}

	resp{
		Code: 200,
		Msg:  "删除完成",
		Data: nil,
	}.response(c)
	return
}

func Create(c *gin.Context) {
	request := &proxyItem{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		resp{
			Code: 400,
			Msg:  "参数解析失败",
			Data: nil,
		}.response(c)
		return
	}
	err = createProxy(request.ListenAddr, request.ListenPort, request.TargetAddr, request.TargetPort)
	if err != nil {
		resp{
			Code: 400,
			Msg:  "创建失败: " + err.Error(),
			Data: nil,
		}.response(c)
		return
	}
	resp{
		Code: 200,
		Msg:  "创建完成",
		Data: nil,
	}.response(c)
	return
}
//go:embed template
var tmpl embed.FS

//go:embed static
var static embed.FS

func StartService() {
	templateFile, _ := template.ParseFS(tmpl, "template/*")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	engine := gin.Default()
	engine.SetHTMLTemplate(templateFile)
	engine.StaticFS("/static", http.FS(static))
	engine.GET("/", Index)
	engine.POST("/delete", Delete)
	engine.POST("/create", Create)
	err := engine.Run("127.0.0.1:57391")
	if err != nil {
		ioutil.WriteFile("err.log", []byte(err.Error()), 0644)
	}
}