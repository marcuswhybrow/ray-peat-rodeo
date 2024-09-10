window.customElements.define("rpr-search", class extends HTMLElement {
  #assets
  #assetBox 
  #autoCompleteBox 
  #searchBar
  #sidebarButton

  static observedAttributes = [ 
    "assets" 
  ];

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
      <style>
        .top {
          display: flex;
          margin-bottom: 1rem;
        }
        .top search-bar {
          flex: 1 1 auto;
        }
        .top sidebar-button {
          margin-left: 1rem;
        }
      </style>
      <div class="top">
        <search-bar></search-bar>
        <sidebar-button></sidebar-button>
      </div>
      <auto-complete-box></auto-complete-box>
      <asset-box>
        <slot></slot>
      </asset-box>
    `;
    this.#assetBox = this.shadowRoot.querySelector("asset-box");
    this.#autoCompleteBox = this.shadowRoot.querySelector("auto-complete-box");
    this.#searchBar = this.shadowRoot.querySelector("search-bar");
    this.#sidebarButton = this.shadowRoot.querySelector("sidebar-button");
  }

  connectedCallback() {

    this.#autoCompleteBox.initItems(this.#assetBox.assets);

    this.shadowRoot.addEventListener("query-changed", event => {
      const query = event.detail.query;
      // this.#autoCompleteBox.queryChanged(query);
      // this.#assetBox.queryChanged(query);
    });
  }

  disconnectedCallback() {

  }

  attributeChangedCallback(name, oldvalue, newValue) {
    switch(name) {
      case "assets":
        const assetList = JSON.parse(newValue);
        const assets = {};
        for (const data of assetList) {
          const asset = document.createElement("rpr-asset");
          asset.id = data.AbsURL;
          asset.kind = data.Kind;
          asset.date = data.Date;
          asset.title = data.Title;
          assets[asset.id] = asset;
        }
        this.#assets = assets;
        this.#assetBox.replaceChildren(...Object.values(this.#assets));
        break;
    }
  }

  get assets() {
    return this.#assets;
  }

  set assets(a) {
    const data = Object.values(a).map(asset => ({
      AbsURL: asset.id,
      Kind: asset.kind,
      Date: asset.date,
      Title: asset.title,
    }));

    this.setAttribute("assets", JSON.stringify(data));
  }
});


// import SearchBar from "./search-bar.js";

// const pagefind = new Promise(async resolve => {
//   const pf = await import("/pagefind/pagefind.js");
//   await pf.options({
//     highlightParam: 'highlight',
//     excerptLength: 60,
//   });
//   resolve(pf);
// });

// const domHasLoaded = new Promise(resolve => {
//   window.addEventListener("DOMContentLoaded", () => {
//     resolve(true)
//   });
// });

// const assetsContainer = new Promise(async resolve => {
//   await domHasLoaded;
//   resolve(document.getElementById("assets"));
// });

// const reactiveFiltersContainer = new Promise(async resolve => {
//   await domHasLoaded;
//   resolve(document.getElementById("rpr-reactive-filters"));
// });

// const assetElements = new Promise(async resolve => {
//   await domHasLoaded;
//   const container = await assetsContainer;
//   const obj = {};
//   for (const assetElement of container.children) {
//     obj[assetElement.id] = assetElement;
//   }
//   resolve(obj);
// });


// const globalFilters = new Promise(async resolve => {
//   const pf = await pagefind;
//   resolve(await pf.filters());
// });

// /*
// const completionElements = new Promise(async resolve => {
//   const [filters, _] = await Promise.all([globalFilters, assetElements, domHasLoaded]);

//   const removeHandler = event => {

//   };

//   const highlightHandler = event => {
//   };

//   let tabindex = 2;

//   const obj = {};
//   for (const [key, values] of Object.entries(filters)) {
//     for (const [value, count] of Object.entries(options)) {
//       const elem = document.createElement("div");
//       elem.id = `${key}/${value}`.replace(" ", "-");
//       elem.tabindex = tabindex += 1;
//       elem.dataset.filterKey = key;
//       elem.dataset.filterValue = value;
//       elem.addEventListener("rpr-remove", removeHandler);
//       elem.addEventListener("rpr-highlight", highlightHandler);
//       elem.className = `
//         cursor-pointer
//         px-2 py-1 mr-2 mb-2 inline-block rounded
//         hover:bg-slate-400 hover:shadow hover:text-white
//       `;
//       switch (name) {
//         case "Contributor":
//           elem.className += ` bg-sky-50 `;
//           break;
//         case "Publisher":
//           elem.className += `bg-purple-50`;
//           break;
//         case "Medium":
//           elem.className += `bg-rose-50`;
//           break;
//         default:
//           elem.className += `bg-slate-50`;
//       }
//       elem.append((() => {
//         const e = document.createElement("span");
//         e.classList.add("key");
//         e.append(key);
//         e.className = ` text-xs uppercase text-slate-500 mr-2 `;
//         return e;
//       })());
//       elem.append((() => {
//         const e = document.createElement("span");
//         e.classList.add("value");
//         e.append(value);
//         return e;
//       })());
//       elem.append((() => {
//         const close = document.createElement("span");
//         close.append("x");
//         return close;
//       })());
//       obj[elem.id] = elem;
//     }
//   }
//   resolve(obj);
// });
// */


// window.addEventListener("rpr-filter-changed", async event => {
//   window.dispatchEvent(new CustomEvent("rpr-trigger-search", {
//     detail: { debounce: 0 }
//   }));
// });

// window.addEventListener("rpr-query-changed", async event => {
//   window.dispatchEvent(new CustomEvent("rpr-trigger-search", {
//     detail: { debounce: 300 }
//   }));
// });

// const searchInputElement = new Promise(async resolve => {
//   await domHasLoaded;
//   resolve(document.getElementById("rpr-search-input"));
// });

// async function getQuery() {
//   const element = await searchInputElement;
//   return element.value;
// }

// async function getIteration() {
//   const element = await searchInputElement;
//   return element.dataset.iteration;
// }

// let iteration = 0;
// const fuzzy = new uFuzzy({
//   intraChars: ".", // Allows any characters between matches
//   intraIns: Infinity, // Allows any amount of characters between matches
// });

// /*
// const autoComplete = new Promise(async resolve => {
//   const haystack = [];
//   const reverseLookup = [];

//   for (const c of Object.values(await completeionElements)) {
//     haystack.push(c.dataset.value);
//     reverseLookup.push(c);
//   }

//   resolve({ haystack, reverseLookup });
// });
// */ 

// /*
// window.addEventListener("rpr-trigger-search", async event => {
//   const debounce = event.detail.debounce || 300;
//   const pf = await pagefind;
//   const query = await getQuery();
//   const rfc = await reactiveFiltersContainer;

//   // Autocomplete
//   const ac = await autoComplete;
//   const [idxs, info, order] = fuzzy.search(ac.haystack, query);
  
//   if (idxs != null && idxs.length > 0) {
//     for (const i in order) {
//       rfc.append((() => {
//         const groupName = groupLookup[info.idx[order[i]]];
//         const text = allFilterNames[info.idx[order[i]]];
//         const id = `${groupName}/${text}`;
//         return globalFilters[id];
//       })());
//     }
//   }

//   const options = { filters };
//   const result = await pf.debouncedSearch(query || null, options, debounce);

//   if (result === null) return;
//   const instigatingIteration = iteration += 1;
//   if (instigatingIteration !== iteration) return;
//   // Debounced searches return null when recalled during the dbounce window.

//   Object.entries(result.filters).forEach(([group, names]) => {
//     Object.entries(names).forEach(async ([name, count]) => {
//       const id = `${group}/${name}`.replaceAll(" ", "-");
//       (await filterElements)[id].dispatchEvent(new CustomEvent("rpr-update-count", {
//         detail: { count },
//       }));
//     });
//   })

//   const ac = await assetsContainer 
//   ac.replaceChildren();

//   result.results.forEach(async result => {
//     if (instigatingIteration !== iteration) return;

//     const data = await result.data();
//     const id = data.raw_url;

//     const assetEl = (await assetElements)[id];
//     if (assetEl === undefined) return;

//     assetEl.setAttribute("tabindex", 10000 + ac.childElementCount);

//     const container = assetEl.getElementsByClassName("excerpt")[0];
//     if (query !== "") {
//       const excerpts = data.sub_results.map(result => {
//         const e = document.createElement("div");
//         e.innerHTML = result.excerpt;
//         return e;
//       });
//       container.replaceChildren(...excerpts);
//       container.classList.remove("hidden");
//     } else {
//       container.classList.add("hidden");
//     }

//     ac.append(assetEl);
//   });
// });
// */


// /*
// document.addEventListener("DOMContentLoaded", async () => {
//   const assets = document.getElementById("assets");
//   const filtersElm = document.getElementById("rpr-search-filters");
//   const input = document.getElementById("rpr-search-input");

//   const search = async debounce => {
//     let query = input.value;

//     const activeFilters = filtersElm.querySelectorAll(".active");
//     let finalFilters = {};
//     for (const f of activeFilters) {
//       const group = f.getAttribute("data-filter-group");
//       const id = f.getAttribute("data-filter-id");
//       finalFilters[group] = [id, ...(finalFilters[group] || [])];
//     }

//     // Tells pagefind to search anyway
//     if (query === "") {
//       query = null;
//     }

//     const result = await (await pagefind).debouncedSearch(query, {
//       filters: finalFilters,
//     }, debounce);

//     // User is still typing
//     if (result === null) return;

//     // Update filter display
//     for (const [group, options] of Object.entries(result.filters)) {
//       for (const [option, count] of Object.entries(options)) {
//         const countId = `${group}/${option}/count`.replaceAll(" ", "-");
//         const countElm = document.getElementById(countId);
//         countElm.textContent = count;

//         const filterId = `${group}/${option}`.replaceAll(" ", "-");
//         const filterElm = document.getElementById(filterId);
//         if (count == 0) {
//           filterElm.classList.add("no-results");
//         } else {
//           filterElm.classList.remove("no-results");
//         }
//       }
//     }

//     // No results
//     if (query === "") {
//       for (const asset of assets.children) {
//         asset.classList.add("hit");
//         asset.classList.remove("has-excerpt");
//       }
//     }

//     // if (result.results.length === 0) return;

//     const data = await Promise.all(result.results.map(r => r.data()));

//     for (const asset of assets.children) {
//       asset.classList.remove("hit");
//       asset.classList.remove("has-excerpt");
//     }

//     for (const d of data) {
//       const asset = document.getElementById(d.raw_url);
//       if (asset === null) continue;
//       asset.classList.add("hit");
//       if (query !== null) {
//         asset.classList.add("has-excerpt");
//         const excerpt = asset.getElementsByClassName("excerpt")[0];
//         let newExcerpts = [];
//         for (const subResult of d.sub_results) {
//           const elem = document.createElement("div");
//           elem.innerHTML = subResult.excerpt;
//           newExcerpts.push(elem);
//         }
//         excerpt.replaceChildren(...newExcerpts);
//       }
//     }
//   };

//   input.addEventListener("keyup", event => {
//     search(300);
//   });

//   const filters = await (await pagefind).filters();

//   for (const [name, options] of Object.entries(filters)) {
//     const heading = document.createElement("h3");
//     heading.textContent = name;
//     heading.className = `
//       text-xs inline-block uppercase mr-2 mb-2 pr-2 py-1
//     `;
//     filtersElm.append(document.createElement("br"));
//     filtersElm.append(heading);

//     for (const [option, count] of Object.entries(options)) {
//       const elem = document.createElement("span");

//       elem.append(`${option} (`);
//       const spanCount = document.createElement("span");
//       spanCount.id = `${name}/${option}/count`.replaceAll(" ", "-");
//       spanCount.textContent = count;
//       elem.append(spanCount);
//       elem.append(")");

//       elem.id = `${name}/${option}`.replaceAll(" ", "-");
//       elem.className = `
//         text-xs inline-block bg-slate-100 rounded mr-2 mb-2
//         text-slate-600 px-2 py-1
//         cursor-pointer
//         [&.no-results:not(.active)]:bg-slate-50 [&.no-results:not(.active)]:text-slate-400
//         [&.active]:bg-slate-400
//       `;
//       elem.setAttribute("data-filter-group", name);
//       elem.setAttribute("data-filter-id", option);
//       elem.addEventListener("click", event => {
//         elem.classList.toggle("active");
//         search(0);
//       });

//       filtersElm.appendChild(elem);
//     }
//   }

//   const handleParam = (paramName, filterGroup) => {
//     const params = new URLSearchParams(document.location.search);
//     const param = params.get(paramName);
//     if (param === null) return;
//     const paramList = param.split(",");
//     if (paramList.length === 1 && paramList[0] === "") return;
//     for (const i in paramList) {
//       if (paramList[i] === null) continue;
//       if (paramList[i] === "") continue;
//       const id = `${filterGroup}/${paramList[i]}`;
//       const elem = document.getElementById(id);
//       if (elem === null) continue;
//       elem.click(); // assumes click adds the .active class
//     }
//   };

//   handleParam("publisher", "Author");
//   handleParam("medium", "Medium");
//   handleParam("todo", "Todo");
//   handleParam("has-issues", "Has-Issues");

//   (() => {
//     const searchParam = params.get("search");
//     input.value = searchParam;
//     if (searchParam === null) return;
//     if (searchParam === "") return;
//     // hacky, deserves using custom events.
//     input.dispatchEvent(new Event("keyup"));
//   })();

//   (() => {
//     const assetId = window.location.pathname;
//     if (!assetId.endsWith("/")) assetId += "/";
//     const assetElem = document.getElementById(assetId);
//     if (assetElem === null) return;

//     // This needs it's own event
//     assetElem.classList.add("clicked");
//   })();
// });
// */

// await import('/pagefind/pagefind-highlight.js');
// new PagefindHighlight({ highlightParam: 'highlight' });

// (await searchInputElement).addEventListener("keyup", async event => {
//   window.dispatchEvent(new CustomEvent("rpr-trigger-search", {
//     detail: { debounce: 300 }
//   }));
// });
