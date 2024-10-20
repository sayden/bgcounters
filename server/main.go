package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/a-h/templ"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters"
	"golang.org/x/net/http2"
)

type bodyInput struct {
	Cwd             string                   `json:"cwd"`
	CounterTemplate counters.CounterTemplate `json:"counter"`
}

type counterImage struct {
	CounterImage string `json:"counter"`
	Id           string `json:"id"`
}

type response []counterImage

type responseMutex struct {
	sync.Mutex
	response
}

var startingFolder, _ = os.Getwd()
var globalResponse responseMutex

func main() {
	log.SetLevel(log.DebugLevel)

	router := gin.Default()
	gin.ForceConsoleColor()

	router.LoadHTMLGlob("./static/*.html")
	router.StaticFile("/main.css", "./static/main.css")
	router.StaticFile("/img.png", "./static/img.png")

	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })

	ch := make(chan bool)
	router.POST("/code", handlerCode(ch))
	router.GET("/render", func(c *gin.Context) {
		c.HTML(http.StatusOK, "render.html", nil)
	})
	router.GET("/listen", handlerListen(ch))
	router.GET("/state", func(c *gin.Context) {
		globalResponse.Lock()
		defer globalResponse.Unlock()
		component := Counters(globalResponse.response)
		c.Header("Cache-Control", "no-cache")
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	})

	// Create a custom HTTP server
	server := &http.Server{Addr: ":8090", Handler: router}

	// Configure HTTP/2
	http2.ConfigureServer(server, &http2.Server{})

	// log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
	log.Fatal(server.ListenAndServe())
}

func handlerListen(ch <-chan bool) func(c *gin.Context) {
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
			case <-ch:
				func() {
					fmt.Fprintf(c.Writer, "event: Grid\ndata:ok\n\n")
					flusher.Flush()
				}()
			}
		}
	}
}

func handlerCode(ch chan<- bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := os.Chdir(startingFolder); err != nil {
				log.Error(err)
			}
		}()

		byt, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err = counters.ValidateSchemaBytes[counters.CounterTemplate](byt); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		buf := new(bytes.Buffer)
		wc := base64.NewEncoder(base64.StdEncoding, buf)
		defer wc.Close()

		response, err := generateCounter(byt)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		globalResponse.Lock()
		defer globalResponse.Unlock()
		globalResponse.response = response
		ch <- true

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func generateCounter(byt []byte) (response, error) {
	// ParseTemplate requires a byte slice, this is because it Unmarshals the JSON on top
	// of a CounterTemplate struct with default values, overriding them with the JSON values
	tempTemplate, err := counters.ParseCounterTemplate(byt)
	if err != nil {
		return nil, err
	}

	newTemplate, err := tempTemplate.ParsePrototype()
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
		err := counter.EncodeCounter(wc, newTemplate)
		if err != nil {
			return nil, err
		}

		counterImage := counterImage{
			CounterImage: "data:image/png;base64," + buf.String(),
			Id:           counter.GetCounterFilename(i, "img", fileNumberPlaceholder, filenamesInUse),
		}

		i++
		fileNumberPlaceholder++

		response = append(response, counterImage)
		wc.Close()
	}

	log.Debug("generateCounters", "finished", "ok")
	return response, nil
}
