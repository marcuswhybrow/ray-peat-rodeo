import { extractPagefindFilters } from "./utils.js";
import fs from "fs";


const html = fs.readFileSync(process.argv[2], "utf8");
console.log(extractPagefindFilters(html));

console.log(extractPagefindFilters(`
  <span data-pagefind-filter="singleName:inlineContent"></span>
  <span data-pagefind-filter="singleName">valueContent</span>
  <span data-pagefind-filter="name1, name2:inlineContent">valueContent</span>
  <span data-pagefind-filter="name1, name2[data-name], name3:inlineContent" data-name="attrValue">valueContent</span>
`));
