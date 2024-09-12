class Pin extends HTMLElement {

  /** @type {string} */
  #key

  /** @type {string} */
  #value

  /** @type {Boolean} */
  #pinned

  /** @type {[Number, Number][]} */
  #highlights = [];

  /** @type {HTMLElement} */
  #keyElement

  /** @type {HTMLElement} */
  #valueElement

  /** attributes that trigger {@link attributeChangedCallback}. */
  static observedAttributes = ["pinned", "matches", "value", "key"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          border: 1px solid #E2E8F0;
          border-radius: 1rem;
          padding: 0.25rem 0.75rem;
          position: relative;

          display: inline-flex;
          flex-direction: row;

          justify-content: left;
          justify-items: center;
          cursor: pointer;
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          font-size: 1rem;
          line-height: 1rem;
        }
        :host(:not([pinned="true"])) > #unpin {
          display: none;
        }
        :host([pinned="true"]) {
          background: #eee;
          border: 1px solid transparent;
        }
        #value {
          color: rgb(51, 65, 85);
          letter-spacing: 0.07em;
          font-size: 0.875rem;
          line-height: 0.9rem;
        }
        :host([pinned="true"]) #value {
          color: #666;
        }
        #key {
          color: rgb(148 163 184);
          margin-right: 0.5rem;
          font-size: 0.6rem;
          letter-spacing: 0.2em;
          line-height: 1.125rem;
          text-transform: uppercase;
        }
        #unpin {
          position: absolute;
          top: 0;
          right: 0;
          display: inline-flex;
          display: none;
          width: 1.25rem;
          height: 1.25rem;
          margin-left: 0.5rem;
          background: rgb(148, 163, 184);
          border-radius: 9999px;
          color: rgb(241, 245, 249);
          text-transform: uppercase;
          font-size: 0.7rem;
          line-height: 1rem;
          font-weight: 700;
          text-align: center;
          vertical-align: middle;
          justify-content: center;
          justify-items: center;
        }
      </style>
      <span id="key"><span class="inner">${this.#key}</span></span>
      <span id="value">${this.#value}</span>
      <span id="unpin"><span>X</span></span>
    `;

    this.#keyElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#key"));
    this.#valueElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#value"));
  }

  connectedCallback() {
  }

  set value(value) {
    this.setAttribute("value", value);
  }

  get value() {
    return this.getAttribute("value") || "";
  }

  set key(key) {
    this.setAttribute("key", key);
  }

  get key() {
    return this.getAttribute("key") || "";
  }

  set pinned(p) {
    this.setAttribute("pinned", p.toString());
  }

  get pinned() {
    return this.#pinned;
  }

  set highlights(newValue) {
    this.setAttribute("matches", newValue.flat(Infinity).join(","));
  }

  get highlights() {
    return this.#highlights;
  }

  hasHighlights() {
    return this.#highlights.length >= 1
  }


  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "pinned":
        this.#pinned = newValue === "true";
        this.dispatchEvent(new CustomEvent(this.pinned ? "pinned" : "unpinned", {
          bubbles: true,
          detail: {
            key: this.#key,
            value: this.#value
          },
        }));
        this.reflowValue();
        break;
      case "matches":
        const numbers = newValue.split(",").map(x => Number(x));

        /** @type {[Number, Number][]} */
        const matches = [];
        for (let i = 0; i + 1 < numbers.length; i += 2) {
          matches.push([numbers[i], numbers[i + 1]]);
        }
        this.#highlights = matches;
        this.reflowValue();
        break;
      case "value":
        this.#value = newValue;
        this.#valueElement.textContent = newValue;
        this.reflowValue();
        break;
      case "key":
        this.#key = newValue;
        this.#keyElement.textContent = this.key;
        this.reflowValue();
        break;
    }
  }

  reflowValue() {
    if (this.#pinned || !this.hasHighlights()) {
      this.#valueElement.textContent = this.value;
      return;
    }

    let pointer = 0;
    const fragment = new DocumentFragment();
    for (const [start, end] of this.highlights) {
      if (start > pointer) {
        fragment.append(this.#value.substring(pointer, start));
      }
      fragment.append((() => {
        const mark = document.createElement("mark");
        mark.append(this.#value.substring(start, end));
        return mark;
      })());
      pointer = end;
    }
    fragment.append(this.#value.substring(pointer));
    this.#valueElement.replaceChildren(fragment);
  }
}

customElements.define("rpr-pin", Pin);
