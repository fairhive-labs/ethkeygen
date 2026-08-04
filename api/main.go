package main

import (
	"cmp"
	"context"
	"crypto/sha256"
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

// The favicon is embedded and immutable, so read it and compute its ETag
// once at startup instead of on every request.
var (
	faviconData []byte
	faviconETag string
)

func init() {
	var err error
	faviconData, err = tfs.ReadFile("assets/favicon.ico")
	if err != nil {
		log.Fatalf("failed to read embedded favicon: %v", err)
	}
	faviconETag = fmt.Sprintf("%x", sha256.Sum256(faviconData))
}

const (
	robotsTxt = `
User-agent: *
Disallow: /members/*
`
	// maxKeys bounds the number of keys the / endpoint will generate per
	// request, preventing resource exhaustion from an arbitrarily large ?l=.
	maxKeys = 1000
)

func generate(c *gin.Context) {
	l := c.DefaultQuery("l", strconv.Itoa(10))
	ln, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// bound the request to avoid CPU/memory exhaustion from huge ?l= values
	ln = max(1, min(ln, maxKeys))

	m, err := key.GenerateN(ln)
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
	for i := range n {
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
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}).ParseFS(tfs, "templates/*"))
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

	addr := ":" + cmp.Or(os.Getenv("PORT"), "8080")

	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Listen for OS termination signals and cancel the context on the first one.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("✅ Listening and serving HTTP on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("👹 HTTP server ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	stop() // stop intercepting signals so a second one force-quits
	log.Printf("🚦 Here we go for a graceful Shutdown...\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ HTTP server Shutdown: %v", err)
	}
	log.Printf("😴 Server stopped")
}

func cors(c *gin.Context) {
	// The service is public and uses no cookies/credentials, so a wildcard
	// origin is fine. Credentials must NOT be combined with "*" (browsers
	// reject it), and the advertised methods must match the real routes.
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func getFavicon(c *gin.Context) {
	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, faviconETag) {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("ETag", faviconETag)
	c.Data(
		http.StatusOK,
		"image/x-icon",
		faviconData,
	)
}
