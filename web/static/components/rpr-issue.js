window.customElements.define("rpr-issue", class extends HTMLElement {
  #id = null;
  #url = null;
  #title = null;

  static observedAttributes = ["issue-id", "url", "title"];
  constructor() {
    super();
    this.attachShadow({ mode: "open" })
    this.shadowRoot.innerHTML = `
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
          background: linear-gradient(135deg, rgb(254, 240, 138) 10%, rgb(253, 230, 138) 100%);
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
          background: linear-gradient(135deg, rgb(254, 249, 195) 70%, rgb(253, 230, 138) 100%);
        }
        #heading {
          color: rgb(113, 63, 18);
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
          color: rgb(133, 77, 14);
        }
      </style>
      <a id="issue">
        <span id="heading"><img src="/assets/images/github-mark.svg" /> #<span id="id"></span></span>
        <span id="title"></span>
      </a>
    `;
  }

  connectedCallback() {

  }

  disonnectedCallback() {

  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch(name) {
      case "issue-id":
        this.#id = newValue;
        this.setAttribute("id", `issue-${newValue}`);
        this.shadowRoot.querySelector("#id").textContent = newValue;
        break;
      case "url":
        this.#url = newValue;
        this.shadowRoot.querySelector("#issue").href = newValue;
        break;
      case "title":
        this.#title = newValue;
        this.shadowRoot.querySelector("#title").textContent = newValue;
        break;
    }
  }

  set id(newValue) {
    this.setAttribute("issue-id", newValue);
  }

  get id() {
    return this.#id
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
})
