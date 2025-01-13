package grammar

type RuleScanner struct {
}

func (rs *RuleScanner) Scan(ruleString string) (string, *Rule, error) {
	var r *Rule
	return "", r, nil
}
