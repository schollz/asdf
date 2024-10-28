package multiply

import (
	"fmt"
	"strconv"
	"strings"
)

// delimiter types
type Delimiter int

const (
	Parentheses Delimiter = iota
	Brackets
)

func Parse(input string, delimiter Delimiter) string {
	placeholders := make(map[string]string)
	placeholderCounter := 0

	// remove any space between a multiplication sign and the thing that preceds it and after it
	for strings.Contains(input, " *") || strings.Contains(input, "* ") {
		input = strings.ReplaceAll(input, " *", "*")
		input = strings.ReplaceAll(input, "* ", "*")
	}

	// Replace bracketed groups with placeholders
	processedInput := replaceBracketedGroups(input, delimiter, &placeholders, &placeholderCounter)

	// Process multiplication on the input with placeholders
	processedInput = processMultiplications(processedInput, delimiter)

	// Replace placeholders with expanded groups
	result := replacePlaceholders(processedInput, placeholders, delimiter)

	return result
}

func replaceBracketedGroups(input string, delimiter Delimiter, placeholders *map[string]string, counter *int) string {
	var builder strings.Builder
	inBracket := 0
	var groupBuilder strings.Builder

	for i := 0; i < len(input); i++ {
		if (input[i] == '(' && delimiter == Parentheses) || (input[i] == '[' && delimiter == Brackets) {
			if inBracket > 0 {
				groupBuilder.WriteByte(input[i])
			}
			inBracket++
		} else if (input[i] == ')' && delimiter == Parentheses) || (input[i] == ']' && delimiter == Brackets) {
			inBracket--
			if inBracket == 0 {
				groupContent := groupBuilder.String()
				groupBuilder.Reset()
				placeholder := fmt.Sprintf("__PLACEHOLDER_%d__", *counter)
				*counter++
				(*placeholders)[placeholder] = groupContent
				builder.WriteString(placeholder)
			} else {
				groupBuilder.WriteByte(input[i])
			}
		} else {
			if inBracket > 0 {
				groupBuilder.WriteByte(input[i])
			} else {
				builder.WriteByte(input[i])
			}
		}
	}

	return builder.String()
}

func processMultiplications(input string, delimiter Delimiter) string {
	tokens := tokenize(input)
	var result []string
	for _, token := range tokens {
		if strings.Contains(token, "*") {
			parts := strings.Split(token, "*")
			entity := parts[0]
			count, _ := strconv.Atoi(parts[1])
			expandedEntity := replicate(entity, count, delimiter)
			result = append(result, expandedEntity)
		} else {
			result = append(result, token)
		}
	}
	return strings.Join(result, " ")
}

func replicate(entity string, count int, delimiter Delimiter) string {
	var expanded []string
	for i := 0; i < count; i++ {
		expanded = append(expanded, entity)
	}
	if delimiter == Brackets {
		return "[" + strings.Join(expanded, " ") + "]"
	} else {
		return "(" + strings.Join(expanded, " ") + ")"
	}
}

func replacePlaceholders(input string, placeholders map[string]string, delimiter Delimiter) string {
	for placeholder, content := range placeholders {
		expansion := Parse(content, delimiter)
		if delimiter == Parentheses {
			input = strings.ReplaceAll(input, placeholder, "("+expansion+")")
		} else {
			input = strings.ReplaceAll(input, placeholder, "["+expansion+"]")
		}
	}
	return input
}

func tokenize(input string) []string {
	var tokens []string
	start := 0
	inBracket := 0
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '[':
			if inBracket == 0 && i > start {
				tokens = append(tokens, strings.TrimSpace(input[start:i]))
			}
			inBracket++
		case ']':
			inBracket--
			if inBracket == 0 {
				tokens = append(tokens, strings.TrimSpace(input[start:i+1]))
				start = i + 1
			}
		case ' ':
			if inBracket == 0 {
				if i > start {
					tokens = append(tokens, strings.TrimSpace(input[start:i]))
				}
				start = i + 1
			}
		}
	}
	if start < len(input) {
		tokens = append(tokens, strings.TrimSpace(input[start:]))
	}
	return tokens
}
