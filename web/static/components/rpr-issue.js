class Issue extends HTMLElement {

  /** @type {Number} */
  #issueId

  /** @type {string} */
  #url

  /** @type {string} */
  #title

  /** @type {HTMLElement} */
  #idElement

  /** @type {HTMLElement} */
  #titleElement

  /** @type {HTMLAnchorElement} */
  #issueElement

  static observedAttributes = ["issue-id", "url", "title"];
  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" })
    shadowRoot.innerHTML = `
      <style>
        #issue {
          z-index: 10;
          display: block;
          transition: all;
          transition-duration: 100ms;
          margin: 0.5rem;
          padding: 1rem;
          box-shadow: 0 20px 25px -5px rgb(133 77 14 / 0.2), 0 8px 10px -6px rgb(133 77 14 / 0.1);
          border-radius: 0.375rem;
          background: linear-gradient(135deg, var(--yellow-200) 10%, var(--amber-200) 100%);
          float: right;
          clear: right;
          font-size: 0.875rem;
          line-height: 1.25rem;
          position: relative;
          letter-spacing: -0.025em;
          text-decoration: none;
          width: 40%;
          margin-right: -4rem;
          vertical-align: top;
        }
        #issue:hover {
          transform: translate(0, 0.25rem); 
          box-shadow: 0 25px 50px -12px rgb(202 138 4 / 0.4);
          background: linear-gradient(135deg, var(--yellow-100) 70%, var(--amber-200) 100%);
        }
        #heading {
          color: var(--yellow-900);
          font-weight: 700;
          margin-right: 0.125rem;
        }
        #heading img {
          height: 1rem;
          width: 1rem;
          display: inline-block;
          position: relative;
          top: 2px;
          margin-right: 0.125rem;
        }
        #title {
          color: var(--yellow-800);
        }
      </style>
      <a id="issue">
        <span id="heading"><img src="/assets/images/github-mark.svg" /> #<span id="id"></span></span>
        <span id="title"></span>
      </a>
    `;

    this.#idElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#id"));
    this.#titleElement = /** @type {HTMLElement} */ (shadowRoot.querySelector("#title"));
    this.#issueElement = /** @type {HTMLAnchorElement} */ (shadowRoot.querySelector("#issue"));
  }

  /**
  * @param {string} name
  * @param {string} _oldValue
  * @param {string} newValue
  */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "issue-id":
        this.#issueId = parseInt(newValue);
        this.setAttribute("id", `issue-${newValue}`);
        this.#idElement.textContent = newValue;
        break;
      case "url":
        this.#url = newValue;
        this.#issueElement.href = newValue;
        break;
      case "title":
        this.#title = newValue;
        this.#titleElement.textContent = newValue;
        break;
    }
  }

  set issueId(newValue) {
    this.setAttribute("issue-id", newValue.toString());
  }

  get issueId() {
    return this.#issueId
  }

  set url(newValue) {
    this.setAttribute("url", newValue);
  }

  get url() {
    return this.#url;
  }

  set title(newValue) {
    this.setAttribute("title", newValue);
  }

  get title() {
    return this.#title;
  }
}

customElements.define("rpr-issue", Issue);
