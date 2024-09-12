class Sidenode extends HTMLElement {
  /** @type {string} */
  #sidenoteId

  static observedAttributes = ["sidenote-id"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        label {
          counter-increment: sidenote;
          font-family: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
        }
        label::after {
          content: counter(sidenote);
          top: -0.25rem;
          left: 0;
          vertical-align: baseline;
          font-size: 0.875rem;
          line-height: 1.25rem;
          position: relative;
          background: white;
          border-radius: 0.375rem;
          box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
          color: rgb(75, 85, 99);
          padding: 0.25rem 0.5rem;
        }
        #sidenote {
          z-index: 10;
          display: block;
          background: white;
          border-radius: 0.365rem;
          box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
          width: 50%;
          margin-right: -4rem;
          float: right;
          clear: right;
          font-size: 0.875rem;
          line-height: 1.25rem;
          position: relative;
          padding: 1rem;
          line-height: 1.25rem;
          vertical-align: middle;
          transition: all;
          transition-duration: 100ms;
        }
        #sidenote::before {
          content: counter(sidenote) ".";
          float: left;
          margin-right: 0.25rem;
          color: rgb(107, 114, 128);
          font-size: 0.875rem;
          line-height: 1.25rem;
        }
        #sidenote ::slotted(img) {
          padding: 0.5rem 0;
          max-width: 100%;
        }
        #sidenote img:last-child {
          padding-bottom: 0;
        }
      </style>
      <label for="sidenote"></label><span id="sidenote"><slot></slot></span>
    `;
  }

  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "sidenote-id":
        this.#sidenoteId = newValue;
        this.setAttribute("id", `sidenote-${newValue}`);
        break;
    }
  }

  set sidenoteId(newValue) {
    this.setAttribute("sidenote-id", newValue);
  }

  get sidenoteId() {
    return this.#sidenoteId;
  }
}

customElements.define("rpr-sidenote", Sidenode);
