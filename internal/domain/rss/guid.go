package rss

type Guid struct {
	Value string `json:"value"`
}

func (g Guid) MarshalText() ([]byte, error) {
	return []byte(g.Value), nil
}

func (g *Guid) UnmarshalText(text []byte) error {
	g.Value = string(text)
	return nil
}
