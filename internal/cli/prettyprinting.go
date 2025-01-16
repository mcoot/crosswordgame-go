package cli

import "fmt"

func PrettyPrint(v interface{}) string {
	return fmt.Sprintf("%v\n", v)
}
