package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MikhailFerapontow/yadro-go/internal/adapters/handler"
	db "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/db"
	stemmer "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/stemmer"
	xkcd "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/xkcd"
	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/internal/core/services"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", ".", "Path to config folder")

	flag.Parse()

	config.MustLoad(configPath)

	client := xkcd.NewCLient(viper.GetString("source_url"))
	stemmer := stemmer.InitStemmer()
	db := db.NewDbApi(viper.GetString("db_file"))

	service := services.NewComicService(db, stemmer, client)

	InitRoutes(service)
}

func InitRoutes(service *services.ComicService) {
	router := http.NewServeMux()

	handler := handler.NewComicHandler(service)

	handler.GetComics(context.Background()) // update on startup

	router.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		ctx, stop := signal.NotifyContext(r.Context(), os.Interrupt)
		defer stop()
		new, total := handler.GetComics(ctx)

		type comicsResponse struct {
			New   int `json:"new"`
			Total int `json:"total"`
		}

		response := comicsResponse{
			New:   new,
			Total: total,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	router.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		comics := handler.Find(r.Context(), search)

		fmt.Println(comics)

		type comicResponse struct {
			Url string
		}

		response := make([]comicResponse, len(comics))
		for i, comic := range comics {
			response[i].Url = comic.Url

		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	server := &http.Server{
		Addr:           ":" + viper.GetString("server.port"),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
