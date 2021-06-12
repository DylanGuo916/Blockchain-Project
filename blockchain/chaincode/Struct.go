package main

type Company struct {
	Id   string `json:"id"`

	CompanyName string `json:"company_name"`

	Legal string `json:"legal"` //法人

	Date string `json:"date"`

	Score []float64 `json:"score"` //得分

	Rank     string    `json:"rank"`     //信用等级
}
