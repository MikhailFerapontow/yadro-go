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
	var configPath string
	var searchQuery string
	var searchByIndex bool
	flag.StringVar(&configPath, "c", ".", "Path to config file")
	flag.StringVar(&searchQuery, "s", "", "Query for searching comics")
	flag.BoolVar(&searchByIndex, "i", false, "Enable search by index")

	flag.Parse()

	config.MustLoad(configPath)

	db := database.NewDbApi(viper.GetString("db_file"))
	client := xkcd.NewCLient(viper.GetString("source_url"))
	app := app.InitApp(db, client, viper.GetInt("parallel"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.GetComics(ctx)
	if searchQuery == "" {
		return
	}
	app.Find(searchQuery, searchByIndex)
}
