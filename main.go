package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"infoftex/infra"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	rand.Seed(time.Now().UnixNano())
}

/**
Filter params:
	--id=tt0010442
*/
func main() {
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)

	opts := infra.InitOptions()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opts.MaxRunTime))
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	imdb := infra.NewImdbws(opts.ImdbApiKey, opts.MaxApiRequests)
	analyzer := infra.NewAnalyzer(&wg, &imdb, opts.Filters, opts.ThreadsCount)
	file := infra.NewFileReader(&wg, opts.FilePath)

	lineCans := analyzer.InitChan(ctxTimeout)

	err := file.Open()
	if err != nil {
		infra.ShowErrorAndExit(err)
	}

	err = file.Read(ctxTimeout, lineCans)
	if err != nil {
		infra.ShowErrorAndExit(err)
	}

	wg.Done()
	wg.Wait()

	err = file.Close()
	if err != nil {
		infra.ShowErrorAndExit(err)
	}

	fmt.Println("lines count:", analyzer.Count)
	fmt.Println("timer:", time.Since(start))

	analyzer.Close(lineCans)
}
