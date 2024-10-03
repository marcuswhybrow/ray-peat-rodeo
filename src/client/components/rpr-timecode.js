class Timecode extends HTMLElement {
  /** @type {string} */
  #externalUrl

  /** @type {string} */
  #time

  /** @type {HTMLAnchorElement} */
  #linkElement

  static observedAttributes = ["external-url", "time"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        #timecode {
          text-align: right;
        }
        #timecode #link {
          font-size: 0.875rem;
          line-height: 1.25rem;
          padding: 0.25rem 0.5rem;
          border-radius: 0.5rem;
          text-decoration: none;
        }
        :host([primary="true"]) #timecode #link {
          background: rgb(209, 213, 219);
          color: rgb(249, 250, 251);
        }
        :host([primary="true"]) #timecode #link:hover {
          background: rgb(107, 114, 128);
        }
        :host(:not([primary="true"])) #timecode #link {
          background: rgb(125, 211, 252);
          color: rgb(240, 249, 255);
        }
        :host(:not([primary="true"])) #timecode #link:hover {
          background: rgb(14, 165, 233);
        }
      </style>
      <span id="timecode" data-pagefind-ignore>
        <a id="link" href=""></a>
      </span>
    `;

    this.#linkElement = /** @type {HTMLAnchorElement} */ (shadowRoot.querySelector("#link"));
  }

  /** 
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "external-url":
        this.#externalUrl = newValue;
        this.#linkElement.href = newValue;
        break;
      case "time":
        this.#time = newValue;
        this.#linkElement.textContent = newValue;
        break;
    }
  }

  get externalUrl() {
    return this.#externalUrl;
  }

  set externalUrl(newValue) {
    this.setAttribute("external-url", newValue);
  }

  get time() {
    return this.#time;
  }

  set time(newValue) {
    this.setAttribute("time", newValue);
  }
}

customElements.define("rpr-timecode", Timecode);
