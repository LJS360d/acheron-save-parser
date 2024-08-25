package main

type PC struct {
	currentBox uint8
	pokemon    [420]Pokemon
}

func (pc *PC) new(section []byte) {
	pc.currentBox = section[0]
	for i := 0; i < len(pc.pokemon); i++ {
		ix := 0x0001 + i*80
		pc.pokemon[i].newBoxed(section[ix : ix+80])
	}
}
