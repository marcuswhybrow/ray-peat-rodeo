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

window.customElements.define("rpr-search", class Search extends HTMLElement {
  #query = "";
  #filters = {};
  #prevFilters = {};
  #searchCount = 0;
  #assets = {};
  #assetsNewestFirst = [];
  #deck = null;
  #pins = {};
  #interactive = false;
  #touch = false;

  static observedAttributes = ["query", "filters", "interactive"];
  static #domParser = new DOMParser();

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
      <style>
        form {
          width: 100%;
          display: flex;
          gap: 1rem;
        }
        form input {
          flex: 1 1 auto;
          border: 1px solid rgb(148,163,184);
          border-radius: 0.25rem;
          padding: 0.5rem 1rem 0.5rem;
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
        }
        #pins {
          margin-bottom: 1rem;
        }
      </style>
      <div id="header">
        <form>
          <input 
            id="input" 
            type="text" 
            placeholder="Search Ray Peat Rodeo"
            tabIndex="1"
            autocomplete="off"
            value="${this.#query}"
          />
        </form>
        <rpr-pins id="pins"></rpr-pins>
      </div>
      <rpr-deck id="deck"></rpr-deck>
    `;

    const input = this.shadowRoot.querySelector("#input");
    input.addEventListener("keyup", () => {
      this.query = input.value;
    });

    this.focus = () => input.focus();
    this.select = () => input.select();

    window.addEventListener("keydown", event => {
      if (event.key === "/") {
        if (this !== document.activeElement) {
          this.dispatchEvent(new CustomEvent("sidebar", {
            bubbles: true,
            detail: "open",
          }));
          input.focus();
          input.select();
          event.preventDefault();
        }
      } else if (event.key === "Escape" && this === document.activeElement) {
        input.blur();
      }
    });

    this.#deck = this.shadowRoot.querySelector("#deck");

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
      const assets = {};
      for (const asset of assetList) {
        const a = document.createElement("rpr-asset");
        a.path = asset.Path;
        a.title = asset.Title;
        a.series = asset.Series;
        a.date = asset.Date;
        a.kind = asset.Kind;
        a.tabIndex = tabIndex++;

        if (a.deriveActive()) {
          activeAsset = a;
        }

        a.addEventListener("pick", async event => {
          const response = await fetch(a.path);

          if (!response.ok) {
            console.error(`Failed to fetch ${a.path}`);
            return null;
          }

          const text = await response.text();
          const doc = Search.#domParser.parseFromString(text, "text/html");
          const grabbed = doc.querySelector("#reading-pane");
          const replace = document.querySelector("#reading-pane");
          replace.replaceWith(grabbed);
          window.scrollTo({ top: 0, behavior: "instant" });

          let link = a.path; 
          if (window.location.search) link += window.location.search;
          if (window.location.hash) link += window.location.hash;
          history.pushState({}, "", link);
        });

        assets[a.path] = a;
      }

      const newsetFirst = Object.values(assets)
        .sort((a, b) => b.date.localeCompare(a.date));

      if (activeAsset === null) {
        activeAsset = newsetFirst[0];
        activeAsset.active = true;
      }

      this.#deck.replace(...newsetFirst);

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
      for(const [key, values] of Object.entries(filters)) {
        for (const value of Object.keys(values)) {
          const pin = document.createElement("rpr-pin");
          pin.key = key;
          pin.value = value;

          pin.addEventListener("click", async event => {
            if (await this.togglePin(pin)) {
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

      // for (const asset of Object.values(await this.#assets)) {
      //   const pin = document.createElement("rpr-pin");
      //   pin.key = "Asset";
      //   pin.value = asset.title;
      //   pin.onPinnedCallback = () => {
      //     asset.click();
      //     const elem = this.shadowRoot.querySelector("#pins");
      //     elem.reset(false);
      //     elem.reflowPins();
      //     return false;
      //   };
      //   pins.push(pin);
      // }

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

  static sortFiltersAlphabetically(a, b) {
    const key = a[0].localeCompare(b[0]);
    if (key !== 0) {
      return a;
    }
    const val = a[1].localeCompare(b[1]);
    return val;
  }

  static sanitiseFilters(newFilters) {
    const result = {};
    for (const key of Object.keys(newFilters).sort()) {
      result[key] = newFilters[key].sort();
    }
    return result;
  }

  static filtersMatch(a, b) {
    const aKeys = Object.keys(a);
    const bKeys = Object.keys(b);
    if (aKeys.length !== bKeys.length) {
      return false;
    }
    return aKeys.every((key, i) => {
      const aValues = a[key];
      const bValues = b[key];
      if (aValues.length !== bValues.length) {
        return false;
      }
      return aValues.every((val, i) => val === bValues[i]);
    });
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch(name) {
      case "query":
        this.#query = newValue;
        this.shadowRoot.querySelector("#input").value = newValue;
        if (newValue !== oldValue && this.interactive) {
          this.#search();
        }
        break;
      case "filters":
        this.#filters = Search.sanitiseFilters(JSON.parse(newValue));
        if (newValue !== oldValue && this.interactive) {
          if (!Search.filtersMatch(this.#filters, this.#prevFilters)) {
            this.#search();
          }
        }
        this.#prevFilters = structuredClone(this.#filters);
        break;
      case "interactive":
        this.#interactive = newValue === "true";
        break;
    }
  }

  async #search() {
    const searchCount = ++this.#searchCount;
    this.#deck.replace();

    // 🪟🔗 Update Window Location
    
    (async () => {
      if (!this.interactive) {
        return;
      }

      const params = new URLSearchParams();
      for (const [key, values] of Object.entries(await this.filters)) {
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
      const pinsComponent = this.shadowRoot.querySelector("#pins");
      const query = this.query.replaceAll('"', "");

      const pins = [];
      for (const pin of await this.#pins) {
        const isPinned = await this.isPinned(pin);
        pin.pinned = isPinned;
        pins.push({ 
          pin, 
          pinned: isPinned,
          order: -1
        });
      }

      const values = pins.map(pin => pin.pin.value);
      const [unordered, info, ordered] = fuzzy.search(values, query);

      if (unordered === null || unordered.length <= 0) {
        pinsComponent.replaceUnpinned();
        const pinned = pins.filter(pin => pin.pinned).map(pin => pin.pin);
        pinsComponent.replacePinned(...pinned);
        return;
      }

      for (const [orderedIndex, unorderedIndex] of ordered.entries()) {
        const pin = pins[unordered[unorderedIndex]];
        pin.pin.matches = info.ranges[unorderedIndex] || [];
        pin.order = orderedIndex;
      }

      if (this.interactive) {
        const unpinned = pins
          .filter(p => !p.pinned && p.order >= 0)
          .sort((a, b) => a.order - b.order)
          .map(m => m.pin);
        pinsComponent.replaceUnpinned(...unpinned);
      }

      const pinned = pins.filter(m => m.pinned).map(m => m.pin);
      pinsComponent.replacePinned(...pinned);
    })();


    // 📃 Filter and search assets

    if (this.query === "" && Search.filtersMatch(this.filters, {})) {
      this.#deck.replace(...await this.getDefaultAssets());
      return false;
    }

    const forceQuery = this.query || null;
    const sort = {};
    if (this.query === "") {
      sort["date"] = "desc";
    }
    const result = await (await pagefind).search(
      forceQuery,
      { filters: this.filters, sort }
    );

    if (searchCount !== this.#searchCount) {
      return false;
    }

    if (result === null) {
      this.#deck.replace(...await this.getDefaultAssets());
      return false;
    }

    const trailingSlashes = /\/+$/;
    let tabIndex = this.#pins.length + 10;
    let activeAsset = null;
    const assets = await this.#assets;
    const toElement = data => {
      const key = data.raw_url.replace(trailingSlashes, "");
      const asset = assets[key];
      asset.tabIndex = tabIndex++;

      if (asset.active) {
        activeAsset = asset;
      }

      if (typeof asset === "undefined") {
        return null;
      }

      asset.replaceResults(...data.sub_results);
      return asset;
    };

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
          this.#deck.append(asset);
        }
        resolve(asset);
      });
      await interrupt();
    }

    await Promise.all(promises);

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
    this.#deck.append(...remainingAssets);

    if (!this.interactive && activeAsset !== null) {
      activeAsset.scrollIntoView({
        block: "center",
        behavior: "instant",
      });
    }

    return true;
  }

  reset(refocus = false) {
    this.query = "";
    if (refocus) this.focus();
    document.getElementById("assets").reset();
  }

  recheck(refocus = false) {
    this.search(this.query);
    document.getElementById("pins").queryChanged(this.query);
    if (refocus) {
      this.focus();
      this.select();
    }
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
    this.setAttribute("interactive", newValue);
  }

  async getDefaultAssets() {
    const assets = await this.#assetsNewestFirst;

    for (const asset of assets) {
      asset.replaceResults();
    }

    return assets;
  }

  async togglePin(pin) {
    const isPinned = await this.isPinned(pin);
    isPinned ? await this.unpin(pin) : await this.pin(pin);
    return !isPinned;
  }

  async pin(pin) {
    const filters = await this.filters;
    let values = filters[pin.key] || [];
    values.push(pin.value);
    values = [...new Set(values)];
    filters[pin.key] = values.sort();
    this.filters = filters;
  }

  async unpin(pin) {
    const filters = await this.filters;
    const values = filters[pin.key] || [];
    const index = values.indexOf(pin.value);
    if (index >= 0) {
      values.splice(index, 1);
      if (values.length > 0) {
        filters[pin.key] = values;
      } else {
        delete filters[pin.key];
      }
      this.filters = filters;
    }
  }

  async isPinned(pin) {
    const filters = await this.filters;
    return filters.hasOwnProperty(pin.key) && filters[pin.key].includes(pin.value);
  }
});

