/** @param {Context} context */
export default function(context) {

  /** @param {import("marked").Token} token */
  return token => {
    if (context.asset.slug === "tribute-to-dr-raymond-peat") {
      console.log(token.type, token);
    }
  };
}
