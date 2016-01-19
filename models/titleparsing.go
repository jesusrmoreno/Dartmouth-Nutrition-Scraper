package models

import "strings"

// TitleToProps returns stuff ...
func TitleToProps(title string) []string {
	props := []string{}
	if strings.Contains(title, "[") && strings.Contains(title, "]") {
		openDietary := strings.Index(title, "[")
		closeDietary := strings.Index(title, "]")
		dietary := strings.Replace(title[openDietary+1:closeDietary], ".", ",", -1)
		dietary = strings.Replace(dietary, " ", "", -1)
		for _, item := range strings.Split(dietary, ",") {
			props = append(props, strings.ToLower(string(item)))
		}
	}

	if strings.Contains(title, "(") && strings.Contains(title, ")") {
		openAllergens := strings.Index(title, "(")
		closeAllergens := strings.Index(title, ")")
		allergens := strings.Replace(title[openAllergens+1:closeAllergens], ".", ",", -1)
		allergens = strings.Replace(allergens, " ", "", -1)
		for _, item := range strings.Split(allergens, ",") {
			props = append(props, strings.ToLower(string(item)))
		}
	}

	return props
}

// minInteger
func minInteger(a, b int) int {
	if a == -1 && a < b {
		return b
	}
	if b == -1 && b < a {
		return a
	}
	if a < b {
		return a
	}
	if b < a {
		return b
	}
	return 0
}

// RemoveMetaData ...
func RemoveMetaData(title string) string {
	openDietary := strings.Index(title, "[")
	openAllergens := strings.Index(title, "(")
	if openDietary != -1 || openAllergens != -1 {
		return title[:minInteger(openAllergens, openDietary)]
	}
	return title
}
