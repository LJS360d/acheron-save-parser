package data

var SpeciesName = map[uint16]string{
	0:    "None",
	1:    "Bulbasaur",
	2:    "Ivysaur",
	3:    "Venusaur",
	4:    "Charmander",
	5:    "Charmeleon",
	6:    "Charizard",
	7:    "Squirtle",
	8:    "Wartortle",
	9:    "Blastoise",
	10:   "Caterpie",
	11:   "Metapod",
	12:   "Butterfree",
	13:   "Weedle",
	14:   "Kakuna",
	15:   "Beedrill",
	16:   "Pidgey",
	17:   "Pidgeotto",
	18:   "Pidgeot",
	19:   "Rattata",
	20:   "Raticate",
	21:   "Spearow",
	22:   "Fearow",
	23:   "Ekans",
	24:   "Arbok",
	25:   "Pikachu",
	26:   "Raichu",
	27:   "Sandshrew",
	28:   "Sandslash",
	29:   "Nidoran-F",
	30:   "Nidorina",
	31:   "Nidoqueen",
	32:   "Nidoran-M",
	33:   "Nidorino",
	34:   "Nidoking",
	35:   "Clefairy",
	36:   "Clefable",
	37:   "Vulpix",
	38:   "Ninetales",
	39:   "Jigglypuff",
	40:   "Wigglytuff",
	41:   "Zubat",
	42:   "Golbat",
	43:   "Oddish",
	44:   "Gloom",
	45:   "Vileplume",
	46:   "Paras",
	47:   "Parasect",
	48:   "Venonat",
	49:   "Venomoth",
	50:   "Diglett",
	51:   "Dugtrio",
	52:   "Meowth",
	53:   "Persian",
	54:   "Psyduck",
	55:   "Golduck",
	56:   "Mankey",
	57:   "Primeape",
	58:   "Growlithe",
	59:   "Arcanine",
	60:   "Poliwag",
	61:   "Poliwhirl",
	62:   "Poliwrath",
	63:   "Abra",
	64:   "Kadabra",
	65:   "Alakazam",
	66:   "Machop",
	67:   "Machoke",
	68:   "Machamp",
	69:   "Bellsprout",
	70:   "Weepinbell",
	71:   "Victreebel",
	72:   "Tentacool",
	73:   "Tentacruel",
	74:   "Geodude",
	75:   "Graveler",
	76:   "Golem",
	77:   "Ponyta",
	78:   "Rapidash",
	79:   "Slowpoke",
	80:   "Slowbro",
	81:   "Magnemite",
	82:   "Magneton",
	83:   "Farfetchd",
	84:   "Doduo",
	85:   "Dodrio",
	86:   "Seel",
	87:   "Dewgong",
	88:   "Grimer",
	89:   "Muk",
	90:   "Shellder",
	91:   "Cloyster",
	92:   "Gastly",
	93:   "Haunter",
	94:   "Gengar",
	95:   "Onix",
	96:   "Drowzee",
	97:   "Hypno",
	98:   "Krabby",
	99:   "Kingler",
	100:  "Voltorb",
	101:  "Electrode",
	102:  "Exeggcute",
	103:  "Exeggutor",
	104:  "Cubone",
	105:  "Marowak",
	106:  "Hitmonlee",
	107:  "Hitmonchan",
	108:  "Lickitung",
	109:  "Koffing",
	110:  "Weezing",
	111:  "Rhyhorn",
	112:  "Rhydon",
	113:  "Chansey",
	114:  "Tangela",
	115:  "Kangaskhan",
	116:  "Horsea",
	117:  "Seadra",
	118:  "Goldeen",
	119:  "Seaking",
	120:  "Staryu",
	121:  "Starmie",
	122:  "Mr-Mime",
	123:  "Scyther",
	124:  "Jynx",
	125:  "Electabuzz",
	126:  "Magmar",
	127:  "Pinsir",
	128:  "Tauros",
	129:  "Magikarp",
	130:  "Gyarados",
	131:  "Lapras",
	132:  "Ditto",
	133:  "Eevee",
	134:  "Vaporeon",
	135:  "Jolteon",
	136:  "Flareon",
	137:  "Porygon",
	138:  "Omanyte",
	139:  "Omastar",
	140:  "Kabuto",
	141:  "Kabutops",
	142:  "Aerodactyl",
	143:  "Snorlax",
	144:  "Articuno",
	145:  "Zapdos",
	146:  "Moltres",
	147:  "Dratini",
	148:  "Dragonair",
	149:  "Dragonite",
	150:  "Mewtwo",
	151:  "Mew",
	152:  "Chikorita",
	153:  "Bayleef",
	154:  "Meganium",
	155:  "Cyndaquil",
	156:  "Quilava",
	157:  "Typhlosion",
	158:  "Totodile",
	159:  "Croconaw",
	160:  "Feraligatr",
	161:  "Sentret",
	162:  "Furret",
	163:  "Hoothoot",
	164:  "Noctowl",
	165:  "Ledyba",
	166:  "Ledian",
	167:  "Spinarak",
	168:  "Ariados",
	169:  "Crobat",
	170:  "Chinchou",
	171:  "Lanturn",
	172:  "Pichu",
	173:  "Cleffa",
	174:  "Igglybuff",
	175:  "Togepi",
	176:  "Togetic",
	177:  "Natu",
	178:  "Xatu",
	179:  "Mareep",
	180:  "Flaaffy",
	181:  "Ampharos",
	182:  "Bellossom",
	183:  "Marill",
	184:  "Azumarill",
	185:  "Sudowoodo",
	186:  "Politoed",
	187:  "Hoppip",
	188:  "Skiploom",
	189:  "Jumpluff",
	190:  "Aipom",
	191:  "Sunkern",
	192:  "Sunflora",
	193:  "Yanma",
	194:  "Wooper",
	195:  "Quagsire",
	196:  "Espeon",
	197:  "Umbreon",
	198:  "Murkrow",
	199:  "Slowking",
	200:  "Misdreavus",
	201:  "Unown",
	202:  "Wobbuffet",
	203:  "Girafarig",
	204:  "Pineco",
	205:  "Forretress",
	206:  "Dunsparce",
	207:  "Gligar",
	208:  "Steelix",
	209:  "Snubbull",
	210:  "Granbull",
	211:  "Qwilfish",
	212:  "Scizor",
	213:  "Shuckle",
	214:  "Heracross",
	215:  "Sneasel",
	216:  "Teddiursa",
	217:  "Ursaring",
	218:  "Slugma",
	219:  "Magcargo",
	220:  "Swinub",
	221:  "Piloswine",
	222:  "Corsola",
	223:  "Remoraid",
	224:  "Octillery",
	225:  "Delibird",
	226:  "Mantine",
	227:  "Skarmory",
	228:  "Houndour",
	229:  "Houndoom",
	230:  "Kingdra",
	231:  "Phanpy",
	232:  "Donphan",
	233:  "Porygon2",
	234:  "Stantler",
	235:  "Smeargle",
	236:  "Tyrogue",
	237:  "Hitmontop",
	238:  "Smoochum",
	239:  "Elekid",
	240:  "Magby",
	241:  "Miltank",
	242:  "Blissey",
	243:  "Raikou",
	244:  "Entei",
	245:  "Suicune",
	246:  "Larvitar",
	247:  "Pupitar",
	248:  "Tyranitar",
	249:  "Lugia",
	250:  "Ho-Oh",
	251:  "Celebi",
	252:  "Treecko",
	253:  "Grovyle",
	254:  "Sceptile",
	255:  "Torchic",
	256:  "Combusken",
	257:  "Blaziken",
	258:  "Mudkip",
	259:  "Marshtomp",
	260:  "Swampert",
	261:  "Poochyena",
	262:  "Mightyena",
	263:  "Zigzagoon",
	264:  "Linoone",
	265:  "Wurmple",
	266:  "Silcoon",
	267:  "Beautifly",
	268:  "Cascoon",
	269:  "Dustox",
	270:  "Lotad",
	271:  "Lombre",
	272:  "Ludicolo",
	273:  "Seedot",
	274:  "Nuzleaf",
	275:  "Shiftry",
	276:  "Taillow",
	277:  "Swellow",
	278:  "Wingull",
	279:  "Pelipper",
	280:  "Ralts",
	281:  "Kirlia",
	282:  "Gardevoir",
	283:  "Surskit",
	284:  "Masquerain",
	285:  "Shroomish",
	286:  "Breloom",
	287:  "Slakoth",
	288:  "Vigoroth",
	289:  "Slaking",
	290:  "Nincada",
	291:  "Ninjask",
	292:  "Shedinja",
	293:  "Whismur",
	294:  "Loudred",
	295:  "Exploud",
	296:  "Makuhita",
	297:  "Hariyama",
	298:  "Azurill",
	299:  "Nosepass",
	300:  "Skitty",
	301:  "Delcatty",
	302:  "Sableye",
	303:  "Mawile",
	304:  "Aron",
	305:  "Lairon",
	306:  "Aggron",
	307:  "Meditite",
	308:  "Medicham",
	309:  "Electrike",
	310:  "Manectric",
	311:  "Plusle",
	312:  "Minun",
	313:  "Volbeat",
	314:  "Illumise",
	315:  "Roselia",
	316:  "Gulpin",
	317:  "Swalot",
	318:  "Carvanha",
	319:  "Sharpedo",
	320:  "Wailmer",
	321:  "Wailord",
	322:  "Numel",
	323:  "Camerupt",
	324:  "Torkoal",
	325:  "Spoink",
	326:  "Grumpig",
	327:  "Spinda",
	328:  "Trapinch",
	329:  "Vibrava",
	330:  "Flygon",
	331:  "Cacnea",
	332:  "Cacturne",
	333:  "Swablu",
	334:  "Altaria",
	335:  "Zangoose",
	336:  "Seviper",
	337:  "Lunatone",
	338:  "Solrock",
	339:  "Barboach",
	340:  "Whiscash",
	341:  "Corphish",
	342:  "Crawdaunt",
	343:  "Baltoy",
	344:  "Claydol",
	345:  "Lileep",
	346:  "Cradily",
	347:  "Anorith",
	348:  "Armaldo",
	349:  "Feebas",
	350:  "Milotic",
	351:  "Castform",
	352:  "Kecleon",
	353:  "Shuppet",
	354:  "Banette",
	355:  "Duskull",
	356:  "Dusclops",
	357:  "Tropius",
	358:  "Chimecho",
	359:  "Absol",
	360:  "Wynaut",
	361:  "Snorunt",
	362:  "Glalie",
	363:  "Spheal",
	364:  "Sealeo",
	365:  "Walrein",
	366:  "Clamperl",
	367:  "Huntail",
	368:  "Gorebyss",
	369:  "Relicanth",
	370:  "Luvdisc",
	371:  "Bagon",
	372:  "Shelgon",
	373:  "Salamence",
	374:  "Beldum",
	375:  "Metang",
	376:  "Metagross",
	377:  "Regirock",
	378:  "Regice",
	379:  "Registeel",
	380:  "Latias",
	381:  "Latios",
	382:  "Kyogre",
	383:  "Groudon",
	384:  "Rayquaza",
	385:  "Jirachi",
	386:  "Deoxys",
	387:  "Turtwig",
	388:  "Grotle",
	389:  "Torterra",
	390:  "Chimchar",
	391:  "Monferno",
	392:  "Infernape",
	393:  "Piplup",
	394:  "Prinplup",
	395:  "Empoleon",
	396:  "Starly",
	397:  "Staravia",
	398:  "Staraptor",
	399:  "Bidoof",
	400:  "Bibarel",
	401:  "Kricketot",
	402:  "Kricketune",
	403:  "Shinx",
	404:  "Luxio",
	405:  "Luxray",
	406:  "Budew",
	407:  "Roserade",
	408:  "Cranidos",
	409:  "Rampardos",
	410:  "Shieldon",
	411:  "Bastiodon",
	412:  "Burmy",
	413:  "Wormadam",
	414:  "Mothim",
	415:  "Combee",
	416:  "Vespiquen",
	417:  "Pachirisu",
	418:  "Buizel",
	419:  "Floatzel",
	420:  "Cherubi",
	421:  "Cherrim",
	422:  "Shellos",
	423:  "Gastrodon",
	424:  "Ambipom",
	425:  "Drifloon",
	426:  "Drifblim",
	427:  "Buneary",
	428:  "Lopunny",
	429:  "Mismagius",
	430:  "Honchkrow",
	431:  "Glameow",
	432:  "Purugly",
	433:  "Chingling",
	434:  "Stunky",
	435:  "Skuntank",
	436:  "Bronzor",
	437:  "Bronzong",
	438:  "Bonsly",
	439:  "Mime-Jr",
	440:  "Happiny",
	441:  "Chatot",
	442:  "Spiritomb",
	443:  "Gible",
	444:  "Gabite",
	445:  "Garchomp",
	446:  "Munchlax",
	447:  "Riolu",
	448:  "Lucario",
	449:  "Hippopotas",
	450:  "Hippowdon",
	451:  "Skorupi",
	452:  "Drapion",
	453:  "Croagunk",
	454:  "Toxicroak",
	455:  "Carnivine",
	456:  "Finneon",
	457:  "Lumineon",
	458:  "Mantyke",
	459:  "Snover",
	460:  "Abomasnow",
	461:  "Weavile",
	462:  "Magnezone",
	463:  "Lickilicky",
	464:  "Rhyperior",
	465:  "Tangrowth",
	466:  "Electivire",
	467:  "Magmortar",
	468:  "Togekiss",
	469:  "Yanmega",
	470:  "Leafeon",
	471:  "Glaceon",
	472:  "Gliscor",
	473:  "Mamoswine",
	474:  "Porygon-Z",
	475:  "Gallade",
	476:  "Probopass",
	477:  "Dusknoir",
	478:  "Froslass",
	479:  "Rotom",
	480:  "Uxie",
	481:  "Mesprit",
	482:  "Azelf",
	483:  "Dialga",
	484:  "Palkia",
	485:  "Heatran",
	486:  "Regigigas",
	487:  "Giratina",
	488:  "Cresselia",
	489:  "Phione",
	490:  "Manaphy",
	491:  "Darkrai",
	492:  "Shaymin",
	493:  "Arceus",
	494:  "Victini",
	495:  "Snivy",
	496:  "Servine",
	497:  "Serperior",
	498:  "Tepig",
	499:  "Pignite",
	500:  "Emboar",
	501:  "Oshawott",
	502:  "Dewott",
	503:  "Samurott",
	504:  "Patrat",
	505:  "Watchog",
	506:  "Lillipup",
	507:  "Herdier",
	508:  "Stoutland",
	509:  "Purrloin",
	510:  "Liepard",
	511:  "Pansage",
	512:  "Simisage",
	513:  "Pansear",
	514:  "Simisear",
	515:  "Panpour",
	516:  "Simipour",
	517:  "Munna",
	518:  "Musharna",
	519:  "Pidove",
	520:  "Tranquill",
	521:  "Unfezant",
	522:  "Blitzle",
	523:  "Zebstrika",
	524:  "Roggenrola",
	525:  "Boldore",
	526:  "Gigalith",
	527:  "Woobat",
	528:  "Swoobat",
	529:  "Drilbur",
	530:  "Excadrill",
	531:  "Audino",
	532:  "Timburr",
	533:  "Gurdurr",
	534:  "Conkeldurr",
	535:  "Tympole",
	536:  "Palpitoad",
	537:  "Seismitoad",
	538:  "Throh",
	539:  "Sawk",
	540:  "Sewaddle",
	541:  "Swadloon",
	542:  "Leavanny",
	543:  "Venipede",
	544:  "Whirlipede",
	545:  "Scolipede",
	546:  "Cottonee",
	547:  "Whimsicott",
	548:  "Petilil",
	549:  "Lilligant",
	550:  "Basculin",
	551:  "Sandile",
	552:  "Krokorok",
	553:  "Krookodile",
	554:  "Darumaka",
	555:  "Darmanitan",
	556:  "Maractus",
	557:  "Dwebble",
	558:  "Crustle",
	559:  "Scraggy",
	560:  "Scrafty",
	561:  "Sigilyph",
	562:  "Yamask",
	563:  "Cofagrigus",
	564:  "Tirtouga",
	565:  "Carracosta",
	566:  "Archen",
	567:  "Archeops",
	568:  "Trubbish",
	569:  "Garbodor",
	570:  "Zorua",
	571:  "Zoroark",
	572:  "Minccino",
	573:  "Cinccino",
	574:  "Gothita",
	575:  "Gothorita",
	576:  "Gothitelle",
	577:  "Solosis",
	578:  "Duosion",
	579:  "Reuniclus",
	580:  "Ducklett",
	581:  "Swanna",
	582:  "Vanillite",
	583:  "Vanillish",
	584:  "Vanilluxe",
	585:  "Deerling",
	586:  "Sawsbuck",
	587:  "Emolga",
	588:  "Karrablast",
	589:  "Escavalier",
	590:  "Foongus",
	591:  "Amoonguss",
	592:  "Frillish",
	593:  "Jellicent",
	594:  "Alomomola",
	595:  "Joltik",
	596:  "Galvantula",
	597:  "Ferroseed",
	598:  "Ferrothorn",
	599:  "Klink",
	600:  "Klang",
	601:  "Klinklang",
	602:  "Tynamo",
	603:  "Eelektrik",
	604:  "Eelektross",
	605:  "Elgyem",
	606:  "Beheeyem",
	607:  "Litwick",
	608:  "Lampent",
	609:  "Chandelure",
	610:  "Axew",
	611:  "Fraxure",
	612:  "Haxorus",
	613:  "Cubchoo",
	614:  "Beartic",
	615:  "Cryogonal",
	616:  "Shelmet",
	617:  "Accelgor",
	618:  "Stunfisk",
	619:  "Mienfoo",
	620:  "Mienshao",
	621:  "Druddigon",
	622:  "Golett",
	623:  "Golurk",
	624:  "Pawniard",
	625:  "Bisharp",
	626:  "Bouffalant",
	627:  "Rufflet",
	628:  "Braviary",
	629:  "Vullaby",
	630:  "Mandibuzz",
	631:  "Heatmor",
	632:  "Durant",
	633:  "Deino",
	634:  "Zweilous",
	635:  "Hydreigon",
	636:  "Larvesta",
	637:  "Volcarona",
	638:  "Cobalion",
	639:  "Terrakion",
	640:  "Virizion",
	641:  "Tornadus",
	642:  "Thundurus",
	643:  "Reshiram",
	644:  "Zekrom",
	645:  "Landorus",
	646:  "Kyurem",
	647:  "Keldeo",
	648:  "Meloetta",
	649:  "Genesect",
	650:  "Chespin",
	651:  "Quilladin",
	652:  "Chesnaught",
	653:  "Fennekin",
	654:  "Braixen",
	655:  "Delphox",
	656:  "Froakie",
	657:  "Frogadier",
	658:  "Greninja",
	659:  "Bunnelby",
	660:  "Diggersby",
	661:  "Fletchling",
	662:  "Fletchinder",
	663:  "Talonflame",
	664:  "Scatterbug",
	665:  "Spewpa",
	666:  "Vivillon",
	667:  "Litleo",
	668:  "Pyroar",
	669:  "Flabebe",
	670:  "Floette",
	671:  "Florges",
	672:  "Skiddo",
	673:  "Gogoat",
	674:  "Pancham",
	675:  "Pangoro",
	676:  "Furfrou",
	677:  "Espurr",
	678:  "Meowstic",
	679:  "Honedge",
	680:  "Doublade",
	681:  "Aegislash",
	682:  "Spritzee",
	683:  "Aromatisse",
	684:  "Swirlix",
	685:  "Slurpuff",
	686:  "Inkay",
	687:  "Malamar",
	688:  "Binacle",
	689:  "Barbaracle",
	690:  "Skrelp",
	691:  "Dragalge",
	692:  "Clauncher",
	693:  "Clawitzer",
	694:  "Helioptile",
	695:  "Heliolisk",
	696:  "Tyrunt",
	697:  "Tyrantrum",
	698:  "Amaura",
	699:  "Aurorus",
	700:  "Sylveon",
	701:  "Hawlucha",
	702:  "Dedenne",
	703:  "Carbink",
	704:  "Goomy",
	705:  "Sliggoo",
	706:  "Goodra",
	707:  "Klefki",
	708:  "Phantump",
	709:  "Trevenant",
	710:  "Pumpkaboo",
	711:  "Gourgeist",
	712:  "Bergmite",
	713:  "Avalugg",
	714:  "Noibat",
	715:  "Noivern",
	716:  "Xerneas",
	717:  "Yveltal",
	718:  "Zygarde-50%",
	719:  "Diancie",
	720:  "Hoopa",
	721:  "Volcanion",
	722:  "Rowlet",
	723:  "Dartrix",
	724:  "Decidueye",
	725:  "Litten",
	726:  "Torracat",
	727:  "Incineroar",
	728:  "Popplio",
	729:  "Brionne",
	730:  "Primarina",
	731:  "Pikipek",
	732:  "Trumbeak",
	733:  "Toucannon",
	734:  "Yungoos",
	735:  "Gumshoos",
	736:  "Grubbin",
	737:  "Charjabug",
	738:  "Vikavolt",
	739:  "Crabrawler",
	740:  "Crabominable",
	741:  "Oricorio",
	742:  "Cutiefly",
	743:  "Ribombee",
	744:  "Rockruff",
	745:  "Lycanroc",
	746:  "Wishiwashi",
	747:  "Mareanie",
	748:  "Toxapex",
	749:  "Mudbray",
	750:  "Mudsdale",
	751:  "Dewpider",
	752:  "Araquanid",
	753:  "Fomantis",
	754:  "Lurantis",
	755:  "Morelull",
	756:  "Shiinotic",
	757:  "Salandit",
	758:  "Salazzle",
	759:  "Stufful",
	760:  "Bewear",
	761:  "Bounsweet",
	762:  "Steenee",
	763:  "Tsareena",
	764:  "Comfey",
	765:  "Oranguru",
	766:  "Passimian",
	767:  "Wimpod",
	768:  "Golisopod",
	769:  "Sandygast",
	770:  "Palossand",
	771:  "Pyukumuku",
	772:  "Type-Null",
	773:  "Silvally",
	774:  "Minior-Red",
	775:  "Komala",
	776:  "Turtonator",
	777:  "Togedemaru",
	778:  "Mimikyu",
	779:  "Bruxish",
	780:  "Drampa",
	781:  "Dhelmise",
	782:  "Jangmo-O",
	783:  "Hakamo-O",
	784:  "Kommo-O",
	785:  "Tapu-Koko",
	786:  "Tapu-Lele",
	787:  "Tapu-Bulu",
	788:  "Tapu-Fini",
	789:  "Cosmog",
	790:  "Cosmoem",
	791:  "Solgaleo",
	792:  "Lunala",
	793:  "Nihilego",
	794:  "Buzzwole",
	795:  "Pheromosa",
	796:  "Xurkitree",
	797:  "Celesteela",
	798:  "Kartana",
	799:  "Guzzlord",
	800:  "Necrozma",
	801:  "Magearna",
	802:  "Marshadow",
	803:  "Poipole",
	804:  "Naganadel",
	805:  "Stakataka",
	806:  "Blacephalon",
	807:  "Zeraora",
	808:  "Meltan",
	809:  "Melmetal",
	810:  "Grookey",
	811:  "Thwackey",
	812:  "Rillaboom",
	813:  "Scorbunny",
	814:  "Raboot",
	815:  "Cinderace",
	816:  "Sobble",
	817:  "Drizzile",
	818:  "Inteleon",
	819:  "Skwovet",
	820:  "Greedent",
	821:  "Rookidee",
	822:  "Corvisquire",
	823:  "Corviknight",
	824:  "Blipbug",
	825:  "Dottler",
	826:  "Orbeetle",
	827:  "Nickit",
	828:  "Thievul",
	829:  "Gossifleur",
	830:  "Eldegoss",
	831:  "Wooloo",
	832:  "Dubwool",
	833:  "Chewtle",
	834:  "Drednaw",
	835:  "Yamper",
	836:  "Boltund",
	837:  "Rolycoly",
	838:  "Carkol",
	839:  "Coalossal",
	840:  "Applin",
	841:  "Flapple",
	842:  "Appletun",
	843:  "Silicobra",
	844:  "Sandaconda",
	845:  "Cramorant",
	846:  "Arrokuda",
	847:  "Barraskewda",
	848:  "Toxel",
	849:  "Toxtricity",
	850:  "Sizzlipede",
	851:  "Centiskorch",
	852:  "Clobbopus",
	853:  "Grapploct",
	854:  "Sinistea",
	855:  "Polteageist",
	856:  "Hatenna",
	857:  "Hattrem",
	858:  "Hatterene",
	859:  "Impidimp",
	860:  "Morgrem",
	861:  "Grimmsnarl",
	862:  "Obstagoon",
	863:  "Perrserker",
	864:  "Cursola",
	865:  "Sirfetchd",
	866:  "Mr-Rime",
	867:  "Runerigus",
	868:  "Milcery",
	869:  "Alcremie-Vanilla-Cream",
	870:  "Falinks",
	871:  "Pincurchin",
	872:  "Snom",
	873:  "Frosmoth",
	874:  "Stonjourner",
	875:  "Eiscue",
	876:  "Indeedee",
	877:  "Morpeko",
	878:  "Cufant",
	879:  "Copperajah",
	880:  "Dracozolt",
	881:  "Arctozolt",
	882:  "Dracovish",
	883:  "Arctovish",
	884:  "Duraludon",
	885:  "Dreepy",
	886:  "Drakloak",
	887:  "Dragapult",
	888:  "Zacian",
	889:  "Zamazenta",
	890:  "Eternatus",
	891:  "Kubfu",
	892:  "Urshifu",
	893:  "Zarude",
	894:  "Regieleki",
	895:  "Regidrago",
	896:  "Glastrier",
	897:  "Spectrier",
	898:  "Calyrex",
	899:  "Wyrdeer",
	900:  "Kleavor",
	901:  "Ursaluna",
	902:  "Basculegion",
	903:  "Sneasler",
	904:  "Overqwil",
	905:  "Enamorus",
	906:  "Venusaur-Mega",
	907:  "Charizard-Mega-X",
	908:  "Charizard-Mega-Y",
	909:  "Blastoise-Mega",
	910:  "Beedrill-Mega",
	911:  "Pidgeot-Mega",
	912:  "Alakazam-Mega",
	913:  "Slowbro-Mega",
	914:  "Gengar-Mega",
	915:  "Kangaskhan-Mega",
	916:  "Pinsir-Mega",
	917:  "Gyarados-Mega",
	918:  "Aerodactyl-Mega",
	919:  "Mewtwo-Mega-X",
	920:  "Mewtwo-Mega-Y",
	921:  "Ampharos-Mega",
	922:  "Steelix-Mega",
	923:  "Scizor-Mega",
	924:  "Heracross-Mega",
	925:  "Houndoom-Mega",
	926:  "Tyranitar-Mega",
	927:  "Sceptile-Mega",
	928:  "Blaziken-Mega",
	929:  "Swampert-Mega",
	930:  "Gardevoir-Mega",
	931:  "Sableye-Mega",
	932:  "Mawile-Mega",
	933:  "Aggron-Mega",
	934:  "Medicham-Mega",
	935:  "Manectric-Mega",
	936:  "Sharpedo-Mega",
	937:  "Camerupt-Mega",
	938:  "Altaria-Mega",
	939:  "Banette-Mega",
	940:  "Absol-Mega",
	941:  "Glalie-Mega",
	942:  "Salamence-Mega",
	943:  "Metagross-Mega",
	944:  "Latias-Mega",
	945:  "Latios-Mega",
	946:  "Lopunny-Mega",
	947:  "Garchomp-Mega",
	948:  "Lucario-Mega",
	949:  "Abomasnow-Mega",
	950:  "Gallade-Mega",
	951:  "Audino-Mega",
	952:  "Diancie-Mega",
	953:  "Rayquaza-Mega",
	954:  "Kyogre-Primal",
	955:  "Groudon-Primal",
	956:  "Rattata-Alola",
	957:  "Raticate-Alola",
	958:  "Raichu-Alola",
	959:  "Sandshrew-Alola",
	960:  "Sandslash-Alola",
	961:  "Vulpix-Alola",
	962:  "Ninetales-Alola",
	963:  "Diglett-Alola",
	964:  "Dugtrio-Alola",
	965:  "Meowth-Alola",
	966:  "Persian-Alola",
	967:  "Geodude-Alola",
	968:  "Graveler-Alola",
	969:  "Golem-Alola",
	970:  "Grimer-Alola",
	971:  "Muk-Alola",
	972:  "Exeggutor-Alola",
	973:  "Marowak-Alola",
	974:  "Meowth-Galar",
	975:  "Ponyta-Galar",
	976:  "Rapidash-Galar",
	977:  "Slowpoke-Galar",
	978:  "Slowbro-Galar",
	979:  "Farfetchd-Galar",
	980:  "Weezing-Galar",
	981:  "Mr-Mime-Galar",
	982:  "Articuno-Galar",
	983:  "Zapdos-Galar",
	984:  "Moltres-Galar",
	985:  "Slowking-Galar",
	986:  "Corsola-Galar",
	987:  "Zigzagoon-Galar",
	988:  "Linoone-Galar",
	989:  "Darumaka-Galar",
	990:  "Darmanitan-Galar",
	991:  "Yamask-Galar",
	992:  "Stunfisk-Galar",
	993:  "Growlithe-Hisui",
	994:  "Arcanine-Hisui",
	995:  "Voltorb-Hisui",
	996:  "Electrode-Hisui",
	997:  "Typhlosion-Hisui",
	998:  "Qwilfish-Hisui",
	999:  "Sneasel-Hisui",
	1000: "Samurott-Hisui",
	1001: "Lilligant-Hisui",
	1002: "Zorua-Hisui",
	1003: "Zoroark-Hisui",
	1004: "Braviary-Hisui",
	1005: "Sliggoo-Hisui",
	1006: "Goodra-Hisui",
	1007: "Avalugg-Hisui",
	1008: "Decidueye-Hisui",
	1009: "Pikachu-Cosplay",
	1010: "Pikachu-Rock-Star",
	1011: "Pikachu-Belle",
	1012: "Pikachu-Pop-Star",
	1013: "Pikachu-Phd",
	1014: "Pikachu-Libre",
	1015: "Pikachu-Original",
	1016: "Pikachu-Hoenn",
	1017: "Pikachu-Sinnoh",
	1018: "Pikachu-Unova",
	1019: "Pikachu-Kalos",
	1020: "Pikachu-Alola",
	1021: "Pikachu-Partner",
	1022: "Pikachu-World",
	1023: "Pichu-Spiky-Eared",
	1024: "Unown-B",
	1025: "Unown-C",
	1026: "Unown-D",
	1027: "Unown-E",
	1028: "Unown-F",
	1029: "Unown-G",
	1030: "Unown-H",
	1031: "Unown-I",
	1032: "Unown-J",
	1033: "Unown-K",
	1034: "Unown-L",
	1035: "Unown-M",
	1036: "Unown-N",
	1037: "Unown-O",
	1038: "Unown-P",
	1039: "Unown-Q",
	1040: "Unown-R",
	1041: "Unown-S",
	1042: "Unown-T",
	1043: "Unown-U",
	1044: "Unown-V",
	1045: "Unown-W",
	1046: "Unown-X",
	1047: "Unown-Y",
	1048: "Unown-Z",
	1049: "Unown-Exclamation",
	1050: "Unown-Question",
	1051: "Castform-Sunny",
	1052: "Castform-Rainy",
	1053: "Castform-Snowy",
	1054: "Deoxys-Attack",
	1055: "Deoxys-Defense",
	1056: "Deoxys-Speed",
	1057: "Burmy-Sandy",
	1058: "Burmy-Trash",
	1059: "Wormadam-Sandy",
	1060: "Wormadam-Trash",
	1061: "Cherrim-Sunshine",
	1062: "Shellos-East",
	1063: "Gastrodon-East",
	1064: "Rotom-Heat",
	1065: "Rotom-Wash",
	1066: "Rotom-Frost",
	1067: "Rotom-Fan",
	1068: "Rotom-Mow",
	1069: "Dialga-Origin",
	1070: "Palkia-Origin",
	1071: "Giratina-Origin",
	1072: "Shaymin-Sky",
	1073: "Arceus-Fighting",
	1074: "Arceus-Flying",
	1075: "Arceus-Poison",
	1076: "Arceus-Ground",
	1077: "Arceus-Rock",
	1078: "Arceus-Bug",
	1079: "Arceus-Ghost",
	1080: "Arceus-Steel",
	1081: "Arceus-Fire",
	1082: "Arceus-Water",
	1083: "Arceus-Grass",
	1084: "Arceus-Electric",
	1085: "Arceus-Psychic",
	1086: "Arceus-Ice",
	1087: "Arceus-Dragon",
	1088: "Arceus-Dark",
	1089: "Arceus-Fairy",
	1090: "Basculin-Blue-Striped",
	1091: "Basculin-White-Striped",
	1092: "Darmanitan-Zen",
	1093: "Darmanitan-Galar-Zen",
	1094: "Deerling-Summer",
	1095: "Deerling-Autumn",
	1096: "Deerling-Winter",
	1097: "Sawsbuck-Summer",
	1098: "Sawsbuck-Autumn",
	1099: "Sawsbuck-Winter",
	1100: "Tornadus-Therian",
	1101: "Thundurus-Therian",
	1102: "Landorus-Therian",
	1103: "Enamorus-Therian",
	1104: "Kyurem-White",
	1105: "Kyurem-Black",
	1106: "Keldeo-Resolute",
	1107: "Meloetta-Pirouette",
	1108: "Genesect-Douse",
	1109: "Genesect-Shock",
	1110: "Genesect-Burn",
	1111: "Genesect-Chill",
	1112: "Greninja-Bond",
	1113: "Greninja-Ash",
	1114: "Vivillon-Polar",
	1115: "Vivillon-Tundra",
	1116: "Vivillon-Continental",
	1117: "Vivillon-Garden",
	1118: "Vivillon-Elegant",
	1119: "Vivillon-Meadow",
	1120: "Vivillon-Modern",
	1121: "Vivillon-Marine",
	1122: "Vivillon-Archipelago",
	1123: "Vivillon-High-Plains",
	1124: "Vivillon-Sandstorm",
	1125: "Vivillon-River",
	1126: "Vivillon-Monsoon",
	1127: "Vivillon-Savanna",
	1128: "Vivillon-Sun",
	1129: "Vivillon-Ocean",
	1130: "Vivillon-Jungle",
	1131: "Vivillon-Fancy",
	1132: "Vivillon-Pokeball",
	1133: "Flabebe-Yellow",
	1134: "Flabebe-Orange",
	1135: "Flabebe-Blue",
	1136: "Flabebe-White",
	1137: "Floette-Yellow",
	1138: "Floette-Orange",
	1139: "Floette-Blue",
	1140: "Floette-White",
	1141: "Floette-Eternal",
	1142: "Florges-Yellow",
	1143: "Florges-Orange",
	1144: "Florges-Blue",
	1145: "Florges-White",
	1146: "Furfrou-Heart",
	1147: "Furfrou-Star",
	1148: "Furfrou-Diamond",
	1149: "Furfrou-Debutante",
	1150: "Furfrou-Matron",
	1151: "Furfrou-Dandy",
	1152: "Furfrou-La-Reine",
	1153: "Furfrou-Kabuki",
	1154: "Furfrou-Pharaoh",
	1155: "Meowstic-Female",
	1156: "Aegislash-Blade",
	1157: "Pumpkaboo-Small",
	1158: "Pumpkaboo-Large",
	1159: "Pumpkaboo-Super",
	1160: "Gourgeist-Small",
	1161: "Gourgeist-Large",
	1162: "Gourgeist-Super",
	1163: "Xerneas-Active",
	1164: "Zygarde-10%",
	// power construct
	1165: "Zygarde-10%",
	1166: "Zygarde-50%",
	1167: "Zygarde-Complete",
	1168: "Hoopa-Unbound",
	1169: "Oricorio-Pom-Pom",
	1170: "Oricorio-Pau",
	1171: "Oricorio-Sensu",
	1172: "Rockruff-Own-Tempo",
	1173: "Lycanroc-Midnight",
	1174: "Lycanroc-Dusk",
	1175: "Wishiwashi-School",
	1176: "Silvally-Fighting",
	1177: "Silvally-Flying",
	1178: "Silvally-Poison",
	1179: "Silvally-Ground",
	1180: "Silvally-Rock",
	1181: "Silvally-Bug",
	1182: "Silvally-Ghost",
	1183: "Silvally-Steel",
	1184: "Silvally-Fire",
	1185: "Silvally-Water",
	1186: "Silvally-Grass",
	1187: "Silvally-Electric",
	1188: "Silvally-Psychic",
	1189: "Silvally-Ice",
	1190: "Silvally-Dragon",
	1191: "Silvally-Dark",
	1192: "Silvally-Fairy",
	1193: "Minior-Orange",
	1194: "Minior-Yellow",
	1195: "Minior-Green",
	1196: "Minior-Blue",
	1197: "Minior-Indigo",
	1198: "Minior-Violet",
	1199: "Minior-Core",
	1200: "Minior-Core-Orange",
	1201: "Minior-Core-Yellow",
	1202: "Minior-Core-Green",
	1203: "Minior-Core-Blue",
	1204: "Minior-Core-Indigo",
	1205: "Minior-Core-Violet",
	1206: "Mimikyu-Busted",
	1207: "Necrozma-Dusk-Mane",
	1208: "Necrozma-Dawn-Wings",
	1209: "Necrozma-Ultra",
	1210: "Magearna-Original",
	1211: "Cramorant-Gulping",
	1212: "Cramorant-Gorging",
	1213: "Toxtricity-Low-Key",
	1214: "Sinistea-Antique",
	1215: "Polteageist-Antique",
	1216: "Alcremie-Ruby-Cream",
	1217: "Alcremie-Matcha-Cream",
	1218: "Alcremie-Mint-Cream",
	1219: "Alcremie-Lemon-Cream",
	1220: "Alcremie-Salted-Cream",
	1221: "Alcremie-Ruby-Swirl",
	1222: "Alcremie-Caramel-Swirl",
	1223: "Alcremie-Rainbow-Swirl",
	1224: "Eiscue-Noice",
	1225: "Indeedee-Female",
	1226: "Morpeko-Hangry",
	1227: "Zacian-Crowned",
	1228: "Zamazenta-Crowned",
	1229: "Eternatus-Eternamax",
	1230: "Urshifu-Rapid-Strike",
	1231: "Zarude-Dada",
	1232: "Calyrex-Ice",
	1233: "Calyrex-Shadow",
	1234: "Basculegion-Female",
	1235: "Alcremie-Berry",
	1236: "Alcremie-Berry-Ruby-Cream",
	1237: "Alcremie-Berry-Matcha-Cream",
	1238: "Alcremie-Berry-Mint-Cream",
	1239: "Alcremie-Berry-Lemon-Cream",
	1240: "Alcremie-Berry-Salted-Cream",
	1241: "Alcremie-Berry-Ruby-Swirl",
	1242: "Alcremie-Berry-Caramel-Swirl",
	1243: "Alcremie-Berry-Rainbow-Swirl",
	1244: "Alcremie-Love",
	1245: "Alcremie-Love-Ruby-Cream",
	1246: "Alcremie-Love-Matcha-Cream",
	1247: "Alcremie-Love-Mint-Cream",
	1248: "Alcremie-Love-Lemon-Cream",
	1249: "Alcremie-Love-Salted-Cream",
	1250: "Alcremie-Love-Ruby-Swirl",
	1251: "Alcremie-Love-Caramel-Swirl",
	1252: "Alcremie-Love-Rainbow-Swirl",
	1253: "Alcremie-Star",
	1254: "Alcremie-Star-Ruby-Cream",
	1255: "Alcremie-Star-Matcha-Cream",
	1256: "Alcremie-Star-Mint-Cream",
	1257: "Alcremie-Star-Lemon-Cream",
	1258: "Alcremie-Star-Salted-Cream",
	1259: "Alcremie-Star-Ruby-Swirl",
	1260: "Alcremie-Star-Caramel-Swirl",
	1261: "Alcremie-Star-Rainbow-Swirl",
	1262: "Alcremie-Clover",
	1263: "Alcremie-Clover-Ruby-Cream",
	1264: "Alcremie-Clover-Matcha-Cream",
	1265: "Alcremie-Clover-Mint-Cream",
	1266: "Alcremie-Clover-Lemon-Cream",
	1267: "Alcremie-Clover-Salted-Cream",
	1268: "Alcremie-Clover-Ruby-Swirl",
	1269: "Alcremie-Clover-Caramel-Swirl",
	1270: "Alcremie-Clover-Rainbow-Swirl",
	1271: "Alcremie-Flower",
	1272: "Alcremie-Flower-Ruby-Cream",
	1273: "Alcremie-Flower-Matcha-Cream",
	1274: "Alcremie-Flower-Mint-Cream",
	1275: "Alcremie-Flower-Lemon-Cream",
	1276: "Alcremie-Flower-Salted-Cream",
	1277: "Alcremie-Flower-Ruby-Swirl",
	1278: "Alcremie-Flower-Caramel-Swirl",
	1279: "Alcremie-Flower-Rainbow-Swirl",
	1280: "Alcremie-Ribbon",
	1281: "Alcremie-Ribbon-Ruby-Cream",
	1282: "Alcremie-Ribbon-Matcha-Cream",
	1283: "Alcremie-Ribbon-Mint-Cream",
	1284: "Alcremie-Ribbon-Lemon-Cream",
	1285: "Alcremie-Ribbon-Salted-Cream",
	1286: "Alcremie-Ribbon-Ruby-Swirl",
	1287: "Alcremie-Ribbon-Caramel-Swirl",
	1288: "Alcremie-Ribbon-Rainbow-Swirl",
	1289: "Sprigatito",
	1290: "Floragato",
	1291: "Meowscarada",
	1292: "Fuecoco",
	1293: "Crocalor",
	1294: "Skeledirge",
	1295: "Quaxly",
	1296: "Quaxwell",
	1297: "Quaquaval",
	1298: "Lechonk",
	1299: "Oinkologne",
	1300: "Oinkologne-Female",
	1301: "Tarountula",
	1302: "Spidops",
	1303: "Nymble",
	1304: "Lokix",
	1305: "Pawmi",
	1306: "Pawmo",
	1307: "Pawmot",
	1308: "Tandemaus",
	1309: "Maushold",
	1310: "Maushold-Four",
	1311: "Fidough",
	1312: "Dachsbun",
	1313: "Smoliv",
	1314: "Dolliv",
	1315: "Arboliva",
	1316: "Squawkabilly",
	1317: "Squawkabilly-Blue",
	1318: "Squawkabilly-Yellow",
	1319: "Squawkabilly-White",
	1320: "Nacli",
	1321: "Naclstack",
	1322: "Garganacl",
	1323: "Charcadet",
	1324: "Armarouge",
	1325: "Ceruledge",
	1326: "Tadbulb",
	1327: "Bellibolt",
	1328: "Wattrel",
	1329: "Kilowattrel",
	1330: "Maschiff",
	1331: "Mabosstiff",
	1332: "Shroodle",
	1333: "Grafaiai",
	1334: "Bramblin",
	1335: "Brambleghast",
	1336: "Toedscool",
	1337: "Toedscruel",
	1338: "Klawf",
	1339: "Capsakid",
	1340: "Scovillain",
	1341: "Rellor",
	1342: "Rabsca",
	1343: "Flittle",
	1344: "Espathra",
	1345: "Tinkatink",
	1346: "Tinkatuff",
	1347: "Tinkaton",
	1348: "Wiglett",
	1349: "Wugtrio",
	1350: "Bombirdier",
	1351: "Finizen",
	1352: "Palafin",
	1353: "Palafin-Hero",
	1354: "Varoom",
	1355: "Revavroom",
	1356: "Cyclizar",
	1357: "Orthworm",
	1358: "Glimmet",
	1359: "Glimmora",
	1360: "Greavard",
	1361: "Houndstone",
	1362: "Flamigo",
	1363: "Cetoddle",
	1364: "Cetitan",
	1365: "Veluza",
	1366: "Dondozo",
	1367: "Tatsugiri",
	1368: "Tatsugiri-Droopy",
	1369: "Tatsugiri-Stretchy",
	1370: "Annihilape",
	1371: "Clodsire",
	1372: "Farigiraf",
	1373: "Dudunsparce",
	1374: "Dudunsparce-Three-Segment",
	1375: "Kingambit",
	1376: "Great Tusk",
	1377: "Scream Tail",
	1378: "Brute Bonnet",
	1379: "Flutter Mane",
	1380: "Slither Wing",
	1381: "Sandy Shocks",
	1382: "Iron Treads",
	1383: "Iron Bundle",
	1384: "Iron Hands",
	1385: "Iron Jugulis",
	1386: "Iron Moth",
	1387: "Iron Thorns",
	1388: "Frigibax",
	1389: "Arctibax",
	1390: "Baxcalibur",
	1391: "Gimmighoul",
	1392: "Gimmighoul-Roaming",
	1393: "Gholdengo",
	1394: "Wo-Chien",
	1395: "Chien-Pao",
	1396: "Ting-Lu",
	1397: "Chi-Yu",
	1398: "Roaring Moon",
	1399: "Iron Valiant",
	1400: "Koraidon",
	1401: "Miraidon",
	1402: "Tauros-Paldea-Combat",
	1403: "Tauros-Paldea-Blaze",
	1404: "Tauros-Paldea-Aqua",
	1405: "Wooper-Paldea",
	1406: "Walking Wake",
	1407: "Iron Leaves",
	1408: "Dipplin",
	1409: "Poltchageist",
	1410: "Poltchageist-Artisan",
	1411: "Sinistcha",
	1412: "Sinistcha-Masterpiece",
	1413: "Okidogi",
	1414: "Munkidori",
	1415: "Fezandipiti",
	1416: "Ogerpon",
	1417: "Ogerpon-Wellspring",
	1418: "Ogerpon-Hearthflame",
	1419: "Ogerpon-Cornerstone",
	1420: "Ogerpon-Teal-Tera",
	1421: "Ogerpon-Wellspring-Tera",
	1422: "Ogerpon-Hearthflame-Tera",
	1423: "Ogerpon-Cornerstone-Tera",
	1424: "Ursaluna-Bloodmoon",
	1425: "Archaludon",
	1426: "Hydrapple",
	1427: "Gouging Fire",
	1428: "Raging Bolt",
	1429: "Iron Boulder",
	1430: "Iron Crown",
	1431: "Terapagos",
	1432: "Terapagos-Terastal",
	1433: "Terapagos-Stellar",
	1434: "Pecharunt",
	1435: "Lugia-Shadow",
	// really stupid extra forms for mons that shouldn't have them but they decided to add them anyway because stupid
	1436: "Mothim",
	1437: "Mothim",
	1438: "Scatterbug",
	1439: "Scatterbug",
	1440: "Scatterbug",
	1441: "Scatterbug",
	1442: "Scatterbug",
	1443: "Scatterbug",
	1444: "Scatterbug",
	1445: "Scatterbug",
	1446: "Scatterbug",
	1447: "Scatterbug",
	1448: "Scatterbug",
	1449: "Scatterbug",
	1450: "Scatterbug",
	1451: "Scatterbug",
	1452: "Scatterbug",
	1453: "Scatterbug",
	1454: "Scatterbug",
	1455: "Scatterbug",
	1456: "Scatterbug",
	1457: "Spewpa",
	1458: "Spewpa",
	1459: "Spewpa",
	1460: "Spewpa",
	1461: "Spewpa",
	1462: "Spewpa",
	1463: "Spewpa",
	1464: "Spewpa",
	1465: "Spewpa",
	1466: "Spewpa",
	1467: "Spewpa",
	1468: "Spewpa",
	1469: "Spewpa",
	1470: "Spewpa",
	1471: "Spewpa",
	1472: "Spewpa",
	1473: "Spewpa",
	1474: "Spewpa",
	1475: "Spewpa",
	1476: "Raticate-Alola-Totem",
	1477: "Gumshoos-Totem",
	1478: "Vikavolt-Totem",
	1479: "Lurantis-Totem",
	1480: "Salazzle-Totem",
	1481: "Mimikyu-Totem",
	1482: "Kommo-O_Totem",
	1483: "Marowak-Alola-Totem",
	1484: "Ribombee-Totem",
	1485: "Araquanid-Totem",
	1486: "Togedemaru-Totem",
	1487: "Pikachu-Starter",
	1488: "Eevee-Starter",
	1489: "Venusaur-Gmax",
	1490: "Blastoise-Gmax",
	1491: "Charizard-Gmax",
	1492: "Butterfree-Gmax",
	1493: "Pikachu-Gmax",
	1494: "Meowth-Gmax",
	1495: "Machamp-Gmax",
	1496: "Gengar-Gmax",
	1497: "Kingler-Gmax",
	1498: "Lapras-Gmax",
	1499: "Eevee-Gmax",
	1500: "Snorlax-Gmax",
	1501: "Garbodor-Gmax",
	1502: "Melmetal-Gmax",
	1503: "Rillaboom-Gmax",
	1504: "Cinderace-Gmax",
	1505: "Inteleon-Gmax",
	1506: "Corviknight-Gmax",
	1507: "Orbeetle-Gmax",
	1508: "Drednaw-Gmax",
	1509: "Coalossal-Gmax",
	1510: "Flapple-Gmax",
	1511: "Appletun-Gmax",
	1512: "Sandaconda-Gmax",
	1513: "Toxtricity-Amped-Gmax",
	1514: "Toxtricity-Low-Key-Gmax",
	1515: "Centiskorch-Gmax",
	1516: "Hatterene-Gmax",
	1517: "Grimmsnarl-Gmax",
	1518: "Alcremie-Gmax",
	1519: "Copperajah-Gmax",
	1520: "Duraludon-Gmax",
	1521: "Urshifu-Gmax",
	1522: "Urshifu-Rapid-Strike-Gmax",
	1523: "Mimikyu-Totem-Busted",
}
