import { writeFileSync, readFileSync } from 'node:fs';

const rhh = JSON.parse(readFileSync('./build/species_rhh.json', 'utf8'));
const our = JSON.parse(readFileSync('./build/species_our.json', 'utf8'));

our.forEach(mon => {
  const rhhMon = rhh.find(m => m.id === mon.id);
  mon.old = rhhMon;
});

writeFileSync('./build/pokemon.json', JSON.stringify(our, null, 2));
writeFileSync('./../docs/src/data/pokemon.json', JSON.stringify(our, null, 2));