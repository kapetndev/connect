package recovery_test

const panicMessage = "very bad thing happened"

// nilPanic is used to prevent the static code analysis tool from warning of
// a panic with a nil value.
var nilPanic any = nil

func returnPanics(m string) {
	switch m {
	case "panic":
		panic(panicMessage)
	case "nilPanic":
		panic(nilPanic)
	}
}
