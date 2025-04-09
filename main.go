package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SabianF/ghasp/src/pages"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//? GHASP Stack
//? - Go
//? - HTMX (Front end) + Templ
//? - AlpineJS (Front-end interactivity)
//? - Shoelace (UI component library) (TailwindCSS)
//? - PostgreSQL (Database)

func main() {
	handleSigTerm()
	loadEnvironmentVariables()
	router := initRouter()
	listenAndServe(router)
}

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if (err != nil) {
		log.Fatal("Error loading .env")
	}
}

// Allows graceful shutdown
func handleSigTerm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received SIGTERM. Exiting...")
		os.Exit(1)
	}()
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRoot)
	serveStaticFiles(router)
	router.Use(middlewareLogRequests)
	http.Handle("/", router)
	return router
}

func serveStaticFiles(router *mux.Router) {
	const defaultDir string = "./assets"

	var dir string
	flag.StringVar(&dir, "dir", defaultDir, "the directory to serve files from.")
	flag.Parse()

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(dir)))
	assetsRouter := router.PathPrefix("/assets/")

	assetsRouter.Handler(assetsHandler)
}

func middlewareLogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			currentTime := time.Now().Format(time.RFC3339Nano)
			log.Printf("%s %s %s (%s)", currentTime, r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w, r)
		},
	)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	rootPageComponent := pages.RootPage()

	err := rootPageComponent.Render(r.Context(), w)
	if (err != nil) {
		log.Fatalf("Failed to render component: %v\n", err)
	}
}

func listenAndServe(router *mux.Router) {
	port := ":" + os.Getenv("APP_PORT")
	log.Printf("Starting server on port %s ...\n", port)

	err := http.ListenAndServe(port, router)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")

	} else if err != nil {
		log.Fatalf("Error with server: %s\n", err)
	}
}
