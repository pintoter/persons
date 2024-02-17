package entity

const (
	Male   = "male"
	Female = "female"
)

type Nationality struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Person struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Surname     string        `json:"surname"`
	Patronymic  string        `json:"patronymic,omitempty"`
	Age         int           `json:"age"`
	Gender      string        `json:"gender"`
	Nationalize []Nationality `json:"nationalize"`
}
