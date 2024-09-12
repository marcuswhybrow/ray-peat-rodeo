class Contributor extends HTMLElement {
  static observedAttributes = ["avatar", "name"];

  /** 
   * Absolute path to the avatar image for the person speaking or writing.
   *
   * @type {string} 
   */
  #avatar

  /** 
   * Full name of the person speaking or writing.
   *
   * @type {string} 
   */
  #name

  /** @type {HTMLImageElement} */
  #avatarElement

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          display: inline-block;
          margin
        }
        span {
          display: flex;
          width: 1.5rem;
          height: 1.5rem;
          overflow: hidden;
          border-radius: 9999px;
        }
        span img {
          display: block;
          height: 100%;
        }
      </style>
      <span><img /></span>
    `;

    this.#avatarElement = /** @type {HTMLImageElement} */ (shadowRoot.querySelector("img"));
  }

  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "avatar":
        this.#avatar = newValue;
        this.#avatarElement.src = newValue;
        break;
      case "name":
        this.#name = newValue;
        this.#avatarElement.title = newValue;
        break;
    }
  }

  set avatar(a) {
    this.setAttribute("avatar", a);
  }

  set name(n) {
    this.setAttribute("name", n);
  }

  get avatar() {
    return this.#avatar;
  }

  get name() {
    return this.#name;
  }
}

customElements.define("rpr-contributor", Contributor);
