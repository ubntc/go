package boxart

import (
	"fmt"
	"strings"
)

func Row(lcr ...string) (l, c, r string) {
	if len(lcr) == 1 {
		lcr = strings.Split(lcr[0], "")
	}

	n := len(lcr)
	if n < 3 {
		panic(fmt.Sprintf("bad box strings: %s", lcr))
	}
	return lcr[0], strings.Join(lcr[1:n-1], ""), lcr[n-1]
}

func BlockToString(name string, chars map[string]string) string {
	if s, ok := chars[name]; ok {
		return s
	}
	return "?"
}
