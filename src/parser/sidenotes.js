import { JSDOM } from "jsdom";
import he from "he";
import { lastSection } from "./utils.js";

/**
 * @param {Context} context 
 * @returns {import("marked").TokenizerAndRendererExtension} 
 */
export default ({ fetcher, asset }) => ({
  name: "sidenote",
  level: "inline",

  start(src) {
    return src.match(/\{/)?.index;
  },

  tokenizer(src, _tokens) {
    const rule = /^\{(.*?)\}/;
    const match = rule.exec(src);
    if (match) {
      const text = match[1].trim();
      const issue = match[1].trim().match(/^\#([0-9]+)$/);

      if (issue) {
        const id = issue[1];
        const url = `https://github.com/marcuswhybrow/ray-peat-rodeo/issues/${id}`;

        /** @type {AsyncToken} */
        const asyncToken = {
          type: "sidenote",
          raw: match[0],
          issueId: id,
          issueTitle: fetcher.fetch(url, "title", async response => {
            const dom = new JSDOM(await response.text());
            const title = dom.window.document.querySelector("h1 bdi")?.textContent;
            if (!title) asset.errors.push(`Failed to find title for issue "${id}".`);
            return title || "";
          }),
          resolveAsync: async token => {
            token.issueTitle = await token.issueTitle;
          }
        };
        return asyncToken;
      } else {
        return {
          type: "sidenote",
          raw: match[0],
          tokens: this.lexer.inlineTokens(text),
        };
      }
    }
  },

  renderer(token) {
    if (token.issueId) {

      const section = asset.sections.length === 0
        ? null
        : lastSection(asset.sections[asset.sections.length - 1]);

      /** @type {Issue} */
      const issue = {
        id: parseInt(token.issueId),
        title: token.issueTitle
      };

      if (section) {
        section.issues.push(issue);
      } else {
        asset.issues.push(issue);
      }

      return `<rpr-issue id="issue-${token.issueId}" issueid="${token.issueId}" issuetitle="${he.escape(token.issueTitle)}"></rpr-issue>`;
    } else {
      return `
        <rpr-sidenote>${this.parser.parseInline(token.tokens || [])}</rpr-sidenote>
      `;
    }
  },
});

