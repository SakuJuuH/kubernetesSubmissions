package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	imageDirectory     = os.Getenv("IMAGE_DIR")
	port               = os.Getenv("PORT")
	imageUrl           = getEnvWithDefault("IMAGE_URL", "https://picsum.photos/300")
	allowedOrigins     = getEnvWithDefault("ALLOWED_ORIGINS", "*")
	cachedImageName    = getEnvWithDefault("CACHED_IMAGE_NAME", "current_image.jpg")
	imageCacheDuration = getCacheDuration()
)

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getCacheDuration() time.Duration {
	durationStr := os.Getenv("CACHE_DURATION_MINUTES")
	if durationStr == "" {
		return 10 * time.Minute
	}

	minutes, err := strconv.Atoi(durationStr)
	if err != nil || minutes <= 0 {
		fmt.Printf("Invalid CACHE_DURATION value '%s', defaulting to 10 minutes\n", durationStr)
		return 10 * time.Minute
	}

	return time.Duration(minutes) * time.Minute
}

type ImageInfo struct {
	Path     string    `json:"path"`
	CachedAt time.Time `json:"cached_at"`
}

func main() {
	if port == "" {
		fmt.Println("$PORT must be set")
		os.Exit(1)
	}

	if imageDirectory == "" {
		cwd, _ := os.Getwd()
		parentDir := filepath.Dir(cwd)
		imageDirectory = filepath.Join(parentDir, "image")
	}

	err := os.MkdirAll(imageDirectory, 0755)
	if err != nil {
		fmt.Printf("Error creating image directory: %v\n", err)
		os.Exit(1)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.Static("/api/image/files", imageDirectory)

	router.GET("/api/image", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Image Service API",
			"status":  http.StatusOK,
		})
	})

	router.GET("/api/image/current", func(c *gin.Context) {
		imageInfo, err := getCachedImage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, imageInfo)
	})

	router.POST("/api/image/shutdown", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Shutting down server..."})
		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}()
	})

	fmt.Printf("Server started in port %s\n", port)

	err = router.Run(":" + port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func getCachedImage() (*ImageInfo, error) {
	imagePath := filepath.Join(imageDirectory, cachedImageName)

	if fileInfo, err := os.Stat(imagePath); err == nil {
		cacheAge := time.Since(fileInfo.ModTime())
		if cacheAge < imageCacheDuration {
			fmt.Printf("Serving cached image (cached %v ago)\n", cacheAge)
			return &ImageInfo{
				Path:     "/files/" + cachedImageName,
				CachedAt: fileInfo.ModTime(),
			}, nil
		}
		fmt.Printf("Cache expired (%v old), downloading new image...\n", cacheAge)
	} else {
		fmt.Println("No cached image found, downloading new image...")
	}

	return downloadNewImage()
}

func downloadNewImage() (*ImageInfo, error) {
	imagePath := filepath.Join(imageDirectory, cachedImageName)

	fmt.Printf("Downloading new image from Lorem Picsum...\n")

	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create image file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %v: %v\n", file.Name(), err)
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %v", err)
	}

	now := time.Now()
	imageInfo := &ImageInfo{
		Path:     "/files/" + cachedImageName,
		CachedAt: now,
	}

	fmt.Printf("Successfully downloaded and cached new image\n")
	return imageInfo, nil
}
