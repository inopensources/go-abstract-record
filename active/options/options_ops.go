package options

type Options struct {
	// True: Every object is rendered with all of its relationships
	// False: Objects are rendered without their relationships
	DeepQuery bool
}

//Using variadic functions here, as it would be a bummer to pass
//far too many values.
func NewOptions(values ...bool) *Options {
	return &Options{
		DeepQuery: values[0],
	}
}