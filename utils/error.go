package utils

import (
	"log"
	"net/http"

	"forum/models"
)

func RenderError(app *models.App, w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	err := app.Templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"Message": message,
		"Status":  status,
	})
	if err != nil {
		log.Printf("Failed to render error.html: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
