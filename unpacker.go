// Unpacker for Dean Edward's p.a.c.k.e.r
// Ported from: https://github.com/beautifier/js-beautify/blob/main/python/jsbeautifier/unpackers/packer.py
package unpacker

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type DEUnpacker struct {
	Source string
	Prefix string
	Suffix string
}

func NewDEUnpacker(source string) (*DEUnpacker, bool) {
	deu := &DEUnpacker{
		Source: source,
	}
	beginOffset := -1
	indexes := regexp.MustCompile(`eval[ ]*\([ ]*function[ ]*\([ ]*p[ ]*,[ ]*a[ ]*,[ ]*c[ ]*,[ ]*k[ ]*,[ ]*e[ ]*,[ ]*`).FindStringIndex(source)
	if indexes != nil {
		beginOffset = indexes[0]
		deu.Prefix = source[:beginOffset]
	}
	if beginOffset == -1 {
		return nil, false
	}
	sourceEnd := source[beginOffset:]
	if arr := strings.SplitN(sourceEnd, "')))", 2); arr[0] == sourceEnd {
		if arr2 := strings.SplitN(sourceEnd, "}))", 2); len(arr2) > 1 {
			deu.Suffix = arr2[1]
		}
	} else {
		deu.Suffix = arr[1]
	}
	return deu, true
}

func filterArgs(source string) (string, []string, int, int, error) {
	juicers := []string{
		`(?s)}\('(.*)', *(\d+|\[\]), *(\d+), *'(.*)'\.split\('\|'\), *(\d+), *(.*)\)\)`,
		`(?s)}\('(.*)', *(\d+|\[\]), *(\d+), *'(.*)'\.split\('\|'\)`,
	}
	for _, juicer := range juicers {
		args := regexp.MustCompile(juicer).FindStringSubmatch(source)
		if args != nil {
			if args[2] == "[]" {
				args[2] = "62"
			}
			payload := args[1]
			symtab := strings.Split(args[4], "|")
			base, _ := strconv.Atoi(args[2])
			count, _ := strconv.Atoi(args[3])
			return payload, symtab, base, count, nil
		}
	}
	return "", nil, 0, 0, errors.New("unexpected code structure")
}

func replaceStrings(source, prefix, suffix string) string {
	match := regexp.MustCompile(`(?s)var *(_\w+)\=\["(.*?)"\];`).FindStringSubmatch(source)
	if match == nil {
		return prefix + source + suffix
	}
	startpoint := len(match[0])
	lookup := strings.Split(match[2], `","`)
	variable := match[1] + `[%d]`
	for index, value := range lookup {
		source = strings.ReplaceAll(source, fmt.Sprintf(variable, index), `"`+value+`"`)
	}
	return source[startpoint:]
}

func (deu *DEUnpacker) Unpack() (string, error) {
	payload, symtab, base, count, err := filterArgs(deu.Source)
	if err != nil {
		return "", err
	}
	if count != len(symtab) {
		return "", errors.New("malformed p.a.c.k.e.r. symtab")
	}
	ub, err := newUnbaser(base)
	if err != nil {
		return "", err
	}
	payload = strings.ReplaceAll(payload, "\\\\", "\\")
	payload = strings.ReplaceAll(payload, "\\'", "'")
	re := regexp.MustCompile(`\b\w+\b`)
	deu.Source = re.ReplaceAllStringFunc(payload, func(word string) string {
		tab, _ := strconv.Atoi(ub.unbase(word))
		if tab >= len(symtab) || symtab[tab] == "" {
			return word
		}
		return symtab[tab]
	})
	return replaceStrings(deu.Source, deu.Prefix, deu.Suffix), nil
}

type unbaser struct {
	base       int
	alphabet   map[int]string
	dictionary map[rune]int
	unbase     func(string) string
}

func newUnbaser(base int) (*unbaser, error) {
	ub := &unbaser{
		base: base,
		alphabet: map[int]string{
			62: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			95: ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`,
		},
	}
	if 36 < base && base < 62 {
		ub.alphabet[base] = ub.alphabet[62][:base]
	}
	if 2 <= base && base <= 36 {
		ub.unbase = func(str string) string {
			val, _ := strconv.ParseInt(str, base, 64)
			return strconv.FormatInt(val, 10)
		}
	} else {
		if _, present := ub.alphabet[base]; !present {
			return nil, errors.New("unsupported base encoding")
		}
		ub.dictionary = make(map[rune]int)
		for index, char := range ub.alphabet[base] {
			ub.dictionary[char] = index
		}
		ub.unbase = func(str string) string {
			res := 0
			for i := range str {
				c := rune(str[len(str)-i-1])
				res += ub.dictionary[c] * int(math.Pow(float64(ub.base), float64(i)))
			}
			return strconv.Itoa(res)
		}
	}
	return ub, nil
}
