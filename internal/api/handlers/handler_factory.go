package handlers

import database "testProject/internal/data"

type Handler struct {
	DB *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{
		DB: db,
	}
}