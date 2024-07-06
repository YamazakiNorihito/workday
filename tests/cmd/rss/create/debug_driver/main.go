package main

type shower interface {
	getWater() []shower
}

type display struct {
	SubDisplay *display
}

func (d display) getWater() []shower {
	return []shower{display{}, d.SubDisplay}
}

func main() {
	// SubDisplayはnullで初期化されます
	s := display{}
	// water := []shower{nil}
	water := s.getWater()
	for _, x := range water {
		if x == nil {
			panic("すべて正常、nilが見つかりました")
		}

		// 最初のイテレーションではdisplay{}はnilではないため
		// 正常に動作しますが、2回目のイテレーションで
		// xはnilになり、getWaterはpanicします。
		x.getWater()
	}
}
