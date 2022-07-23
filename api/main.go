package main

import (
	"context"
	"crypto/md5"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	key "github.com/fairhive-labs/ethkeygen/pkg"
	"github.com/gin-gonic/gin"
)

type signature struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	Signature  string `json:"signature"`
}

//go:embed assets templates
var tfs embed.FS

const (
	robotsTxt = `
User-agent: *
Disallow: /members/*
`
)

func generate(c *gin.Context) {
	m, err := key.GenerateN(10)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "index.html", m)
}

func bulkSignature(c *gin.Context) {
	// get message
	var m struct {
		Message string
	}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nQ := c.DefaultQuery("n", "1")
	// control iterations
	n, err := strconv.Atoi(nQ)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if n <= 0 || n > 10000 {
		n = 1
	}

	signatures := make([]signature, n)
	for i := 0; i < n; i++ {
		prk, a, err := key.Generate()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		s, err := key.SignMessage(prk, m.Message)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		signatures[i] = signature{prk, a, s}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"total":      n,
		"signatures": signatures,
		"message":    m.Message})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	t := template.Must(template.ParseFS(tfs, "templates/*"))
	r.SetHTMLTemplate(t)

	r.Use(cors)
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/favicon.ico", getFavicon)
	r.GET("/robots.txt", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte(robotsTxt))
	})
	r.PUT("/bulk-signature", bulkSignature)
	r.GET("/", generate)
	return r
}

func main() {
	r := setupRouter()

	var addr string
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	} else {
		addr = ":8080" // default port
	}

	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-quit
		log.Printf("🚨 Shutdown signal \"%v\" received\n", s)

		log.Printf("🚦 Here we go for a graceful Shutdown...\n")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("⚠️ HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("✅ Listening and serving HTTP on %s\n", addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("👹 HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("😴 Server stopped")

}

func cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func getFavicon(c *gin.Context) {
	file, err := tfs.ReadFile("assets/favicon.ico")
	etag := fmt.Sprintf("%x", md5.Sum(file))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("ETag", etag)
	c.Data(
		http.StatusOK,
		"image/x-icon",
		file,
	)
}
