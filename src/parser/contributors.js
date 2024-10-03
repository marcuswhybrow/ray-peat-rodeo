import path from "path";

/** 
  * @param {Context} context 
 *  @returns {import("marked").TokenizerAndRendererExtension} 
 */
export default ({ asset: parsed, avatars }) => ({
  name: "contribution",
  level: "block",

  start: src => src.match(/^[a-zA-Z0-9]+: /)?.index,

  tokenizer(src, _tokens) {
    const rule = /^([a-zA-Z0-9]+): /;
    const match = rule.exec(src);
    if (match) {
      const initials = match[1];
      const token = {
        type: "contribution",
        raw: match[0],
        initials,
      };

      const contributors = parsed.frontMatter.speakers;

      const contributorMapExists = typeof contributors !== "undefined";
      const speakerExists = contributors.hasOwnProperty(token.initials);

      if (!contributorMapExists || !speakerExists) {
        parsed.errors.push(`Cannot find "speakers.${token.initials}" in front matter.`);
        return;
      }

      let name = contributors[token.initials];
      token.filterable = true;
      if (name.startsWith("-")) {
        name = name.substring(1);
        token.filterable = false;
      }
      token.name = name;

      const slug = name.toLowerCase().replaceAll(" ", "-");

      const avatar = avatars.find(avatar => path.parse(avatar).name === slug);
      token.avatar = avatar ? `/avatars/${avatar}` : "";

      if (!parsed.contributors.some(contributor => contributor.name === name)) {
        parsed.contributors.push({
          name,
          filterable: token.filterable,
          initials: token.initials,
          avatar: token.avatar,
        });
      }
      return token;
    }
  },

  renderer({ name, initials, avatar, filterable }) {
    return `
      <rpr-contribution
        initials="${initials}" 
        name="${name}"
        ${avatar ? `avatar="${avatar}"` : ""} 
        ${filterable ? `filterable` : ""}
      ></rpr-contribution>
    `;
  },
});
