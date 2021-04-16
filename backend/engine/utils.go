package engine

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

func CopyJsonMap(m map[string]map[string]interface{}) map[string]map[string]interface{} {
	cp := make(map[string]map[string]interface{})
	for k, v := range m {
		cp[k] = v
	}

	return cp
}

func RemoveFromSliceMapStringInterface(s []map[string]interface{}, i int) []map[string]interface{} {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxF(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinF(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}


func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return a * -1
	} else {
		return a
	}
}

func Constraint(num float32, lower float32, upper float32) float32 {
	return MinF(MaxF(num, lower), upper)
}
