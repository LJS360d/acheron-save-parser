package jsutils

import (
	"acheron-save-parser/sav"
	"reflect"
	"syscall/js"
)

func ToJsArray(value interface{}) js.Value {
	// Get the reflection value of the input
	refValue := reflect.ValueOf(value)

	// Check if the value is a slice or an array
	if refValue.Kind() != reflect.Slice && refValue.Kind() != reflect.Array {
		panic("provided value is neither a slice nor an array")
	}

	// Create a JavaScript array with the length of the Go slice or array
	jsArray := js.Global().Get("Array").New(refValue.Len())

	// Fill the JavaScript array with the values from the Go slice or array
	for i := 0; i < refValue.Len(); i++ {
		elem := refValue.Index(i).Interface()
		jsArray.SetIndex(i, js.ValueOf(elem))
	}

	return jsArray
}

func PCToJs(pc *sav.PC) js.Value {
	jsMons := make([]js.Value, len(pc.Pokemon))
	for i, value := range pc.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return js.ValueOf(map[string]interface{}{
		"currentBox": pc.CurrentBox,
		"pokemon":    ToJsArray(jsMons),
		"boxNames":   pc.BoxNames,
	})
}

func TeamToJs(t *sav.Team) js.Value {
	jsMons := make([]js.Value, t.Size)
	for i, value := range t.Pokemon {
		jsMons[i] = PokemonToJs(&value)
	}
	return ToJsArray(jsMons)
}

func PokemonToJs(p *sav.Pokemon) js.Value {
	if p.SpeciesName() == "None" {
		return js.Null()
	}
	return js.ValueOf(map[string]interface{}{
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
	return js.ValueOf(map[string]interface{}{
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
