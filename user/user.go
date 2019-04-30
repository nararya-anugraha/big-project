package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// user.handler

type UserModuleType struct {
	db *sqlx.DB
}

func Wire(router *httprouter.Router, db *sqlx.DB) {
	userModule := UserModuleType{
		db: db,
	}
	router.POST("/api/user", userModule.getUsersHandler)
}

// user.data

type User struct {
	UserID     int64      `db:"user_id" json:"user_id"`
	UserName   *string    `db:"user_name" json:"user_name"`
	MSISDN     *string    `db:"msisdn" json:"msisdn"`
	UserEmail  *string    `db:"user_email" json:"user_email"`
	BirthDate  *time.Time `db:"birth_date" json:"birth_date"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	UpdateTime *time.Time `db:"update_time" json:"update_time"`
	UserAge    *string    `db:"user_age" json:"user_age"`
}

type PagedUsers struct {
	PageSize    int64   `json:"page_size"`
	CurrentPage int64   `json:"current_page"`
	TotalPages  int64   `json:"total_pages"`
	Users       []*User `json:"users"`
	Filter      string  `json:"filter"`
}

func (userModule *UserModuleType) getUsersHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pagedUsers := PagedUsers{}

	err := json.NewDecoder(request.Body).Decode(&pagedUsers)

	if err != nil {
		handleError(writer, err)
	}

	err = userModule.getUser(&pagedUsers)
	if err != nil {
		handleError(writer, err)
	}

	err = json.NewEncoder(writer).Encode(pagedUsers)
	if err != nil {
		handleError(writer, err)
	}
}

func handleError(writer http.ResponseWriter, err error) {
	encoder := json.NewEncoder(writer)
	encoder.Encode(err)

	writer.Header().Add("status", "500")
	writer.Header().Add("content-type", "application/json")

}

const query = `
select 
	user_id, 
	user_name,
	msisdn, 
	user_email, 
	birth_date, 
	create_time,
	update_time,
	to_char(age(now(), birth_date), 'YY "years," MM "months," DD "days"') user_age
from ws_user
`

type queryParametersType struct {
	Limit    int64  `db:"limit"`
	Offset   int64  `db:"offset"`
	UserName string `db:"user_name"`
}

func (userModule *UserModuleType) getUser(pagedUsers *PagedUsers) error {
	users := make([]*User, 0)

	queryParameters := queryParametersType{}

	filteredQuery := query
	if len(pagedUsers.Filter) > 0 {
		filteredQuery += " where user_name ~* :user_name"
		queryParameters.UserName = pagedUsers.Filter
	}

	// Get Total Pages
	countQuery := "select count(1) from (" + filteredQuery + " ) a "
	rows, err := userModule.db.NamedQuery(countQuery, queryParameters)
	defer rows.Close()
	if err != nil {
		return err
	}

	var totalRows int64
	rows.Next()
	rows.Scan(&totalRows)

	if totalRows == 0 {
		return nil
	}

	pageSize := pagedUsers.PageSize

	if totalRows%pageSize > 0 {
		totalRows += (pageSize - (totalRows % pageSize))
	}

	pagedUsers.TotalPages = (totalRows / pagedUsers.PageSize)

	// Get Users
	pagedQuery := filteredQuery + " limit :limit  offset :offset"
	statement, err := userModule.db.PrepareNamed(pagedQuery)
	if err != nil {
		return err
	}

	limit := pagedUsers.PageSize
	offset := (pagedUsers.CurrentPage - 1) * pagedUsers.PageSize

	queryParameters.Limit = limit
	queryParameters.Offset = offset

	userRows, err := statement.Queryx(queryParameters)
	defer userRows.Close()
	if err != nil {
		return err
	}

	for userRows.Next() {
		user := User{}
		err = userRows.StructScan(&user)
		if err != nil {
			return err
		}
		users = append(users, &user)
	}

	pagedUsers.Users = users
	return nil
}
