package data_repos

import (
	"log"
	"net/http"

	domain_repos "github.com/SabianF/ghasp/src/common/domain/repositories"
	"github.com/SabianF/ghasp/src/common/presentation/pages"
)

const HtmxExamplesUrl string = "/htmx-examples"

func HtmxExamplesHandleRequest(w http.ResponseWriter, r *http.Request) {

	htmxExamplesPageProps := pages.HtmxExamplesPageProps{
		LayoutProps: domain_repos.NewLayoutPropsDefault(),
	}

	htmxExamplesPage := pages.HtmxExamplesPage(htmxExamplesPageProps)
	err := htmxExamplesPage.Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
	}
}
