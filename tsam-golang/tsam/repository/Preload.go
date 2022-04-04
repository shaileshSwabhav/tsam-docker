package repository

// Preload contains preload schema and query details.
type Preload struct {
	Schema          string
	Queryprocessors []QueryProcessor
}
