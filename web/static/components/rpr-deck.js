class Deck extends HTMLElement {
  #stageElement

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        #stage {
          display: flex;
          flex-direction: column;
          gap: 0.5rem;
        }

        #stage > *[score="-1"] {
          display: none;
        }
      </style>
      <div id="stage"></div>
    `;

    this.#stageElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#stage"));

    /**
      * @type {(event: Event) => void}
      * @param {PickEvent} event
      */
    function pickHandler(event) {
      /** @type {Asset|null} */
      const current = shadowRoot.querySelector(`[active="true"]`);

      const asset = event.detail.asset;

      if (current !== asset) {
        if (current !== null) {
          current.active = false;
        }
        asset.active = true;
        this.dispatchEvent(new CustomEvent("pick", {
          bubbles: true,
          detail: event.detail,
        }));
      }

    }

    shadowRoot.addEventListener("pick", pickHandler);
  }

  /** 
  * @param {...HTMLElement} elements
  */
  append(...elements) {
    this.#stageElement.append(...elements);
  }

  /** 
  * @param {...HTMLElement} elements
  */
  replace(...elements) {
    this.#stageElement.replaceChildren(...elements);
  }
}

customElements.define("rpr-deck", Deck);
