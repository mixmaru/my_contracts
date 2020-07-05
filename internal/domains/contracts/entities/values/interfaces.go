package values

type IValidation interface {
	Validate(value interface{}) (validErrors []int, err error)
}
