package general

// Performer conatins details about each talent to be displayed on leader board.
type Performer struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	TotalScore uint16  `json:"totalScore"`
	Rank       uint16  `json:"rank"`
	Image      *string `json:"image"`
}

// LeaderBoradDTO conatins multiple performers.
type LeaderBoradDTO struct {
	SelfPerformer Performer   `json:"selfPerformer"`
	AllPerformers []Performer `json:"allPerformers"`
}

// TableName defines table name of the struct.
func (*Performer) TableName() string {
	return "programming_questions"
}
