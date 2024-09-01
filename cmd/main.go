package main

import (
	"fmt"
	"syscall/js"
)

func DecodeSaveFile(this js.Value, args []js.Value) interface{} {
	fileBuffer := args[0]
	fmt.Println(fileBuffer.Get("length").Int())
	data := make([]byte, fileBuffer.Get("length").Int())
	js.CopyBytesToGo(data, fileBuffer)

	save := DecodeSaveData(data)
	return js.ValueOf(map[string]interface{}{
		// "trainer": save.Trainer.toJS(),
		// "pokedex":  save.Pokedex.toJS(),
		"team": save.Team.toJS(),
		// "bag":  save.Bag.toJS(),
		"pc": save.PC.toJS(),
	})
}

func main() {
	js.Global().Set("DecodeSaveData", js.FuncOf(DecodeSaveFile))
	select {}
}
