package main

import "strings"

func prepareWildcard(s string) func(string) bool {
	if len(s) == 0 {
		return emptyMatcher
	}

	parts := removeEmptyStrings(strings.Split(s, "*"))

	// Only * characters
	if len(parts) == 0 {
		return alwaysMatcher
	}
	if len(parts) == 1 && len(parts[0]) == len(s) {
		return equalsMatcher(parts[0])
	}

	start := 0
	end := len(parts) - 1
	strictPrefix := s[0] != '*'
	strictSuffix := s[len(s)-1] != '*'
	if strictPrefix {
		start++
	}
	if strictSuffix {
		end--
	}

	return func(s string) bool {
		if strictPrefix {
			prefix := parts[0]
			if !strings.HasPrefix(s, prefix) {
				return false
			}
			s = s[len(prefix):]
		}
		if strictSuffix {
			suffix := parts[len(parts)-1]
			if !strings.HasSuffix(s, suffix) {
				return false
			}
			s = s[:len(s)-len(suffix)]
		}

		for i := start; i < end; i++ {
			subStringIndex := strings.Index(s, parts[i])
			if subStringIndex == -1 {
				// Either the substring doesn't exist or we've got overlap, we don't want that.
				return false
			}
			s = s[subStringIndex:]
		}

		return true
	}
}

func removeEmptyStrings(s []string) []string {
	i := 0
	for _, str := range s {
		if str != "" {
			s[i] = str
			i++
		}
	}
	return s[:i]
}

// Some shortcuts which should make things a bit faster in some special-case scenarios

func emptyMatcher(s string) bool {
	return s == ""
}

func alwaysMatcher(_ string) bool {
	return true
}

func equalsMatcher(a string) func(string) bool {
	return func(b string) bool { return a == b }
}
