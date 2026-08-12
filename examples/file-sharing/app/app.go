package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const maxFileSize = 10 * 1024 * 1024 // 10 MB

// Start runs the web service.
func Start() {
	Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "127.0.0.1:8080"
	}

	log.Printf("file sharing service listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")

	switch r.Method {
	case http.MethodPost:
		upload(w, r, key)
	case http.MethodGet:
		download(w, r, key)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func upload(w http.ResponseWriter, r *http.Request, key string) {
	if key == "" {
		http.Error(w, "no file path given", http.StatusBadRequest)
		return
	}

	// Reject if key already exists
	_, err := Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		http.Error(w, "file already exists", http.StatusConflict)
		return
	}

	// Read body (limited to max size)
	limited := http.MaxBytesReader(w, r.Body, maxFileSize)
	data, err := io.ReadAll(limited)
	if err != nil {
		http.Error(w, "file too large (max 10 MB) or upload failed", http.StatusRequestEntityTooLarge)
		return
	}
	if len(data) == 0 {
		http.Error(w, "empty file", http.StatusBadRequest)
		return
	}

	_, err = Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Printf("s3 put failed: %v", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	log.Printf("uploaded %s (%d bytes)", key, len(data))
	w.WriteHeader(http.StatusCreated)
}

func download(w http.ResponseWriter, r *http.Request, key string) {
	if key == "" {
		http.Error(w, "no file path given", http.StatusBadRequest)
		return
	}

	out, err := Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *s3types.NoSuchKey
		if errors.As(err, &notFound) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		log.Printf("s3 get failed: %v", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	defer out.Body.Close()

	if out.ContentType != nil {
		w.Header().Set("Content-Type", *out.ContentType)
	}
	if out.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", *out.ContentLength))
	}
	io.Copy(w, out.Body)
}
