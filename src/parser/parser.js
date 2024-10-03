import fs from "fs";
import { Marked } from "marked";
import timecodes from "./timecodes.js";
import sections from "./sections.js";
import contributors from "./contributors.js";
import sidenotes from "./sidenotes.js";
import mentions from "./mentions.js";
import { Fetcher } from "./fetcher.js";
import yaml from "js-yaml";

/**
  * @param {Promise<Asset>[]} assetPromiseList 
  * @returns {Promise<Asset>[]}
  */
export function parse(assetPromiseList) {
  const cacheYaml = fs.readFileSync("assets/data/cache.yml", "utf8");
  const cacheData = /** @type {import("./fetcher.js").FetcherCache} */ (yaml.load(cacheYaml, {}));
  const fetcher = new Fetcher(cacheData);
  const avatars = fs.readdirSync("./src/public/avatars");

  return assetPromiseList.map(async inputAssetPromise => {
    const asset = await inputAssetPromise;
    if (asset.html || !asset.markdown) return asset;

    /** @type {Context} */
    const context = { fetcher, asset, avatars };

    const marked = new Marked({
      extensions: [
        timecodes(context),
        mentions(context),
        sidenotes(context),
        contributors(context),
        sections(context),
      ],
      async: true,
      walkTokens: async token => {
        const asyncToken = /** @type {AsyncToken} */ (token);
        if (asyncToken.resolveAsync) await asyncToken.resolveAsync(token, context);
      }
    });

    asset.html = await marked.parse(asset.markdown);

    return asset;
  });
}

export default { parse };

