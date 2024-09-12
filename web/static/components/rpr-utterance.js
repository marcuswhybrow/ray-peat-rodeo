class Utterance extends HTMLElement {
  /**
   * Full name of the person uttering this utterance.
   *
   * @type {string} 
   */
  #by

  /** 
   * Absolute path to the avatar image for this person.
   *
   * @type {string} 
   */
  #avatar

  /** 
   * Is the person speaking a "primary" speaker?
   *
   * @type {Boolean} 
   */
  #primary

  /** 
   * Should this utterance be displyed as a short utterance? The avatar and 
   * byline will be hidden, the width may shrink, and it will overlap with 
   * the previous utterance.
   *
   * @type {Boolean} 
   */
  #short

  /** @type {HTMLElement} */
  #nameElement

  /** @type {HTMLElement} */
  #contentElement

  /** @type {HTMLImageElement} */
  #avatarImgElement

  static observedAttributes = ["by", "avatar", "primary", "short"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
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
          color: var(--gray-400);
        }
        :host(:not([primary="true"])) #name {
          color: var(--sky-400);
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
          color: var(--gray-900);
          background: var(--gray-100);
        }
        :host(:not([primary="true"])) #content {
          color: var(--sky-900);
          background: linear-gradient(135deg, var(--sky-100) 0%, var(--blue-200) 100%);
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
        <div id="name"></div>
        <div id="content">
          <slot></slot>
        </div>
      </div>
    `;

    this.#nameElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#name"));
    this.#avatarImgElement = /** @type {HTMLImageElement} */ (shadowRoot.querySelector("#avatar img"));
  }

  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "by":
        this.#by = newValue;
        this.#nameElement.textContent = newValue;
        this.#avatarImgElement.alt = newValue;
        this.#avatarImgElement.title = newValue;
        break;
      case "avatar":
        this.#avatar = newValue;
        this.#avatarImgElement.src = newValue;
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
    this.setAttribute("primary", newValue.toString());
  }

  get primary() {
    return this.#primary;
  }

  set short(newValue) {
    this.setAttribute("short", newValue.toString());
  }

  get short() {
    return this.#short;
  }
}

customElements.define("rpr-utterance", Utterance);
