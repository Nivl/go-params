package params

// Scanner is an interface used to allow injecting data into a custom type
type Scanner interface {
	ScanString(date string) error
}
