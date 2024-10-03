import * as nodePath from "path";
import fs from "fs";
import fm from "front-matter";
import { parse as parseHTML } from "node-html-parser";

/**
 * @param {string} assetFilename
 * @returns {{date:string, slug:string, ext:string}}
 */
export function parseAssetFilename(assetFilename) {
  const path = nodePath.parse(assetFilename);
  const date = path.name.substring(0, "0000-00-00".length);
  const slug = path.name.substring("0000-00-00-".length);
  const ext = path.ext;
  return { date, slug, ext };
}

/**
 * @param {string} assetFilename
 * @returns {Asset}
 */
export function createAssetStub(assetFilename) {
  const content = fs.readFileSync(assetFilename, "utf8");
  const { date, slug } = parseAssetFilename(assetFilename);

  const {
    attributes: frontMatter,
    body: markdown
  } = /** @type{import("front-matter").FrontMatterResult<FrontMatter>} */ (fm(content));

  return {
    date, slug, markdown, frontMatter,
    filename: assetFilename,
    html: "",
    issues: [],
    sections: [],
    errors: [],
    contributors: [],
  };
}

/**
  * Dummy template literal for IDE HTML syntax highlighting.
  *
  * @param {TemplateStringsArray} strings
  * @param {string[]} values
  */
export function html(strings, ...values) {
  let str = "";
  strings.forEach((string, i) => {
    str += string + (values[i] || "");
  });
  return str;
}

/**
 * Dummy template literal for IDE JavaScript syntax highlighting.
 */
export const js = html;

/**
  * @param {string[]} columnList
  * @returns {string}
  */
export function tableRow(...columnList) {
  return html`<tr>${columnList.map(col => html`<td>${col}</td>`).join("")}</tr>`
}

/** 
 * Parses data-pagefind-filter attributes to exctract all Pagefind filters.
 *
 * Pagefind's node wrapper lib doesn't say which filters it discovered, forcing
 * filter lookup in the client lib, slowing time to first render of filters.
 * This functions gets around that by parsing the HTML again ourselves looking
 * for the same pagefind HTML element attributes which pagefind itself does to 
 * reconstruct the same data that [pagefind.filters()] would return.
 *
 * This is an upcomming feature of Pagefind, so this approach will soon be 
 * obsolete. See reference issues. This implementation is a best guess effort
 * following the Pagefind docs, it may not perfectly match edge cases in filter 
 * names or values.
 *
 * # Reference
 *
 * - https://pagefind.app/docs/filtering/
 * - https://github.com/CloudCannon/pagefind/issues/715
 * - https://github.com/CloudCannon/pagefind/issues/371
 *
 * # Example
 * ```js 
 * import assert from "assert";
 * assert.deepEqual(extractPagefindFilters(`
 *   <span data-pagefind-filter="singleName:inlineContent"></span>
 *   <span data-pagefind-filter="singleName">valueContent</span>
 *   <span data-pagefind-filter="name1, name2:inlineContent">valueContent</span>
 *   <span data-pagefind-filter="name1, name2[data-name], name3:inlineContent" data-name="attrValue">valueContent</span>
 * `), {
 *   singleName: { inlineContent: 1, valueContent: 1 },
 *   name1: { valueContent: 2 },
 *   name2: { inlineContent: 1, attrValue: 1 },
 *   name3: { inlineContent: 1 }
 * }
 * ```
 *
 * @param {string} html 
 * @returns {PagefindFilters}
 */
export function extractPagefindFilters(html) {
  /** @type {PagefindFilters} */
  const pagefindFilters = {};

  parseHTML(html).querySelectorAll("[data-pagefind-filter]").forEach(element => {
    let signature = element.getAttribute("data-pagefind-filter");
    if (!signature) return;

    let filters = [];

    let chars = signature.split("");
    let name = "";
    let mod = "";

    chars.forEach(char => {
      switch (char) {
        case ',':
          if (mod[0] === ":") mod += char;
          else {
            filters.push([name, mod]);
            name = ""; mod = "";
          }
          break;
        case '[':
        case ':':
          mod += char;
          break;
        case ']':
        default:
          if (mod) mod += char;
          else name += char;
      }
    });

    if (name || mod) filters.push([name, mod]);

    filters = filters.map(([name, mod]) => {
      name = name.trim();
      mod = mod.trim();
      if (mod[0] === ":") {
        return [name, mod.substring(1).trim()];
      } else if (mod[0] === "[") {
        return [name, element.getAttribute(mod.substring(1, mod.length - 1))?.trim() || ""];
      } else {
        return [name, element.textContent?.trim() || ""];
      }
    });

    filters.forEach(([name, value]) => {
      if (!pagefindFilters.hasOwnProperty(name)) pagefindFilters[name] = {};
      if (!pagefindFilters[name].hasOwnProperty(value))
        pagefindFilters[name][value] = 1;
      else pagefindFilters[name][value]++;
    });
  });

  return pagefindFilters;
}
