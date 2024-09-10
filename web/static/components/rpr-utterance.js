window.customElements.define("rpr-utterance", class extends HTMLElement {
  #by = null;
  #avatar = null;
  #primary = false;
  #short = false;

  static observedAttributes = ["by", "avatar", "primary", "short"];

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
      <style>
        :host(:first) #utterance {
          margin-top: 0;
        }
        #utterance {
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          position: relative;
        }
        :host([primary="true"]) #utterance {
          margin-left: 0.25rem;
          margin-right: 4rem;
        }
        :host(:not([primary="true"])) #utterance {
          margin-left: 4rem;
          margin-right: 0.25;
        }
        :host([short="true"]) #utterance {
          margin-top: -1rem;
        }
        :host(:not([short="true"])) #utterance {
          margin-top: 1rem;
        }
        :host(:not([avatar])) #avatar {
          display: none;
        }
        #avatar {
          width: 2rem;
          height: 2rem;
          border-radius: 9999px;
          display: inline-block;
          box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
          float: left;
          margin-right: 1rem;
          margin-bottom: 0;
          overflow: hidden;
          position: absolute;
          left: -3rem;
          top: -0.25rem;
        }
        #avatar > div {
          width: 9999px;
        }
        #avatar img {
          height: 2rem;
        }
        :host([short="true"]) #avatar {
          display: none;
        }
        :host([short="true"]) #name {
          display: none;
        }
        :host(:not([short="true"])) #name {
          font-size: 0.75rem;
          line-height: 1rem;
          margin-top: 2rem;
          margin-bottom: 1rem;
          display: block;
        }
        :host([primary="true"]) #name {
          color: #9CA3AF;
        }
        :host(:not([primary="true"])) #name {
          color: #38BDF8;
        }
        #content {
          padding: 2rem;
          border-radius: 0.25rem;
          box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
        }
        #content ::slotted(p) {
          margin-top: 0;
          margin-bottom: 1.5rem;
        }
        #content ::slotted(p:last-child) {
          margin-bottom: 0 !important;
        }
        #content > blockquote {
          padding-left: 1rem;
          font-size: 0.75rem;
          line-height: 1rem;
        }
        :host([primary="true"]) #content {
          color: #111827;
          background: #F3F4F6;
        }
        :host(:not([primary="true"])) #content {
          color: #0C4A6E;
          background: linear-gradient(135deg, rgb(224, 242, 254) 0%, rgb(191, 219, 254) 100%);
        }
        :host([short="true"]) #content {
          display: inline-block;
        }
        :host(:not([short="true"])) #content {
          display: block;
        }
      </style>
      <div id="utterance">
        <div id="avatar">
          <div>
            <img src="" alt="" />
          </div>
        </div>
        <div id="name" data-pagefind-ignore></div>
        <div id="content">
          <slot></slot>
        </div>
      </div>
    `;
  }

  connectedCallback() {

  }

  disonnectedCallback() {

  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "by":
        this.#by = newValue;
        this.shadowRoot.querySelector("#name").textContent = newValue;
        const img = this.shadowRoot.querySelector("#avatar img");
        img.alt = newValue;
        img.title = newValue;
        break;
      case "avatar":
        this.#avatar = newValue;
        this.shadowRoot.querySelector("#avatar img").src = newValue;
        break;
      case "primary":
        this.#primary = newValue === "true";
        break;
      case "short":
        this.#short = newValue === "true";
        break;
    }
  }

  set by(newValue) {
    this.setAttribute("by", newValue);
  }

  get by() {
    return this.#by;
  }

  set avatar(newValue) {
    this.setAttribute("avatar", newValue);
  }

  get avatar() {
    return this.#avatar;
  }

  set primary(newValue) {
    this.setAttribute("primary", newValue);
  }

  get primary() {
    return this.#primary;
  }

  set short(newValue) {
    this.setAttribute("short", newValue);
  }

  get short() {
    return this.#short;
  }
});
