class Toolbar extends HTMLElement {


  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" })
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          background: white;
          box-shadow: 0 0 20px #ddd;
          border-radius: 0.5rem;
          display: flex;
          flex-direction: row;
          gap: 2rem;
          justify-content: center;
        }

      </style> 

    `;

  }
}

customElements.define("rpr-toolbar", Toolbar);
