package frontpage

var (
	strHandler func(string)
)

func HandleString(f func(string)) {
	strHandler = f
}
