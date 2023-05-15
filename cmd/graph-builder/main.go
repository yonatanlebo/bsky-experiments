package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/bluesky-social/indigo/events"
	intEvents "github.com/ericvolp12/bsky-experiments/pkg/events"
	"github.com/ericvolp12/bsky-experiments/pkg/graph"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Count is used for sorting and storing mention counts
type Count struct {
	Handle string
	Count  int
}

const (
	maxBackoff       = 30 * time.Second
	maxBackoffFactor = 2
)

var tracer trace.Tracer

func main() {
	ctx := context.Background()
	rawlog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %+v\n", err)
	}
	defer func() {
		err := rawlog.Sync()
		if err != nil {
			fmt.Printf("failed to sync logger on teardown: %+v", err.Error())
		}
	}()

	log := rawlog.Sugar().With("source", "graph_builder_main")

	log.Info("starting graph builder...")

	// Registers a tracer Provider globally if the exporter endpoint is set
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		log.Info("initializing tracer...")
		shutdown, err := installExportPipeline(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// Replace with the WebSocket URL you want to connect to.
	u := url.URL{Scheme: "wss", Host: "bsky.social", Path: "/xrpc/com.atproto.sync.subscribeRepos"}

	includeLinks := os.Getenv("INCLUDE_LINKS") == "true"

	workerCount := 5

	graphFile := os.Getenv("BINARY_GRAPH_FILE")
	if graphFile == "" {
		graphFile = "social-graph.bin"
	}

	postRegistryEnabled := false
	dbConnectionString := os.Getenv("REGISTRY_DB_CONNECTION_STRING")
	if dbConnectionString != "" {
		postRegistryEnabled = true
	}

	sentimentAnalysisEnabled := false
	sentimentServiceHost := os.Getenv("SENTIMENT_SERVICE_HOST")
	if sentimentServiceHost != "" {
		sentimentAnalysisEnabled = true
	}

	log.Info("initializing BSky Event Handler...")
	bsky, err := intEvents.NewBSky(
		ctx,
		log,
		includeLinks, postRegistryEnabled, sentimentAnalysisEnabled,
		dbConnectionString, sentimentServiceHost,
		workerCount,
	)
	if err != nil {
		log.Fatal(err)
	}

	binReaderWriter := graph.BinaryGraphReaderWriter{}

	log.Info("reading social graph from binary file...")
	resumedGraph, err := binReaderWriter.ReadGraph(graphFile)
	if err != nil {
		log.Infof("error reading social graph from binary: %s\n", err)
		log.Info("creating a new graph for this session...")
	} else {
		log.Info("social graph resumed successfully")
		bsky.SocialGraph = resumedGraph
	}

	graphTicker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	wg := &sync.WaitGroup{}

	// Run a routine that dumps graph data to a file every 5 minutes
	wg.Add(1)
	go func() {
		log = log.With("source", "graph_dump")
		log.Info("starting graph dump routine...")
		ctx := context.Background()
		for {
			select {
			case <-graphTicker.C:
				saveGraph(ctx, log, bsky, &binReaderWriter, graphFile)
			case <-quit:
				graphTicker.Stop()
				wg.Done()
				return
			}
		}
	}()

	// Server for pprof and prometheus via promhttp
	go func() {
		log = log.With("source", "pprof_server")
		log.Info("starting pprof and prometheus server...")
		// Create a handler to write out the plaintext graph
		http.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
			log.Info("writing graph to HTTP Response...")
			bsky.SocialGraphMux.Lock()
			defer bsky.SocialGraphMux.Unlock()

			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", "attachment; filename=social-graph.txt")
			w.Header().Set("Content-Transfer-Encoding", "binary")
			w.Header().Set("Expires", "0")
			w.Header().Set("Cache-Control", "must-revalidate")
			w.Header().Set("Pragma", "public")

			err := bsky.SocialGraph.Write(w)
			if err != nil {
				log.Errorf("error writing graph: %s", err)
			} else {
				log.Info("graph written to HTTP Response successfully")
			}
		})

		http.Handle("/metrics", promhttp.Handler())
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// Run a routine that handles the events from the WebSocket
	log.Info("starting repo sync routine...")
	err = handleRepoStreamWithRetry(ctx, bsky, u, &events.RepoStreamCallbacks{
		RepoCommit: bsky.HandleRepoCommit,
		RepoInfo:   intEvents.HandleRepoInfo,
		Error:      intEvents.HandleError,
	})

	if err != nil {
		log.Errorf("Error: %v\n", err)
	}

	log.Info("waiting for routines to finish...")
	close(quit)
	wg.Wait()
	log.Info("routines finished, exiting...")
}

func getHalfHourFileName(baseName string) string {
	now := time.Now()
	min := now.Minute()
	halfHourSuffix := "00"
	if min >= 30 {
		halfHourSuffix = "30"
	}

	fileExt := filepath.Ext(baseName)
	fileName := strings.TrimSuffix(baseName, fileExt)

	return fmt.Sprintf("%s-%s_%s%s", fileName, now.Format("2006_01_02_15"), halfHourSuffix, fileExt)
}

func saveGraph(
	ctx context.Context,
	log *zap.SugaredLogger,
	bsky *intEvents.BSky,
	binReaderWriter *graph.BinaryGraphReaderWriter,
	graphFileFromEnv string,
) {
	tracer := otel.Tracer("graph-builder")
	ctx, span := tracer.Start(ctx, "SaveGraph")
	defer span.End()

	timestampedGraphFilePath := getHalfHourFileName(graphFileFromEnv)

	// Acquire a lock on the graph before we write it
	span.AddEvent("saveGraph:AcquireGraphLock")
	bsky.SocialGraphMux.RLock()
	span.AddEvent("saveGraph:GraphLockAcquired")

	log.Infof("writing social graph to binary file, last updated: %s",
		bsky.SocialGraph.LastUpdate.Format("02.01.06 15:04:05"),
		"graph_last_updated_at", bsky.SocialGraph.LastUpdate,
	)

	// Write the graph to a timestamped file
	err := binReaderWriter.WriteGraph(bsky.SocialGraph, timestampedGraphFilePath)
	if err != nil {
		log.Errorf("error writing social graph to binary file: %s", err)
	}

	// Release the lock after we're done writing the inital file
	span.AddEvent("saveGraph:ReleaseGraphLock")
	bsky.SocialGraphMux.RUnlock()

	// Copy the file to the "latest" path, we don't need to lock the graph for this
	err = copyFile(timestampedGraphFilePath, graphFileFromEnv)
	if err != nil {
		log.Errorf("error copying binary file: %s", err)
	}

	log.Info("social graph written to binary file successfully")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

func installExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tracerProvider := newTraceProvider(exporter)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider.Shutdown, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("BSkyGraphBuilder"),
		),
	)

	if err != nil {
		panic(err)
	}

	// initialize the traceIDRatioBasedSampler
	traceIDRatioBasedSampler := sdktrace.TraceIDRatioBased(1)

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(traceIDRatioBasedSampler),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func getNextBackoff(currentBackoff time.Duration) time.Duration {
	if currentBackoff == 0 {
		return time.Second
	}

	return currentBackoff + time.Duration(rand.Int63n(int64(maxBackoffFactor*currentBackoff)))
}

func handleRepoStreamWithRetry(
	ctx context.Context,
	bsky *intEvents.BSky,
	u url.URL,
	callbacks *events.RepoStreamCallbacks,
) error {
	var backoff time.Duration

	for {
		log.Println("connecting to BSky WebSocket...")
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("failed to connect to websocket: %v", err)
		}
		defer c.Close()

		// Create a new context with a cancel function
		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Start a timer to check for graph updates
		updateCheckDuration := 30 * time.Second
		updateCheckTimer := time.NewTicker(updateCheckDuration)
		defer updateCheckTimer.Stop()

		// Run a goroutine to handle the graph update check
		go func() {
			for {
				select {
				case <-updateCheckTimer.C:
					bsky.SocialGraphMux.RLock()
					if bsky.SocialGraph.LastUpdate.Add(updateCheckDuration).Before(time.Now()) {
						log.Printf("The graph hasn't been updated in the past %v seconds, closing and reopening the WebSocket...", updateCheckDuration)
						bsky.SocialGraphMux.RUnlock()
						// Cancel the context and close the WebSocket connection
						// Reset the backoff timer
						backoff = 0
						// Set the lastupdate to now so we don't trigger this again
						bsky.SocialGraphMux.Lock()
						bsky.SocialGraph.LastUpdate = time.Now()
						bsky.SocialGraphMux.Unlock()
						cancel()
						c.Close()
						return
					}
					bsky.SocialGraphMux.RUnlock()
				case <-streamCtx.Done():
					return
				}
			}
		}()

		err = events.HandleRepoStream(streamCtx, c, callbacks)
		if err != nil {
			log.Printf("Error in event handler routine: %v\n", err)
			backoff = getNextBackoff(backoff)
			if backoff > maxBackoff {
				return fmt.Errorf("maximum backoff of %v reached, giving up", maxBackoff)
			}

			select {
			case <-time.After(backoff):
				log.Printf("Reconnecting after %v...\n", backoff)
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			return nil
		}
	}
}
