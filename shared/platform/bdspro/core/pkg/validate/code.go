package validate

type Code string

const (
	CodeRequired     Code = "required"
	CodeTooLong      Code = "too_long"
	CodeTooShort     Code = "too_short"
	CodeOutOfRange   Code = "out_of_range"
	CodeInvalid      Code = "invalid"
	CodeInvalidDate  Code = "invalid_date"
	CodeInvalidValue Code = "invalid_value"
)

func (c Code) String() string {
	if c == "" {
		return "invalid"
	}
	return string(c)
}
