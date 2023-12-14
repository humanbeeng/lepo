package testpackage

type FuncTestStruct struct{}

func (ft *FuncTestStruct) Hello() (string, error) {
	return "From testpackage", nil
}
