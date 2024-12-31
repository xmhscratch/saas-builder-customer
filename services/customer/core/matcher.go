package core

import (
	"regexp"
)

const (
	LowerFlat = "LowerFlat"
	UpperFlat = "UpperFlat"
	LowerCamel = "LowerCamel"
	UpperCamel = "UpperCamel"
	Snake = "Snake"
	ScreamingSnake = "ScreamingSnake"
	CamelSnake = "CamelSnake"
	Kebab = "Kebab"
	ScreamingKebab = "ScreamingKebab"
	Train = "Train"
)

const (
    LowerFlatRegex = "[a-z0-9]+"
    UpperFlatRegex ="[A-Z0-9]+"
    LowerCamelRegex = "[a-z]+(?:[A-Z0-9]+[a-z0-9]+[A-Za-z0-9]*)+"
    UpperCamelRegex = "(?:[A-Z][a-z0-9]+)(?:[A-Z]+[a-z0-9]*)+"
    SnakeRegex = "[a-z0-9]+(?:_[a-z0-9]+)+"
    ScreamingSnakeRegex = "[A-Z0-9]+(?:_[A-Z0-9]+)+"
    CamelSnakeRegex = "[A-Z][a-z0-9]+(?:_[A-Z]+[a-z0-9]*)+"
    KebabRegex = "[a-z0-9]+(?:-[a-z0-9]+)+"
    ScreamingKebabRegex = "[A-Z0-9]+(?:-[A-Z0-9]+)+"
    TrainRegex = "[A-Z][a-z0-9]+(?:-[A-Z]+[a-z0-9]*)+"
)

func matchCase(casePattern string, requestedCase string) bool {
	match, _ := regexp.MatchString(casePattern, requestedCase)
    return match
}

func contains(slices []string, term string) bool {
	for _, val := range slices {
		if val == term {
			return true
		}
	}
	return false
}

func containsAll(slices []string, terms... string) bool {
	if len(terms) == 0 {
		return false
	}
	for _, term := range terms {
		if !contains(slices, term) {
			return false
		}
	}
	return true
}

func containsExcluding(slices []string, include string, excludes []string) bool {
	for _, exclude := range excludes {
		if contains(slices, exclude) {
			return false
		}
	}
	return contains(slices, include)
}

func appendOnce(slices []string, val string) []string {
	if !contains(slices, val) {
		slices = append(slices, val)
	}
	return slices
}

func getTags(requestedCase string) []string {
	tags := []string{}

	/**
	* Represents matching for flat case, e.g. 'flatcase'
	*
	* Can be also referred as: lower flat case
	*/
	if matchCase(LowerFlatRegex, requestedCase) {
		tags = appendOnce(tags, "lower")
		tags = appendOnce(tags, "flat")
	}

	/**
	* Represents matching for upper flat case, e.g. 'UPPERFLATCASE'
	*/
	if matchCase(UpperFlatRegex, requestedCase) {
		tags = appendOnce(tags, "upper")
		tags = appendOnce(tags, "flat")
	}

	/**
	* Represents matching for camel case, e.g. 'camelCase'
	*
	* Can be also referred as: lower camel case, dromedary case
	*/
	if matchCase(LowerCamelRegex, requestedCase) {
		tags = appendOnce(tags, "lower")
		tags = appendOnce(tags, "camel")
	}

	/**
	* Represents matching for upper camel case, e.g. 'UpperCamelCase'
	*
	* Can be also referred as: pascal case, studly case
	*/
	if matchCase(UpperCamelRegex, requestedCase) {
		tags = appendOnce(tags, "upper")
		tags = appendOnce(tags, "camel")
	}

	/**
	* Represents matching for snake case, e.g. 'snake_case'
	*
	* Can be also referred as: lower snake case, pothole case
	*/
	if matchCase(SnakeRegex, requestedCase) {
		tags = appendOnce(tags, "snake")
	}

	/**
	* Represents matching for screaming snake case, e.g. 'SCREAMING_SNAKE_CASE'
	*
	* Can be also referred as: upper snake case, macro case, constant case
	*/
	if matchCase(ScreamingSnakeRegex, requestedCase) {
		tags = appendOnce(tags, "screaming")
		tags = appendOnce(tags, "snake")
	}

	/**
	* Represents matching for camel snake case, e.g. 'Camel_Snake_Case'
	*/
	if matchCase(CamelSnakeRegex, requestedCase) {
		tags = appendOnce(tags, "camel")
		tags = appendOnce(tags, "snake")
	}

	/**
	* Represents matching for kebab case, e.g. 'kebab-case'
	*
	* Can be also referred as: lower kebab case, dash case, lisp case
	*/
	if matchCase(KebabRegex, requestedCase) {
		tags = appendOnce(tags, "kebab")
	}

	/**
	* Represents matching for screaming kebab case, e.g. 'SCREAMING-KEBAB-CASE'
	*
	* Can be also referred as: upper kebab case, cobol case
	*/
	if matchCase(ScreamingKebabRegex, requestedCase) {
		tags = appendOnce(tags, "screaming")
		tags = appendOnce(tags, "kebab")
	}

	/**
	* Represents matching for train case, e.g. 'Train-Case'
	*/
	if matchCase(TrainRegex, requestedCase) {
		tags = appendOnce(tags, "train")
	}

	return tags
}

func GetCaseName(requestedCase string) string {
	if len(requestedCase) == 0 {
		return ""
	}

	tags := getTags(requestedCase)
	switch (true) {
		case containsAll(tags, "train"):
			return "Train"
		case containsAll(tags, "upper", "kebab") || containsAll(tags, "screaming", "kebab") || containsAll(tags, "cobol"):
			return "ScreamingKebab"
		case containsExcluding(tags, "kebab", []string{"upper", "screaming"}) || containsAll(tags, "lower", "kebab") || containsAll(tags, "dash") || containsAll(tags, "lisp"):
			return "Kebab"
		case containsAll(tags, "camel", "snake"):
			return "CamelSnake"
		case containsAll(tags, "upper", "snake") || containsAll(tags, "screaming", "snake") || containsAll(tags, "macro") || containsAll(tags, "constant"):
			return "ScreamingSnake"
		case containsExcluding(tags, "snake", []string{"upper", "screaming", "camel"}) || containsAll(tags, "lower", "snake") || containsAll(tags, "pothole"):
			return "Snake"
		case containsExcluding(tags, "camel", []string{"upper", "snake"}) || containsAll(tags, "lower", "camel") || containsAll(tags, "dromedary"):
			return "LowerCamel"
		case containsAll(tags, "upper", "camel") || containsAll(tags, "pascal") || containsAll(tags, "studly"):
			return "UpperCamel"
		case containsAll(tags, "upper", "flat"):
			return "UpperFlat"
		case containsExcluding(tags, "flat", []string{"upper"}) || containsAll(tags, "lower", "flat"):
			return "LowerFlat"
		default: return ""
	}
	return ""
}
