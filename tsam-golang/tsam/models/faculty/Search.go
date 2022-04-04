package faculty

// Search Contain Email, Contatc, FirstName, City
type Search struct {
	Email     string `json:"email"`
	Contact   string `json:"contact"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
}
