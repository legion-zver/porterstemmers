package porterstemmers

import (
    "regexp"
)

// Pattern for search and replace
type Pattern struct {
    Rx *regexp.Regexp
    To string
}