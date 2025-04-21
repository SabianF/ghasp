package data_repos

import (
	"log"
	"net/http"

	"github.com/SabianF/ghasp/src/common/presentation/pages"
)

const RootUrl string = "/"

func RootHandleRequest(w http.ResponseWriter, r *http.Request) {
	rootPageComponent := pages.RootPage()

	err := rootPageComponent.Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
	}
}
