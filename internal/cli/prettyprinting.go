package cli

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mcoot/crosswordgame-go/internal/apitypes"
)

func PrettyPrint(v interface{}) {
	if ok := apiPrettyPrint(v); !ok {
		spew.Dump(v)
	}
}

func apiPrettyPrint(v interface{}) bool {
	switch v := v.(type) {
	case *apitypes.HealthcheckResponse:
		printHealthcheckResponse(v)
		return true
	case []interface{}:
		spew.Dump(v)
		return true
	}
	return false
}

func printHealthcheckResponse(v *apitypes.HealthcheckResponse) {
	fmt.Printf(`Server health:
  Start Time: %s
`, v.StartTime)
}
