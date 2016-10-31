package porterstemmers

import (    
    "regexp"
    "strings"
	"strconv"    
)

// Tokenizer - токенизер
type Tokenizer struct {
}

// Tokenize - токенизация
func (t *Tokenizer) Tokenize(text string) (words map[string]int64) {
    words = make(map[string]int64) 
    rusPS := RussianPorterStemmer{}
    
    // Разбиваем на слова
    for _, w := range strings.Split(replaceWSpacesRx.ReplaceAllString(strings.ToLower(text), " ")," ") {
        if len(w) > 2 {
            if _, err := strconv.Atoi(w); err != nil {
                stem := rusPS.StemString(w)                
                if len(stem) > 2 {
                    words[stem]++                    
                }
            }
        }
    }
    return
}

var replaceWSpacesRx = regexp.MustCompile(`[\s|\n]+`)