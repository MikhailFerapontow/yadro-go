package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/MikhailFerapontow/yadro-go/internal/adapters/handler"
	"github.com/MikhailFerapontow/yadro-go/internal/adapters/repository"
	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/internal/core/services"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var (
	service *services.ComicService
)

func main() {
	var configPath string
	var searchQuery string
	var searchByIndex bool
	flag.StringVar(&configPath, "c", ".", "Path to config folder")
	flag.StringVar(&searchQuery, "s", "", "Query for searching comics")
	flag.BoolVar(&searchByIndex, "i", false, "Enable search by index")

	flag.Parse()

	config.MustLoad(configPath)

	client := repository.NewCLient(viper.GetString("source_url"))
	db := repository.NewDbApi(viper.GetString("db_file"))
	stemmer := repository.InitStemmer()

	service = services.NewComicService(db, stemmer, client)

	InitRoutes()
}

func InitRoutes() {
	router := http.NewServeMux()

	handler := handler.NewComicHandler(service)

	router.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		ctx, stop := signal.NotifyContext(r.Context(), os.Interrupt)
		defer stop()
		handler.GetComics(ctx)
	})

	router.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		handler.Find(r.Context(), search)
	})

	go func() {
		if err := http.ListenAndServe(viper.GetString("server.port"), router); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	StartCroneJob(handler)
}

func StartCroneJob(handler *handler.ComicHandler) {
	c := cron.New()

	c.AddFunc(viper.GetString("server.cron"), func() {
		handler.GetComics(context.Background())
	})

	c.Run()
	defer c.Stop()
}
