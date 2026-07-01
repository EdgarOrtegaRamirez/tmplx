package tmplx

import (
	"strings"
	"testing"
)

func renderTemplate(t *testing.T, tmpl string, vars map[string]interface{}) string {
	t.Helper()
	e := New()
	e.SetVariables(vars)
	result, err := e.Render(tmpl)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	return result
}

func TestEvaluatorVariable(t *testing.T) {
	result := renderTemplate(t, "Hello {{ name }}!", map[string]interface{}{"name": "World"})
	if result != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got %q", result)
	}
}

func TestEvaluatorIntVariable(t *testing.T) {
	result := renderTemplate(t, "Count: {{ count }}", map[string]interface{}{"count": 42})
	if result != "Count: 42" {
		t.Errorf("Expected 'Count: 42', got %q", result)
	}
}

func TestEvaluatorBoolVariable(t *testing.T) {
	result := renderTemplate(t, "Active: {{ active }}", map[string]interface{}{"active": true})
	if result != "Active: true" {
		t.Errorf("Expected 'Active: true', got %q", result)
	}
}

func TestEvaluatorDottedVariable(t *testing.T) {
	result := renderTemplate(t, "Hello {{ user.name }}!", map[string]interface{}{
		"user": map[string]interface{}{"name": "Alice"},
	})
	if result != "Hello Alice!" {
		t.Errorf("Expected 'Hello Alice!', got %q", result)
	}
}

func TestEvaluatorFilterUpper(t *testing.T) {
	result := renderTemplate(t, "{{ name | upper }}", map[string]interface{}{"name": "hello"})
	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got %q", result)
	}
}

func TestEvaluatorFilterLower(t *testing.T) {
	result := renderTemplate(t, "{{ name | lower }}", map[string]interface{}{"name": "HELLO"})
	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}
}

func TestEvaluatorFilterCapitalize(t *testing.T) {
	result := renderTemplate(t, "{{ text | capitalize }}", map[string]interface{}{"text": "hello"})
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got %q", result)
	}
}

func TestEvaluatorFilterTitle(t *testing.T) {
	result := renderTemplate(t, "{{ text | title }}", map[string]interface{}{"text": "hello world"})
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", result)
	}
}

func TestEvaluatorFilterTruncate(t *testing.T) {
	result := renderTemplate(t, "{{ text | truncate(5) }}", map[string]interface{}{"text": "Hello World"})
	if result != "Hello..." {
		t.Errorf("Expected 'Hello...', got %q", result)
	}
}

func TestEvaluatorFilterTruncateWords(t *testing.T) {
	result := renderTemplate(t, "{{ text | truncate_words(2) }}", map[string]interface{}{"text": "hello world foo bar"})
	if result != "hello world..." {
		t.Errorf("Expected 'hello world...', got %q", result)
	}
}

func TestEvaluatorFilterJoin(t *testing.T) {
	result := renderTemplate(t, "{{ items | join(', ') }}", map[string]interface{}{"items": []interface{}{"a", "b", "c"}})
	if result != "a, b, c" {
		t.Errorf("Expected 'a, b, c', got %q", result)
	}
}

func TestEvaluatorFilterDefault(t *testing.T) {
	result := renderTemplate(t, "{{ missing | default('N/A') }}", map[string]interface{}{})
	if result != "N/A" {
		t.Errorf("Expected 'N/A', got %q", result)
	}
}

func TestEvaluatorFilterLength(t *testing.T) {
	result := renderTemplate(t, "{{ name | length }}", map[string]interface{}{"name": "hello"})
	if result != "5" {
		t.Errorf("Expected '5', got %q", result)
	}
}

func TestEvaluatorFilterSnakeCase(t *testing.T) {
	result := renderTemplate(t, "{{ name | snake_case }}", map[string]interface{}{"name": "helloWorld"})
	if result != "hello_world" {
		t.Errorf("Expected 'hello_world', got %q", result)
	}
}

func TestEvaluatorFilterCamelCase(t *testing.T) {
	result := renderTemplate(t, "{{ name | camel_case }}", map[string]interface{}{"name": "hello_world"})
	if result != "helloWorld" {
		t.Errorf("Expected 'helloWorld', got %q", result)
	}
}

func TestEvaluatorFilterKebabCase(t *testing.T) {
	result := renderTemplate(t, "{{ name | kebab_case }}", map[string]interface{}{"name": "helloWorld"})
	if result != "hello-world" {
		t.Errorf("Expected 'hello-world', got %q", result)
	}
}

func TestEvaluatorFilterPascalCase(t *testing.T) {
	result := renderTemplate(t, "{{ text | pascal_case }}", map[string]interface{}{"text": "hello_world"})
	if result != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', got %q", result)
	}
}

func TestEvaluatorFilterReverse(t *testing.T) {
	result := renderTemplate(t, "{{ name | reverse }}", map[string]interface{}{"name": "hello"})
	if result != "olleh" {
		t.Errorf("Expected 'olleh', got %q", result)
	}
}

func TestEvaluatorFilterTrim(t *testing.T) {
	result := renderTemplate(t, "{{ name | trim }}", map[string]interface{}{"name": "  hello  "})
	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}
}

func TestEvaluatorFilterReplace(t *testing.T) {
	result := renderTemplate(t, "{{ text | replace('world', 'Go') }}", map[string]interface{}{"text": "hello world"})
	if result != "hello Go" {
		t.Errorf("Expected 'hello Go', got %q", result)
	}
}

func TestEvaluatorFilterRepeat(t *testing.T) {
	result := renderTemplate(t, "{{ text | repeat(3) }}", map[string]interface{}{"text": "ab"})
	if result != "ababab" {
		t.Errorf("Expected 'ababab', got %q", result)
	}
}

func TestEvaluatorFilterAbs(t *testing.T) {
	result := renderTemplate(t, "{{ value | abs }}", map[string]interface{}{"value": -5})
	if result != "5" {
		t.Errorf("Expected '5', got %q", result)
	}
}

func TestEvaluatorFilterRound(t *testing.T) {
	result := renderTemplate(t, "{{ value | round(2) }}", map[string]interface{}{"value": 3.14159})
	if result != "3.14" {
		t.Errorf("Expected '3.14', got %q", result)
	}
}

func TestEvaluatorFilterCeil(t *testing.T) {
	result := renderTemplate(t, "{{ value | ceil }}", map[string]interface{}{"value": 3.1})
	if result != "4" {
		t.Errorf("Expected '4', got %q", result)
	}
}

func TestEvaluatorFilterFloor(t *testing.T) {
	result := renderTemplate(t, "{{ value | floor }}", map[string]interface{}{"value": 3.9})
	if result != "3" {
		t.Errorf("Expected '3', got %q", result)
	}
}

func TestEvaluatorFilterSum(t *testing.T) {
	result := renderTemplate(t, "{{ nums | sum }}", map[string]interface{}{"nums": []interface{}{1, 2, 3, 4}})
	if result != "10" {
		t.Errorf("Expected '10', got %q", result)
	}
}

func TestEvaluatorFilterMin(t *testing.T) {
	result := renderTemplate(t, "{{ nums | min }}", map[string]interface{}{"nums": []interface{}{5, 3, 8, 1, 9}})
	if result != "1" {
		t.Errorf("Expected '1', got %q", result)
	}
}

func TestEvaluatorFilterMax(t *testing.T) {
	result := renderTemplate(t, "{{ nums | max }}", map[string]interface{}{"nums": []interface{}{5, 3, 8, 1, 9}})
	if result != "9" {
		t.Errorf("Expected '9', got %q", result)
	}
}

func TestEvaluatorFilterContains(t *testing.T) {
	result := renderTemplate(t, "{{ text | contains('llo') }}", map[string]interface{}{"text": "hello"})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterStartsWith(t *testing.T) {
	result := renderTemplate(t, "{{ text | starts_with('hel') }}", map[string]interface{}{"text": "hello"})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterEndsWith(t *testing.T) {
	result := renderTemplate(t, "{{ text | ends_with('llo') }}", map[string]interface{}{"text": "hello"})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterUnique(t *testing.T) {
	result := renderTemplate(t, "{{ nums | unique | join(',') }}", map[string]interface{}{"nums": []interface{}{1, 2, 2, 3, 3, 3}})
	if result != "1,2,3" {
		t.Errorf("Expected '1,2,3', got %q", result)
	}
}

func TestEvaluatorFilterSort(t *testing.T) {
	result := renderTemplate(t, "{{ nums | sort | join(',') }}", map[string]interface{}{"nums": []interface{}{3, 1, 2}})
	if result != "1,2,3" {
		t.Errorf("Expected '1,2,3', got %q", result)
	}
}

func TestEvaluatorFilterFlatten(t *testing.T) {
	result := renderTemplate(t, "{{ nested | flatten | join(',') }}", map[string]interface{}{
		"nested": []interface{}{
			[]interface{}{1, 2},
			[]interface{}{3, 4},
		},
	})
	if result != "1,2,3,4" {
		t.Errorf("Expected '1,2,3,4', got %q", result)
	}
}

func TestEvaluatorFilterBase64Encode(t *testing.T) {
	result := renderTemplate(t, "{{ text | base64_encode }}", map[string]interface{}{"text": "hello"})
	if result != "aGVsbG8=" {
		t.Errorf("Expected 'aGVsbG8=', got %q", result)
	}
}

func TestEvaluatorFilterBase64Decode(t *testing.T) {
	result := renderTemplate(t, "{{ text | base64_decode }}", map[string]interface{}{"text": "aGVsbG8="})
	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}
}

func TestEvaluatorFilterUrlEncode(t *testing.T) {
	result := renderTemplate(t, "{{ text | url_encode }}", map[string]interface{}{"text": "hello world"})
	if result != "hello%20world" {
		t.Errorf("Expected 'hello%%20world', got %q", result)
	}
}

func TestEvaluatorFilterUrlDecode(t *testing.T) {
	result := renderTemplate(t, "{{ text | url_decode }}", map[string]interface{}{"text": "hello%20world"})
	if result != "hello world" {
		t.Errorf("Expected 'hello world', got %q", result)
	}
}

func TestEvaluatorFilterPadLeft(t *testing.T) {
	result := renderTemplate(t, "{{ text | pad_left(5, '0') }}", map[string]interface{}{"text": "42"})
	if result != "00042" {
		t.Errorf("Expected '00042', got %q", result)
	}
}

func TestEvaluatorFilterPadRight(t *testing.T) {
	result := renderTemplate(t, "{{ text | pad_right(5, '.') }}", map[string]interface{}{"text": "hi"})
	if result != "hi..." {
		t.Errorf("Expected 'hi...', got %q", result)
	}
}

func TestEvaluatorFilterHtmlEscape(t *testing.T) {
	result := renderTemplate(t, "{{ text | html_escape }}", map[string]interface{}{"text": "<b>bold</b>"})
	if !strings.Contains(result, "&lt;b&gt;") {
		t.Errorf("Expected HTML escaped text, got %q", result)
	}
}

func TestEvaluatorFilterTypeOf(t *testing.T) {
	result := renderTemplate(t, "{{ value | type_of }}", map[string]interface{}{"value": "hello"})
	if result != "string" {
		t.Errorf("Expected 'string', got %q", result)
	}
}

func TestEvaluatorFilterClamp(t *testing.T) {
	result := renderTemplate(t, "{{ value | clamp(0, 100) }}", map[string]interface{}{"value": 150})
	if result != "100" {
		t.Errorf("Expected '100', got %q", result)
	}
}

func TestEvaluatorFilterCount(t *testing.T) {
	result := renderTemplate(t, "{{ items | count }}", map[string]interface{}{"items": []interface{}{"a", "b", "c"}})
	if result != "3" {
		t.Errorf("Expected '3', got %q", result)
	}
}

func TestEvaluatorFilterKeys(t *testing.T) {
	result := renderTemplate(t, "{{ obj | keys | sort | join(',') }}", map[string]interface{}{
		"obj": map[string]interface{}{"c": 3, "a": 1, "b": 2},
	})
	if result != "a,b,c" {
		t.Errorf("Expected 'a,b,c', got %q", result)
	}
}

func TestEvaluatorFilterValues(t *testing.T) {
	result := renderTemplate(t, "{{ obj | values | sort | join(',') }}", map[string]interface{}{
		"obj": map[string]interface{}{"a": 3, "b": 1, "c": 2},
	})
	if result != "1,2,3" {
		t.Errorf("Expected '1,2,3', got %q", result)
	}
}

func TestEvaluatorFilterToJSON(t *testing.T) {
	result := renderTemplate(t, "{{ obj | to_json }}", map[string]interface{}{
		"obj": map[string]interface{}{"name": "test"},
	})
	if !strings.Contains(result, `"name"`) {
		t.Errorf("Expected JSON with name, got %q", result)
	}
}

func TestEvaluatorFilterWordCount(t *testing.T) {
	result := renderTemplate(t, "{{ text | word_count }}", map[string]interface{}{"text": "hello world foo"})
	if result != "3" {
		t.Errorf("Expected '3', got %q", result)
	}
}

func TestEvaluatorFilterCharCount(t *testing.T) {
	result := renderTemplate(t, "{{ text | char_count }}", map[string]interface{}{"text": "hello"})
	if result != "5" {
		t.Errorf("Expected '5', got %q", result)
	}
}

func TestEvaluatorFilterFirst(t *testing.T) {
	result := renderTemplate(t, "{{ items | first }}", map[string]interface{}{"items": []interface{}{"a", "b", "c"}})
	if result != "a" {
		t.Errorf("Expected 'a', got %q", result)
	}
}

func TestEvaluatorFilterLast(t *testing.T) {
	result := renderTemplate(t, "{{ items | last }}", map[string]interface{}{"items": []interface{}{"a", "b", "c"}})
	if result != "c" {
		t.Errorf("Expected 'c', got %q", result)
	}
}

func TestEvaluatorFilterSlice(t *testing.T) {
	result := renderTemplate(t, "{{ items | slice(1, 3) | join(',') }}", map[string]interface{}{
		"items": []interface{}{"a", "b", "c", "d"},
	})
	if result != "b,c" {
		t.Errorf("Expected 'b,c', got %q", result)
	}
}

func TestEvaluatorFilterMap(t *testing.T) {
	result := renderTemplate(t, "{{ users | map('name') | join(',') }}", map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice"},
			map[string]interface{}{"name": "Bob"},
		},
	})
	if result != "Alice,Bob" {
		t.Errorf("Expected 'Alice,Bob', got %q", result)
	}
}

func TestEvaluatorFilterSelect(t *testing.T) {
	result := renderTemplate(t, "{{ users | select('active') | length }}", map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "active": true},
			map[string]interface{}{"name": "Bob", "active": false},
			map[string]interface{}{"name": "Charlie", "active": true},
		},
	})
	if result != "2" {
		t.Errorf("Expected '2', got %q", result)
	}
}

func TestEvaluatorFilterReject(t *testing.T) {
	result := renderTemplate(t, "{{ users | reject('active') | length }}", map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "active": true},
			map[string]interface{}{"name": "Bob", "active": false},
			map[string]interface{}{"name": "Charlie", "active": true},
		},
	})
	if result != "1" {
		t.Errorf("Expected '1', got %q", result)
	}
}

func TestEvaluatorFilterGroupBy(t *testing.T) {
	// group_by groups items; verify it doesn't error and produces output
	e := New()
	e.SetVariable("items", []interface{}{
		map[string]interface{}{"name": "a", "type": "x"},
		map[string]interface{}{"name": "b", "type": "y"},
		map[string]interface{}{"name": "c", "type": "x"},
	})
	result, err := e.Render("{{ items | group_by('type') }}")
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got %q", result)
	}
}

func TestEvaluatorFilterSortBy(t *testing.T) {
	result := renderTemplate(t, "{{ items | sort_by('name') | map('name') | join(',') }}", map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "c"},
			map[string]interface{}{"name": "a"},
			map[string]interface{}{"name": "b"},
		},
	})
	if result != "a,b,c" {
		t.Errorf("Expected 'a,b,c', got %q", result)
	}
}

func TestEvaluatorFilterReverseCollection(t *testing.T) {
	result := renderTemplate(t, "{{ items | reverse_collection | join(',') }}", map[string]interface{}{
		"items": []interface{}{"a", "b", "c"},
	})
	if result != "c,b,a" {
		t.Errorf("Expected 'c,b,a', got %q", result)
	}
}

func TestEvaluatorFilterRegexFind(t *testing.T) {
	result := renderTemplate(t, "{{ text | regex_find('\\d+') }}", map[string]interface{}{"text": "abc123def"})
	if result != "123" {
		t.Errorf("Expected '123', got %q", result)
	}
}

func TestEvaluatorFilterRegexMatch(t *testing.T) {
	result := renderTemplate(t, "{{ text | regex_match('^[a-z]+$') }}", map[string]interface{}{"text": "hello"})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterNot(t *testing.T) {
	result := renderTemplate(t, "{{ flag | not }}", map[string]interface{}{"flag": false})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterTernary(t *testing.T) {
	result := renderTemplate(t, "{{ flag | ternary('yes', 'no') }}", map[string]interface{}{"flag": true})
	if result != "yes" {
		t.Errorf("Expected 'yes', got %q", result)
	}
}

func TestEvaluatorFilterIsNil(t *testing.T) {
	result := renderTemplate(t, "{{ text | is_nil }}", map[string]interface{}{})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterIsEmpty(t *testing.T) {
	result := renderTemplate(t, "{{ text | is_empty }}", map[string]interface{}{"text": ""})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterIsBlank(t *testing.T) {
	result := renderTemplate(t, "{{ text | is_blank }}", map[string]interface{}{"text": "   "})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorIfBlock(t *testing.T) {
	result := renderTemplate(t, "{{# if active }}Yes{{/ if }}", map[string]interface{}{"active": true})
	if result != "Yes" {
		t.Errorf("Expected 'Yes', got %q", result)
	}
}

func TestEvaluatorIfBlockFalse(t *testing.T) {
	result := renderTemplate(t, "{{# if active }}Yes{{/ if }}", map[string]interface{}{"active": false})
	if result != "" {
		t.Errorf("Expected '', got %q", result)
	}
}

func TestEvaluatorIfElseBlock(t *testing.T) {
	result := renderTemplate(t, "{{# if active }}Yes{{ else }}No{{/ if }}", map[string]interface{}{"active": false})
	if result != "No" {
		t.Errorf("Expected 'No', got %q", result)
	}
}

func TestEvaluatorUnlessBlock(t *testing.T) {
	result := renderTemplate(t, "{{# unless hidden }}Visible{{/ unless }}", map[string]interface{}{"hidden": false})
	if result != "Visible" {
		t.Errorf("Expected 'Visible', got %q", result)
	}
}

func TestEvaluatorUnlessBlockTrue(t *testing.T) {
	result := renderTemplate(t, "{{# unless hidden }}Visible{{/ unless }}", map[string]interface{}{"hidden": true})
	if result != "" {
		t.Errorf("Expected '', got %q", result)
	}
}

func TestEvaluatorEachBlock(t *testing.T) {
	result := renderTemplate(t, "{{# each items }}{{ item }} {{/ each }}", map[string]interface{}{
		"items": []interface{}{"a", "b", "c"},
	})
	if result != "a b c " {
		t.Errorf("Expected 'a b c ', got %q", result)
	}
}

func TestEvaluatorEachBlockWithVariable(t *testing.T) {
	result := renderTemplate(t, "{{# each users }}{{ item.name }} {{/ each }}", map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice"},
			map[string]interface{}{"name": "Bob"},
		},
	})
	if result != "Alice Bob " {
		t.Errorf("Expected 'Alice Bob ', got %q", result)
	}
}

func TestEvaluatorWithBlock(t *testing.T) {
	result := renderTemplate(t, "{{# with user }}{{ name }} ({{ email }}){{/ with }}", map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
		},
	})
	if result != "Alice (alice@example.com)" {
		t.Errorf("Expected 'Alice (alice@example.com)', got %q", result)
	}
}

func TestEvaluatorNestedIfEach(t *testing.T) {
	result := renderTemplate(t,
		"{{# if show }}{{# each items }}{{ item }} {{/ each }}{{/ if }}",
		map[string]interface{}{
			"show":  true,
			"items": []interface{}{"x", "y"},
		})
	if result != "x y " {
		t.Errorf("Expected 'x y ', got %q", result)
	}
}

func TestEvaluatorPartial(t *testing.T) {
	e := New()
	e.SetPartial("header", "Hello {{ name }}!")
	e.SetVariable("name", "World")
	result, err := e.Render("{{> header }}")
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if result != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got %q", result)
	}
}

func TestEvaluatorComment(t *testing.T) {
	result := renderTemplate(t, "Hello {{! this is a comment }} World", map[string]interface{}{})
	if result != "Hello  World" {
		t.Errorf("Expected 'Hello  World', got %q", result)
	}
}

func TestEvaluatorMultipleFilters(t *testing.T) {
	result := renderTemplate(t, "{{ name | upper | truncate(3) }}", map[string]interface{}{"name": "hello"})
	// upper("hello") = "HELLO", truncate(3, "HELLO") = "HEL..."
	if result != "HEL..." {
		t.Errorf("Expected 'HEL...', got %q", result)
	}
}

func TestEvaluatorFilterRegexFindAll(t *testing.T) {
	result := renderTemplate(t, "{{ text | regex_find_all('\\d+') | join(',') }}", map[string]interface{}{"text": "a1b2c3"})
	if result != "1,2,3" {
		t.Errorf("Expected '1,2,3', got %q", result)
	}
}

func TestEvaluatorFilterRegexReplace(t *testing.T) {
	result := renderTemplate(t, "{{ text | regex_replace('\\d+', 'NUM') }}", map[string]interface{}{"text": "abc123def456"})
	if result != "abcNUMdefNUM" {
		t.Errorf("Expected 'abcNUMdefNUM', got %q", result)
	}
}

func TestEvaluatorFilterCompact(t *testing.T) {
	result := renderTemplate(t, "{{ items | compact | join(',') }}", map[string]interface{}{
		"items": []interface{}{"a", "", "b", nil, "c"},
	})
	if result != "a,b,c" {
		t.Errorf("Expected 'a,b,c', got %q", result)
	}
}

func TestEvaluatorFilterToString(t *testing.T) {
	result := renderTemplate(t, "{{ value | string }}", map[string]interface{}{"value": 42})
	if result != "42" {
		t.Errorf("Expected '42', got %q", result)
	}
}

func TestEvaluatorFilterToInt(t *testing.T) {
	result := renderTemplate(t, "{{ value | int }}", map[string]interface{}{"value": "42"})
	if result != "42" {
		t.Errorf("Expected '42', got %q", result)
	}
}

func TestEvaluatorFilterToFloat(t *testing.T) {
	result := renderTemplate(t, "{{ value | float }}", map[string]interface{}{"value": "3.14"})
	if result != "3.14" {
		t.Errorf("Expected '3.14', got %q", result)
	}
}

func TestEvaluatorFilterToBool(t *testing.T) {
	result := renderTemplate(t, "{{ value | bool }}", map[string]interface{}{"value": "hello"})
	if result != "true" {
		t.Errorf("Expected 'true', got %q", result)
	}
}

func TestEvaluatorFilterPercentage(t *testing.T) {
	result := renderTemplate(t, "{{ value | percentage }}", map[string]interface{}{"value": 0.75})
	if result != "75.0%" {
		t.Errorf("Expected '75.0%%', got %q", result)
	}
}

func TestEvaluatorFilterFormatNumber(t *testing.T) {
	result := renderTemplate(t, "{{ value | format_number(3) }}", map[string]interface{}{"value": 3.14159})
	if result != "3.142" {
		t.Errorf("Expected '3.142', got %q", result)
	}
}

func TestEvaluatorFilterRegexSplit(t *testing.T) {
	result := renderTemplate(t, "{{ text | regex_split(',') | length }}", map[string]interface{}{"text": "a,b,c"})
	if result != "3" {
		t.Errorf("Expected '3', got %q", result)
	}
}

func TestEvaluatorFilterSize(t *testing.T) {
	result := renderTemplate(t, "{{ items | size }}", map[string]interface{}{"items": []interface{}{"a", "b", "c"}})
	if result != "3" {
		t.Errorf("Expected '3', got %q", result)
	}
}
