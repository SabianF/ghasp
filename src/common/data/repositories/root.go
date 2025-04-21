package data_repos

import (
	"log"
	"net/http"

	domain_repos "github.com/SabianF/ghasp/src/common/domain/repositories"
	"github.com/SabianF/ghasp/src/common/presentation/pages"
)

func RootPageHandleRequest(w http.ResponseWriter, r *http.Request) {

	rootPageProps := pages.RootPageProps{
		LayoutProps: domain_repos.NewLayoutPropsDefault(),
	}

	rootPageComponent := pages.RootPage(rootPageProps)

	err := rootPageComponent.Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
	}
}
