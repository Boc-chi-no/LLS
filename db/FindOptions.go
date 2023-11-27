package db

type FindOptions struct {
	// A document specifying the inclusive lower bound for a specific index. The default value is 0, which means that
	// there is no minimum value.
	Min interface{}

	// A document specifying the exclusive upper bound for a specific index. The default value is nil, which means that
	// there is no maximum value.
	Max interface{}

	// The maximum number of documents to return. The default value is 0, which means that all documents matching the
	// filter will be returned. A negative limit specifies that the resulting documents should be returned in a single
	// batch. The default value is 0.
	Limit int64

	// The number of documents to skip before adding documents to the result. The default value is 0.
	Skip int64

	Key string

	PrefixScans bool
}

func Find() *FindOptions {
	return &FindOptions{}
}

// SetMin sets the value for the Min field.
func (f *FindOptions) SetMin(min interface{}) *FindOptions {
	f.Min = min
	return f
}

// SetMax sets the value for the Max field.
func (f *FindOptions) SetMax(max interface{}) *FindOptions {
	f.Max = max
	return f
}

// SetLimit sets the value for the Limit field.
func (f *FindOptions) SetLimit(i int64) *FindOptions {
	f.Limit = i
	return f
}

// SetSkip sets the value for the Skip field.
func (f *FindOptions) SetSkip(i int64) *FindOptions {
	f.Skip = i
	return f
}

// SetKey sets the value for the Key field.
func (f *FindOptions) SetKey(i string) *FindOptions {
	f.Key = i
	return f
}

// SetPrefixScans sets the value for the PrefixScans field.
func (f *FindOptions) SetPrefixScans(i bool) *FindOptions {
	f.PrefixScans = i
	return f
}
