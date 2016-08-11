package main

import (
	"github.com/gin-gonic/gin"
	rdrr "github.com/mattlgy/rdrr/lib"
	"strings"
)

type RedirTplOpts struct {
	Dest string
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.POST("/", func(c *gin.Context) {
		dest := c.PostForm("dest")
		if len(dest) > 256 {
			c.String(400, "URL too long")
			return
		}

		str := c.PostForm("words")
		if len(str) > 256 {
			c.String(400, "Text too long")
			return
		}
		words := strings.Split(str, " ")
		if len(words) > 32 {
			c.String(400, "Text too long")
			return
		}

		slug := rdrr.AddDest(dest, words)

		c.String(200, "http://rdrr.at/"+slug)
	})

	r.GET("/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		co, s := rdrr.Get(slug)
		switch {
		case co == nil, co.GetCount() < 0:
			c.String(404, "oops")
		case co.GetCount() == 0:
			c.HTML(200, "redir.html", &gin.H{
				"dest": co.GetURL(),
				"word": co.GetWord(),
			})
		case co.GetCount() > 0:
			c.HTML(200, "redir.html", &gin.H{
				"dest": "/" + s,
				"word": co.GetWord(),
			})
		}
	})

	r.Run(":80")
}
