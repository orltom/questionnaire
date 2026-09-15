package domain

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

func (v Visibility) Valid() bool {
	return v == Public || v == Private
}
