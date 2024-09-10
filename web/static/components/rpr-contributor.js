window.customElements.define("rpr-contributor", class extends HTMLElement {
  static observedAttributes = [ "avatar", "name" ];

  #avatar = null;
  #name = null;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
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
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "avatar":
        this.#avatar = newValue;
        this.shadowRoot.querySelector("img").src = newValue;
        break;
      case "name":
        this.#name = newValue;
        this.shadowRoot.querySelector("img").title = newValue;
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
});
