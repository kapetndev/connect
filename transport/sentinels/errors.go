package sentinels

const (
	// MarshalError defines an error used when unable to marshal a payload.
	MarshalError = error("failed to marshal payload")
	// UnmarshalError defines an error used when unable to unmarshal a payload.
	UnmarshalError = error("failed to unmarshal payload")
)

// error is a simple wrapper for a sentinel error type. It is favoured over the
// use of `errors.New` because the latter is not a compile time constant.
type error string

func (e error) Error() string {
	return string(e)
}
