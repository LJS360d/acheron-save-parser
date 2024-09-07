package main

import (
	"acheron-save-parser/export/jsutils"
	"acheron-save-parser/gba"
	"acheron-save-parser/sav"
	"fmt"
	"syscall/js"
)

type JSON = map[string]interface{}

func ParseSavBytes(this js.Value, args []js.Value) interface{} {
	fileBuffer := args[0]
	data := make([]byte, fileBuffer.Get("length").Int())
	js.CopyBytesToGo(data, fileBuffer)

	savData := sav.ParseSavBytes(data)
	return js.ValueOf(map[string]interface{}{
		// "trainer": save.Trainer.ToJs(),
		// "pokedex":  save.Pokedex.ToJs(),
		"team": TeamToJs(&savData.Team),
		// "bag":  save.Bag.ToJs(),
		"pc": PCToJs(&savData.PC),
	})
}

func ParseGbaBytes(this js.Value, args []js.Value) interface{} {
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
	return js.ValueOf(jsutils.ToJsonValue(gbaData))
}

func main() {
	js.Global().Set("ParseSavBytes", js.FuncOf(ParseSavBytes))
	js.Global().Set("ParseGbaBytes", js.FuncOf(ParseGbaBytes))
	select {}
}

func SavToJs(s *sav.SavData) js.Value {
	return jsutils.ToJsonValue(s)
}

func PCToJs(pc *sav.PC) js.Value {
	jsMons := make([]js.Value, len(pc.Pokemon))
	for i, value := range pc.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return js.ValueOf(map[string]interface{}{
		"currentBox": pc.CurrentBox,
		"pokemon":    jsutils.ToJsArray(jsMons),
		"boxNames":   jsutils.ToJsArray(pc.BoxNames),
	})
}

func TeamToJs(t *sav.Team) js.Value {
	jsMons := make([]js.Value, t.Size)
	for i, value := range t.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return jsutils.ToJsArray(jsMons)
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
		"toSDExportFormat": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return p.SDExportFormat()
		}),
	})
}

func TrainerToJs(t *sav.Trainer) js.Value {
	return js.ValueOf(JSON{
		"name":   t.Name(),
		"gender": t.Gender(),
		// "saveWarpFlags":   t.saveWarpFlags,
		"id": t.Id,
		// "publicId":        t.publicId,
		// "privateId":       t.privateId,
		"playtimeHours":   t.PlaytimeHours,
		"playtimeMinutes": t.PlaytimeMinutes,
		"playtimeSeconds": t.PlaytimeSeconds,
		// "playTimeVBlanks": t.playTimeVBlanks,
		// "options":         t.options,
		// "securityKey":     t.securityKey,
	})
}
