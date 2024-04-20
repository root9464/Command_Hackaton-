package models

type Homework struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	SubjectName string `json:"subject_name"`
	Task        string `json:"task"`
}

type Test struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	SubjectName string `json:"subject_name"`
	Task        string `json:"task"`
}
