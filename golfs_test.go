package logrus

type Person struct {
	Name    string   `golf:"name"`
	Alias   string   `golf:"alias"`
	Hideout *Hideout `golf:"hideout"`
}

func (p *Person) PlayGolf() bool {
	return true
}

type Hideout struct {
	Name        string `golf:"name"`
	DimensionId int    `golf:"dimensionId"`
}
