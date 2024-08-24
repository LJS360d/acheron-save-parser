package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Gen3Pokemon struct {
	name        string
	data        []byte
	personality uint32
	trainer     struct {
		id   uint32
		name string
	}
	species struct {
		id   uint16
		name string
		nid  int
	}
	exp      uint32
	expgroup string
	level    int
	moves    []struct {
		id   uint16
		name string
		pp   uint8
	}
	nature string
	ivs    map[string]int
	evs    map[string]uint8
}

func (g *Gen3Pokemon) kantoid(id int) int {
	/* if id <= 251 {
		return id
	}
	if id >= 413 {
		return 201
	}
	kanto := []int{252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 290, 291, 292, 276, 277, 285, 286, 327, 278, 279, 283, 284, 320, 321, 300, 301, 352, 343, 344, 299, 324, 302, 339, 340, 370, 341, 342, 349, 350, 318, 319, 328, 329, 330, 296, 297, 309, 310, 322, 323, 363, 364, 365, 331, 332, 361, 362, 337, 338, 298, 325, 326, 311, 312, 303, 307, 308, 333, 334, 360, 355, 356, 315, 287, 288, 289, 316, 317, 357, 293, 294, 295, 366, 367, 368, 359, 353, 354, 336, 335, 369, 304, 305, 306, 351, 313, 314, 345, 346, 347, 348, 280, 281, 282, 371, 372, 373, 374, 375, 376, 377, 378, 379, 382, 383, 384, 380, 381, 385, 386, 358}
	ix := id - 277
	if ix < len(kanto) {
		return kanto[ix]
	}
	return 0 */
	return id
}

func (g *Gen3Pokemon) speciesname(id int) string {
	/* names := []string{"Bulbasaur", "Ivysaur", "Venusaur", "Charmander", "Charmeleon", "Charizard", "Squirtle", "Wartortle", "Blastoise", "Caterpie", "Metapod", "Butterfree", "Weedle", "Kakuna", "Beedrill", "Pidgey", "Pidgeotto", "Pidgeot", "Rattata", "Raticate", "Spearow", "Fearow", "Ekans", "Arbok", "Pikachu", "Raichu", "Sandshrew", "Sandslash", "Nidoran (F)", "Nidorina", "Nidoqueen", "Nidoran (M)", "Nidorino", "Nidoking", "Clefairy", "Clefable", "Vulpix", "Ninetales", "Jigglypuff", "Wigglytuff", "Zubat", "Golbat", "Oddish", "Gloom", "Vileplume", "Paras", "Parasect", "Venonat", "Venomoth", "Diglett", "Dugtrio", "Meowth", "Persian", "Psyduck", "Golduck", "Mankey", "Primeape", "Growlithe", "Arcanine", "Poliwag", "Poliwhirl", "Poliwrath", "Abra", "Kadabra", "Alakazam", "Machop", "Machoke", "Machamp", "Bellsprout", "Weepinbell", "Victreebel", "Tentacool", "Tentacruel", "Geodude", "Graveler", "Golem", "Ponyta", "Rapidash", "Slowpoke", "Slowbro", "Magnemite", "Magneton", "Farfetch'd", "Doduo", "Dodrio", "Seel", "Dewgong", "Grimer", "Muk", "Shellder", "Cloyster", "Gastly", "Haunter", "Gengar", "Onix", "Drowzee", "Hypno", "Krabby", "Kingler", "Voltorb", "Electrode", "Exeggcute", "Exeggutor", "Cubone", "Marowak", "Hitmonlee", "Hitmonchan", "Lickitung", "Koffing", "Weezing", "Rhyhorn", "Rhydon", "Chansey", "Tangela", "Kangaskhan", "Horsea", "Seadra", "Goldeen", "Seaking", "Staryu", "Starmie", "Mr. Mime", "Scyther", "Jynx", "Electabuzz", "Magmar", "Pinsir", "Tauros", "Magikarp", "Gyarados", "Lapras", "Ditto", "Eevee", "Vaporeon", "Jolteon", "Flareon", "Porygon", "Omanyte", "Omastar", "Kabuto", "Kabutops", "Aerodactyl", "Snorlax", "Articuno", "Zapdos", "Moltres", "Dratini", "Dragonair", "Dragonite", "Mewtwo", "Mew", "Chikorita", "Bayleef", "Meganium", "Cyndaquil", "Quilava", "Typhlosion", "Totodile", "Croconaw", "Feraligatr", "Sentret", "Furret", "Hoothoot", "Noctowl", "Ledyba", "Ledian", "Spinarak", "Ariados", "Crobat", "Chinchou", "Lanturn", "Pichu", "Cleffa", "Igglybuff", "Togepi", "Togetic", "Natu", "Xatu", "Mareep", "Flaaffy", "Ampharos", "Bellossom", "Marill", "Azumarill", "Sudowoodo", "Politoed", "Hoppip", "Skiploom", "Jumpluff", "Aipom", "Sunkern", "Sunflora", "Yanma", "Wooper", "Quagsire", "Espeon", "Umbreon", "Murkrow", "Slowking", "Misdreavus", "Unown", "Wobbuffet", "Girafarig", "Pineco", "Forretress", "Dunsparce", "Gligar", "Steelix", "Snubbull", "Granbull", "Qwilfish", "Scizor", "Shuckle", "Heracross", "Sneasel", "Teddiursa", "Ursaring", "Slugma", "Magcargo", "Swinub", "Piloswine", "Corsola", "Remoraid", "Octillery", "Delibird", "Mantine", "Skarmory", "Houndour", "Houndoom", "Kingdra", "Phanpy", "Donphan", "Porygon2", "Stantler", "Smeargle", "Tyrogue", "Hitmontop", "Smoochum", "Elekid", "Magby", "Miltank", "Blissey", "Raikou", "Entei", "Suicune", "Larvitar", "Pupitar", "Tyranitar", "Lugia", "Ho-Oh", "Celebi", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Treecko", "Grovyle", "Sceptile", "Torchic", "Combusken", "Blaziken", "Mudkip", "Marshtomp", "Swampert", "Poochyena", "Mightyena", "Zigzagoon", "Linoone", "Wurmple", "Silcoon", "Beautifly", "Cascoon", "Dustox", "Lotad", "Lombre", "Ludicolo", "Seedot", "Nuzleaf", "Shiftry", "Nincada", "Ninjask", "Shedinja", "Taillow", "Swellow", "Shroomish", "Breloom", "Spinda", "Wingull", "Pelipper", "Surskit", "Masquerain", "Wailmer", "Wailord", "Skitty", "Delcatty", "Kecleon", "Baltoy", "Claydol", "Nosepass", "Torkoal", "Sableye", "Barboach", "Whiscash", "Luvdisc", "Corphish", "Crawdaunt", "Feebas", "Milotic", "Carvanha", "Sharpedo", "Trapinch", "Vibrava", "Flygon", "Makuhita", "Hariyama", "Electrike", "Manectric", "Numel", "Camerupt", "Spheal", "Sealeo", "Walrein", "Cacnea", "Cacturne", "Snorunt", "Glalie", "Lunatone", "Solrock", "Azurill", "Spoink", "Grumpig", "Plusle", "Minun", "Mawile", "Meditite", "Medicham", "Swablu", "Altaria", "Wynaut", "Duskull", "Dusclops", "Roselia", "Slakoth", "Vigoroth", "Slaking", "Gulpin", "Swalot", "Tropius", "Whismur", "Loudred", "Exploud", "Clamperl", "Huntail", "Gorebyss", "Absol", "Shuppet", "Banette", "Seviper", "Zangoose", "Relicanth", "Aron", "Lairon", "Aggron", "Castform", "Volbeat", "Illumise", "Lileep", "Cradily", "Anorith", "Armaldo", "Ralts", "Kirlia", "Gardevoir", "Bagon", "Shelgon", "Salamence", "Beldum", "Metang", "Metagross", "Regirock", "Regice", "Registeel", "Kyogre", "Groudon", "Rayquaza", "Latias", "Latios", "Jirachi", "Deoxys", "Chimecho", "Pok\u00e9mon Egg", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown", "Unown"}
	if id > len(names) {
		return ""
	} */
	return strconv.Itoa(id)
}

func (g *Gen3Pokemon) movename(id int) string {
	/* moves := []string{"Pound", "Karate Chop", "Double Slap", "Comet Punch", "Mega Punch", "Pay Day", "Fire Punch", "Ice Punch", "Thunder Punch", "Scratch", "Vice Grip", "Guillotine", "Razor Wind", "Swords Dance", "Cut", "Gust", "Wing Attack", "Whirlwind", "Fly", "Bind", "Slam", "Vine Whip", "Stomp", "Double Kick", "Mega Kick", "Jump Kick", "Rolling Kick", "Sand Attack", "Headbutt", "Horn Attack", "Fury Attack", "Horn Drill", "Tackle", "Body Slam", "Wrap", "Take Down", "Thrash", "Double-Edge", "Tail Whip", "Poison Sting", "Twineedle", "Pin Missile", "Leer", "Bite", "Growl", "Roar", "Sing", "Supersonic", "Sonic Boom", "Disable", "Acid", "Ember", "Flamethrower", "Mist", "Water Gun", "Hydro Pump", "Surf", "Ice Beam", "Blizzard", "Psybeam", "Bubble Beam", "Aurora Beam", "Hyper Beam", "Peck", "Drill Peck", "Submission", "Low Kick", "Counter", "Seismic Toss", "Strength", "Absorb", "Mega Drain", "Leech Seed", "Growth", "Razor Leaf", "Solar Beam", "Poison Powder", "Stun Spore", "Sleep Powder", "Petal Dance", "String Shot", "Dragon Rage", "Fire Spin", "Thunder Shock", "Thunderbolt", "Thunder Wave", "Thunder", "Rock Throw", "Earthquake", "Fissure", "Dig", "Toxic", "Confusion", "Psychic", "Hypnosis", "Meditate", "Agility", "Quick Attack", "Rage", "Teleport", "Night Shade", "Mimic", "Screech", "Double Team", "Recover", "Harden", "Minimize", "Smokescreen", "Confuse Ray", "Withdraw", "Defense Curl", "Barrier", "Light Screen", "Haze", "Reflect", "Focus Energy", "Bide", "Metronome", "Mirror Move", "Self-Destruct", "Egg Bomb", "Lick", "Smog", "Sludge", "Bone Club", "Fire Blast", "Waterfall", "Clamp", "Swift", "Skull Bash", "Spike Cannon", "Constrict", "Amnesia", "Kinesis", "Soft-Boiled", "High Jump Kick", "Glare", "Dream Eater", "Poison Gas", "Barrage", "Leech Life", "Lovely Kiss", "Sky Attack", "Transform", "Bubble", "Dizzy Punch", "Spore", "Flash", "Psywave", "Splash", "Acid Armor", "Crabhammer", "Explosion", "Fury Swipes", "Bonemerang", "Rest", "Rock Slide", "Hyper Fang", "Sharpen", "Conversion", "Tri Attack", "Super Fang", "Slash", "Substitute", "Struggle", "Sketch", "Triple Kick", "Thief", "Spider Web", "Mind Reader", "Nightmare", "Flame Wheel", "Snore", "Curse", "Flail", "Conversion 2", "Aeroblast", "Cotton Spore", "Reversal", "Spite", "Powder Snow", "Protect", "Mach Punch", "Scary Face", "Feint Attack", "Sweet Kiss", "Belly Drum", "Sludge Bomb", "Mud-Slap", "Octazooka", "Spikes", "Zap Cannon", "Foresight", "Destiny Bond", "Perish Song", "Icy Wind", "Detect", "Bone Rush", "Lock-On", "Outrage", "Sandstorm", "Giga Drain", "Endure", "Charm", "Rollout", "False Swipe", "Swagger", "Milk Drink", "Spark", "Fury Cutter", "Steel Wing", "Mean Look", "Attract", "Sleep Talk", "Heal Bell", "Return", "Present", "Frustration", "Safeguard", "Pain Split", "Sacred Fire", "Magnitude", "Dynamic Punch", "Megahorn", "Dragon Breath", "Baton Pass", "Encore", "Pursuit", "Rapid Spin", "Sweet Scent", "Iron Tail", "Metal Claw", "Vital Throw", "Morning Sun", "Synthesis", "Moonlight", "Hidden Power", "Cross Chop", "Twister", "Rain Dance", "Sunny Day", "Crunch", "Mirror Coat", "Psych Up", "Extreme Speed", "Ancient Power", "Shadow Ball", "Future Sight", "Rock Smash", "Whirlpool", "Beat Up", "Fake Out", "Uproar", "Stockpile", "Spit Up", "Swallow", "Heat Wave", "Hail", "Torment", "Flatter", "Will-O-Wisp", "Memento", "Facade", "Focus Punch", "Smelling Salts", "Follow Me", "Nature Power", "Charge", "Taunt", "Helping Hand", "Trick", "Role Play", "Wish", "Assist", "Ingrain", "Superpower", "Magic Coat", "Recycle", "Revenge", "Brick Break", "Yawn", "Knock Off", "Endeavor", "Eruption", "Skill Swap", "Imprison", "Refresh", "Grudge", "Snatch", "Secret Power", "Dive", "Arm Thrust", "Camouflage", "Tail Glow", "Luster Purge", "Mist Ball", "Feather Dance", "Teeter Dance", "Blaze Kick", "Mud Sport", "Ice Ball", "Needle Arm", "Slack Off", "Hyper Voice", "Poison Fang", "Crush Claw", "Blast Burn", "Hydro Cannon", "Meteor Mash", "Astonish", "Weather Ball", "Aromatherapy", "Fake Tears", "Air Cutter", "Overheat", "Odor Sleuth", "Rock Tomb", "Silver Wind", "Metal Sound", "Grass Whistle", "Tickle", "Cosmic Power", "Water Spout", "Signal Beam", "Shadow Punch", "Extrasensory", "Sky Uppercut", "Sand Tomb", "Sheer Cold", "Muddy Water", "Bullet Seed", "Aerial Ace", "Icicle Spear", "Iron Defense", "Block", "Howl", "Dragon Claw", "Frenzy Plant", "Bulk Up", "Bounce", "Mud Shot", "Poison Tail", "Covet", "Volt Tackle", "Magical Leaf", "Water Sport", "Calm Mind", "Leaf Blade", "Dragon Dance", "Rock Blast", "Shock Wave", "Water Pulse", "Doom Desire", "Psycho Boost"}
	return moves[id-1] */
	return strconv.Itoa(id)

}

func (g *Gen3Pokemon) naturename(id int) string {
	natures := []string{"Hardy", "Lonely", "Brave", "Adamant", "Naughty", "Bold", "Docile", "Relaxed", "Impish", "Lax", "Timid", "Hasty", "Serious", "Jolly", "Naive", "Modest", "Mild", "Quiet", "Bashful", "Rash", "Calm", "Gentle", "Sassy", "Careful", "Quirky"}
	return natures[id]
}

func (g *Gen3Pokemon) movetype(id int) string {
	return ""
}

func (g *Gen3Pokemon) _expgroup(id int) string {
	/* groups := []string{"Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Fast", "Fast", "Medium Fast", "Medium Fast", "Fast", "Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Slow", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Slow", "Slow", "Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Fast", "Fast", "Fast", "Fast", "Medium Fast", "Slow", "Slow", "Medium Fast", "Fast", "Fast", "Fast", "Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Fast", "Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Fast", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Fast", "Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Fast", "Fast", "Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Slow", "Fast", "Medium Fast", "Medium Fast", "Fast", "Slow", "Slow", "Slow", "Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Slow", "Fast", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Medium Slow", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Erratic", "Erratic", "Erratic", "Medium Slow", "Medium Slow", "Fluctuating", "Fluctuating", "Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Fluctuating", "Fluctuating", "Fast", "Fast", "Medium Slow", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Fast", "Medium Fast", "Fast", "Fluctuating", "Fluctuating", "Erratic", "Erratic", "Slow", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Fluctuating", "Fluctuating", "Slow", "Slow", "Medium Fast", "Medium Fast", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Medium Fast", "Medium Fast", "Fast", "Fast", "Fast", "Fast", "Fast", "Medium Fast", "Medium Fast", "Fast", "Medium Fast", "Medium Fast", "Erratic", "Erratic", "Medium Fast", "Fast", "Fast", "Medium Slow", "Slow", "Slow", "Slow", "Fluctuating", "Fluctuating", "Slow", "Medium Slow", "Medium Slow", "Medium Slow", "Erratic", "Erratic", "Erratic", "Medium Slow", "Fast", "Fast", "Fluctuating", "Erratic", "Slow", "Slow", "Slow", "Slow", "Medium Fast", "Erratic", "Fluctuating", "Erratic", "Erratic", "Erratic", "Erratic", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Slow", "Fast"}
	return groups[id-1] */
	return strconv.Itoa(id)
}

func (g *Gen3Pokemon) _level(expgroup string, exp uint32) int {
	var levels []uint32
	level := 0

	switch expgroup {
	case "Fast":
		levels = []uint32{}
	case "Medium Fast":
		levels = []uint32{}
	case "Medium Slow":
		levels = []uint32{}
	case "Slow":
		levels = []uint32{}
	case "Fluctuating":
		levels = []uint32{}
	case "Erratic":
		levels = []uint32{}
	default:
		levels = []uint32{}
	}

	for i := 0; i < 100; i++ {
		if exp > levels[i] {
			level = i
		} else {
			break
		}
	}
	return level
}

func NewGen3Pokemon(pkm []byte) *Gen3Pokemon {
	orders := []string{"GAEM", "GAME", "GEAM", "GEMA", "GMAE", "GMEA", "AGEM", "AGME", "AEGM", "AEMG", "AMGE", "AMEG", "EGAM", "EGMA", "EAGM", "EAMG", "EMGA", "EMAG", "MGAE", "MGEA", "MAGE", "MAEG", "MEGA", "MEAG"}
	g := &Gen3Pokemon{}
	g.name = g.readstring(pkm[8:18])
	g.data = pkm
	trainer := g.readstring(pkm[20:27])
	if trainer == "" && g.name == "" {
		return nil
	}
	g.personality = binary.LittleEndian.Uint32(pkm[0:4])
	g.trainer.id = binary.LittleEndian.Uint32(pkm[4:8])
	g.trainer.name = trainer
	key := g.trainer.id ^ g.personality

	data := pkm[32:]
	order := g.personality % 24
	orderstring := orders[order]
	sections := make(map[byte][]byte)
	for i := 0; i < 4; i++ {
		section := orderstring[i]
		sectiondata := data[i*12 : (i+1)*12]
		decr := g.decryptsubsection(sectiondata, key)
		sections[section] = decr
	}

	var decrypted bytes.Buffer
	for _, section := range []byte{'G', 'A', 'E', 'M'} {
		decrypted.Write(sections[section])
	}
	g.data = append(pkm[:32], decrypted.Bytes()...)

	g.species.id = binary.LittleEndian.Uint16(sections['G'][0:2])
	g.species.name = g.speciesname(int(g.species.id))
	g.species.nid = g.kantoid(int(g.species.id))
	g.exp = binary.LittleEndian.Uint32(sections['G'][4:8])
	g.expgroup = g._expgroup(int(g.species.id))
	g.level = g._level(g.expgroup, g.exp)
	if g.name == "" {
		g.name = strings.ToUpper(g.species.name)
	}

	for i := 0; i < 4; i++ {
		moveID := binary.LittleEndian.Uint16(sections['A'][i*2 : (i+1)*2])
		if moveID == 0 {
			continue
		}
		g.moves = append(g.moves, struct {
			id   uint16
			name string
			pp   uint8
		}{
			id:   moveID,
			name: g.movename(int(moveID)),
			pp:   sections['A'][i+8],
		})
	}

	g.nature = g.naturename(int(g.personality % 25))
	g.ivs = g.getivs(binary.LittleEndian.Uint32(sections['M'][4:8]))
	g.evs = g.getevs(sections['E'])

	return g
}

func (g *Gen3Pokemon) getivs(value uint32) map[string]int {
	iv := make(map[string]int)
	bitstring := fmt.Sprintf("%032b", value)
	iv["hp"] = int(binary.BigEndian.Uint32([]byte(bitstring[27:32] + "000")))
	iv["attack"] = int(binary.BigEndian.Uint32([]byte(bitstring[22:27] + "000")))
	iv["defence"] = int(binary.BigEndian.Uint32([]byte(bitstring[17:22] + "000")))
	iv["speed"] = int(binary.BigEndian.Uint32([]byte(bitstring[12:17] + "000")))
	iv["spatk"] = int(binary.BigEndian.Uint32([]byte(bitstring[7:12] + "000")))
	iv["spdef"] = int(binary.BigEndian.Uint32([]byte(bitstring[2:7] + "000")))
	return iv
}

func (g *Gen3Pokemon) getevs(section []byte) map[string]uint8 {
	ev := make(map[string]uint8)
	ev["hp"] = section[0]
	ev["attack"] = section[1]
	ev["defence"] = section[2]
	ev["speed"] = section[3]
	ev["spatk"] = section[4]
	ev["spdef"] = section[5]
	ev["cool"] = section[6]
	ev["beauty"] = section[7]
	ev["cute"] = section[8]
	ev["smart"] = section[9]
	ev["tough"] = section[10]
	ev["feel"] = section[11]
	return ev
}

func (g *Gen3Pokemon) decryptsubsection(data []byte, key uint32) []byte {
	if len(data) != 12 {
		return []byte{}
	}
	a := binary.LittleEndian.Uint32(data[0:4]) ^ key
	b := binary.LittleEndian.Uint32(data[4:8]) ^ key
	c := binary.LittleEndian.Uint32(data[8:12]) ^ key
	result := make([]byte, 12)
	binary.LittleEndian.PutUint32(result[0:4], a)
	binary.LittleEndian.PutUint32(result[4:8], b)
	binary.LittleEndian.PutUint32(result[8:12], c)
	return result
}

func (g *Gen3Pokemon) readstring(text []byte) string {
	chars := "0123456789!?.-         ,  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var ret strings.Builder
	for _, i := range text {
		c := int(i) - 161
		if c < 0 || c >= len(chars) {
			ret.WriteRune(' ')
		} else {
			ret.WriteByte(chars[c])
		}
	}
	return strings.TrimSpace(ret.String())
}
