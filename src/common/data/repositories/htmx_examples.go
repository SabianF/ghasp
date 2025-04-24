package data_repos

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SabianF/ghasp/src/common/data/sources"
	"github.com/SabianF/ghasp/src/common/domain/entities"
	domain_repos "github.com/SabianF/ghasp/src/common/domain/repositories"
	"github.com/SabianF/ghasp/src/common/presentation/components"
	"github.com/SabianF/ghasp/src/common/presentation/pages"
)

func HtmxExamplesPageHandleRequest(w http.ResponseWriter, r *http.Request) {

	htmxExamplesPageProps := pages.HtmxExamplesPageProps{
		LayoutProps: domain_repos.NewLayoutPropsDefault(),
		TestTableProps: components.TableProps{
			Page: strconv.Itoa(sources.TestTablePageCurrent),
			Headings: sources.TestTableHeadings,
			RowsAndColumns: sources.TestTableData,
		},
	}

	htmxExamplesPage := pages.HtmxExamplesPage(htmxExamplesPageProps)
	err := htmxExamplesPage.Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
	}
}

func HtmxExamplesAddEntryHandleRequest(w http.ResponseWriter, r *http.Request) {

	// TODO: Validate inputs from data repo
	nameFirst := r.FormValue("name_first")
	nameLast := r.FormValue("name_last")
	email := r.FormValue("email")

	newEntry, err := entities.NewUser(
		nameFirst,
		nameLast,
		email,
	)
	if (err != nil) {
		log.Printf("Failed to create new user: %s\n", err.Error())

		w.Header().Add("HX-Target", "form-error")
		w.WriteHeader(http.StatusBadRequest)

		components.NotificationError(
			components.NotificationErrorProps{
				Message: err.Error(),
			},
		).Render(r.Context(), w)

		return
	}

	newRow := []string{
		newEntry.User().Name_first,
		newEntry.User().Name_last,
		newEntry.User().Email,
	}

	sources.TestTableData = append(sources.TestTableData, newRow)

	tableRowComponent := components.TableRow(components.TableRowProps{
		Columns: newRow,
	})

	tableRowComponent.Render(r.Context(), w)
}
