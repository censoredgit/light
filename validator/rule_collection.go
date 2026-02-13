package validator

type RuleCollection map[string][]Rule

func NewRuleCollection() *RuleCollection {
	v := make(RuleCollection)
	return &v
}

func (v *RuleCollection) AddRule(field string, rules ...Rule) *RuleCollection {
	(*v)[field] = append((*v)[field], rules...)
	return v
}
