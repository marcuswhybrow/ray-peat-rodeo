/** 
 * @param {Context} _context
 * @returns {import("marked").TokenizerAndRendererExtension} 
 */
export default (_context) => ({
  name: "mention",
  level: "inline",

  start(src) {
    return src.match(/\[\[/)?.index;
  },

  tokenizer(src, _tokens) {
    const rule = /^\[\[(.*?)(\|.*?)?]\]/;
    const match = rule.exec(src);
    if (match) {
      const signature = match[1].trim();
      const label = match[2]?.substring(1) || signature;
      const token = {
        type: "mention",
        raw: match[0],
        signature,
        tokens: this.lexer.inlineTokens(label.trim()),
      }
      return token;
    }
  },

  renderer(token) {
    return this.parser.parseInline(token.tokens || []);
  },
});
