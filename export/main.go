package main

import (
	jsconvert "acheron-save-parser/export/js"
	"acheron-save-parser/gba"
	"acheron-save-parser/sav"
	"fmt"
	"syscall/js"
)

type JSON = map[string]any

func ParseSavBytes(this js.Value, args []js.Value) any {
	fileBuffer := args[0]
	data := make([]byte, fileBuffer.Get("length").Int())
	js.CopyBytesToGo(data, fileBuffer)

	savData := sav.ParseSavBytes(data)
	return js.ValueOf(map[string]any{
		"trainer": jsconvert.ToJsValue(&savData.Trainer),
		// "pokedex":  save.Pokedex.ToJs(),
		"team": TeamToJs(&savData.Team),
		// "bag":  save.Bag.ToJs(),
		"pc": PCToJs(&savData.PC),
	})
}

func ParseGbaBytes(this js.Value, args []js.Value) any {
	fileBuffer := args[0]
	data := make([]byte, fileBuffer.Get("length").Int())
	js.CopyBytesToGo(data, fileBuffer)

	gbaData := gba.ParseGbaBytes(data)
	fmt.Printf("Loaded GBA Rom data\nSpecies: %d\nItems: %d\nMoves: %d\nAbilities: %d\nNatures: %d\n",
		len(gba.Species),
		len(gba.Items),
		len(gba.Moves),
		len(gba.Abilities),
		len(gba.Natures),
	)
	return jsconvert.ToJsValue(gbaData)
}

func main() {
	js.Global().Set("ParseSavBytes", js.FuncOf(ParseSavBytes))
	js.Global().Set("ParseGbaBytes", js.FuncOf(ParseGbaBytes))
	select {}
}

func PCToJs(pc *sav.PC) js.Value {
	jsMons := make([]js.Value, len(pc.Pokemon))
	for i, value := range pc.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return js.ValueOf(map[string]any{
		"currentBox": pc.CurrentBox,
		"pokemon":    jsconvert.ToJsArray(jsMons),
		"boxNames":   jsconvert.ToJsArray(pc.BoxNames),
	})
}

func TeamToJs(t *sav.Team) js.Value {
	jsMons := make([]js.Value, t.Size)
	for i, value := range t.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return jsconvert.ToJsArray(jsMons)
}

func PokemonToJs(p *sav.Pokemon) js.Value {
	if p.PokemonData.Species == 0 {
		return js.Null()
	}
	return js.ValueOf(JSON{
		"nickname": p.Nickname(),
		"species":  p.SpeciesName(),
		"item":     p.ItemName(),
		"level":    p.Level(),
		"toSDExportFormat": js.FuncOf(func(this js.Value, args []js.Value) any {
			return p.SDExportFormat()
		}),
	})
}
