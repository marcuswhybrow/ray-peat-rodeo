import { walkTokens } from "marked";
import slugify from "slugify";
import he from "he";
import { findLastSection, findSectionParent } from "./utils.js";

const prefixRegex = /^(([0-9]+.)+) /;
const slugPurgeRegex = /['"()!:@.~+*`]/g;

/** 
 * @param {Context} _context
 * @returns {import("marked").RendererExtension} 
 */
export default ({ asset }) => ({
  name: "heading",
  renderer: ({ depth, raw, tokens }) => {
    const tag = `h${depth}`;

    if (depth === 1) {
      asset.errors.push(`Level one heading "${raw}" found. The level one heading is reserved for the page title. Consider using heading levels 2 through 6.`);
      return;
    }

    let [title, timecode] = ["", ""];
    walkTokens(tokens || [], t => {
      if (t.type === "text") {
        title += t.text;
      } else if (t.type === "timecode") {
        timecode = t.timecode;
      }
    });
    title = title.trim();

    const id = slugify(he.decode(title), {
      lower: true,
      trim: true,
      remove: slugPurgeRegex,
    });

    let prefix = "";
    const prefixMatch = title.match(prefixRegex);

    if (prefixMatch) {
      prefix = prefixMatch[1];
      title = title.substring(prefixMatch[0].length);
    }

    const sectionParent = findSectionParent(asset, depth);

    if (!sectionParent) {
      const lastSection = findLastSection(asset);
      if (!lastSection) {
        asset.errors.push(`The first heading "${raw}" is a level ${depth} heading, but must be a level 2 heading.`);
      } else {
        asset.errors.push(`Level ${depth} heading "${raw}" follows a level ${lastSection.depth} heading. Must be a level ${lastSection.depth - 1} heading or less.`);
      }
      return;
    }

    sectionParent.push({
      depth, title, id, prefix, timecode,
      issues: [],
      subsections: [],
      excerpt: null,
    });

    return `
      <rpr-section 
        level="${depth}"
        title="${title}"
        timecode="${timecode}"
        prefix="${prefix}"
        section-id="${id}"
      ><${tag} id="${id}">${title}</${tag}></rpr-section>
    `;
  },
});
