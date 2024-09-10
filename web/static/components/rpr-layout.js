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


        .side { 
          height: 100vh;
          width: 33%;
          overflow-y: scroll;
          position: fixed;
          top: 0;
          left: 0;
          padding: 0 2rem 2rem;
          box-sizing: border-box;
          transition: 200ms ease-in-out;
        }
        :host(:not([sidebar="true"])) .side {
          transform: translateX(-100%);
        }

        .side rpr-search {
          margin-bottom: 8rem;
        }
        rpr-toolbar {
          position: fixed;
          bottom: 2rem;
          left: 2rem;
          width: calc(33% - 4rem);
          transition: 300ms;
          transition-timing-function: ease-in-out;
        }
        :host(:not([sidebar="true"])) rpr-toolbar {
          transform: translateX(calc(50vw - 50%));
        }

        .main::slotted(*) {
          padding-left: calc(33% + 2rem);
          padding-right: 2rem;
          padding-top: 2rem;
          padding-bottom: 2rem;
          transition: 300ms ease-in-out;
        }
        :host(:not([sidebar="true"])) .main::slotted(*) {
          padding-left: 2rem;
        }

        .sidebar-icon {
          position: fixed;
          top: 1rem;
          left: 1rem;
          width: 2.2rem;
          opacity: 0.4;
          cursor: pointer;
        }
        :host([sidebar="true"]) .sidebar-icon {
          display: none;
        }

      </style>
      <div class="side">
				<rpr-search></rpr-search>
      </div>

      <slot class="main"></slot>

      <!-- <rpr-ad></rpr-ad> -->

      <rpr-toolbar></rpr-toolbar>
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
      switch (event.detail) {
        case "open":
          this.sidebar = true;
          break;
        case "close":
          this.sidebar = false;
          break;
        case "toggle":
          this.sidebar = !this.sidebar;
          break;
      }
    });

    this.shadowRoot.addEventListener("search", () => {
      this.sidebar = true;
      const search = this.shadowRoot.querySelector("rpr-search");
      search.focus();
      search.select();
    });

    const params = new URLSearchParams(window.location.search);
    if (params.get("sidebar") === "false") {
      this.sidebar = false;
    } else {
      this.sidebar = true;
    }
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
