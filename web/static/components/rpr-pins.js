class Pins extends HTMLElement {

  /** @type {Number} */
  #tabIndexStart = 2;

  /** @type {Number} */
  #tabIndexEnd = 2

  /** @type {Pin[]} */
  #pinned = [];

  /** @type {Pin[]} */
  #unpinned = [];

  /** @type {HTMLElement} */
  #pinnedElement

  /** @type {HTMLElement} */
  #unpinnedElement

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          display: flex;
          flex-direction: column;
          gap: 0.5rem;
        }

        #unpinned {
          display: flex;
          flex-direction: row;
          flex-wrap: wrap;
          justify-content: left; 
          gap: 0.5rem;
        }

        #pinned {
          display: flex;
          flex-direction: row;
          flex-wrap: wrap;
          justify-content: left;
          gap: 0.5rem;
        }
      </style>
      <div id="unpinned"></div>
      <div id="pinned"></div>
    `;

    this.#pinnedElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#pinned"));
    this.#unpinnedElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#unpinned"));
  }

  /**
  * @param {...Pin} pins
  */
  replacePinned(...pins) {
    this.#pinned = pins;
    this.#pinnedElement.replaceChildren(...pins);
    this.#recalculateTabIndexes();
  }

  /**
  * @param {...Pin} pins
  */
  replaceUnpinned(...pins) {
    this.#unpinned = pins;
    this.#unpinnedElement.replaceChildren(...pins);
    this.#recalculateTabIndexes();
  }

  #recalculateTabIndexes() {
    let tabIndex = this.#tabIndexStart;
    for (const pin of this.#unpinned) pin.tabIndex = tabIndex++;
    for (const pin of this.#pinned) pin.tabIndex = tabIndex++;
    this.#tabIndexEnd = tabIndex;
  }

  get tabIndexEnd() {
    return this.#tabIndexEnd;
  }
}

customElements.define("rpr-pins", Pins);
