package template

import (
	"strings"
)

// Parse text and replace all the placeholders that match with params keys
func Parse(text string, params map[string]string) string {
	if len(params) == 0 {
		return text
	}

	var tokens []string

	for key := range params {
		if len(tokens) == 0 {
			tokens = tokenize(text, key)
			continue
		}

		i := 0
		for i < len(tokens) {
			token := tokens[i]
			if _, ok := params[token]; ok {
				i++
				continue
			}

			stokens := tokenize(token, key)
			if len(stokens) == 1 {
				i++
				continue
			}

			tokens = concat(tokens[0:i], stokens, tokens[i+1:])
			i += len(stokens)
		}
	}

	for i, token := range tokens {
		if v, ok := params[token]; ok {
			tokens[i] = v
		}
	}


	return strings.Join(tokens, "")
}

func concat(arrs ...[]string) []string {
	var r []string

	for _, arr := range arrs {
		r = append(r, arr...)
	}

	return r
}

func tokenize(text string, needle string) []string {
	if text == "" {
		return []string{""}
	}

	parts := strings.Split(text, needle)
	var tokens []string

	for i, token := range parts {
		if token != "" {
			tokens = append(tokens, token)
		}

		if i < len(parts) - 1 {
			tokens = append(tokens, needle)
		}
	}

	return tokens
}
