package forms

type errors map[string][]string

// Add func增加error message
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get func return 第一個error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}