package conditions

func IfThen(condition bool, optionA interface{}) interface{} {
	if condition {
		return optionA
	}
	return nil
}

func IfThenElse(condition bool, optionA interface{}, optionB interface{}) interface{} {
	if condition {
		return optionA
	}
	return optionB
}
