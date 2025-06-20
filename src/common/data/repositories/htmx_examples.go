package data_repos

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SabianF/ghasp/src/common/data/models"
	"github.com/SabianF/ghasp/src/common/domain/entities"
	domain_repos "github.com/SabianF/ghasp/src/common/domain/repositories"
	"github.com/SabianF/ghasp/src/common/presentation/components"
	"github.com/SabianF/ghasp/src/common/presentation/pages"
)

func HtmxExamplesPageHandleRequest(w http.ResponseWriter, r *http.Request) {
	pageString := r.FormValue("page")
	if (pageString == "") {
		pageString = "1"
	}
	_, errParseInt := strconv.ParseInt(
		pageString, 10, 64,
	)
	if (errParseInt != nil) {
		errMsg := "valid page not provided: " + errParseInt.Error()
		log.Println(errMsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errMsg))
		return
	}

	userTableHeadings := models.GetUserFieldNames()

	rowsAndColumns, errGetTestTableData := getTestTableData()
	if (errGetTestTableData != nil) {
		errMsg := "cannot get user data: " + errGetTestTableData.Error()
		log.Println(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}

	htmxExamplesPageProps := pages.HtmxExamplesPageProps{
		LayoutProps: domain_repos.NewLayoutPropsDefault(),
		TestTableProps: components.TableProps{
			Page: pageString,
			Headings: userTableHeadings,
			RowsAndColumns: rowsAndColumns,
		},
	}

	htmxExamplesPage := pages.HtmxExamplesPage(htmxExamplesPageProps)
	errRender := htmxExamplesPage.Render(r.Context(), w)
	if errRender != nil {
		log.Printf("Failed to render page: %v\n", errRender)
		return
	}
}

func HtmxExamplesAddEntryHandleRequest(w http.ResponseWriter, r *http.Request) {
	nameFirst := r.FormValue("name_first")
	nameLast := r.FormValue("name_last")
	email := r.FormValue("email")

	newUser, errCreateUser := entities.NewUser(
		nameFirst,
		nameLast,
		email,
	)
	if (errCreateUser != nil) {
		log.Println("Failed to create new user:" + errCreateUser.Error())
		sendErrorNotification(w, errCreateUser, r)
		return
	}

	newRow := []string{
		newUser.User().Name_first,
		newUser.User().Name_last,
		newUser.User().Email,
	}

	// Update user data
	_, err := models.CreateUserModel(newUser).Create()
	if (err != nil) {
		errMsg := "failed to add entry: " + err.Error()
		log.Println(errMsg)
		sendErrorNotification(w, err, r)
		return
	}

	tableRowComponent := components.TableRow(components.TableRowProps{
		Columns: newRow,
	})

	errRender := tableRowComponent.Render(r.Context(), w)
	if (errRender != nil) {
		log.Println("Failed to render table row component: " + errRender.Error())
	}
}

func getTestTableData() ([][]string, error) {
	userData, errGetAllUsers := models.GetAllUsers()
	if (errGetAllUsers != nil) {
		return nil, errGetAllUsers
	}

	var rowsAndColumns = [][]string{}

	for _, user := range userData {
		rowsAndColumns = append(rowsAndColumns, []string{
			user.User().Name_first,
			user.User().Name_last,
			user.User().Email,
		})
	}
	return rowsAndColumns, nil
}

func sendErrorNotification(w http.ResponseWriter, errCreateUser error, r *http.Request) {
	w.Header().Add("HX-Target", "form-error")
	w.WriteHeader(http.StatusBadRequest)

	components.NotificationError(
		components.NotificationErrorProps{
			Message: errCreateUser.Error(),
		},
	).Render(r.Context(), w)
}
