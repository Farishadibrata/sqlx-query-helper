package models

type ExTable struct {
	Id           int    `json:"id" db:"id"`
	Description  string `json:"description" db:"description"`
	ProjectCode  string `json:"project_code" db:"project_code"`
	UserID       int    `json:"user_id" db:"user_id"`
	RandomNumber int    `json:"random_number" db:"random_number"`
}

type UserCredential struct {
	ProjectCode string `json:"project_code" db:"project_code"`
}
