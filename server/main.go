package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters"
	"github.com/sayden/counters/output"
	"github.com/sayden/counters/transform"
	"golang.org/x/net/http2"
)

type bodyInput struct {
	Cwd             string                   `json:"cwd"`
	CounterTemplate counters.CounterTemplate `json:"counter"`
}

func main() {
	log.SetLevel(log.DebugLevel)

	router := gin.Default()
	gin.ForceConsoleColor()

	router.LoadHTMLGlob("./static/*.html")
	router.StaticFile("/main.css", "./static/main.css")
	router.StaticFile("/img.png", "./static/img.png")

	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	// router.POST("/counter", handlerCounter)

	// Those 3 endpoins are for the same purpose: to update an image in the browser using SSE
	// with the content that arrives as a POST request into the /code endpoint
	ch := make(chan *response)
	router.POST("/code", handlerCode(ch))
	router.GET("/render", func(c *gin.Context) { c.HTML(http.StatusOK, "render.html", nil) })
	router.Any("/listen", handlerListen(ch))

	// Create a custom HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Configure HTTP/2
	http2.ConfigureServer(server, &http2.Server{})

	// log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
	log.Fatal(server.ListenAndServe())
}

func handlerListen(ch <-chan *response) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Set headers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		if c.Request.ProtoMajor == 2 {
			fmt.Println("Client is using HTTP/2")
		} else {
			fmt.Println("Client is using HTTP/1.x")
		}

		// Create a channel to notify of client disconnect
		clientChan := make(chan bool)
		go func() {
			<-c.Request.Context().Done()
			clientChan <- true
		}()

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			http.Error(c.Writer, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Send events
		for {
			select {
			case <-clientChan:
				fmt.Println("Client disconnected")
				return
			case data := <-ch:
				byt, err := json.Marshal(data)
				if err != nil {
					log.Error(err)
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				fmt.Fprintf(c.Writer, "data:%s\n\n", string(byt))
				flusher.Flush()
			}
		}
	}
}

type counterImage struct {
	CounterImage string `json:"counter"`
	Id           string `json:"id"`
}

type response []counterImage

func handlerCode(ch chan<- *response) func(c *gin.Context) {
	return func(c *gin.Context) {
		byt, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// FIXME: Validate the schema
		// if err = counters.ValidateSchemaBytes(byt); err != nil {
		// 	log.Error(err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		buf := new(bytes.Buffer)
		wc := base64.NewEncoder(base64.StdEncoding, buf)
		defer wc.Close()

		response, err := generateCounter(byt, "/home/mcastro/projects/prototypes/ukraine")
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ch <- response

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// func handlerCounter(c *gin.Context) {
// 	// Read request body
// 	body := bodyInput{}
// 	err := c.BindJSON(&body)
// 	if err != nil {
// 		log.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	response, err := generateCounter(&body.CounterTemplate, body.Cwd)
// 	if err != nil {
// 		log.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	byt, err := json.Marshal(response)
// 	if err != nil {
// 		log.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	fmt.Fprintf(c.Writer, "data:%s\n\n", string(byt))
//
// 	c.JSON(http.StatusOK, gin.H{"status": "ok"})
// }

func generateCounter(byt []byte, wd string) (*response, error) {
	// Capture current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Request body contains the current working directory to use
	// This is relevant because we need to use relavite paths
	if err = os.Chdir(wd); err != nil {
		return nil, err
	}
	// Restore working directory after the function ends
	defer func() {
		if err = os.Chdir(cwd); err != nil {
			log.Error(err)
		}
	}()

	// ParseTemplate requires a byte slice, this is because it Unmarshals the JSON on top
	// of a CounterTemplate struct with default values, overriding them with the JSON values
	tempTemplate, err := counters.ParseTemplate(byt)
	if err != nil {
		return nil, err
	}

	newTemplate, err := transform.ParsePrototypedTemplate(tempTemplate)
	if err != nil {
		return nil, err
	}

	response := response(make([]counterImage, 0, len(newTemplate.Counters)))

	i := 0
	fileNumberPlaceholder := 0
	filenamesInUse := make(map[string]bool)
	for _, counter := range newTemplate.Counters {
		buf := new(bytes.Buffer)
		wc := base64.NewEncoder(base64.StdEncoding, buf)

		// get a canvas with the rendered counter. The canvas can be written to a io.Writer
		gc, err := output.GetCounterCanvas(&counter, newTemplate)
		if err != nil {
			return nil, err
		}

		if err = gc.EncodePNG(wc); err != nil {
			return nil, err
		}

		counterImage := counterImage{
			CounterImage: "data:image/png;base64," + buf.String(),
			Id:           counter.GetCounterFilename(i, "img", fileNumberPlaceholder, filenamesInUse),
		}

		i++
		fileNumberPlaceholder++

		response = append(response, counterImage)
		log.Debug("counterImage", "id", counterImage.Id)
		wc.Close()
	}

	return &response, nil
}
