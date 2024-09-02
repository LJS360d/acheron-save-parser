package main

import (
	"acheron-save-parser/export/jsutils"
	"acheron-save-parser/sav"
	"fmt"
	"syscall/js"
)

func DecodeSaveFile(this js.Value, args []js.Value) interface{} {
	fileBuffer := args[0]
	fmt.Println(fileBuffer.Get("length").Int())
	data := make([]byte, fileBuffer.Get("length").Int())
	js.CopyBytesToGo(data, fileBuffer)

	save := sav.DecodeSaveData(data)
	return js.ValueOf(map[string]interface{}{
		// "trainer": save.Trainer.ToJs(),
		// "pokedex":  save.Pokedex.ToJs(),
		"team": jsutils.TeamToJs(&save.Team),
		// "bag":  save.Bag.ToJs(),
		"pc": jsutils.PCToJs(&save.PC),
	})
}

func main() {
	js.Global().Set("DecodeSaveData", js.FuncOf(DecodeSaveFile))
	select {}
}
