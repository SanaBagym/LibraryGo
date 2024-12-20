package internal

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
}
