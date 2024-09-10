window.customElements.define("rpr-pins", class extends HTMLElement {
  #pinned = null;
  #unpinned = null;
  #tabIndexStart = 2;
  #tabIndexEnd = null;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
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
          gap: 0.25rem;
        }

        #pinned {
          display: flex;
          flex-direction: row;
          flex-wrap: wrap;
          justify-content: left;
          gap: 0.25rem;
        }
      </style>
      <div id="unpinned"></div>
      <div id="pinned"></div>
    `;
    this.#pinned = this.shadowRoot.querySelector("#pinned");
    this.#unpinned = this.shadowRoot.querySelector("#unpinned");
  }

  connectedCallback() {
  }

  replacePinned(...pins) {
    this.#pinned.replaceChildren(...pins);
    this.#recalculateTabIndexes();
  }

  replaceUnpinned(...pins) {
    this.#unpinned.replaceChildren(...pins);
    this.#recalculateTabIndexes();
  }

  #recalculateTabIndexes() {
    let tabIndex = this.#tabIndexStart;
    for (const pin of this.#unpinned.children) pin.tabIndex = tabIndex++;
    for (const pin of this.#pinned.children) pin.tabIndex = tabIndex++;
    this.#tabIndexEnd = tabIndex;
  }

  get tabIndexEnd() {
    return this.#tabIndexEnd;
  }
});
