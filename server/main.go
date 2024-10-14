package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters"
	"github.com/sayden/counters/output"
)

type bodyInput struct {
	Cwd             string                   `json:"cwd"`
	CounterTemplate counters.CounterTemplate `json:"counter"`
}

func main() {
	r := gin.Default()
	gin.ForceConsoleColor()

	r.LoadHTMLGlob("./static/*.html")
	r.StaticFile("/main.css", "./static/main.css")
	r.StaticFile("/img.png", "./static/img.png")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/counter", func(c *gin.Context) {
		// Read request body
		body := bodyInput{}
		err := c.BindJSON(&body)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Capture current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// Restore working directory after the function ends
		defer func() {
			if err = os.Chdir(cwd); err != nil {
				log.Error(err)
			}
		}()
		// Request body contains the current working directory to use
		// This is relevant because we need to use relavite paths
		if err = os.Chdir(body.Cwd); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ParseTemplate requires a byte slice, this is because it Unmarshals the JSON on top
		// of a CounterTemplate struct with default values, overriding them with the JSON values
		byt, err := json.Marshal(body.CounterTemplate)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		newTemplate, err := counters.ParseTemplate(byt)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// get a canvas with the rendered counter. The canvas can be written to a io.Writer
		counterCanvas, err := output.GetCounterCanvas(&newTemplate.Counters[0], newTemplate)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Avoid caching the image in the browser or the image won't be updated after
		// the first request
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Writer.Header().Set("Cache-Control", "no-cache")

		// Write the image to the response body
		err = counterCanvas.EncodePNG(c.Writer)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/code", func(c *gin.Context) {
		// Read request body
		counterTemplate := counters.CounterTemplate{}
		err := c.BindJSON(&counterTemplate)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		f, err := os.Create("./static/img.png")
		if err != nil {
			log.Error(err)
			return
		}
		defer f.Close()

		if err := generateCounter(&counterTemplate, "/home/mcastro/projects/prototypes/ukraine", f); err != nil {
			log.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	})

	r.GET("/render", func(c *gin.Context) {
		c.HTML(http.StatusOK, "render.html", nil)
	})

	r.Run(":8080")
}

func generateCounter(template *counters.CounterTemplate, wd string, writer io.Writer) error {
	if writer == nil {
		return errors.New("writer is nil")
	}

	// Capture current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// Request body contains the current working directory to use
	// This is relevant because we need to use relavite paths
	if err = os.Chdir(wd); err != nil {
		return err
	}
	// Restore working directory after the function ends
	defer func() {
		if err = os.Chdir(cwd); err != nil {
			log.Error(err)
		}
	}()

	// ParseTemplate requires a byte slice, this is because it Unmarshals the JSON on top
	// of a CounterTemplate struct with default values, overriding them with the JSON values
	byt, err := json.Marshal(template)
	if err != nil {
		return err
	}
	newTemplate, err := counters.ParseTemplate(byt)
	if err != nil {
		return err
	}

	// get a canvas with the rendered counter. The canvas can be written to a io.Writer
	res, err := output.GetCounterCanvas(&newTemplate.Counters[0], newTemplate)
	if err != nil {
		return err
	}

	return res.EncodePNG(writer)
}

func logError(c *gin.Context, err error) {
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
