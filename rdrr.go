package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rdrr "github.com/mattlgy/rdrr/lib"
)

type RedirTplOpts struct {
	Dest string
}

func main() {
	slug := rdrr.AddDest("http://example.com", []string{"this", "is", "a", "test"})
	fmt.Println(slug)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/x/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		co, s := rdrr.Get(slug)
		switch {
		case co == nil, co.GetCount() < 0:
			c.String(404, "oops")
		case co.GetCount() == 0:
			c.Redirect(303, co.GetURL())
		case co.GetCount() > 0:
			c.Redirect(303, "/x/"+s)
		}
	})
	r.GET("/r/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		co, s := rdrr.Get(slug)
		switch {
		case co == nil, co.GetCount() < 0:
			c.String(404, "oops")
		case co.GetCount() == 0:
			c.HTML(200, "redir.html", gin.H{
				"dest": co.GetURL(),
				"word": co.GetWord(),
			})
		case co.GetCount() > 0:
			c.HTML(200, "redir.html", gin.H{
				"dest": "/r/" + s,
				"word": co.GetWord(),
			})
		}
	})

	r.Run() // listen and server on 0.0.0.0:8080
}
