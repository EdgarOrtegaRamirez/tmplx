package tmplx

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// FilterFunc is a function that transforms a value.
type FilterFunc func(value interface{}, args []interface{}) (interface{}, error)

// FilterRegistry holds all available filters.
type FilterRegistry struct {
	filters map[string]FilterFunc
}

// NewFilterRegistry creates a registry with built-in filters.
func NewFilterRegistry() *FilterRegistry {
	r := &FilterRegistry{
		filters: make(map[string]FilterFunc),
	}
	r.registerBuiltins()
	return r
}

// Register adds a custom filter.
func (r *FilterRegistry) Register(name string, fn FilterFunc) {
	r.filters[name] = fn
}

// Get returns a filter by name.
func (r *FilterRegistry) Get(name string) (FilterFunc, bool) {
	fn, ok := r.filters[name]
	return fn, ok
}

func (r *FilterRegistry) registerBuiltins() {
	r.Register("upper", filterUpper)
	r.Register("lower", filterLower)
	r.Register("title", filterTitle)
	r.Register("capitalize", filterCapitalize)
	r.Register("trim", filterTrim)
	r.Register("trim_left", filterTrimLeft)
	r.Register("trim_right", filterTrimRight)
	r.Register("ltrim", filterTrimLeft)
	r.Register("rtrim", filterTrimRight)
	r.Register("strip", filterTrim)
	r.Register("lstrip", filterTrimLeft)
	r.Register("rstrip", filterTrimRight)
	r.Register("truncate", filterTruncate)
	r.Register("truncate_words", filterTruncateWords)
	r.Register("replace", filterReplace)
	r.Register("replace_first", filterReplaceFirst)
	r.Register("replace_last", filterReplaceLast)
	r.Register("repeat", filterRepeat)
	r.Register("reverse", filterReverse)
	r.Register("pad_left", filterPadLeft)
	r.Register("pad_right", filterPadRight)
	r.Register("center", filterCenter)
	r.Register("wrap", filterWrap)
	r.Register("snake_case", filterSnakeCase)
	r.Register("camel_case", filterCamelCase)
	r.Register("kebab_case", filterKebabCase)
	r.Register("pascal_case", filterPascalCase)
	r.Register("html_escape", filterHTMLEscape)
	r.Register("html_unescape", filterHTMLUnescape)
	r.Register("url_encode", filterURLEncode)
	r.Register("url_decode", filterURLDecode)
	r.Register("base64_encode", filterBase64Encode)
	r.Register("base64_decode", filterBase64Decode)
	r.Register("abs", filterAbs)
	r.Register("ceil", filterCeil)
	r.Register("floor", filterFloor)
	r.Register("round", filterRound)
	r.Register("clamp", filterClamp)
	r.Register("min", filterMin)
	r.Register("max", filterMax)
	r.Register("sum", filterSum)
	r.Register("average", filterAverage)
	r.Register("mean", filterAverage)
	r.Register("percentage", filterPercentage)
	r.Register("format_number", filterFormatNumber)
	r.Register("int", filterToInt)
	r.Register("float", filterToFloat)
	r.Register("str", filterToString)
	r.Register("bool", filterToBool)
	r.Register("string", filterToString)
	r.Register("number", filterToFloat)
	r.Register("length", filterLength)
	r.Register("len", filterLength)
	r.Register("size", filterLength)
	r.Register("first", filterFirst)
	r.Register("last", filterLast)
	r.Register("slice", filterSlice)
	r.Register("sort", filterSort)
	r.Register("sort_by", filterSortBy)
	r.Register("reverse_collection", filterReverseCollection)
	r.Register("flatten", filterFlatten)
	r.Register("unique", filterUnique)
	r.Register("uniq", filterUnique)
	r.Register("compact", filterCompact)
	r.Register("join", filterJoin)
	r.Register("map", filterMap)
	r.Register("filter", filterFilterFunc)
	r.Register("select", filterFilterFunc)
	r.Register("reject", filterRejectFunc)
	r.Register("group_by", filterGroupBy)
	r.Register("count", filterCount)
	r.Register("keys", filterKeys)
	r.Register("values", filterValues)
	r.Register("to_json", filterToJSON)
	r.Register("to_yaml", filterToYAML)
	r.Register("to_csv", filterToCSV)
	r.Register("to_lines", filterToLines)
	r.Register("default", filterDefault)
	r.Register("d", filterDefault)
	r.Register("ternary", filterTernary)
	r.Register("not", filterNot)
	r.Register("is_empty", filterIsEmpty)
	r.Register("is_blank", filterIsBlank)
	r.Register("is_nil", filterIsNil)
	r.Register("is_none", filterIsNil)
	r.Register("regex_match", filterRegexMatch)
	r.Register("regex_find", filterRegexFind)
	r.Register("regex_find_all", filterRegexFindAll)
	r.Register("regex_replace", filterRegexReplace)
	r.Register("regex_split", filterRegexSplit)
	r.Register("contains", filterContains)
	r.Register("starts_with", filterStartsWith)
	r.Register("ends_with", filterEndsWith)
	r.Register("words", filterWords)
	r.Register("lines", filterLines)
	r.Register("char_count", filterCharCount)
	r.Register("word_count", filterWordCount)
	r.Register("line_count", filterLineCount)
	r.Register("type_of", filterTypeOf)
	r.Register("now", filterNow)
	r.Register("timestamp", filterTimestamp)
}

func toStringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to number", v)
	}
}

func toInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func toSlice(v interface{}) ([]interface{}, bool) {
	switch val := v.(type) {
	case []interface{}:
		return val, true
	case []string:
		result := make([]interface{}, len(val))
		for i, s := range val {
			result[i] = s
		}
		return result, true
	case []int:
		result := make([]interface{}, len(val))
		for i, n := range val {
			result[i] = n
		}
		return result, true
	default:
		return nil, false
	}
}

func toMap(v interface{}) (map[string]interface{}, bool) {
	switch val := v.(type) {
	case map[string]interface{}:
		return val, true
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v2 := range val {
			result[fmt.Sprintf("%v", k)] = v2
		}
		return result, true
	default:
		return nil, false
	}
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	case int:
		return val != 0
	case float64:
		return val != 0
	default:
		return true
	}
}

func filterUpper(value interface{}, args []interface{}) (interface{}, error) {
	return strings.ToUpper(toStringVal(value)), nil
}

func filterLower(value interface{}, args []interface{}) (interface{}, error) {
	return strings.ToLower(toStringVal(value)), nil
}

func filterTitle(value interface{}, args []interface{}) (interface{}, error) {
	return strings.Title(toStringVal(value)), nil
}

func filterCapitalize(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	if s == "" {
		return s, nil
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:], nil
}

func filterTrim(value interface{}, args []interface{}) (interface{}, error) {
	return strings.TrimSpace(toStringVal(value)), nil
}

func filterTrimLeft(value interface{}, args []interface{}) (interface{}, error) {
	return strings.TrimLeftFunc(toStringVal(value), unicode.IsSpace), nil
}

func filterTrimRight(value interface{}, args []interface{}) (interface{}, error) {
	return strings.TrimRightFunc(toStringVal(value), unicode.IsSpace), nil
}

func filterTruncate(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	length := 30
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			length = n
		}
	}
	suffix := "..."
	if len(args) > 1 {
		suffix = toStringVal(args[1])
	}
	if utf8.RuneCountInString(s) <= length {
		return s, nil
	}
	runes := []rune(s)
	return string(runes[:length]) + suffix, nil
}

func filterTruncateWords(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	count := 10
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			count = n
		}
	}
	suffix := "..."
	if len(args) > 1 {
		suffix = toStringVal(args[1])
	}
	words := strings.Fields(s)
	if len(words) <= count {
		return s, nil
	}
	return strings.Join(words[:count], " ") + suffix, nil
}

func filterReplace(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	return strings.ReplaceAll(toStringVal(value), toStringVal(args[0]), toStringVal(args[1])), nil
}

func filterReplaceFirst(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	return strings.Replace(toStringVal(value), toStringVal(args[0]), toStringVal(args[1]), 1), nil
}

func filterReplaceLast(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	s := toStringVal(value)
	old := toStringVal(args[0])
	newStr := toStringVal(args[1])
	idx := strings.LastIndex(s, old)
	if idx == -1 {
		return s, nil
	}
	return s[:idx] + newStr + s[idx+len(old):], nil
}

func filterRepeat(value interface{}, args []interface{}) (interface{}, error) {
	count := 1
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			count = n
		}
	}
	return strings.Repeat(toStringVal(value), count), nil
}

func filterReverse(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func filterPadLeft(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	length := 10
	pad := " "
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			length = n
		}
	}
	if len(args) > 1 {
		pad = toStringVal(args[1])
	}
	runeLen := utf8.RuneCountInString(s)
	if runeLen >= length {
		return s, nil
	}
	padLen := length - runeLen
	return strings.Repeat(pad, padLen) + s, nil
}

func filterPadRight(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	length := 10
	pad := " "
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			length = n
		}
	}
	if len(args) > 1 {
		pad = toStringVal(args[1])
	}
	runeLen := utf8.RuneCountInString(s)
	if runeLen >= length {
		return s, nil
	}
	padLen := length - runeLen
	return s + strings.Repeat(pad, padLen), nil
}

func filterCenter(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	length := 10
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			length = n
		}
	}
	runeLen := utf8.RuneCountInString(s)
	if runeLen >= length {
		return s, nil
	}
	totalPad := length - runeLen
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad), nil
}

func filterWrap(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	width := 80
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			width = n
		}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return s, nil
	}
	var lines []string
	var current strings.Builder
	for _, word := range words {
		if current.Len() > 0 && current.Len()+1+len(word) > width {
			lines = append(lines, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteString(" ")
		}
		current.WriteString(word)
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return strings.Join(lines, "\n"), nil
}

func filterSnakeCase(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func filterCamelCase(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var result strings.Builder
	for i, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(word)
		if i == 0 {
			result.WriteRune(unicode.ToLower(runes[0]))
			if len(runes) > 1 {
				result.WriteString(string(runes[1:]))
			}
		} else {
			result.WriteRune(unicode.ToUpper(runes[0]))
			if len(runes) > 1 {
				result.WriteString(string(runes[1:]))
			}
		}
	}
	return result.String(), nil
}

func filterKebabCase(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func filterPascalCase(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(word)
		result.WriteRune(unicode.ToUpper(runes[0]))
		if len(runes) > 1 {
			result.WriteString(string(runes[1:]))
		}
	}
	return result.String(), nil
}

func filterHTMLEscape(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s, nil
}

func filterHTMLUnescape(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s, nil
}

func filterURLEncode(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
			result.WriteByte(c)
		} else {
			result.WriteString(fmt.Sprintf("%%%02X", c))
		}
	}
	return result.String(), nil
}

func filterURLDecode(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			b, err := strconv.ParseUint(s[i+1:i+3], 16, 8)
			if err == nil {
				result.WriteByte(byte(b))
				i += 2
				continue
			}
		}
		if s[i] == '+' {
			result.WriteByte(' ')
		} else {
			result.WriteByte(s[i])
		}
	}
	return result.String(), nil
}

func filterBase64Encode(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	var result strings.Builder
	b := []byte(s)
	for i := 0; i < len(b); i += 3 {
		var b0, b1, b2 byte
		b0 = b[i]
		if i+1 < len(b) {
			b1 = b[i+1]
		}
		if i+2 < len(b) {
			b2 = b[i+2]
		}
		result.WriteByte(alphabet[(b0>>2)&0x3F])
		result.WriteByte(alphabet[((b0&0x3)<<4)|((b1>>4)&0xF)])
		if i+1 < len(b) {
			result.WriteByte(alphabet[((b1&0xF)<<2)|((b2>>6)&0x3)])
		} else {
			result.WriteByte('=')
		}
		if i+2 < len(b) {
			result.WriteByte(alphabet[b2&0x3F])
		} else {
			result.WriteByte('=')
		}
	}
	return result.String(), nil
}

func filterBase64Decode(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	s = strings.TrimRight(s, "=")
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, 0, len(s)*3/4)
	for i := 0; i < len(s); i += 4 {
		var vals [4]int
		for j := 0; j < 4 && i+j < len(s); j++ {
			for k, c := range alphabet {
				if s[i+j] == byte(c) {
					vals[j] = k
					break
				}
			}
		}
		result = append(result, byte((vals[0]<<2)|(vals[1]>>4)))
		if len(s)-i > 2 && s[i+2] != '=' {
			result = append(result, byte((vals[1]<<4)|(vals[2]>>2)))
		}
		if len(s)-i > 3 && s[i+3] != '=' {
			result = append(result, byte((vals[2]<<6)|vals[3]))
		}
	}
	return string(result), nil
}

func filterAbs(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	return math.Abs(f), nil
}

func filterCeil(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	return math.Ceil(f), nil
}

func filterFloor(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	return math.Floor(f), nil
}

func filterRound(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	precision := 0
	if len(args) > 0 {
		if p, err := toInt(args[0]); err == nil {
			precision = p
		}
	}
	multiplier := math.Pow(10, float64(precision))
	return math.Round(f*multiplier) / multiplier, nil
}

func filterClamp(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	minVal, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	maxVal, err := toFloat(args[1])
	if err != nil {
		return nil, err
	}
	if f < minVal {
		return minVal, nil
	}
	if f > maxVal {
		return maxVal, nil
	}
	return f, nil
}

func filterMin(value interface{}, args []interface{}) (interface{}, error) {
	if s, ok := toSlice(value); ok {
		if len(s) == 0 {
			return nil, nil
		}
		minVal, err := toFloat(s[0])
		if err != nil {
			return s[0], nil
		}
		for _, item := range s[1:] {
			f, err := toFloat(item)
			if err == nil && f < minVal {
				minVal = f
			}
		}
		return minVal, nil
	}
	return value, nil
}

func filterMax(value interface{}, args []interface{}) (interface{}, error) {
	if s, ok := toSlice(value); ok {
		if len(s) == 0 {
			return nil, nil
		}
		maxVal, err := toFloat(s[0])
		if err != nil {
			return s[0], nil
		}
		for _, item := range s[1:] {
			f, err := toFloat(item)
			if err == nil && f > maxVal {
				maxVal = f
			}
		}
		return maxVal, nil
	}
	return value, nil
}

func filterSum(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("sum requires a list")
	}
	total := 0.0
	for _, item := range s {
		f, err := toFloat(item)
		if err == nil {
			total += f
		}
	}
	return total, nil
}

func filterAverage(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("average requires a list")
	}
	if len(s) == 0 {
		return 0.0, nil
	}
	total := 0.0
	for _, item := range s {
		f, err := toFloat(item)
		if err == nil {
			total += f
		}
	}
	return total / float64(len(s)), nil
}

func filterPercentage(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	precision := 1
	if len(args) > 0 {
		if p, err := toInt(args[0]); err == nil {
			precision = p
		}
	}
	return fmt.Sprintf("%.*f%%", precision, f*100), nil
}

func filterFormatNumber(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	precision := 2
	if len(args) > 0 {
		if p, err := toInt(args[0]); err == nil {
			precision = p
		}
	}
	return fmt.Sprintf("%.*f", precision, f), nil
}

func filterToInt(value interface{}, args []interface{}) (interface{}, error) {
	return toInt(value)
}

func filterToFloat(value interface{}, args []interface{}) (interface{}, error) {
	return toFloat(value)
}

func filterToString(value interface{}, args []interface{}) (interface{}, error) {
	return toStringVal(value), nil
}

func filterToBool(value interface{}, args []interface{}) (interface{}, error) {
	return isTruthy(value), nil
}

func filterLength(value interface{}, args []interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return utf8.RuneCountInString(v), nil
	case []interface{}:
		return len(v), nil
	case map[string]interface{}:
		return len(v), nil
	default:
		return 0, nil
	}
}

func filterFirst(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok || len(s) == 0 {
		return nil, nil
	}
	return s[0], nil
}

func filterLast(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok || len(s) == 0 {
		return nil, nil
	}
	return s[len(s)-1], nil
}

func filterSlice(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("slice requires a list")
	}
	start := 0
	end := len(s)
	if len(args) > 0 {
		if n, err := toInt(args[0]); err == nil {
			start = n
		}
	}
	if len(args) > 1 {
		if n, err := toInt(args[1]); err == nil {
			end = n
		}
	}
	if start < 0 {
		start = len(s) + start
	}
	if end < 0 {
		end = len(s) + end
	}
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	return s[start:end], nil
}

func filterSort(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("sort requires a list")
	}
	result := make([]interface{}, len(s))
	copy(result, s)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			a, b := toStringVal(result[i]), toStringVal(result[j])
			if a > b {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result, nil
}

func filterSortBy(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return filterSort(value, args)
	}
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("sort_by requires a list")
	}
	key := toStringVal(args[0])
	result := make([]interface{}, len(s))
	copy(result, s)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			a, _ := toMap(result[i])
			b, _ := toMap(result[j])
			aVal := toStringVal(a[key])
			bVal := toStringVal(b[key])
			if aVal > bVal {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result, nil
}

func filterReverseCollection(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("reverse requires a list")
	}
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result, nil
}

func filterFlatten(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("flatten requires a list")
	}
	var result []interface{}
	for _, item := range s {
		if nested, ok := toSlice(item); ok {
			result = append(result, nested...)
		} else {
			result = append(result, item)
		}
	}
	return result, nil
}

func filterUnique(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("unique requires a list")
	}
	seen := make(map[string]bool)
	var result []interface{}
	for _, item := range s {
		key := toStringVal(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return result, nil
}

func filterCompact(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("compact requires a list")
	}
	var result []interface{}
	for _, item := range s {
		if item != nil && item != "" {
			result = append(result, item)
		}
	}
	return result, nil
}

func filterJoin(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("join requires a list")
	}
	sep := ", "
	if len(args) > 0 {
		sep = toStringVal(args[0])
	}
	strs := make([]string, len(s))
	for i, item := range s {
		strs[i] = toStringVal(item)
	}
	return strings.Join(strs, sep), nil
}

func filterMap(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("map requires a list")
	}
	if len(args) == 0 {
		return value, nil
	}
	key := toStringVal(args[0])
	var result []interface{}
	for _, item := range s {
		if m, ok := toMap(item); ok {
			result = append(result, m[key])
		}
	}
	return result, nil
}

func filterFilterFunc(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("filter requires a list")
	}
	if len(args) == 0 {
		return value, nil
	}
	key := toStringVal(args[0])
	var result []interface{}
	for _, item := range s {
		if m, ok := toMap(item); ok {
			if isTruthy(m[key]) {
				result = append(result, item)
			}
		}
	}
	return result, nil
}

func filterRejectFunc(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("reject requires a list")
	}
	if len(args) == 0 {
		return value, nil
	}
	key := toStringVal(args[0])
	var result []interface{}
	for _, item := range s {
		if m, ok := toMap(item); ok {
			if !isTruthy(m[key]) {
				result = append(result, item)
			}
		}
	}
	return result, nil
}

func filterGroupBy(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("group_by requires a list")
	}
	if len(args) == 0 {
		return value, nil
	}
	key := toStringVal(args[0])
	groups := make(map[string][]interface{})
	for _, item := range s {
		if m, ok := toMap(item); ok {
			groupKey := toStringVal(m[key])
			groups[groupKey] = append(groups[groupKey], item)
		}
	}
	return groups, nil
}

func filterCount(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return nil, fmt.Errorf("count requires a list")
	}
	if len(args) == 0 {
		return len(s), nil
	}
	target := args[0]
	count := 0
	for _, item := range s {
		if toStringVal(item) == toStringVal(target) {
			count++
		}
	}
	return count, nil
}

func filterKeys(value interface{}, args []interface{}) (interface{}, error) {
	m, ok := toMap(value)
	if !ok {
		return nil, fmt.Errorf("keys requires a map")
	}
	keys := make([]interface{}, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys, nil
}

func filterValues(value interface{}, args []interface{}) (interface{}, error) {
	m, ok := toMap(value)
	if !ok {
		return nil, fmt.Errorf("values requires a map")
	}
	vals := make([]interface{}, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals, nil
}

func filterToJSON(value interface{}, args []interface{}) (interface{}, error) {
	indent := false
	if len(args) > 0 {
		indent = isTruthy(args[0])
	}
	return toJSONStr(value, indent), nil
}

func toJSONStr(v interface{}, indent bool) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case []interface{}:
		items := make([]string, len(val))
		for i, item := range val {
			items[i] = toJSONStr(item, indent)
		}
		if indent {
			return "[\n  " + strings.Join(items, ",\n  ") + "\n]"
		}
		return "[" + strings.Join(items, ",") + "]"
	case map[string]interface{}:
		items := make([]string, 0, len(val))
		for k, v2 := range val {
			items = append(items, fmt.Sprintf("%q: %s", k, toJSONStr(v2, indent)))
		}
		if indent {
			return "{\n  " + strings.Join(items, ",\n  ") + "\n}"
		}
		return "{" + strings.Join(items, ",") + "}"
	default:
		return fmt.Sprintf("%q", toStringVal(val))
	}
}

func filterToYAML(value interface{}, args []interface{}) (interface{}, error) {
	return toYAMLStr(value, 0), nil
}

func toYAMLStr(v interface{}, indent int) string {
	prefix := strings.Repeat("  ", indent)
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		if strings.ContainsAny(val, "\n:\"'") {
			return fmt.Sprintf("|-\n%s%s", prefix, strings.ReplaceAll(val, "\n", "\n"+prefix))
		}
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		var sb strings.Builder
		for _, item := range val {
			sb.WriteString(fmt.Sprintf("\n%s- %s", prefix, toYAMLStr(item, indent+1)))
		}
		return sb.String()
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
		var sb strings.Builder
		for k, v2 := range val {
			sb.WriteString(fmt.Sprintf("\n%s%s: %s", prefix, k, toYAMLStr(v2, indent+1)))
		}
		return sb.String()
	default:
		return toStringVal(val)
	}
}

func filterToCSV(value interface{}, args []interface{}) (interface{}, error) {
	s, ok := toSlice(value)
	if !ok {
		return toStringVal(value), nil
	}
	if len(s) == 0 {
		return "", nil
	}
	var sb strings.Builder
	for _, item := range s {
		if m, ok := toMap(item); ok {
			if sb.Len() == 0 {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				sb.WriteString(strings.Join(keys, ","))
				sb.WriteString("\n")
			}
			vals := make([]string, 0, len(m))
			for _, v := range m {
				vals = append(vals, toStringVal(v))
			}
			sb.WriteString(strings.Join(vals, ","))
			sb.WriteString("\n")
		} else {
			sb.WriteString(toStringVal(item))
			sb.WriteString("\n")
		}
	}
	return sb.String(), nil
}

func filterToLines(value interface{}, args []interface{}) (interface{}, error) {
	s := toStringVal(value)
	return strings.Split(s, "\n"), nil
}

func filterDefault(value interface{}, args []interface{}) (interface{}, error) {
	if !isTruthy(value) && len(args) > 0 {
		return args[0], nil
	}
	return value, nil
}

func filterTernary(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	if isTruthy(value) {
		return args[0], nil
	}
	return args[1], nil
}

func filterNot(value interface{}, args []interface{}) (interface{}, error) {
	return !isTruthy(value), nil
}

func filterIsEmpty(value interface{}, args []interface{}) (interface{}, error) {
	if value == nil {
		return true, nil
	}
	switch v := value.(type) {
	case string:
		return v == "", nil
	case []interface{}:
		return len(v) == 0, nil
	case map[string]interface{}:
		return len(v) == 0, nil
	default:
		return false, nil
	}
}

func filterIsBlank(value interface{}, args []interface{}) (interface{}, error) {
	if value == nil {
		return true, nil
	}
	return strings.TrimSpace(toStringVal(value)) == "", nil
}

func filterIsNil(value interface{}, args []interface{}) (interface{}, error) {
	return value == nil, nil
}

func filterRegexMatch(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return false, nil
	}
	s := toStringVal(value)
	pattern := toStringVal(args[0])
	return regexp.MatchString(pattern, s)
}

func filterRegexFind(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return "", nil
	}
	s := toStringVal(value)
	pattern := toStringVal(args[0])
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.FindString(s), nil
}

func filterRegexFindAll(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return []interface{}{}, nil
	}
	s := toStringVal(value)
	pattern := toStringVal(args[0])
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllString(s, -1)
	result := make([]interface{}, len(matches))
	for i, m := range matches {
		result[i] = m
	}
	return result, nil
}

func filterRegexReplace(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return value, nil
	}
	s := toStringVal(value)
	pattern := toStringVal(args[0])
	replacement := toStringVal(args[1])
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.ReplaceAllString(s, replacement), nil
}

func filterRegexSplit(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return value, nil
	}
	s := toStringVal(value)
	pattern := toStringVal(args[0])
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	parts := re.Split(s, -1)
	result := make([]interface{}, len(parts))
	for i, p := range parts {
		result[i] = p
	}
	return result, nil
}

func filterContains(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return false, nil
	}
	return strings.Contains(toStringVal(value), toStringVal(args[0])), nil
}

func filterStartsWith(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return false, nil
	}
	return strings.HasPrefix(toStringVal(value), toStringVal(args[0])), nil
}

func filterEndsWith(value interface{}, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return false, nil
	}
	return strings.HasSuffix(toStringVal(value), toStringVal(args[0])), nil
}

func filterWords(value interface{}, args []interface{}) (interface{}, error) {
	return strings.Fields(toStringVal(value)), nil
}

func filterLines(value interface{}, args []interface{}) (interface{}, error) {
	return strings.Split(toStringVal(value), "\n"), nil
}

func filterCharCount(value interface{}, args []interface{}) (interface{}, error) {
	return utf8.RuneCountInString(toStringVal(value)), nil
}

func filterWordCount(value interface{}, args []interface{}) (interface{}, error) {
	return len(strings.Fields(toStringVal(value))), nil
}

func filterLineCount(value interface{}, args []interface{}) (interface{}, error) {
	return strings.Count(toStringVal(value), "\n") + 1, nil
}

func filterTypeOf(value interface{}, args []interface{}) (interface{}, error) {
	if value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%T", value), nil
}

func filterNow(value interface{}, args []interface{}) (interface{}, error) {
	return strconv.FormatInt(time.Now().Unix(), 10), nil
}

func filterTimestamp(value interface{}, args []interface{}) (interface{}, error) {
	f, err := toFloat(value)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(f), 0)
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		format = toStringVal(args[0])
	}
	return t.Format(format), nil
}
