import { LitElement, css, html } from "lit";
import { activeFiltersContext, availableFiltersContext } from "./app-root.js";
import { ContextConsumer } from "@lit/context";

export class Filter extends LitElement {
  static properties = {
    key: { type: String, reflect: true },
    value: { type: String, reflect: true },
    hideKey: { type: Boolean, reflect: true },
    hideCount: { type: Boolean, reflect: true },
    active: { type: Boolean, reflect: true },
    count: { type: Number, reflect: true },
    pertinant: { type: Boolean, reflect: true },
  };

  static styles = css`
    :host(*) {
      position: relative;

      display: inline-flex;
      flex-direction: row;

      justify-content: left;
      gap: 0.5rem;
      cursor: pointer;
      font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
    }
    .wrapper {
      cursor: pointer;
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
    }
    .key {
      color: rgb(148 163 184);
      text-transform: capitalize;
      font-size: var(--font-size-xs);
      line-height: var(--line-height-xs);
    }
    .checkbox {
      cursor: pointer;
      margin: 0;
    }
    .value {
      color: var(--slate-600);
      cursor: pointer;
      word-break: break-all;
    }
    .count {
      color: var(--slate-400);
      text-align: right;
    }
  `;

  _handleChange(event) {
    this.active = event.target.checked;
    this.dispatchEvent(new FilterClickEvent(this.key, this.value, this.active));
  }

  render() {
    const filterId = `${this.key}/${this.value}`;
    const activeFilters = this.activeFilters.value || [];
    const active = activeFilters.some(f => f.key === this.key && f.value === this.value);

    this.count = (() => {
      if (!this.availableFilters.value) return this.count;
      const filters = this.availableFilters.value;
      if (!filters.hasOwnProperty(this.key)) return 0;
      if (!filters[this.key].hasOwnProperty(this.value)) return 0;
      return filters[this.key][this.value];
    })();

    this.active = active;
    this.pertinant = this.count > 0 || active;

    return html`
      <label for="${filterId}" part="wrapper" class="wrapper">
        ${this.hideKey
        ? html``
        : html`<span part="key" class="key"><span class="inner">${this.key}</span></span>`
      }

        <input 
          class="checkbox"
          type="checkbox" 
          id="${filterId}" 
          name="${filterId}" 
          part="checkbox"
          @change=${this._handleChange}
          .checked="${active}"
        >

        <span part="value" class="value">${this.value}</span>

        ${this.hideCount
        ? html``
        : html`<span part="count" class="count">${this.count}</span>`
      }
      </label>
    `;
  }

  /**
   * @param {string} key 
   * @param {string} value 
   * @param {Number} count
   * @param {Boolean} hideKey
   * @param {Boolean} hideCount
   */
  constructor(key, value, count, hideKey = false, hideCount = false) {
    super();

    this.key = key;
    this.value = value;
    this.count = count;
    this.hideKey = hideKey;
    this.hideCount = hideCount;

    this.activeFilters = new ContextConsumer(this, {
      context: activeFiltersContext,
      subscribe: true,
    });

    this.availableFilters = new ContextConsumer(this, {
      context: availableFiltersContext,
      subscribe: true,
    });
  }
}

customElements.define("rpr-filter", Filter);

export class FilterClickEvent extends Event {
  /**
  * @param {string} key
  * @param {string} value
  * @param {Boolean|null} force
  */
  constructor(key, value, force = null) {
    super("filter-click", { bubbles: true, composed: true });
    this.key = key;
    this.value = value;
    this.force = force;
  }
}

