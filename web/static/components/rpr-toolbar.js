customElements.define("rpr-toolbar", class Toolbar extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" })
    this.shadowRoot.innerHTML = `
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
        button {
          flex-grow: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;
          gap: 0.25rem;

          background: transparent;
          border: none;

          cursor: pointer;
          text-align: center;
          padding: 1rem 0;
        }
        button .caption {
          font-size: 0.725rem;
          line-height: 1rem;
          text-transform: lowercase;
          letter-spacing: 0.1em;
          color: #aaa;
        }

        button svg {
          height: 2rem;
        }
        button svg * {
          stroke: #aaa;
          fill: transparent;
        }
        button svg .filled {
          fill: #aaa;
        }

        button.book svg {
          transform: translateX(3px);
        }

      </style> 

      <button class="search">
        <svg viewBox="0 0 100 100">
          <circle 
            class="border" 
            cx="50" cy="50" r="43" 
            stroke-width="5"
            stroke="#F00"
            fill="transparent"
          />
          <line x1="80" y1="80" x2="100" y2="100" stroke="#F00" stroke-width="5" /> 
        </svg>
        <div class="caption">Search</div>
      </button>

      <button class="book">
        <svg viewBox="0 0 100 100">
          <rect 
            x="3" y="3" 
            width="70" height="94"
            rx="10" ry="10" 
            stroke-width="5"
          />
        </svg>
        <div class="caption">Book</div>
      </button>

      <button class="help">
        <svg viewBox="0 0 110 100">
          <polygon points="3,97 107,97 55,3" stroke-width="5" />
          <circle class="filled" cx="72" cy="78" r="5" />
          <circle class="filled" cx="55" cy="78" r="5" />
          <circle class="filled" cx="38" cy="78" r="5" />
        </svg>
        <div class="caption">Help</div>
      </button>

      <button class="sidebar">
        <svg viewBox="0 0 100 100">
          <rect 
            class="border" 
            x="3" y="3" 
            width="94" height="94"
            rx="10" ry="10" 
            stroke-width="5"
            stroke="#F00"
            fill="transparent"
          />
          <line x1="25" y1="0" x2="25" y2="100" stroke="#F00" stroke-width="5" /> 
        </svg>
        <div class="caption">Sidebar</div>
      </button>
    `;

    this.shadowRoot.querySelector("button.sidebar").addEventListener("click", () => {
      this.dispatchEvent(new CustomEvent("sidebar", {
        bubbles: true,
        detail: "toggle",
      }));
    });

    this.shadowRoot.querySelector("button.search").addEventListener("click", () => {
      this.dispatchEvent(new CustomEvent("search", {
        bubbles: true,
      }));
    });
  }

  connectedCallback() {

  }

  attributeChangedCallback(name, oldValue, newValue) {

  }
});
