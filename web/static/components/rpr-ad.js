class Ad extends HTMLElement {
  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          position: fixed;
          left: 2.5rem; 
          bottom: 2rem;
          width: calc(33% - 5rem);
          background: linear-gradient(184deg, #FEF9C3 0%, #FACC15 100%);
          border-radius: 0.5rem;
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          cursor: pointer;
        }
        #wrapper {
          padding: 0.25rem 1rem;
        }
        #wrapper:hover #book {
          transform: rotate(10deg);
        }
        h3 {
          font-size: 1rem;
          font-weight: 400;
          letter-spacing: 0.05em;
          text-transform: uppercase;
          color: #A16207;
          margin-left: 1rem;
        }
        #book {
          width: 4rem;
          height: 6rem;
          background: linear-gradient(10deg, #64748B 0%, #94A3B8 100%);
          position: absolute;
          bottom: 0.5rem;
          right: 1rem;
          transform: rotate(5deg);
          transition: all;
          transition-duration: 100ms;
        }
      </style>
      <div id="wrapper">
        <h3>£9.99 <i>The Complete Works eBook</i></h3>
        <div id="book"></div>
      </div>
    `;
  }
}

customElements.define("rpr-ad", Ad);
