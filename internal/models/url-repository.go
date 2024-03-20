package models

type URLRepository interface {
	Initialize() error

	Get(shortKey string) (string, bool)

	Save(events []Event) error
}
