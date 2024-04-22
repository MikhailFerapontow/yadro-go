package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/pkg/app"
	"github.com/MikhailFerapontow/yadro-go/pkg/database"
	"github.com/MikhailFerapontow/yadro-go/pkg/xkcd"
	"github.com/spf13/viper"
)

func main() {
	var config_path string
	var search_query string
	flag.StringVar(&config_path, "c", ".", "path to config file")
	flag.StringVar(&search_query, "s", "", "Query for searching comics")

	flag.Parse()

	config.MustLoad(config_path)

	db := database.NewDbApi(viper.GetString("db_file"))
	client := xkcd.NewCLient(viper.GetString("source_url"))
	app := app.InitApp(db, client, viper.GetInt("parallel"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.GetComics(ctx)
	if search_query == "" {
		return
	}
	app.Find(search_query)
}
