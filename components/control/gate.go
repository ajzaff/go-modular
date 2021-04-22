package control

type Gate V

func (g *Gate) On() bool { return *g != 0 }
func (g *Gate) SetOn(on bool) {
	if !on {
		*g = 0
		return
	}
	*g = 1
}
func (g *Gate) Read(vs []V) (n int, err error) {
	if !g.On() {
		for i := range vs {
			vs[i] = V(0)
		}
	}
	return len(vs), nil
}
