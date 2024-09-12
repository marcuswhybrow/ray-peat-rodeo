const pagefind = new Promise(resolve => {
  import("/pagefind/pagefind.js").then(pagefind => {
    pagefind.options({
      highlightParam: 'highlight',
      excerptLength: 60,
      showSubResults: true,
    }).then(() => {
      resolve(pagefind);
    });
  });
});

const fuzzy = new uFuzzy({
  intraChars: ".", // Allows any characters between matches
  intraIns: 3, // Allows any amount of characters between matches
});

const interrupt = () => new Promise(resolve => setZeroTimeout(resolve));

class Search extends HTMLElement {
  /** @type {string} */
  #query = "";

  /** @type {Object.<string, string[]>} */
  #filters = {};

  /** @type {Object.<string, string[]>} */
  #prevFilters = {};

  #searchCount = 0;

  /** @type {Promise<Object.<string, Asset>>} */
  #assets

  /** @type {Promise<Asset[]>} */
  #assetsNewestFirst;

  /** @type {Deck} */
  #deckElement

  /** @type {Promise<Pin[]>} */
  #pins

  /** @type {Boolean} */
  #interactive

  /** @type {HTMLElement} */
  #readingPaneElement

  /** @type {HTMLInputElement} */
  #inputElement

  /** @type {Pins} */
  #pinsElement

  /** @type {HTMLButtonElement} */
  #sidebarButton

  static observedAttributes = ["query", "filters", "interactive"];
  static #domParser = new DOMParser();

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
        }
        form {
          width: 100%;
          display: flex;
          flex-direction: column;
          gap: 1rem;
        }
        form input {
          flex: 1 1 auto;
          border: 2px solid #bbb;
          border-radius: 0.25rem;
          padding: 0.5rem 1rem 0.5rem;
          font-size: 1.1rem;
          line-height: 2rem;
        }
        form input::placeholder {
          color: #777;
        }

        .under {
          font-size: 0.75rem;
          line-height: 0.75rem;
          color: #999;
          cursor: pointer;
          text-align: center;
          text-transform: lowercase;
          display: flex;
          flex-direction: row;
          gap: 1rem;
        }
        .under > *:hover {
          color: #111;
        }

        :host(:not([interactive="true"])) #pins {
          opacity: 0.2;
        }
        :host(:not([interactive="true"])) #deck {
          opacity: 0.2;
        }
        #header {
          padding-top: 2rem;
          background: white;
          display: flex;
          flex-direction: column;
          gap: 1rem;
          position: sticky;
          top: 0;
        }
        #pins {
          margin-bottom: 1rem;
        }

        .buttons {
          display: flex;
          flex-direction: row;
        }

        button {
          flex-grow: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;
          justify-items: center;
          gap: 0.25rem;

          background: transparent;
          border: none;

          cursor: pointer;
          text-align: center;
          padding: 1rem 0;
          justify-content: center;
          align-items: center;
        }
        button .caption {
          font-size: 0.725rem;
          line-height: 1rem;
          text-transform: lowercase;
          letter-spacing: 0.1em;
          color: #aaa;
        }

        button svg {
          height: 2rem;
        }
        button svg * {
          stroke: #aaa;
          fill: transparent;
        }
        button svg .filled {
          fill: #aaa;
        }

        button.book svg {
          transform: translateX(3px);
        }

        button.ray-peat img {
          clip-path: circle(50% at 20% 40%);
          aspect-ratio: 1/1;
          width: 3rem;
          transform: translateX(0.5rem) translateY(-0.5rem);
        }
        button.ray-peat .img-wrapper {
          height:2rem;

        }

      </style>
      <div id="header">
        <div class="buttons">
          <button class="ray-peat" tabindex="-10">
            <!-- <div class="img-wrapper"><img src="/assets/images/ray-peat-head.webp"/></div> -->
            <svg viewBox="0 0 100 100">
              <circle cx="50" cy="50" r="47" stroke-width="5" />
              </svg>
            <div class="caption">Ray Peat</div>
          </button>

          <button class="help" tabindex="-9">
            <svg viewBox="0 0 110 100">
              <polygon points="3,97 107,97 55,3" stroke-width="5" />
              <circle class="filled" cx="72" cy="78" r="5" />
              <circle class="filled" cx="55" cy="78" r="5" />
              <circle class="filled" cx="38" cy="78" r="5" />
            </svg>
            <div class="caption">Help</div>
          </button>

          <button class="sidebar" tabindex="-8">
            <svg viewBox="0 0 100 100">
              <rect 
                class="border" 
                x="3" y="3" 
                width="94" height="94"
                rx="10" ry="10" 
                stroke-width="5"
                stroke="#F00"
                fill="transparent"
              />
              <line x1="25" y1="0" x2="25" y2="100" stroke="#F00" stroke-width="5" /> 
            </svg>
            <div class="caption">Sidebar</div>
          </button>

        </div>
        <form>
          <input 
            id="input" 
            type="text" 
            placeholder="Search"
            tabIndex="1"
            autocomplete="off"
            value="${this.#query}"
          />
          <div class="under">
            <div class="advanced-search">show all filters</div>
            <div>share this search</div>
            <div>help</div>
          </div>
        </form>
        <rpr-pins id="pins"></rpr-pins>
      </div>
      <rpr-deck id="deck"></rpr-deck>
    `;

    this.#deckElement = /** @type {Deck} */ (shadowRoot.querySelector("#deck"));
    this.#readingPaneElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#reading-pane"));
    this.#inputElement = /** @type {HTMLInputElement} */ (shadowRoot.querySelector("#input"));
    this.#pinsElement = /** @type {Pins} */ (shadowRoot.querySelector("#pins"));
    this.#sidebarButton = /** @type {HTMLButtonElement} */ (shadowRoot.querySelector("button.sidebar"));
    const advancedSearch = /** @type {HTMLElement} */ (shadowRoot.querySelector(".advanced-search"));
    const rayPeatElement = /** @type {HTMLElement} */ (shadowRoot.querySelector(".ray-peat"));

    rayPeatElement.addEventListener("click", () => {
      const pane = document.querySelector("#reading-pane");
      if (pane) {
        pane.replaceChildren(new RayPeat());
      }
    });

    advancedSearch.addEventListener("click", async () => {
      const advancedSearchPage = new AdvancedSearch();
      advancedSearchPage.pins = /** @type {Pin[]} */ ((await this.#pins).map(pin => pin.cloneNode()));
      advancedSearchPage.filters = this.filters;
      const pane = document.querySelector("#reading-pane");
      if (pane) {
        pane.replaceChildren(advancedSearchPage);
      }
    });

    this.#sidebarButton.addEventListener("click", event => {
      this.dispatchEvent(new CustomEvent("sidebar", {
        bubbles: true,
        detail: "toggle",
      }));
      event.stopPropagation();
    });

    this.#inputElement.addEventListener("keyup", () => {
      this.query = this.#inputElement.value;
    });

    this.focus = () => this.#inputElement.focus();
    this.select = () => this.#inputElement.select();

    window.addEventListener("keydown", event => {
      if (event.key === "/") {
        if (this !== document.activeElement) {
          this.dispatchEvent(new CustomEvent("sidebar", {
            bubbles: true,
            detail: "open",
          }));
          this.#inputElement.focus();
          this.#inputElement.select();
          event.preventDefault();
        }
      } else if (event.key === "Escape" && this === document.activeElement) {
        this.#inputElement.blur();
      }
    });


    this.#assets = new Promise(async resolve => {
      const response = await fetch("/search.json");
      if (!response.ok) {
        console.error(`Failed to fetch /search.json`);
        return null;
      }

      const assetList = await response.json();
      if (assetList.length <= 0) {
        console.error(`Found 0 assets listed at /search.json`);
        return null;
      }

      let tabIndex = 10000;
      let activeAsset = null;

      /** @type {Object.<string, Asset>} */
      const assets = {};

      for (const asset of assetList) {
        const a = new Asset();
        a.path = asset.Path;
        a.title = asset.Title;
        a.series = asset.Series;
        a.date = asset.Date;
        a.kind = asset.Kind;
        a.issues = asset.Issues;
        a.tabIndex = tabIndex++;

        if (a.deriveActive()) {
          activeAsset = a;
        }

        /**
          * @type {(event: Event) => Promise<void>}
          * @param {PickEvent} event
          */
        async function pickHandler(event) {
          const response = await fetch(a.path);
          const issue = event.detail.issue;

          if (!response.ok) {
            console.error(`Failed to fetch ${a.path}`);
            return;
          }

          const text = await response.text();
          const doc = Search.#domParser.parseFromString(text, "text/html");

          const grabbed = doc.querySelector("#reading-pane");
          if (!grabbed) {
            throw new Error(`Failed to find #reading-pane.`);
          }

          // #reading-pane must be re-selected every time because it may be replaced in the DOM.
          const readingPane = document.querySelector("#reading-pane");
          if (!readingPane) {
            throw new Error("Failed to find #reading-pane");
          }

          readingPane.replaceWith(grabbed);

          if (issue === null) {
            window.scrollTo({ top: 0, behavior: "instant" });
          } else {
            const issueBubble = /** @type {Issue|null} */ (document.querySelector(`#issue-${issue}`));
            if (issueBubble) {
              window.scrollTo({ top: issueBubble.offsetTop, behavior: "instant" });
            }
          }

          let link = a.path;
          if (window.location.search) link += window.location.search;
          if (window.location.hash) link += window.location.hash;
          history.pushState({}, "", link);
        }

        a.addEventListener("pick", pickHandler);
        assets[a.path] = a;
      }

      const newsetFirst = Object.values(assets)
        .sort((a, b) => b.date.localeCompare(a.date));

      if (activeAsset === null) {
        activeAsset = newsetFirst[0];
        activeAsset.active = true;
      }

      this.#deckElement.replace(...newsetFirst);

      if (!this.#interactive) {
        activeAsset.scrollIntoView({
          block: "center",
          behavior: "instant",
        });
      }

      resolve(assets);
    });

    this.#pins = new Promise(async resolve => {
      const pf = await pagefind;
      const filters = await pf.filters();

      const pins = [];
      for (const [key, values] of Object.entries(filters)) {
        for (const value of Object.keys(values)) {
          const pin = new Pin();
          pin.key = key;
          pin.value = value;

          pin.addEventListener("click", async () => {
            if (this.togglePin(pin)) {
              this.query = "";
            }

            const isTouchScreen = window.matchMedia("(pointer: coarse)").matches;
            if (!isTouchScreen) {
              this.select();
            }
          });

          pin.addEventListener("keyup", event => {
            if (event.key === "Enter") pin.click()
          });

          pins.push(pin);
        }
      }

      resolve(pins);
    });

    this.#assetsNewestFirst = new Promise(async resolve => {
      const elements = Object.values(await this.#assets);
      const sorted = elements.sort((a, b) => b.date.localeCompare(a.date));
      resolve(sorted);
    });
  }

  async connectedCallback() {
    const params = new URLSearchParams(window.location.search);
    const search = params.get("search");
    if (search !== null) {
      params.delete("search");
      this.query = search;
    }

    if (params.toString() === "") {
      if (this.query !== "") {
        await this.#search();
      } else {
        this.interactive = true;
      }
    }

    /** @type {Object.<string, string[]>} */
    const filters = {};

    for (const pin of (await this.#pins)) {
      if (params.getAll(pin.key).includes(pin.value)) {
        filters[pin.key] = filters[pin.key] || [];
        filters[pin.key].push(pin.value);
      }
      this.filters = filters;
    }

    if (this.query !== "" || !Search.filtersMatch(filters, {})) {
      await this.#search();
    } else {
      this.interactive = true;
    }
  }

  /**
  * @param {[string, string]} a
  * @param {[string, string]} b
  */
  static sortFiltersAlphabetically(a, b) {
    const key = a[0].localeCompare(b[0]);
    if (key !== 0) {
      return a;
    }
    const val = a[1].localeCompare(b[1]);
    return val;
  }

  /**
  * @param {Object.<string, string[]>} newFilters
  * @returns {Object.<string, string[]>}
  */
  static sanitiseFilters(newFilters) {
    /** @type {Object.<string, string[]>} */
    const result = {};

    for (const key of Object.keys(newFilters).sort()) {
      result[key] = newFilters[key].sort();
    }
    return result;
  }


  /** 
  * @param {Object.<string, string[]>} a
  * @param {Object.<string, string[]>} b
  */
  static filtersMatch(a, b) {
    const aKeys = Object.keys(a);
    const bKeys = Object.keys(b);
    if (aKeys.length !== bKeys.length) {
      return false;
    }
    return aKeys.every(key => {
      const aValues = a[key];
      const bValues = b[key];
      if (aValues.length !== bValues.length) {
        return false;
      }
      return aValues.every((val, i) => val === bValues[i]);
    });
  }

  /**
  * @param {string} name
  * @param {string} oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "query":
        this.#query = newValue;
        this.#inputElement.value = newValue;
        if (newValue !== oldValue && this.interactive) {
          this.#search();
        }
        break;
      case "filters":
        if (newValue !== oldValue) {
          this.#filters = Search.sanitiseFilters(JSON.parse(newValue));
          if (this.interactive) {
            if (!Search.filtersMatch(this.#filters, this.#prevFilters)) {
              this.#search();
            }
          }
          this.#prevFilters = structuredClone(this.#filters);
        }
        break;
      case "interactive":
        this.#interactive = newValue === "true";
        break;
    }
  }

  async #search() {
    const searchCount = ++this.#searchCount;
    this.#deckElement.replace();

    // 🪟🔗 Update Window Location

    (async () => {
      if (!this.interactive) {
        return;
      }

      const params = new URLSearchParams();
      for (const [key, values] of Object.entries(this.filters)) {
        for (const value of values) {
          params.append(key, value);
        }
      }

      if (this.query !== "") {
        params.append("search", this.query);
      }

      let link = window.location.pathname;

      const search = params.toString();
      if (search !== "") link += `?${search}`;

      if (window.location.hash !== "") link += window.location.hash;

      history.replaceState({}, "", link);
    })();

    // 📌 Filter pins

    (async () => {
      const query = this.query.replaceAll('"', "");

      const pins = [];
      for (const pin of await this.#pins) {
        pin.pinned = this.isPinned(pin);
        pins.push({
          pin,
          pinned: pin.pinned,
          order: -1
        });
      }

      const values = pins.map(pin => pin.pin.value);
      const [unordered, info, ordered] = fuzzy.search(values, query);

      if (unordered === null || unordered.length <= 0) {
        this.#pinsElement.replaceUnpinned();
        const pinned = pins.filter(pin => pin.pinned).map(pin => pin.pin);
        this.#pinsElement.replacePinned(...pinned);
        return;
      }

      for (const [orderedIndex, unorderedIndex] of ordered.entries()) {
        const pin = pins[unordered[unorderedIndex]];
        pin.pin.highlights = info.ranges[unorderedIndex] || [];
        pin.order = orderedIndex;
      }

      if (this.interactive) {
        const unpinned = pins
          .filter(p => !p.pinned && p.order >= 0)
          .sort((a, b) => a.order - b.order)
          .map(m => m.pin);
        this.#pinsElement.replaceUnpinned(...unpinned);
      }

      const pinned = pins.filter(m => m.pinned).map(m => m.pin);
      this.#pinsElement.replacePinned(...pinned);
    })();


    // 📃 Filter and search assets

    if (this.query === "" && Search.filtersMatch(this.filters, {})) {
      this.#deckElement.replace(...await this.getDefaultAssets());
      return false;
    }

    const forceQuery = this.query || null;
    const sort = {};
    if (this.query === "") {
      sort["date"] = "desc";
    }
    const filters = this.filters;
    const result = await (await pagefind).search(
      forceQuery,
      { filters, sort }
    );

    if (searchCount !== this.#searchCount) {
      return false;
    }

    if (result === null) {
      this.#deckElement.replace(...await this.getDefaultAssets());
      return false;
    }

    const trailingSlashes = /\/+$/;
    let tabIndex = (await this.#pins).length + 10;

    /** @type {Asset|null} */
    let activeAsset = null;

    const assets = await this.#assets;

    /**
     * @typedef {Object} PagefindResult
     * @param {string} raw_url
     * @param {PagefindSubResult} sub_results
     */

    /** 
      * @param {PagefindResult} data
      * @returns {Asset}
      */
    function toElement(data) {
      const key = data.raw_url.replace(trailingSlashes, "");
      const asset = assets[key];
      asset.tabIndex = tabIndex++;

      asset.showIssues = filters.hasOwnProperty("Issues")
        && filters["Issues"].includes("Has Issues");

      if (typeof asset === "undefined") {
        throw new Error("Asset not found");
      }

      asset.replaceResults(...data.sub_results);
      return asset;
    }

    // Show a usefull amount of results as quick as possible
    const burstSize = 10;
    const realSize = Math.min(burstSize, result.results.length);
    const burst = result.results.splice(0, realSize);

    const promises = new Array(realSize);
    for (const i in burst) {
      promises[i] = new Promise(async resolve => {
        const data = await burst[i].data();
        const asset = toElement(data);
        if (searchCount === this.#searchCount) {
          this.#deckElement.append(asset);
        }
        resolve(asset);
      });
      await interrupt();
    }


    // A great moment to make the UI interactive;
    this.interactive = true;

    if (searchCount !== this.#searchCount) {
      return false;
    }

    // Slow down the remainder a little to prevent UI hanging
    const remainingPromises = new Array(result.results.length);
    for (const i in result.results) {
      remainingPromises[i] = new Promise(async resolve => {
        const data = await result.results[i].data();
        resolve(toElement(data));
      });
      await interrupt();
    }

    const remainingAssets = await Promise.all(remainingPromises);

    if (searchCount !== this.#searchCount) return false;
    this.#deckElement.append(...remainingAssets);

    const allAssets = await Promise.all(promises);
    allAssets.push(remainingAssets);

    for (const asset of allAssets) {
      if (asset.active) {
        activeAsset = asset;
      }
    }

    if (activeAsset && !this.interactive) {
      activeAsset.scrollIntoView({
        block: "center",
        behavior: "instant",
      });
    }

    return true;
  }

  get query() {
    return this.#query;
  }

  set query(newValue) {
    this.setAttribute("query", newValue);
  }

  get filters() {
    return this.#filters;
  }

  set filters(newValue) {
    this.setAttribute("filters", JSON.stringify(newValue));
  }

  get interactive() {
    return this.#interactive;
  }

  set interactive(newValue) {
    this.setAttribute("interactive", newValue.toString());
  }

  async getDefaultAssets() {
    const assets = await this.#assetsNewestFirst;

    let tabIndex = this.#pinsElement.tabIndexEnd;
    for (const asset of assets) {
      asset.replaceResults();
      asset.showIssues = false;
      asset.tabIndex = ++tabIndex;
    }

    return assets;
  }

  /**
  * @param {Pin} pin
  * @returns {Boolean}
  */
  togglePin(pin) {
    const isPinned = this.isPinned(pin);
    isPinned ? this.unpin(pin) : this.pin(pin);
    return !isPinned;
  }

  /**
  * @param {Pin} pin
  */
  pin(pin) {
    const filters = structuredClone(this.filters);
    let values = filters[pin.key] || [];
    values.push(pin.value);
    values = [...new Set(values)];
    filters[pin.key] = values.sort();
    this.filters = filters;
  }

  /**
  * @param {Pin} pin
  */
  unpin(pin) {
    const filters = structuredClone(this.filters);
    const values = filters[pin.key] || [];
    const index = values.indexOf(pin.value);
    if (index >= 0) {
      values.splice(index, 1);
      if (values.length > 0) {
        filters[pin.key] = values;
      } else {
        delete filters[pin.key];
      }
    }
    this.filters = filters;
  }

  /**
  * @param {Pin} pin
  * @returns {Boolean}
  */
  isPinned(pin) {
    return this.filters.hasOwnProperty(pin.key) && this.filters[pin.key].includes(pin.value);
  }
}

customElements.define("rpr-search", Search);

