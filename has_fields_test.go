package logrus

type Person struct {
	Name          string
	Alias         string
	Hideout       *Hideout
	useTypeFields bool
	except        []string
}

func (p *Person) Fields() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name, "alias": p.Alias, "hideout": p.Hideout}
}

func (p *Person) Flatten() bool {
	return true
}

func (p *Person) UseTypeFields() bool {
	return p.useTypeFields
}

func (p *Person) ExceptFields() []string {
	return p.except
}

type Hideout struct {
	Name          string
	DimensionId   int
	useTypeFields bool
	except        []string
}

func (h *Hideout) Fields() map[string]interface{} {
	return map[string]interface{}{
		"name": h.Name, "dimensionId": h.DimensionId}
}

func (h *Hideout) Flatten() bool {
	return false
}

func (h *Hideout) UseTypeFields() bool {
	return h.useTypeFields
}

func (h *Hideout) ExceptFields() []string {
	return h.except
}
