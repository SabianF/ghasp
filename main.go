package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SabianF/ghasp/src/pages"

	_ "github.com/lib/pq"
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
	initDb()
	router := initRouter()
	listenAndServe(router)
}

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if (err != nil) {
		log.Fatal("Error loading .env")
	}
}

func initDb() {
	log.Println("Opening DB connection...")

	host := os.Getenv("DB_HOST")
	if (host == "") {
		log.Fatal("Failed to get DB_HOST from .env")
	}

	port := os.Getenv("DB_PORT")
	if (port == "") {
		log.Fatal("Failed to get DB_PORT from .env")
	}

	user := os.Getenv("DB_USER")
	if (host == "") {
		log.Fatal("Failed to get DB_USER from .env")
	}

	pass := os.Getenv("DB_PASS")
	if (host == "") {
		log.Fatal("Failed to get DB_PASS from .env")
	}

	name := os.Getenv("DB_NAME")
	if (host == "") {
		log.Fatal("Failed to get DB_NAME from .env")
	}

	dataSource := fmt.Sprintf(
		// The '' around password is to include any spaces
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
    host, port, user, pass, name,
	)

	db, err := sql.Open("postgres", dataSource)
	if (err != nil) {
		log.Fatal("Failed to open DB: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if (err != nil) {
		log.Fatal("Failed to establish connection to DB: ", err)
	}

	log.Printf("Successfully opened DB connection: %d", db.Stats().OpenConnections)
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
