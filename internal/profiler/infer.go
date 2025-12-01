package profiler

import "strconv"

func InferType(v string) string {

	if isInt(v) {
		return "int"
	}
	if isBool(v) {
		return "bool"
	}
	if isFloat(v) {
		return "float"
	}
	return "string"
}

func isInt(v string) bool {
	_, err := strconv.Atoi(v)

	return err == nil
}

func isBool(v string) bool {
	_, err := strconv.ParseBool(v)
	return err == nil
}

func isFloat(v string) bool {
	_, err := strconv.ParseFloat(v, 64)
	return err == nil
}
