package models

import (
	"database/sql"
	"html/template"
)

type App struct {
	DB        *sql.DB
	Templates *template.Template
}
