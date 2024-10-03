/** 
 * @param {Context} _context
 * @returns {import("marked").TokenizerAndRendererExtension} 
 */
export default (_context) => ({
  name: "timecode",
  level: "inline",

  start: src => src.match(/\[[0-9]+/)?.index,

  tokenizer(src, _tokens) {
    const rule = /^\[([0-9]+\:)?([0-9]{1,2})\:([0-9]{1,2})\]/;
    const match = rule.exec(src);
    if (match) {
      const hours = match[1]?.substring(0, match[1].length - 1)?.padStart(2, "0") || "00";
      const minutes = match[2].padStart(2, "0");
      const seconds = match[3].padStart(2, "0");
      return {
        type: "timecode",
        raw: match[0],
        timecode: `${hours}:${minutes}:${seconds}`,
      };
    }
  },

  renderer({ timecode }) {
    return `
      <rpr-timecode 
        time="${timecode}" 
      ></rpr-timecode>
    `;
  },
});
