class AdvancedSearch extends HTMLElement {

  /** @type {Object.<string, string[]>} */
  #filters = {};

  /** @type {Object.<string, Pin[]>} */
  #pins = {};

  static observedAttributes = ["filters"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          display: block;
          padding: 2rem;
        }
        h1, h2 {
          color: #bbb;
          font-weight: 400;
        }
        .pin-groups .pins {
          display: flex;
          flex-direction: row;
          flex-wrap: wrap;
          gap: 0.5rem;
        }
      </style>
      <h1>Advanced Search</h1>
      <div class="pin-groups">

        <div data-key="Medium">
          <h2 class="title"></h2>
          <p class="description">Ray Peat wrote articles, appeared on podcasts, and expressed himself in many mediums.</p>
          <div class="pins"></div>
        </div>

        <div data-key="State">
          <h2 class="title"></h2>
          <p class="description">Assets range from awaiting content, to AI audio transcripts, to verified correct text.</p>
          <div class="pins"></div>
        </div>

        <div data-key="Completion">
          <h2 class="title"></h2>
          <p class="description">Assets range from awaiting content, to AI audio transcripts, to verified correct text.</p>
          <div class="pins"></div>
        </div>

        <div data-key="Publisher">
          <h2 class="title"></h2>
          <p class="description">Ray appeared on various shows and published his own newsletter.</p>
          <div class="pins"></div>
        </div>

        <div data-key="Participant">
          <h2 class="title"></h2>
          <p class="description">Many people participated in conversations with Ray Peat.</p>
          <div class="pins"></div>
        </div>

        <div data-key="Issues">
          <h2 class="title"></h2>
          <p class="description">Opportunities to improve an asset's text are tracked in "GitHub issues".</p>
          <div class="pins"></div>
        </div>

        <div data-key="Mention">
          <h2 class="title"></h2>
          <p class="description">People and concepts mentioned by participants.</p>
          <div class="pins"></div>
        </div>

        <div data-key="Also On">
          <h2 class="title"></h2>
          <p class="description">Assets were gathered from the following websites.</p>
          <div class="pins"></div>
        </div>

      </div>
    `;
  }

  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "filters":
        this.#filters = JSON.parse(newValue);

        for (const [key, pins] of Object.entries(this.pins)) {
          const group = this.shadowRoot?.querySelector(`[data-key="${key}"]`);

          const titleElement = group?.querySelector(".title");
          if (titleElement) {
            titleElement.textContent = key;
          }

          const pinsElement = group?.querySelector(".pins");
          if (pinsElement) {
            pinsElement.replaceChildren(...pins);
          }
        }

        break;
    }
  }

  get filters() {
    return this.#filters;
  }

  set filters(newValue) {
    this.setAttribute("filters", JSON.stringify(newValue));
  }

  /**
   * @returns {Object.<string, Pin[]>}
   */
  get pins() {
    return this.#pins;
  }

  /** 
  * @param {Pin[]} newValue
  */
  set pins(newValue) {
    /** @type {Object.<string, Pin[]>} */
    const pins = {};
    for (const pin of newValue) {
      pins[pin.key] = pins[pin.key] || [];
      pins[pin.key].push(pin);
      pin.addEventListener("click", () => {
        pin.pinned = !pin.pinned; // just to demo interactivity
      });
    }
    this.#pins = pins;
  }
}

customElements.define("rpr-advanced-search", AdvancedSearch);
