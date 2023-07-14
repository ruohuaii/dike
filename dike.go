package dike

func Work(ptr interface{}, sourceFile string) {
	matcher := NewMatcher(ptr)
	result, err := matcher.GetDefined("dike")
	if err != nil {
		panic(err)
	}
	name, short := matcher.GetStructName()
	builder := NewBuilder(name, short, result)
	err = builder.Build(sourceFile)
	if err != nil {
		panic(err)
	}
}
