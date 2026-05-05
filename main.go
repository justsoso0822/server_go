package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HelloReq struct {
	Name string
	Age  int
}

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		var req HelloReq
		if err := r.Parse(&req); err != nil {
			r.Response.Write(err.Error())
			return
		}
		if req.Name == "" {
			r.Response.Write("name should not be empty")
			return
		}
		if req.Age <= 0 {
			r.Response.Write("age should be greater than 0")
			return
		}
		name := req.Name
		age := req.Age
		r.Response.Writef("Hello %s, your age is %d", name, age)
	})
	s.SetPort(8000)
	s.Run()
}
