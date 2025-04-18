package main

import (
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

	"github.com/SabianF/ghasp/src/common/data/models"
	db_postgres "github.com/SabianF/ghasp/src/common/data/sources"
	"github.com/SabianF/ghasp/src/common/presentation/components"
	"github.com/SabianF/ghasp/src/common/presentation/pages"

	"github.com/gorilla/mux"
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
	db_postgres.InitDb()
	defer db_postgres.CloseDb()
	models.CreateUserTable()
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
		db_postgres.CloseDb()
		os.Exit(1)
	}()
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRoot)
	router.HandleFunc("/table", handleTable)
	router.HandleFunc("/db", handleDb)
	router.HandleFunc("/multi-replace", handleMultiReplace)
	router.HandleFunc("/multi-replace-button", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("HX-Trigger", "one, two, three")
	})
	router.HandleFunc("/test-one", func(w http.ResponseWriter, r *http.Request) {
		pages.TestComponent("one: " + time.Now().String()).Render(r.Context(), w)
	})
	router.HandleFunc("/test-two", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		pages.TestComponent("two: " + time.Now().String()).Render(r.Context(), w)
	})
	router.HandleFunc("/test-three", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		pages.TestComponent("three: " + time.Now().String()).Render(r.Context(), w)
	})
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
	rootPageComponent := pages.RootPage()

	err := rootPageComponent.Render(r.Context(), w)
	if (err != nil) {
		log.Fatalf("Failed to render component: %v\n", err)
	}
}

func handleTable(w http.ResponseWriter, r *http.Request) {

	tablePropsRowsAndData := generateTableData(1, 10, 3)

	tableProps := components.TableProps{
		Page: "1",
		Headings: tablePropsHeadings,
		RowsAndColumns: tablePropsRowsAndData,
		Footers: tablePropsFooters,
	}

	tablePageProps := pages.TablePageProps{
		TableProps: tableProps,
	}

	tablePageComponent := pages.TablePage(tablePageProps)

	tablePageComponent.Render(r.Context(), w)
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

func handleMultiReplace(w http.ResponseWriter, r *http.Request) {
	page := pages.MultiReplacePage()
	page.Render(r.Context(), w)
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
