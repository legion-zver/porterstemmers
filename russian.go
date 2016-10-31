package porterstemmers

import (     
    "sort"
    "errors"
    "regexp"
    "strings"    
)

// For test: https://jsfiddle.net/tvjhuynh/

// RussianPorterStemmer - для россии
type RussianPorterStemmer struct {
}

func (s *RussianPorterStemmer) attemptReplacePatterns(token string, patterns []Pattern) (string, error) {    
    for i:=0; i<len(patterns);i++ {
        if patterns[i].Rx != nil {
            if patterns[i].Rx.MatchString(token) {
                return patterns[i].Rx.ReplaceAllString(token, patterns[i].To), nil                
            }
        }
    }
    return "", errors.New("Not replace pattern")
}

func (s *RussianPorterStemmer) perfectiveGerund(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{perfectiveRx1,""}, Pattern{perfectiveRx2,""}})
}

func (s *RussianPorterStemmer) adjective(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{adjectiveRx,""}})
}

func (s *RussianPorterStemmer) participle(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{participleRx1,"$1"}, Pattern{participleRx2,""}})
}

func (s *RussianPorterStemmer) adjectival(token string) (string, error) {
    result, err := s.adjective(token)
    if err != nil {
        pariticipleResult, err := s.participle(result);
        if err == nil {
            result = pariticipleResult
        } else {
            return "", err
        }
    }
    return result, nil
}

func (s *RussianPorterStemmer) reflexive(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{reflexiveRx,""}})
}

func (s *RussianPorterStemmer) verb(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{verbRx1,"$1"}, Pattern{verbRx2,""}})
}

func (s *RussianPorterStemmer) noun(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{nounRx,""}})
}

func (s *RussianPorterStemmer) superlative(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{superlativeRx,""}})
}

func (s *RussianPorterStemmer) derivational(token string) (string, error) {
    return s.attemptReplacePatterns(token,
    []Pattern{Pattern{derivationalRx,""}})
}

// StemString - stem string
func (s *RussianPorterStemmer) StemString(token string) string {
    token = strings.TrimSpace(strings.ToLower(token))
    token = eeRx.ReplaceAllString(token, "е")
    rv := volwesRx.FindStringSubmatch(token)
    if rv == nil || len(rv) < 3 {
        return token
    } 
    head := rv[1]    
    r2 := volwesRx.FindStringSubmatch(rv[2])
    result, err := s.perfectiveGerund(rv[2])
    // Проверка исключений
    if len(r2) < 1 {
        if sort.SearchStrings(specificHeads, head) > 0 {
            return token
        }
    }
    if err != nil {        
        resultReflexive, err := s.reflexive(rv[2])
        if err != nil {
            resultReflexive = rv[2]
        }
        result, err = s.adjectival(resultReflexive)
        if err != nil {
            result, err = s.verb(resultReflexive)
            if err != nil {
                result, err = s.noun(resultReflexive);
                if err != nil {
                    result = resultReflexive
                }
            }
        }
    }
    result = andRx.ReplaceAllString(result, "")
    derivationalResult := result
    if r2 != nil && len(r2) > 2 && len(r2[2]) > 0 {
        derivationalResult, err = s.derivational(r2[2]);
        if err != nil {
            derivationalResult, err = s.derivational(result)
            if err != nil {
                derivationalResult = result
            }
        } else {
            derivationalResult = result
        }
    }
    superlativeResult, err := s.superlative(derivationalResult)
    if err != nil {
        superlativeResult = derivationalResult
    }
    superlativeResult = wnRx.ReplaceAllString(superlativeResult, "$1")
    superlativeResult = mzRx.ReplaceAllString(superlativeResult, "")
    return head + superlativeResult
}

// RX List
var (
    perfectiveRx1  = regexp.MustCompile("[ая]в(ши|шись)$")
    perfectiveRx2  = regexp.MustCompile("(ив|ивши|ившись|ывши|ывшись|ыв)$")
    adjectiveRx    = regexp.MustCompile("(ее|ие|ые|ое|ими|ыми|ей|ий|ый|ой|ем|им|ым|ом|его|ого|ему|ому|их|ых|ую|юю|ая|яя|ою|ею)$")
    participleRx1  = regexp.MustCompile("([ая])(ем|нн|вш|ющ|щ)$")
    participleRx2  = regexp.MustCompile("(ивш|ывш|ующ)$")
    reflexiveRx    = regexp.MustCompile("(ся|сь)$")
    verbRx1        = regexp.MustCompile("([ая])(ла|на|ете|йте|ли|й|л|ем|н|ло|но|ет|ют|ны|ть|ешь|нно)$")
    verbRx2        = regexp.MustCompile("(ила|ыла|ена|ейте|уйте|ите|или|ыли|ей|уй|ил|ыл|им|ым|ен|ило|ыло|ено|ят|ует|ит|ыт|ены|ить|ыть|ишь|ую|ю)$")
    nounRx         = regexp.MustCompile("(а|ев|ов|ие|ье|е|иями|ями|ами|еи|ии|и|ией|ей|ой|ий|й|иям|ям|ием|ем|ам|ом|о|у|ах|иях|ях|ы|ь|ию|ью|ю|ия|ья|я)$")
    superlativeRx  = regexp.MustCompile("(ейш|ейше)$")
    derivationalRx = regexp.MustCompile("(ост|ость)$")
    eeRx           = regexp.MustCompile("ё")
    volwesRx       = regexp.MustCompile("^(.*?[аеиоюяуыиэ])(.*)$")
    andRx          = regexp.MustCompile("и$")
    wnRx           = regexp.MustCompile("(н)н")
    mzRx           = regexp.MustCompile("ь$")
    // Исключительно для анализа некоторых матерных слов
    specificHeads  = []string{"","тво","ху"}
)