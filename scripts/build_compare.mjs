import { writeFileSync, readFileSync } from 'node:fs';

// Parse command-line arguments
const args = process.argv.slice(2);

const config = {
  newPath: '',
  oldPath: '',
  objectId: 'id', // Default value for objectId
  outputs: [],
};

args.forEach((arg) => {
  const [key, value] = arg.split('=');
  switch (key) {
    case '-newPath':
    case '-new':
      config.newPath = value;
      break;
    case '-oldPath':
    case '-old':
      config.oldPath = value;
      break;
    case '-objectId':
    case '-id':
      config.objectId = value || 'id';
      break;
    case '-outputs':
    case '-o':
      config.outputs = value.split(',');
      break;
    default:
      console.error(`Unknown argument: ${key}`);
      process.exit(1);
  }
});

// Validate required arguments
if (!config.newPath || !config.oldPath || config.outputs.length === 0) {
  console.error("Missing required arguments: -new/-newPath, -old/-oldPath, and -o/-outputs are required.");
  process.exit(1);
}

// Read the JSON files
const newArr = JSON.parse(readFileSync(config.newPath, 'utf8'));
const oldArr = JSON.parse(readFileSync(config.oldPath, 'utf8'));

// Modify the 'our' JSON data
newArr.forEach(obj => {
  const oldObj = oldArr.find(m => m[config.objectId] === obj[config.objectId]);
  obj.old = oldObj ?? null;
});

// Write to the output files
config.outputs.forEach(outputPath => {
  writeFileSync(outputPath, JSON.stringify(newArr, null, 2));
});

console.log('Comparison completed and files written to:', config.outputs.join(', '));
