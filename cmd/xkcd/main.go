package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/MikhailFerapontow/yadro-go/internal/adapters/handler"
	"github.com/MikhailFerapontow/yadro-go/internal/adapters/repository"
	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/internal/core/services"
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

	// db := database.NewDbApi(viper.GetString("db_file"))
	// client := xkcd.NewCLient(viper.GetString("source_url"))
	// app := app.InitApp(db, client, viper.GetInt("parallel"))

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	// app.GetComics(ctx)
	// app.Find(searchQuery, searchByIndex)

	InitRoutes()
}

func InitRoutes() {
	router := http.NewServeMux()

	handler := handler.NewComicHandler(service)

	router.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		handler.GetComics(r.Context())
	})

	router.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		handler.Find(r.Context(), search)
	})

	if err := http.ListenAndServe(viper.GetString("server.port"), router); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}
