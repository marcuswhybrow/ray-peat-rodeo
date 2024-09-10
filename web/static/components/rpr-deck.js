window.customElements.define("rpr-deck", class Deck extends HTMLElement {
  #cards = null;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
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

    this.shadowRoot.addEventListener("pick", event => {
      const current = this.shadowRoot.querySelector(`[active="true"]`);
      if (current !== event.detail) {
        if (current !== null) {
          current.active = false;
        }
        event.detail.active = true;
        this.dispatchEvent(new CustomEvent("pick", {
          bubbles: true,
          detail: event.detail,
        }));
      }
    });
  }

  append(...elements) {
    this.shadowRoot.querySelector("#stage").append(...elements);
  }

  replace(...elements) {
    this.shadowRoot.querySelector("#stage").replaceChildren(...elements);
  }

  get cards() {
    return this.#cards;
  }

  set cards(newValue) {
    this.#cards = newValue;
  }
});
