package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

//? GHASP Stack
//? - Go
//? - HTMX (Front end) + Templ
//? - AlpineJS (Front-end interactivity)
//? - Shoelace (UI component library) (TailwindCSS)
//? - PostgreSQL (Database)

func main() {
	handleSigTerm()
	router := initRouter()
	listenAndServe(router)
}

// Allows graceful shutdown
func handleSigTerm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Received SIGTERM. Exiting...")
		os.Exit(1)
	}()
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRoot)
	http.Handle("/", router)
	return router
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func listenAndServe(router *mux.Router) {
	const port string = ":3333"
	fmt.Printf("Starting server on port %s ...\n", port)

	err := http.ListenAndServe(port, router)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")

	} else if err != nil {
		fmt.Printf("Error with server: %s\n", err)
		os.Exit(1)
	}
}
