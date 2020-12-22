package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/CastyLab/cdn.service/config"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

var (
	host *string
	port *int
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	port = flag.Int("port", 5555, "CDN server port")
	host = flag.String("host", "0.0.0.0", "CDN server host")

	configFileName := flag.String("config-file", "config.yaml", "config.yaml file")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if err := config.Load(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	if err := sentry.Init(sentry.ClientOptions{Dsn: config.Map.Secrets.SentryDsn}); err != nil {
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
	if config.Map.App.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	}

	minioClient, err := minio.NewV4(
		config.Map.Secrets.ObjectStorage.Endpoint,
		config.Map.Secrets.ObjectStorage.AccessKey,
		config.Map.Secrets.ObjectStorage.SecretKey,
		false,
	)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalln(err)
	}

	router := gin.New()
	router.GET("/uploads/:bucket/:object_id", func(ctx *gin.Context) {
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
	})

	log.Printf("Server running and listening on port [%s:%d]", *host, *port)
	if err := router.Run(fmt.Sprintf("%s:%d", *host, *port)); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}
