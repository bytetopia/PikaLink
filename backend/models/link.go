package models

import (
    "time"
)

type Link struct {
    ID          int       `json:"id" db:"id"`
    OriginalURL string    `json:"original_url" db:"original_url"`
    ShortCode   string    `json:"short_code" db:"short_code"`
    Title       string    `json:"title" db:"title"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    UserID      int       `json:"user_id" db:"user_id"`
}

type User struct {
    ID       int    `json:"id" db:"id"`
    Username string `json:"username" db:"username"`
    Password string `json:"-" db:"password"`
}

type CreateLinkRequest struct {
    OriginalURL    string `json:"original_url" binding:"required"`
    Title          string `json:"title"`
    CustomCode     string `json:"custom_code"`
    UseCustomCode  bool   `json:"use_custom_code"`
}

type UpdateLinkRequest struct {
    OriginalURL string `json:"original_url"`
    Title       string `json:"title"`
    ShortCode   string `json:"short_code"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required"`
}