package options

type OptionsOps struct {
	// True: Every object is rendered with all of its relationships
	// False: Objects are rendered without their relationships
	DeepQuery bool
}

//Using variadic functions here, as it would be a bummer to pass
//far too many values.
func NewOptionsOps(extraOptions ...bool) OptionsOps {
	if extraOptions != nil {
		return OptionsOps{
			DeepQuery: extraOptions[0],
		}
	}

	return OptionsOps{
		DeepQuery: false,
	}
}