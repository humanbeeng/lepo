package current

import "fmt"

type NiceGreeter struct {
	Name string
}

func (nc NiceGreeter) Greet(msg string) (string, error) {
	fmt.Println("Greeting nicely", msg)
	return msg, nil
}
