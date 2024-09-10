window.customElements.define("rpr-timecode", class extends HTMLElement {
  #externalUrl = null;
  #time = null;

  static observedAttributes = ["external-url", "time"];

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
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
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "external-url":
        this.#externalUrl = newValue;
        this.shadowRoot.querySelector("#link").href = newValue;
        break;
      case "time":
        this.#time = newValue;
        this.shadowRoot.querySelector("#link").textContent = newValue;
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
});
