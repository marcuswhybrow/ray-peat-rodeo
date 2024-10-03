import fs from "fs";
import path from "path";
import * as pf from "pagefind";
import assert from "assert";
import parser from "./parser/parser.js";
import { createAssetStub, extractPagefindFilters, html, js, parseAssetFilename, tableRow } from "./utils.js";
import he from "he";
import childProcess from "child_process";

const start = performance.now();
const since = () => Math.ceil(performance.now() - start);

/** 
 * Buffer filesystem changes and execute them all at the end for less dev 
 * server reloads.
 *
 * @type {(() => Promise<void>|any)[]} 
 */
const writeBuffer = [];

const OUT = process.argv[2] || "build";
const CLEAN = process.argv[3] === "clean";
const BUILD_CACHE = `${OUT}/.build-cache.json`;

/** @param {string} slug */
const outJsonFilename = slug => path.join(OUT, `${slug}.json`);

/** @param {string} slug */
const outMarkdownFilename = slug => path.join(OUT, `${slug}.md`);

/** @param {string} slug */
const outHtmlFilename = slug => path.join(OUT, slug, "index.html");

/** @param {string} slug */
const outPartialFilename = slug => path.join(OUT, slug, "partial.html");

writeBuffer.push(() => {
  if (fs.existsSync(OUT)) {
    if (CLEAN) {
      console.log(`Clean build requested: removing existing build directory "${OUT}"`);
      fs.rmSync(OUT, { recursive: true });
    }
    console.log(`Using existing build directory "${OUT}"`);
  } else {
    if (CLEAN) {
      console.log(`Clean build unnecessary: build directory "${OUT}" doesn't yet exist.`);
    }
    console.log(`Creating build directory "${OUT}"`);
    fs.mkdirSync(OUT, { recursive: true });
  }
  console.log(childProcess.execSync(`copy-static ${OUT}`).toString());
});

const assetFilenameList = fs.readdirSync("assets", { withFileTypes: true })
  .filter(file => file.isFile())
  .map(file => path.join(file.parentPath, file.name))
  .reverse();


const dirtyAssetIndexes = [];
const inputAssetPromiseList = assetFilenameList.map(async (assetFilename, index) => {
  const { mtime: inputModified } = fs.statSync(assetFilename);

  const { slug } = parseAssetFilename(assetFilename);
  const cached = outJsonFilename(slug);
  if (fs.existsSync(cached)) {
    const { mtime: cacheModified } = fs.statSync(cached);
    if (cacheModified >= inputModified) {
      const cachedJson = fs.readFileSync(cached, "utf8");
      return /** @type {Asset} */ (JSON.parse(cachedJson));
    }
  }

  dirtyAssetIndexes.push(index);
  return createAssetStub(assetFilename);
});

/** @type {Promise<Asset>[]} */
const assetPromiseList = parser.parse(inputAssetPromiseList);

const assetJsonDone = assetPromiseList.map(async (assetPromise, index) => {
  if (!dirtyAssetIndexes.includes(index)) return;

  const asset = await assetPromise;
  const out = path.join(OUT, `${asset.slug}.json`);
  writeBuffer.push(() => fs.writeFileSync(out, JSON.stringify(asset)));
});

Promise.all(assetJsonDone).then(() => {
  const all = assetFilenameList.length;
  const dirty = dirtyAssetIndexes.length;
  const cached = all - dirty;
  const outJson = outJsonFilename("ASSET-SLUG");
  console.log(`Found ${all} asset(s) of which ${cached} were cached at ${outJson}`);
  console.log(`Therefore, ${dirty} asset(s) were parsed and written to ${outJson} by ${since()}ms`);
});

const assetList = await Promise.all(assetPromiseList);

/** @type {Page[]} */
const pageList = assetList.map(asset => {
  const medium = asset.frontMatter.source.kind
    .split(" ")
    .map(word => word[0].toUpperCase() + word.substring(1))
    .join(" ");

  const locations = asset.frontMatter.source.url ? [asset.frontMatter.source.url] : [];
  locations.push(...asset.frontMatter.source.mirrors || []);

  const issuesFilter = (() => {
    const hasIssues = asset.issues.length > 0 || asset.sections.some(s => s.issues.length > 0);
    if (hasIssues) return "Has Issues";
    if (!asset.frontMatter.completion?.issues) return "Unknown";
    return "No Issues";
  })();

  const metaData = (() => {
    let result = "";

    result += tableRow("Medium", html`<span data-pagefind-filter="medium">${medium}</span>`);
    result += tableRow("Date", html`<span data-pagefind-sort="date">${asset.date}</span>`);

    if (asset.contributors.length > 0) {
      const contributorsHtml = asset.contributors.map(contributor => html`
        <li ${contributor.filterable ? `data-pagefind-filter="contributor"` : ""} >
          ${contributor.name}
        </li>
      `).join("");

      result += tableRow(html`Contributors`, html`
        <ul>${contributorsHtml}</ul>
      `);
    }

    if (locations.length > 0) {
      result += tableRow(html`Locations`, (() => {
        const overrideList = [
          ["https://github.com/0x2447196/raypeatarchive", "github.com/0x2447196/raypeatarchive"],
          ["https://data.raypeatforum.com", "raypeatforum.com"],
        ];

        let locationsHtml = "<ul>";

        locations.forEach(location => {
          const url = new URL(location);
          let hostname = url.hostname;

          for (const override of overrideList) {
            if (location.startsWith(override[0])) {
              hostname = override[1];
              break;
            }
          }

          if (hostname.startsWith("www.")) {
            hostname = hostname.substring(4);
          }

          locationsHtml += html`
            <li
              data-pagefind-filter="also on[data-hostname]"
              data-hostname="${hostname}"
            >
              <a href="${location}">${location}</a>
            </li>
          `;
        });

        locationsHtml += "</ul>";
        return html`${locationsHtml}`;
      })());
    }

    if (asset.frontMatter.added) {
      result += tableRow("Added By", asset.frontMatter.added.author);
      result += tableRow("Added On", asset.frontMatter.added.date);
    }

    const support = "support@raypeat.rodeo";

    result += tableRow("Completion", (() => {
      const completion = asset.frontMatter.completion;
      if (!completion) return "Incomplete";

      let completionHtml = "";
      completionHtml += completion.content
        ? `<span data-name="Content Added" data-pagefind-filter="completion[data-name]">Content added</span>`
        : `Content missing`;

      completionHtml += completion["content-verified"]
        ? ` <span data-pagefind-filter="completion[data-name]" data-name="Content Verified">and verified</span>`
        : ` but unverified`;

      completionHtml += completion["speakers-identified"]
        ? `, <span data-pagefind-filter="completion[data-name]" data-name="Contributors Identified">contributors identified</span>`
        : `, contributors not identified`;

      completionHtml += completion.issues
        ? `, <span data-pagefind-filter="completion[data-name], issues[data-issues]" data-name="Issues Identified" data-issues="${issuesFilter}">issues identified</span>`
        : `, issues not identified`;

      completionHtml += completion.notes
        ? `, <span data-pagefind-filter="completion[data-name]" data-name="Notes Identified">notes identified</span>`
        : `, notes not identified`;

      completionHtml += completion.timestamps
        ? `, <span data-pagefind-filter="completion[data-name]" data-name="Timestamps Identified">timestamps identified</span>.`
        : `, timestamps not identified.`;

      return completionHtml;
    })());

    result += tableRow(html`Legal`, (() => {
      switch (asset.frontMatter.source.kind) {
        case "book":
          return html`Copyright for Ray Peat's books belongs to Ray Peat. 
Permission to republish this book has yet to be sought from Ray Peat's estate.
Ray peat's books are currently out of print, and elsewhere freely available
both in part and in full. If you are associated with Ray Peat's estate and 
would like to discuss copyright permission, please 
<a href="mailto:${support}?subject=re: Ray Peat's Books">get in touch</a>.`;
        case "article":
          return html`Copyright for Ray Peat's articles belongs to Ray Peat.
Permission to republish this article has yet to be sought from Ray Peat's 
estate. Ray's articles were already publically available on his website. If 
you are associated with Ray Peat's estate and would like to discuss copyright 
permission, please 
<a href="mailto:${support}?subject=re: Ray Peat's Articles">get in touch</a>.`;
        case "paper":
          return html`Copyright for this journal article belongs to the original
publisher. Permission to republish this article has yet to be sought from the 
original publisher. All articles were elsewhere freely available. If you are
associated with the original publisher of this journal article and would like 
to discuss copyright permission, please
<a href="mailto:${support}?subject=re: Republishing of ${he.escape(asset.frontMatter.source.series)}">get in touch</a>.`;
        case "audio":
        case "video":
          return html`This transcript, although a unique creation, is considered a 
derivative work which may only be published with the permission of the original
work's publisher, which has not yet been sought. If you are the original 
publisher, please 
<a href="mailto:${support}?subject=re: Derivative Works of ${he.escape(asset.frontMatter.source.series)}">get in touch</a>.`;
        default:
          return html`If you are the original publisher of this content and 
would like to discuss copyright permission, please 
<a href="mailto:${support}?subject=re: Copyright for ${he.escape(asset.frontMatter.source.title)}">get in touch</a>.`;
      }
    })());

    const edit = `https://github.com/marcuswhybrow/ray-peat-rodeo/edit/main/assets/${asset.date}-${asset.slug}.md`;
    const download = `https://raypeat.rodeo/${asset.slug}.md`;
    const dataLink = `https://raypeat.rodeo/${asset.slug}.json`;

    result += tableRow(`Edit on GitHub`, html`<a href="${edit}">${edit}</a>`);
    result += tableRow(`Markdown`, html`<a href="/${asset.slug}.md">${download}</a>`);
    result += tableRow(`JSON`, html`<a href="/${asset.slug}.json">${dataLink}</a>`);

    return result;
  })();

  const results = assetList.map((asset, index) => {
    /** @param {Section[]} sections */
    const sectionsToHtml = sections => {
      if (!sections) return [];
      return sections.map(section => html`
        <li class="section" data-id="${section.id}" data-depth="${section.depth.toString()}">
          <a 
            class="link"
            href="/${asset.slug}/#${section.id}"
            class="depth-${section.depth.toString()}"
            ariaExpanded="false"
          >
            <span class="title">${section.title}</span>
          </a>
          <ol class="subsections">
            ${sectionsToHtml(section.subsections).join("")}
          </ol>
        </li>
      `);
    };

    return html`
      <li class="result" data-slug="${asset.slug}" data-score="0">
        <a 
          class="header"
          href="/${asset.slug}/"
          ariaCurrent="${index === 0 ? "true" : "false"}"
        >
          <h3 class="title">${asset.frontMatter.source.title}</h3>
          <p class="details">${asset.date} ${asset.frontMatter.source.series}</p>
        </a>
        <ol class="sections">
          <div class="highlight" data-slug="${asset.slug}"></div>
          <div class="overlight" data-slug="${asset.slug}"></div>
          ${sectionsToHtml(asset.sections).join("")}
        </ol>
      </li>
    `;
  });

  const partial = html`
    <article data-pagefind-body slot="asset">
      <header>
        <h1>${asset.frontMatter.source?.title}</h1>
        <p>${asset.date} <rpr-filter data-pagefind-filter="publisher[value]" key="publisher" value="${asset.frontMatter.source?.series}"></rpr-filter></p>
      </header>
      <div>
        ${asset.html}
      </div>
      <footer>
        <table>
          ${metaData}
        </table>
      </footer>
    </article>
  `;

  const content = html`
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <title>Ray Peat Rodeo</title>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="manifest" href="/favicons/site.webmanifest">
        <link rel="apple-touch-icon" sizes="180x180" href="/favicons/apple-touch-icon.png">
        <link rel="icon" type="image/png" sizes="32x32" href="/favicons/favicon-32x32.png">
        <link rel="icon" type="image/png" sizes="16x16" href="/favicons/favicon-16x16.png">
        <link rel="stylesheet" href="/global.css">
        <script src="/scripts/setZeroTimeout.min.js"></script>
        <script type="module" src="/derived/filters.js"></script>
        <script type="module" src="/pagefind/pagefind.js"></script>
        <script type="module" src="/components/rpr-timecode.js"></script>
        <script type="module" src="/components/rpr-sidenote.js"></script>
        <script type="module" src="/components/rpr-issue.js"></script>
        <script type="module" src="/components/rpr-contribution.js"></script>
        <script type="module" src="/components/rpr-filter.js"></script>
        <script type="module" src="/components/app-root.js"></script>
      </head>
      <body>
        <app-root>
          ${partial}
          <ol class="results" slot="results">
            ${results.join("")}
          </ol>
        </app-root>
      </body>
    </html>
  `;

  /** @type {Page} */
  const page = {
    asset, partial,
    html: content,
    filters: extractPagefindFilters(partial),
  };

  return page;
});

console.log(`Detertimed data for all pages by ${since()}ms`);

/** @type {PagefindFilters} */
let filters = {};

pageList.forEach(page => {
  for (const [name, values] of Object.entries(page.filters)) {
    for (const [value, count] of Object.entries(values)) {
      if (!filters.hasOwnProperty(name)) filters[name] = {};
      if (!filters[name].hasOwnProperty(value)) filters[name][value] = count;
      else filters[name][value] += count;
    }
  }
});

writeBuffer.push(() => {
  fs.mkdirSync(`${OUT}/derived`, { recursive: true });
  fs.writeFileSync(
    `${OUT}/derived/filters.js`,
    `export const FILTERS = ${JSON.stringify(filters)};`
  );
});

filters = Object.fromEntries(Object.entries(filters)
  .sort((a, b) => a[0].localeCompare(b[0]))
  .map(([key, values]) => [
    key,
    Object.fromEntries(Object.entries(values)
      .sort((a, b) => a[0].localeCompare(b[0]))
    )
  ])
);

const { index: pagefind } = await pf.createIndex();
assert(pagefind);

const latestPagePartial = pageList.reduce((latest, page) => {
  const partialFilename = outPartialFilename(page.asset.slug);
  if (!fs.existsSync(partialFilename)) return latest;
  const { mtime } = fs.statSync(partialFilename);
  const value = mtime.valueOf();
  if (value > latest) return value;
  return latest;
}, 0);

let rebuildPagefindIndex = false;

if (fs.existsSync(BUILD_CACHE)) {
  const prevLatestPagePartial = fs.readFileSync(BUILD_CACHE, "utf8");
  console.log("prev", JSON.parse(prevLatestPagePartial));
  console.log("curr", latestPagePartial);
  if (latestPagePartial > JSON.parse(prevLatestPagePartial)) {
    rebuildPagefindIndex = true;
  }
} else {
  rebuildPagefindIndex = true;
}

writeBuffer.push(() => fs.writeFileSync(BUILD_CACHE, JSON.stringify(latestPagePartial)));

let indexingPromiseList = [];
if (rebuildPagefindIndex) {
  indexingPromiseList = pageList.map(async page => {
    await pagefind.addHTMLFile({
      url: `/${page.asset.slug}`,
      content: html`
        <!DOCTYPE html>
        <html lang="en">
          <head><title>${he.escape(page.asset.frontMatter.source.title)}</title><meta charset="UTF-8"></head>
          <body>${page.partial}</body>
        </html>
      `,
    });
  });
}

const writePromiseList = pageList.map(async (page, index) => {
  const out = outHtmlFilename(page.asset.slug);
  const partial = outPartialFilename(page.asset.slug);
  const parsed = path.parse(out);

  writeBuffer.push(() => {
    fs.mkdirSync(parsed.dir, { recursive: true });
    fs.writeFileSync(out, page.html);
  });

  if (fs.existsSync(partial)) {
    const extant = fs.readFileSync(partial, "utf8");

    // Preserving file modified is important for caching.
    if (page.partial !== extant) {
      writeBuffer.push(() => fs.writeFileSync(partial, page.partial));
    }
  } else {
    writeBuffer.push(() => fs.writeFileSync(partial, page.partial));
  }

  writeBuffer.push(() => fs.copyFileSync(page.asset.filename, outMarkdownFilename(page.asset.slug)));

  // Make home page a copy of the latest asset
  if (index === 0) {
    writeBuffer.push(() => fs.writeFileSync(path.join(OUT, "index.html"), page.html));
  }
});

const derivedDone = Promise.all(assetPromiseList).then(assetList => {
  writeBuffer.push(() => fs.mkdirSync(`${OUT}/derived`, { recursive: true }));

  /** @type {ThinAsset[]} */
  const thinAssetList = assetList.map(asset => ({
    title: asset.frontMatter.source.title,
    slug: asset.slug,
    date: asset.date,
    publisher: asset.frontMatter.source.series,
    sections: asset.sections,
    issues: asset.issues,
  }));

  const assetsJson = Object.fromEntries(assetList.map(asset => [asset.slug, `/${asset.slug}.json`]));

  const assetInputs = assetList
    .map(asset => `"${asset.slug}": resolve(__dirname, "${asset.slug}/index.html"),`)
    .join("\n");

  const viteConfigJs = js`
    import {resolve} from "path"
    import {defineConfig} from "vite"

    export default defineConfig({
      build: {
        rollupOptions: {
          input: {
            main: resolve(__dirname, "index.html"),
            ${assetInputs}
          },
        },
      },
    })
  `;

  writeBuffer.push(() => {
    fs.mkdirSync(`${OUT}/public/derived`, { recursive: true });
    fs.writeFileSync(`${OUT}/public/assets.json`, JSON.stringify(assetsJson));
    fs.writeFileSync(`${OUT}/derived/thin-assets.js`, `export default ${JSON.stringify(thinAssetList)};`);
    fs.writeFileSync(path.join(OUT, "vite.config.js"), viteConfigJs);
  });
});

const pagefindDone = Promise.all(indexingPromiseList).then(async () => {
  if (rebuildPagefindIndex) {
    writeBuffer.push(async () => {
      await pagefind.writeFiles({ outputPath: `${OUT}/public/pagefind` })
      fs.mkdirSync(`${OUT}/pagefind`, { recursive: true });
      fs.renameSync(`${OUT}/public/pagefind/pagefind.js`, `${OUT}/pagefind/pagefind.js`);
      console.log(`Wrote Pagefind files to ${OUT}/pagefind/ by ${since()}ms`);
    });
  } else {
    console.log(`Pagefind files already cached at ${OUT}/pagefind/ by ${since()}ms`);
  }
});

await Promise.all([...writePromiseList, pagefindDone, derivedDone]).then(async () => {
  for (const action of writeBuffer) await action();
  console.log(`Completed in ${since()}ms`);
});
