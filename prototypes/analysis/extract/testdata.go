package extract

type InnerStruct struct {
	InnerField int
}

type OuterStruct struct {
	Field1 int
	Field2 string
	Anon   struct {
		InAnon    string
		InnerAnon struct {
			Haha string
		}
	}
}
