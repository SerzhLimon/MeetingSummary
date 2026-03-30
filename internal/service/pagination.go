package service

import (
	"database/sql"
	"fmt"

	"github.com/SerzhLimon/MeetingSummary/internal/models"
)

// Generated Pagination Meta data
func pagination(table string, limit, page int) *models.Pagination {
	var (
		tmpl        = models.Pagination{}
		recordcount int
		DB          *sql.DB
	)

	// Count all record
	sqltable := fmt.Sprintf("SELECT count(id) FROM %s", table)

	DB.QueryRow(sqltable).Scan(&recordcount)

	total := (recordcount / limit)

	// Calculator Total Page
	remainder := (recordcount % limit)
	if remainder == 0 {
		tmpl.TotalPage = total
	} else {
		tmpl.TotalPage = total + 1
	}

	// Set current/record per page meta data
	tmpl.CurrentPage = page
	tmpl.RecordPerPage = limit

	// Calculator the Next/Previous Page
	if page <= 0 {
		tmpl.Next = page + 1
	} else if page < tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = page + 1
	} else if page == tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = 0
	}

	return &tmpl
}
