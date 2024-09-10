window.customElements.define("rpr-layout", class Layout extends HTMLElement {
  #sidebar = true;

  static observedAttributes = ["sidebar"];

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
  }

  connectedCallback() {
    this.shadowRoot.innerHTML = `
      <style>
        :host(rpr-sidebar) {
        }

        ::-webkit-scrollbar {
          width: 8px;
          background: white;
        }
        ::-webkit-scrollbar:hover {
          background: #fafafa;
        }

        ::-webkit-scrollbar-thumb {
          background: #ddd;
          border-radius: 0.25rem;
        }

        :host(:not([sidebar="true"])) .side {
          display: none;
        }
        :host(:not([sidebar="true"])) .main {
          padding-left: 0;
        }

        .side { 
          height: 100vh;
          width: 33%;
          overflow-y: scroll;
          position: fixed;
          top: 0;
          left: 0;
        }
        .side .inner {
          padding: 0 2rem 2rem;
          margin-bottom: 8rem;
        }

        .main {
          padding-left: 33%;
        }
        .main .inner {
          padding: 2rem;
        }

        #sidebar-icon {
          position: fixed;
          top: 1rem;
          left: 1rem;
          width: 2.2rem;
          opacity: 0.4;
          cursor: pointer;
        }
        :host([sidebar="true"]) #sidebar-icon {
          display: none;
        }
      </style>
      <div class="side">
        <div class="inner">
          <slot name="side"></slot>
        </div>
      </div>

      <div class="main">
        <div class="inner">
          <slot name="main"></slot>
        </div>
      </div>

      <div id="advertisement">
        <slot name="advertisement"></slot>
      </div>

      <img id="sidebar-icon" src="/assets/images/interface-layout-left-sidebar-icon.svg" />
    `;

    this.shadowRoot.addEventListener("pick", async event => {
      const asset = event.detail;
      const response = await fetch(asset.path);
      if (!response.ok) {
        console.error(`Failed to fetch asset ${asset.path}`);
        return;
      }
      const doc = parser.parseFromString(await response.text(), "text/html");
      const newPane = doc.querySelector("#reading-pane");
      document.querySelector("#reading-pane").replaceWith(newPane);
      window.scrollTo(0,0);
      history.pushState({}, "", asset.path);
    });

    this.shadowRoot.addEventListener("sidebar", event => {
      this.sidebar = event.detail;
    });

    const params = new URLSearchParams(window.location.search);
    if (params.get("sidebar") === "false") {
      this.sidebar = false;
    } else {
      this.sidebar = true;
    }

    const icon = this.shadowRoot.querySelector("#sidebar-icon");
    icon.addEventListener("click", event => {
      this.sidebar = true;
    });
    icon.addEventListener("keyup", event => {
      if (event.key === "Enter") icon.click();
    });

  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "sidebar":
        this.#sidebar = newValue === "true";

        break;
    }
  }

  set sidebar(newValue) {
    this.setAttribute("sidebar", newValue);
  }

  get sidebar() {
    return this.#sidebar;
  }

});
