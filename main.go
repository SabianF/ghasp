package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/SabianF/ghasp/src/common/presentation/components"
	"github.com/SabianF/ghasp/src/common/data/models"
	"github.com/SabianF/ghasp/src/common/presentation/pages"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

//? GHASP Stack
//? - Go
//? - HTMX (Front end) + Templ
//? - AlpineJS (Front-end interactivity)
//? - Shoelace (UI component library) (TailwindCSS)
//? - PostgreSQL (Database)



var tablePropsHeadings = []string{
	"1 heading 1",
	"2 heading 2",
	"3 heading 3",
}

var tablePropsFooters = []string{
	"1 footer 1",
	"2 footer 2",
	"3 footer 3",
}

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

	db, err := pgx.Connect(context.Background(), dataSource)
	if (err != nil) {
		log.Fatal("Failed to open DB: ", err)
	}
	defer db.Close(context.Background())

	err = db.Ping(context.Background())
	if (err != nil) {
		log.Fatal("Failed to establish connection to DB: ", err)
	}

	log.Printf("Successfully opened DB connection: %s\n", db.Config().Database)

	models.CreateUserTable(db)
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
	router.HandleFunc("/db", handleDb)
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
			log.Printf("%s %s %s (%s)\n", currentTime, r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w, r)
		},
	)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tablePropsRowsAndData := generateTableData(1, 10, 3)

	tableProps := components.TableProps{
		Page: "1",
		Headings: tablePropsHeadings,
		RowsAndColumns: tablePropsRowsAndData,
		Footers: tablePropsFooters,
	}

	rootPageComponent := pages.RootPage(tableProps)

	err := rootPageComponent.Render(r.Context(), w)
	if (err != nil) {
		log.Fatalf("Failed to render component: %v\n", err)
	}
}

func handleDb(w http.ResponseWriter, r *http.Request) {
	page := r.Header.Get("table-page")

	pageInt, err := strconv.Atoi(page)
	if (err != nil) {
		log.Printf("Invalid table-page: [%v]. Using default page.\n", page)
		pageInt = 0
	}

	pageInt += 1
	page = strconv.Itoa(pageInt)
	tableData := generateTableData(pageInt, 10, 3)

	tableDataProps := components.TableDataProps{
		Page: page,
		RowsAndColumns: tableData,
	}

	tableDataComponent := components.TableData(tableDataProps)

	tableDataComponent.Render(r.Context(), w)
}

func generateTableData(page int, numRows int, numColumns int) [][]string {
	result := [][]string{}

	for row := 0; row < numRows; row++ {
		newColumn := []string{}

		for column := 0; column < numColumns; column++ {
			newColumnData := fmt.Sprintf(
				"row %d, column %d",
				row + ((page - 1) * numRows) + 1,
				column + 1,
			)
			newColumn = append(newColumn, newColumnData)
		}

		result = append(result, newColumn)
	}

	return result
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
