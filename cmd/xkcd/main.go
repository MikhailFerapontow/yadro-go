package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MikhailFerapontow/yadro-go/internal/adapters/handler"
	// db "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/db"
	"github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/db/sqlite"
	stemmer "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/stemmer"
	xkcd "github.com/MikhailFerapontow/yadro-go/internal/adapters/repository/xkcd"
	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/internal/core/services"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", ".", "Path to config folder")

	flag.Parse()

	config.MustLoad(configPath)

	client := xkcd.NewCLient(viper.GetString("source_url"))
	stemmer := stemmer.InitStemmer()

	db, err := sqlite.NewSqliteDB()
	if err != nil {
		log.Fatalf("db initialization failed with %s", err)
	}
	defer db.Close()

	err = sqlite.RunMigrations(db)
	if err != nil {
		log.Fatalf("migrations failed with %s", err)
	}
	log.Printf("migrations successful")

	api := sqlite.NewApiSqlite(db)

	service := services.NewComicService(api, stemmer, client)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	InitRoutes(ctx, service)
}

func InitRoutes(mainCtx context.Context, service *services.ComicService) {
	router := http.NewServeMux()

	handler := handler.NewComicHandler(service)

	port := checkPort(viper.GetInt("server.port"))
	log.Printf("Listening on port %d", port)

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    func(l net.Listener) context.Context { return mainCtx },
	}

	g, gCtx := errgroup.WithContext(mainCtx)

	g.Go(func() error {
		handler.GetComics(gCtx) // update on startup
		return nil
	})

	router.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		new, total := handler.GetComics(r.Context())

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

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		return StartCroneJob(mainCtx, handler)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("\nGraceful shutdown: %s \n", err)
	}
}

func StartCroneJob(ctx context.Context, handler *handler.ComicHandler) error {
	c := cron.New()
	c.AddFunc(viper.GetString("server.cron"), func() {
		handler.GetComics(ctx)
	})
	c.Start()

	<-ctx.Done()

	c.Stop()
	return ctx.Err()
}

// Check if given port is availiable
// If not, return first available port
func checkPort(port int) int {
	host := "127.0.0.1" // TODO: make configurable

	conn, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		newPort := listener.Addr().(*net.TCPAddr).Port
		fmt.Printf("Port %d is not available, using %d instead\n", port, newPort)
		return newPort
	}
	defer conn.Close()

	return port
}
