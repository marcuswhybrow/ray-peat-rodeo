import { LitElement, html, css } from "lit";
import { unsafeHTML } from 'lit/directives/unsafe-html.js';
import { createRef, ref } from "lit/directives/ref.js";
import { ContextProvider, createContext } from "@lit/context";
import { FILTERS } from "/derived/filters.js";
import ASSETS from "/derived/thin-assets.js";
import * as pagefind from "/pagefind/pagefind.js";
import { ResultClickEvent } from "../events.js";
import { isVisible, loadResultData, pagefindFilters } from "../utils.js";

/** @type {import("@lit/context").Context<"activeFilters", Filter[]>} */
export const activeFiltersContext = createContext("activeFilters");

/** @type {import("@lit/context").Context<"availableFilters", PagefindFilters>} */
export const availableFiltersContext = createContext("availableFilters");

export class AppRoot extends LitElement {
  /** @type {Result[]} */
  result

  static properties = {
    _results: { state: true },
    // _allFilters: { state: true },

    activeFilters: {
      type: Array,
      state: true,

      /** @type {(newVal:Filter[], oldVal:Filter[]) => Boolean} */
      hasChanged: (newVal, oldVal) => {
        if (!oldVal) return true;
        if (newVal.length !== oldVal.length) return true;
        const alphaSort = (a, b) => `${a.key}${a.value}`.localeCompare(`${b.key}${b.value}`);
        newVal.sort(alphaSort);
        oldVal.sort(alphaSort);
        const itemMissmatch = (filter, i) => filter.key !== oldVal[i].key || filter.value !== oldVal[i].value;
        return newVal.some(itemMissmatch);
      },
    },

    _activeResultSlug: { state: true },
    _activeResultSectionId: { state: true },

    _search: { state: true },
    results: { state: true },
  };

  /**
   * Get the active filters in the form pagefind expects.
   *
   * @returns {Object.<String, String[]>}
   */
  get pagefindFilters() {
    return pagefindFilters(this.activeFilters);
  }

  constructor() {
    super();

    /** @type {PagefindResponse|null} */
    this.pagefindResponse = null;

    /** 
     * This is a hack to build a reverse lookup table from Pagefind's result 
     * url to result id.
     *
     * This won't be necessary once issue 371 is fixed.
     * - https://github.com/CloudCannon/pagefind/issues/371
     *
     * @type {Promise<[url: string, id: string][]>} 
     */
    this.pagefindReverseLookup = new Promise(async resolve => {
      /** @type {PagefindResponse} */
      const response = await pagefind.search(null);

      resolve(await Promise.all(response.results.map(async result => {
        await new Promise(r => setZeroTimeout(r, 0));
        const data = await result.data();
        return [data.url, result.id];
      })));
    });

    this.initialising = true;
    this.searchInputSet = false;
    this.allowHistoryChanges = false;
    this.activeFiltersProvider = new ContextProvider(this, { context: activeFiltersContext });
    this.availableFiltersProvider = new ContextProvider(this, { context: availableFiltersContext });
    this.showInitResults = true;

    /** @type {import("lit-html/directives/ref.js").Ref<HTMLInputElement>} */
    this.searchInput = createRef();


    /** @type {PagefindFilters} */
    this._allFilters = FILTERS;
    this.availableFiltersProvider.setValue(FILTERS);

    /** @type {Result[]} */
    this.results = ASSETS;

    pagefind.filters(); // Required to receive filters henceforth

    const params = new URLSearchParams(window.location.search);
    const activeFilters = [];

    for (const [key, values] of Object.entries(this._allFilters)) {
      for (const [value, _count] of Object.entries(values)) {
        if (params.has(key, value)) {
          activeFilters.push({ key, value });
        }
      }
    }

    this._search = params.get("search") || "";

    this.activeFilters = activeFilters;
    this.activeFiltersProvider.setValue(activeFilters);

    this.initialising = false;

    window.addEventListener("filter-click", async event => {
      const activeFilters = structuredClone(this.activeFilters);
      const index = activeFilters.findIndex(f => f.key === event.key && f.value === event.value);
      if (index >= 0) {
        if (event.force === null) {
          activeFilters.push({ key: event.key, value: event.value });
        } else if (event.force === false) {
          activeFilters.splice(index, 1);
        }
      } else if (event.force === null || event.force === true) {
        activeFilters.push({ key: event.key, value: event.value });
      }

      this.activeFilters = activeFilters
    });

    window.addEventListener("result-click", async event => {
      const url = new URL(window.location.href);
      const currentSlug = url.pathname.replace(/^\/+(.*?)\/*$/, "$1");
      const currentSection = url.hash.substring(1);

      // Update Window Location

      url.pathname = `/${event.slug}/`; // trailing slash prevents 301
      url.hash = event.section || "";

      if (currentSlug !== event.slug || currentSection !== event.section) {
        history.pushState({}, "", url.href);
      }

      // Highlight Result In List

      this.shadowRoot
        ?.querySelectorAll(`.result[data-slug="${event.slug}"] .header`)
        .forEach(element => element.ariaCurrent = "true");

      this.shadowRoot
        ?.querySelectorAll(`.result:not([data-slug="${event.slug}"]) .header`)
        .forEach(element => element.ariaCurrent = "false");


      // Replace Content

      if (currentSlug !== event.slug) {
        const response = await fetch(`/${event.slug}/partial.html`);
        if (!response.ok) throw new Error(`Failed to fetch result content "${url.pathname}".`);
        this.innerHTML = await response.text();
      }

      // Scroll to Section or Top

      if (event.section) {
        const SCROLL_OFFSET = 100;
        const section = this.querySelector(`[id="${event.section}"]`);
        if (!section) throw new Error(`Failed to find section "${event.section}" in result content "${url.pathname}".`);
        window.scrollTo({
          top: section.getBoundingClientRect().top + window.scrollY - SCROLL_OFFSET,
        });
      } else {
        window.scrollTo({ top: 0 });
      }
    });

    window.addEventListener("popstate", event => {
    });

    window.addEventListener("keydown", event => {
      const inputEl = /** @type {HTMLInputElement|null} */ (this.shadowRoot?.querySelector(".search input"));
      if (inputEl) {
        if (event.key === "k" && event.ctrlKey) {
          inputEl.focus();
          inputEl.select();
          event.preventDefault();
          event.stopPropagation();
        } else if (event.key === "Escape") {
          inputEl.blur();
        }
      }
    });
  }

  connectedCallback() {
    super.connectedCallback();

    /** @type {IntersectionObserverCallback} */
    const intersectionHandler = entries => {
      if (this.pagefindResponse) {
        const ids = entries.map(entry => {
          const htmlElement = /** @type {HTMLElement} */ (entry.target);
          return htmlElement.dataset.pagefindResultId || "";
        });

        loadResultData(this.results, this.pagefindResponse.results, ids)
          .then(results => this.results = results);
      }
    };

    this.observer = new IntersectionObserver(intersectionHandler, {
      root: null,
      rootMargin: "200px",
      threshold: 0,
    });

    this.observer?.takeRecords
  }

  /** @param {import("lit").PropertyValues<this>} changedProperties */
  async willUpdate(changedProperties) {
    if (changedProperties.has("activeFilters")) {
      this.activeFiltersProvider.setValue(this.activeFilters);
    }

    const searchChanged = changedProperties.has("_search");
    const searchPrev = changedProperties.get("_search");
    const searchIsInit = typeof searchPrev === "undefined";
    const dirtySearch = searchChanged && (!searchIsInit || this._search !== "");

    const activeFiltersChanges = changedProperties.has("activeFilters");
    const activeFiltersPrev = changedProperties.get("activeFilters");
    const activeFiltersIsInit = typeof activeFiltersPrev === "undefined";
    const dirtyActiveFilters = activeFiltersChanges && (!activeFiltersIsInit || this.activeFilters.length > 0);

    const performSearch = dirtySearch || dirtyActiveFilters;

    // console.log("app-root willUpdate", changedProperties, {
    //   buffering: this.buffering,
    //   _search: this._search,
    //   results: this.results,
    //   activeFilters: this.activeFilters,
    //   _allFilters: this._allFilters,
    // });

    if (performSearch) {
      /** @type {PagefindResponse} */
      const response = await pagefind.search(this._search || null, {
        filters: pagefindFilters(this.activeFilters || []),

        // Sort by date if no search text is provided, otherwise sort by weight (default);
        sort: this._search ? {} : { date: "desc" },
      });

      if (response !== null) {
        this.pagefindResponse = response;
        this.results = await Promise.all(this.results.map(async result => {
          result.loaded = false;

          // Hack until Pagefind issue 371 is fixed.
          if (!result.pagefindResultId) {
            const lookup = await this.pagefindReverseLookup;
            const entry = lookup.find(l => l[0] === `/${result.slug}`);
            if (!entry) return result;
            result.pagefindResultId = entry[1];
          }

          const index = response.results.findIndex(pr => pr.id === result.pagefindResultId);
          result.score = response.results[index]?.score;
          result.order = index;
          return result;
        }));
        this._allFilters = response.filters;
        this.availableFiltersProvider.setValue(response.filters);
        this.showInitResults = false;

        const all = Array.from(document.querySelectorAll(".results .result"));

        // Maybe be unessary after every search, requires investigation
        all.forEach(element => {
          this.observer?.observe(element);
        });

        const ids = all.filter(e => isVisible(e, 100)).map(element => {
          const htmlElement = /** @type {HTMLElement} */ (element);
          return htmlElement.dataset.pagefindResultId || "";
        });

        this.results = await loadResultData(this.results, response.results, ids);
      }
    }
  }

  render() {
    const url = new URL(window.location.href);
    if (this._search) {
      url.searchParams.set("search", this._search);
    } else {
      url.searchParams.delete("search");
    }

    /** 
     * @param {String} key 
     * @param {Object.<String, Number>} values
     */
    const renderFilter = (key, values) => {
      return Object.entries(values).map(([value, count]) => {
        const active = (this.activeFilters || []).some(f => f.key === key && f.value === value);
        const extant = url.searchParams.has(key, value);

        if (active && !extant) {
          url.searchParams.append(key, value);
        } else if (!active && extant) {
          url.searchParams.delete(key, value);
        }

        return html`
          <li>
            <rpr-filter .key="${key}" .value="${value}" count="${count}" hidekey></rpr-filter>
          </li>
        `;
      });
    };

    /** @param {Event} event */
    const toggleFilterMode = event => {
      const elem = /** @type {HTMLElement} */ (event.target);
      if (elem.textContent === "or") {
        elem.textContent = "and";
      } else {
        elem.textContent = "or";
      }
    };

    /** @type {import("lit-html").TemplateResult[]} */
    const filters = Object.entries(this._allFilters).map(([key, values]) => {
      return html`
        <li class="group">
          <div class="header">
            <span class="name">${key}</span>
            <span class="filter-mode" @click="${toggleFilterMode}">or</span>
          </div>
          <ul class="filters">
            ${renderFilter(key, values)}
          </ul>
        </li>
      `;
    });

    if (this.allowHistoryChanges) {
      history.replaceState({}, "", url.href);
    }

    if (this.initialising === false) {
      this.allowHistoryChanges = true;
    }

    /** @type {(event:MouseEvent) => void} */
    const resetFilters = event => {
      this.activeFilters = [];
      this._search = "";
      if (this.searchInput.value) this.searchInput.value.value = "";
      event.preventDefault();
      event.stopPropagation();
    };

    const resetFiltersButton = (() => {
      if (filters.length > 0) {
        return html`
          <li class="reset">
            <a href="" @click="${resetFilters}">Reset Filters</a>
          </li>
        `;
      }
    })();

    const renderIssues = this.activeFilters?.some(f => f.key === "issues" && f.value === "Has Issues") || false;

    /** @type {import("lit-html").TemplateResult[]} */
    const results = this.results.map(result => {
      const slugFromUrl = window.location.pathname.replace(/^\/+(.*?)\/$$/, "$1");
      const current = result.slug === slugFromUrl;

      /** @type {(event:Event) => void} */
      const onClick = event => {
        event.preventDefault();
        event.stopPropagation();
        event.target?.dispatchEvent(new ResultClickEvent(result.slug));
      };

      return html`
        <li 
          class="result" 
          data-slug="${result.slug}" 
          data-score="${result.score}"
          data-pagefind-result-id="${result.pagefindResultId}"
          ?data-loaded="${result.loaded}"
          style="order:${result.order}"
        >
          <a 
            class="header" 
            href="/${result.slug}/" 
            aria-current="${current ? "true" : "false"}"
            @click=${onClick}
          >
            <h3 class="title">${result.title}</h3> 
            <p class="details">${result.date} ${result.publisher}</p>
          </a>

          ${result.excerpt ? html`<div class="excerpt">${unsafeHTML(result.excerpt)}</div>` : ""}

          ${renderIssues && result.issues
          ? html`<ol class="issues">${issuesToHtml(result.slug, result.issues)}</ol>`
          : ""}

          <ol class="sections">
            <div class="highlight" data-slug="${result.slug}"></div>
            <div class="overlight" data-slug="${result.slug}"></div>
            ${sectionsToHtml(result.slug, result.sections, renderIssues)}
          </ol>
        </li>
      `;
    });

    /** @type {(event:InputEvent) => void} */
    const searchKeyUp = event => {
      const inputEl = /** @type {HTMLInputElement} */ (this.searchInput.value);
      this._search = inputEl.value;
    };

    return html`
      <div class="search">
        <div class="search-input">
          <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 30 30">
<path d="M 13 3 C 7.4889971 3 3 7.4889971 3 13 C 3 18.511003 7.4889971 23 13 23 C 15.396508 23 17.597385 22.148986 19.322266 20.736328 L 25.292969 26.707031 A 1.0001 1.0001 0 1 0 26.707031 25.292969 L 20.736328 19.322266 C 22.148986 17.597385 23 15.396508 23 13 C 23 7.4889971 18.511003 3 13 3 z M 13 5 C 17.430123 5 21 8.5698774 21 13 C 21 17.430123 17.430123 21 13 21 C 8.5698774 21 5 17.430123 5 13 C 5 8.5698774 8.5698774 5 13 5 z"></path>
</svg>
          <input 
            type="search" 
            placeholder="Search..." 
            @keyup="${searchKeyUp}" 
            ${ref(this.searchInput)} 
          />
          <kbd>Ctl K</kbd>
        </div>
      </div>

      <div class="all-filters">
        <a href="/" class="home">Ray Peat Rodeo</a>
        <ul>
          ${filters}
          <p style="grid-column:1 / -1">Advanced Search</p>
          <textarea style="grid-column: 1 / -1; min-height: 8rem; resize: vertical"></textarea>
          ${resetFiltersButton}
        </ul>
      </div>

      ${this.results.length > 0 ? html`<ol class="results">${results}</ol>` : ""}
      ${this.showInitResults ? html`<slot name="results"></slot>` : ""}

      <div class="main">
        <div class="content">
          <slot name="asset"></slot>
        </div>
      </div>
    `;
  }

  updated() {
    this.shadowRoot?.querySelectorAll(".results .result").forEach(result => {
      this.observer?.observe(result);
    });

    if (!this.searchInputSet && this.searchInput.value) {
      this.searchInput.value.value = this._search;
      this.searchInputSet = true;
    }
  }

  static styles = css`
    :host(*) {
      font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
      --sidebar: 20rem;
    }

    .search {
      border-bottom: 1px solid var(--slate-200);
      padding: 1rem 2rem;
      top: 0;
      position: sticky;
      z-index: 50;
      background: rgba(255, 255, 255, 0.8);
      -webkit-backdrop-filter: blur(10px);
      backdrop-filter: blur(10px);
      margin-left: calc(var(--sidebar) * 2 + 2px);
    }

    .search .search-input {
      position: relative;
      width: 20rem;
      display: flex;
      flex-direction: row;
      transition: all .05s;
    }
    .search .search-input:has(input:focus),
    .search .search-input:has(input:not(:placeholder-shown)) {
      width: 25rem;
    }

    .search .search-input input {
      padding: 0.5rem 1rem 0.5rem 2.5rem;
      width: 100%;
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      letter-spacing: 0.025em;
      border: none;
      border: 1px solid var(--slate-200);
      border-radius: 9999px;
      background: white;
    }
    .search .search-input input:placeholder {
      opacity: 1;
      color: var(--slate-200);
    }

    .search .search-input kbd {
      position: absolute;
      right: 1rem;
      top: 0;
      bottom: 0;
      display: grid;
      align-items: center;
      color: var(--slate-400);
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
    }

    .search .search-input svg {
      position: absolute;
      top: 0.55rem;
      left: 0.5rem;

      height: 1.25rem;
      width: 1.25rem;
      opacity: 0.5;
    }

    .all-filters {
      position: fixed;
      top: 0;
      left: 0;
      width: var(--sidebar);
      height: 100vh;
      overflow-y: scroll;
      border-right: 1px solid var(--slate-200);
    }

    .results {
      position: fixed;
      top: 0;
      left: calc(var(--sidebar) + 1px);
      width: var(--sidebar);
      height: 100vh;
      overflow-y: scroll;

      border-right: 1px solid var(--slate-200);
    }

    .main {
      margin-left: calc(var(--sidebar) * 2 + 2px);
      position: relative;
    }
    .main::before {
      content: "";
      position: absolute;
      background: radial-gradient(at top, var(--pink-50) 0%, transparent 70%);
      top: -8rem;
      left: 0;
      right: 30%;
      height: 60rem;
      z-index: -10;
    }
    .main::after {
      content: "";
      position: absolute;
      background: radial-gradient(at top, var(--purple-50) 0%, transparent 70%);
      top: -8rem;
      left: 20%;
      right: 20%;
      height: 30rem;
      z-index: -11;
    }


    .main .content {
      width: 40rem;
      margin: 0 auto;
    }

    @media (max-width: 1330px) {
      .main .content {
        width: auto;
        min-width: 20rem;
        margin: 0;
        padding-right: 2rem;
        padding-left: 2rem;
      }
    }

    .all-filters {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    @media (max-width: 1330px) {
      :host {
        --sidebar: 15rem;
      }
    }

    .all-filters a.home {
      font-size: var(--font-size-3xl);
      line-height: var(--line-height-3xl);
      margin: 2rem 2rem 1rem 2rem;
      text-decoration: none;
      font-weight: 700;
      color: var(--slate-900);
      letter-spacing: -0.05em;
    }


    .all-filters ul {
      margin: 0;
      padding: 0;
      list-style: none;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .all-filters > ul {
      padding: 0rem 2rem 2rem 2rem;
      display: grid;
      grid-template-columns: min-content auto min-content;
    }

    .all-filters {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .all-filters > ul:not(:has(rpr-filter[pertinant])):before {
      content: "No results";
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      color: var(--slate-400);
      grid-column: 1 / -1;
    }

    .all-filters > ul:not(:has(rpr-filter[pertinant])) .group {
      display: none;
    }

    .all-filters > ul .reset {
      display: none;
    }

    .all-filters > ul:not(:has(rpr-filter[pertinant])) .reset {
      display: block;
    }
    .all-filters .reset a {
      display: inline-block;
      text-decoration: none;
      background: var(--slate-500);
      color: white;
      font-weight: 400;
      border-radius: 9999px;
      padding: 0.5rem 1rem;
      letter-spacing: 0.025em;
      font-size: var(--font-size-md);
      line-height: var(--line-height-md);
      transition: all .1s;
    }
    .all-filters .reset a:hover {
      box-shadow: 0 5px 30px rgba(0,0,0,0.3);
      transform: translateY(3px);
    }

    .all-filters .group .filters:not(:has(rpr-filter[pertinant])):after {
      content: "No results";
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      color: var(--slate-400);
      grid-column: 1 / -1;
    }

    .all-filters .group {
      width: 100%;
      display: grid;
      grid-template-columns: subgrid;
      grid-column: 1 / -1;
    }

    .all-filters .group .header {
      grid-column: 1 / -1;
      display: grid;
      grid-template-columns: auto max-content;
      padding: 0.5rem 0;
    }

    .all-filters .group .header .name {
      grid-column: 1;
      font-size: var(--font-size-md);
      line-height: var(--line-height-md);
      color: var(--slate-900);
      text-transform: capitalize;
      letter-spacing: -0.025em;
    }

    .all-filters .group .header .filter-mode {
      grid-column: 2;
      color: var(--slate-400);
      background: var(--slate-100);
      border-radius: 9999px;
      font-size: var(--font-size-xs); 
      line-height: var(--line-height-xs);
      text-transform: uppercase;
      letter-spacing: 0.05em;
      padding: 0 0.5rem;
      display: grid;
      align-items: center;
      cursor: pointer;
      transform: translateX(0.5rem);
      user-select: none;
    }
    .all-filters .group .header .filter-mode:hover {
      background: var(--slate-200);
      color: var(--slate-500);
    }

    .all-filters .group ul.filters {
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      display: grid;
      grid-template-columns: subgrid;
      grid-column: 1 / -1;
      row-gap: 0;
      column-gap: 0.5rem;
    }

    .all-filters .group ul.filters li {
      display: grid;
      grid-template-columns: subgrid;
      grid-column: 1 / -1;
    }

    .all-filters .group .filters rpr-filter {
      margin-bottom: 0.5rem;
    }

    .all-filters rpr-filter {
      display: grid;
      grid-template-columns: subgrid;
      grid-column: 1 / -1;
      align-items: start;
    }

    .all-filters rpr-filter::part(wrapper) {
      display: grid;
      grid-template-columns: subgrid;
      grid-column: 1 / -1;
      align-items: start;
    }
    .all-filters rpr-filter::part(checkbox) {
      margin-top: 0.3rem;
      grid-column: 1;
    }
    .all-filters rpr-filter::part(value) {
      grid-column: 2;
    }
    .all-filters rpr-filter::part(count) {
      grid-column: 3;
    }

    .all-filters rpr-filter[count="0"]:not([active]) {
      display: none;
    }

    .results {
      margin: 0;
      padding: 0;

      display: flex;
      flex-direction: column;
      list-style: none;
    }

    .results .result {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      padding: 1rem 2rem;
    }

    .results .result[data-score=""] {
      display: none;
    }

    .results .issues {
      display: grid;
      grid-template-columns: min-content auto;
      column-gap: 0.5rem;
      row-gap: 0.5rem;
    }
    .results .issues .issue {
      grid-column: 1 / -1;
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
    }
    .results .section .issues .issue {
      border-left: 1px solid var(--amber-100);
      padding-left: 1rem;
    }
    .results .section .issues .issue:last-child {
      padding-bottom: 1rem;
    }
    .results .issues .id {
      color: var(--amber-600);
      font-weight: 700;
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      margin-right: 0.25rem;
      letter-spacing: -0.025em;
    }
    .results .issues .id:before {
      content: "#";
    }
    .results .issues .title {
      display: inline;
      border: 0;
      padding: 0;
      color: var(--amber-600);
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
    }

    .results .result .excerpt:not(:empty) {
      font-size: var(--font-size-xs);
      line-height: var(--line-height-xs);
      padding-bottom: 0.5rem;
    }

    .results .result.hide:not(:has([aria-current="true"])) {
      display: none;
    }

    .results .result a.header {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      color: var(--slate-800);
      text-decoration: none;
    }
    .results .result:has(a.header[aria-current]:not([aria-current="false"])) {
      background: var(--slate-50);
    }

    .results .result:has(a.header[aria-current]:not([aria-current="false"])) ol a {
      border-left-color: var(--slate-200);

    }

    .results .result .header .title {
      font-size: var(--font-size-md);
      line-height: var(--line-height-md);
      margin: 0;
      font-weight: 400;
      letter-spacing: 0.025em;
    }

    .results .result .header .details {
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      margin: 0;
    }

    .results .result .sections {
      position: relative;
    }

    .results .result .sections .highlight {
      position: absolute;
      top: 0;
      left: 0;
      right: -1rem;
      height: 0;
      z-index: 10;
      border-top-right-radius: 0.25rem;
      border-bottom-right-radius: 0.25rem;
      background: var(--gray-100);
    }

    .results .result:has(a.header:not([aria-current="true"])) .highlight {
      display: none;
    }

    .results .result .sections .overlight {
      position: absolute;
      top: 0;
      left: 0;
      width: 1px;
      height: 0;
      z-index: 30;
      background: var(--pink-600);
      box-shadow: 0 0 5px var(--pink-400);
    }

    .results .result:has(a.header:not([aria-current="true"])) .overlight {
      display: none;
    }

    .results .result .sections .section {
      z-index: 20;
      position: relative;
    }

    .results .result .sections .section .excerpt {
      font-size: var(--font-size-xs);
      line-height: var(--line-height-xs);
      border-left: 1px solid var(--slate-100);
      padding-left: 1rem;
    }
    .results .result:has(a[aria-current="true"]) .sections .section .excerpt {
      border-left: 1px solid var(--slate-200);
    }

    .results .result .sections .section[data-depth="3"] .excerpt { padding-left: 2rem; }
    .results .result .sections .section[data-depth="4"] .excerpt { padding-left: 3rem; }
    .results .result .sections .section[data-depth="5"] .excerpt { padding-left: 4rem; }
    .results .result .sections .section[data-depth="6"] .excerpt { padding-left: 5rem; }

    .results .result ol {
      list-style: none;
      margin: 0;
      padding: 0;
    }
    .results .result ol a {
      display: block;
      border-left: 1px solid var(--slate-100);
      padding: 0.3rem 0.5rem 0.3rem 0.5rem;
      color: var(--pink-600);
      text-decoration: none;
      font-size: var(--font-size-md);
      line-height: var(--line-height-sm);
      letter-spacing: -0.025em;
    }
    .results .result ol a:hover {
      text-decoration: underline;
    }
    .results .result ol a.depth-2 { padding-left: 1rem; }
    .results .result ol a.depth-3 { padding-left: 2rem; }
    .results .result ol a.depth-4 { padding-left: 3rem; }
    .results .result ol a.depth-5 { padding-left: 4rem; }
    .results .result ol a.depth-6 { padding-left: 5rem; }

    .results .result .section:has(> a[aria-expanded="false"]) .subsections {
      display: none;
    }
  `;
}

customElements.define("app-root", AppRoot);


/** 
 * @param {string} slug
 * @param {Issue[]} issues 
 * @returns {import("lit-html").TemplateResult[]}
 */
function issuesToHtml(slug, issues) {
  return issues.map(issue => {
    const anchorId = `issue-${issue.id}`;
    /** @type {(event:Event) => void} */
    const click = event => {
      event.preventDefault();
      event.stopPropagation();
      event.target?.dispatchEvent(new ResultClickEvent(slug, anchorId));
    };

    return html`
      <li class="issue">
        <span class="id">${issue.id}</span>
        <a 
          class="title" 
          href="/${slug}/#${anchorId}"
          @click="${click}"
        >${issue.title}</a>
      </li>
    `;
  });
}

/**
 * @param {String} slug
 * @param {Section[]} sections
 * @param {Boolean} renderIssues
 * @returns {import("lit-html").TemplateResult[]}
 */
function sectionsToHtml(slug, sections, renderIssues) {
  if (!sections) return [];
  return sections.map(section => {
    /** @type {(event:MouseEvent) => void} */
    const onClick = event => {
      event.stopPropagation();
      event.preventDefault();
      event.target?.dispatchEvent(new ResultClickEvent(slug, section.id));
    };

    return html`
      <li class="section" data-id="${section.id}" data-depth="${section.depth}">
        <a 
          class="link" 
          href="/${slug}/#${section.id}" 
          class="depth-${section.depth}"
          ariaExpanded="false"
          @click="${onClick}"
        >
          <span class="title">${unsafeHTML(section.title)}</span>
        </a>
        ${section.excerpt
        ? html`<div class="excerpt">${unsafeHTML(section.excerpt)}</div>`
        : ""}

        ${renderIssues && section.issues
        ? html`<ol class="issues">${issuesToHtml(slug, section.issues)}</ol>`
        : ""}

        <ol class="subsections">
          ${sectionsToHtml(slug, section.subsections, renderIssues)}
        </ol>
      </li>
    `;
  });
}

// /**
// * Like Element.closest, but returns all matching ancestors, not just the 
// * closest.
// *
// * @param {Element} element The element to start the search
// * @param {string} selector Returned elements must match this CSS selector
// * @returns {Element[]}
// */
// function ancestors(element, selector) {
//   const results = [];

//   /** @type {Element|null|undefined} */
//   let currentElement = element.parentElement?.closest(selector);

//   while (currentElement) {
//     results.push(currentElement);
//     currentElement = currentElement.parentElement?.closest(selector);
//   }

//   return results;
// }
