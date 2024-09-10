window.customElements.define("rpr-asset-excerpt", class extends HTMLElement {
  #text = "";

  static observedAttributes = [ "text" ];

  constructor() {
    super();
  }

  connectedCallback() {
    this.addEventListener("click", event => {
      const asset = this.parentElement;
      fetch(asset.link).then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error: ${response.status}`)
        }
        return response.text();
      }).then(text => {
        const parser = new DOMParser();
        const selectDoc = parser.parseFromString(text, "text/html");
        const select = selectDoc.getElementById("reading-pane");
        const target = document.getElementById("reading-pane");
        target.replaceWith(select);


        const asset = this.parentElement;
        const textParam = encodeURIComponent(this.text);
        const assetLink = asset.link;
        const hash = `:~:text=${textParam}`;
        const link = `${assetLink}#${hash}`;

        const state = {};
        const unused = "";
        history.pushState(state, unused, assetLink);
        location.hash = hash;
      }).catch(error => {
        console.log(error);
      });
    });
    this.update();
  }

  disconnectedCallback() {

  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "text") {
      this.#text = newValue;
      this.update();
    }
  }

  update() {
    if (!this.isConnected) return;
  }

  get text() {
    return this.#text;
  }

  set text(text) {
    this.setAttribute("text", text);
  }
});
