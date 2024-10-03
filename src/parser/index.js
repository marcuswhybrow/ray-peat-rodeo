import fs from "fs";
import path from "path";
import { parse } from "./parser.js";
import { createAssetStub } from "../utils.js";

const OUTPUT_EXT = ".json";

const inputs = process.argv.slice(2);
let toStdout = false;

if (inputs.length === 0 || inputs[0] === "--help") {
  console.log(`Usage: parser file.md
Usage: parser file1.md file2.md file3.md ./out/directory
Usage: parser file*.md ./out/directory`);
  process.exit(1);
}

/** @type {string[]} */
const inputFilenameList = [];

/** @type {string[]} */
const outputFilenameList = [];

if (inputs.length >= 1) {
  if (!fs.lstatSync(inputs[0]).isFile()) {
    console.error("First argument must be a file.");
    process.exit(1);
  }
}

if (inputs.length === 1) {
  inputFilenameList.push(inputs[0]);
  toStdout = true;
} else {
  const outputDirname = inputs[inputs.length - 1];
  inputs
    .slice(0, inputs.length - 1)
    .forEach(input => {
      inputFilenameList.push(input);
      const outputFilename = path.join(outputDirname, `${path.parse(input).name}${OUTPUT_EXT}`);
      outputFilenameList.push(outputFilename);
    });
}

const start = process.hrtime();

/** @type {Promise<Asset>[]} */
const inputAssetList = inputFilenameList.map(async inputFilename => createAssetStub(inputFilename));

const assetList = parse(inputAssetList);

const nullresults = assetList.map((assetPromise, i) => new Promise(async resolve => {
  const asset = await assetPromise;
  const json = JSON.stringify(asset);

  if (toStdout) {
    console.log(json);
    return resolve(null);
  }

  const outputFilename = outputFilenameList[i];
  const parentDir = path.parse(outputFilename).dir;
  fs.mkdirSync(parentDir, { recursive: true });
  fs.writeFileSync(outputFilename, json);

  console.log(outputFilename)
  resolve(null);
}));

await Promise.all(nullresults);
if (!toStdout) {
  const diff = process.hrtime(start);
  const ms = Math.ceil(diff[1] * 1e-6);
  console.log(`Completed in ${ms}ms.`);
}
