package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

var (
	host         *string
	port         *int
	validBuckets = []string{
		"avatars",
		"subtitles",
		"posters",
	}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	port = flag.Int("port", 5555, "CDN server port")
	host = flag.String("host", "0.0.0.0", "CDN server host")

	configFileName := flag.String("config-file", "config.yaml", "config.yaml file")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if err := loadConfig(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	if err := sentry.Init(sentry.ClientOptions{Dsn: config.SentryDsn}); err != nil {
		log.Fatal(fmt.Errorf("could not initilize sentry: %v", err))
	}
}

func main() {

	defer func() {
		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		if ok := sentry.Flush(time.Second * 5); !ok {
			sentry.CaptureMessage("could not Flush sentry")
			log.Println("could not Flush sentry")
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENV") == "dev" {
		gin.SetMode(gin.DebugMode)
	}

	minioClient, err := minio.NewV4(
		config.Endpoint,
		config.AccessKey,
		config.SecretKey,
		config.UseHttps,
	)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalln(err)
	}

	if config.InsecureSkipVerify {
		minioClient.SetCustomTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	}

	router := gin.New()
	router.GET("/uploads/:bucket/:object_id", func(ctx *gin.Context) {

		bucketName := ctx.Param("bucket")
		if IsValidBucketName(bucketName) {
			mCtx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()

			output, err := minioClient.GetObjectWithContext(mCtx, ctx.Param("bucket"), ctx.Param("object_id"), minio.GetObjectOptions{})
			if err != nil {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
			if _, err := io.Copy(ctx.Writer, output); err != nil {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
		} else {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

	})

	log.Printf("Server running and listening on port [%s:%d]", *host, *port)
	if err := router.Run(fmt.Sprintf("%s:%d", *host, *port)); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}

func IsValidBucketName(bucketname string) bool {
	for _, bk := range validBuckets {
		if bucketname == bk {
			return true
		}
	}
	return false
}
