package template

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestConcat(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"4", "5"}
	c := []string{"6", "7", "8", "9"}

	r := concat(a, b, c)

	assert.Len(t, r, len(a) + len(b) + len(c))

	for i, v := range a {
		assert.Equal(t, r[i], v)
	}

	offset := len(a)
	for i, v := range b {
		assert.Equal(t, r[i+offset], v)
	}

	offset = len(a) + len(b)
	for i, v := range c {
		assert.Equal(t, r[i+offset], v)
	}

	r = concat(a, nil, c)
	assert.Len(t, r, len(a) + len(c))

	for i, v := range a {
		assert.Equal(t, r[i], v)
	}

	offset = len(a)
	for i, v := range c {
		assert.Equal(t, r[i+offset], v)
	}
}

func TestTokenize(t *testing.T) {
	tokens := tokenize("%var1% This is a %var2% a %var3% and %var1%", "%var1%")

	assert.Len(t, tokens, 3)

	expected := []string{"%var1%", " This is a %var2% a %var3% and ", "%var1%"}
	assert.Equal(t, expected, tokens)

	tokens = tokenize("This is a test", "%var1%")

	assert.Len(t, tokens, 1)

	expected = []string{"This is a test"}
	assert.Equal(t, expected, tokens)

	tokens = tokenize("", "%var1%")
	assert.Len(t, tokens, 1)
}

func TestParse(t *testing.T) {
	text := "this is a test"
	params := map[string]string{}
	p := Parse("this is a test", params)
	assert.Equal(t, text, p)

	text = "This is a test %var1%"
	params = map[string]string{"%var1%": "Var1"}

	p = Parse(text, params)
	assert.Equal(t, "This is a test Var1", p)

	text = "This is a test %var1% and %var2% and this is going %var2% to %var2%"
	params = map[string]string{"%var1%": "Var1", "%var2%": "Var2"}

	p = Parse(text, params)
	assert.Equal(t, "This is a test Var1 and Var2 and this is going Var2 to Var2", p)

	text = "This is a test %var1% and %var2%"
	params = map[string]string{"%var1%": "%var2%", "%var2%": "Var2"}

	p = Parse(text, params)
	assert.Equal(t, "This is a test %var2% and Var2", p)
}

func ExampleParse() {
	text := "This is a test %var1%"
	params := map[string]string{"%var1%": "Var1"}

	newText := Parse(text, params)
	fmt.Println(newText)
	// Output: This is a test Var1
}
