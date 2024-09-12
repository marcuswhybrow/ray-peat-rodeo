class RayPeat extends HTMLElement {
  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          display: block;
          padding: 2rem;
        }
      </style>
      <h1>Ray Peat</h1>
      <p>Ray Peat (1936 - 2022) was a biologist and teacher.</p>
    `;
  }
}

customElements.define("rpr-ray-peat", RayPeat);
