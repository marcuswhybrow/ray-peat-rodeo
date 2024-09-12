
/**
  * @typedef {object} SidebarEvent
  * @property {"open"|"close"|"toggle"} detail
  */

/**
 * A main content area with a togglable sidebar. Handles the picking of new 
 * {@link Asset|Assets} by fetching it's content and updaing the DOM.
 */
class Layout extends HTMLElement {
  /** Is the sidebar visible? */
  #sidebar = true;

  /** Attributes which trigger {@link attributeChangedCallback} */
  static observedAttributes = ["sidebar"];

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host(*) {
          display: block;
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
          background: white;
          width: calc(33% - 3rem);
          position: fixed;
          top: 1.5rem;
          left: 1.5rem;
          height: calc(100vh - 3rem);

          padding: 0 2rem 0 2rem;
          // box-shadow: 0 0 20px #ddd;

          transition: 200ms ease-in-out;
          overflow-y: scroll;
          border-radius: 0.25rem;
        }
        :host(:not([sidebar="true"])) .side {
          transform: translateX(calc(-100% + 1rem));
        }
        .side::after {
          content: "";
          position: fixed;
          bottom: 0;
          left: 0;
          right: 0;
          height: 8rem;
          background: linear-gradient(180deg, transparent 0%, white 100%);
          pointer-events: none;
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
          /* transform: translateX(calc(50vw - 50%)); */
          display: none;
        }

        .main::slotted(*) {
          background: white;
          margin-left: calc(33% + 3rem);
          transition: 300ms ease-in-out;
          border-radius: 0.25rem;

          display: block;
        }
        :host(:not([sidebar="true"])) .main::slotted(*) {
          margin-left: 4rem;
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
    `;

    const search = /** @type {Search} */ (shadowRoot.querySelector("rpr-search"));
    const readingPane = /** @type {HTMLDivElement} */ (shadowRoot.querySelector("#reading-pane"));
    const sideElement = /** @type {HTMLElement} */ (shadowRoot.querySelector(".side"));
    const parser = new DOMParser();

    sideElement.addEventListener("click", () => {
      if (!this.sidebar) {
        this.sidebar = true;
      }
    });

    /**
     * @type {(event: Event) => Promise<void>}
     * @param {PickEvent} event
      */
    async function pickHandler(event) {
      const asset = event.detail.asset;
      const response = await fetch(asset.path);
      if (!response.ok) {
        console.error(`Failed to fetch asset ${asset.path}`);
        return;
      }
      const doc = parser.parseFromString(await response.text(), "text/html");

      const newPane = doc.querySelector("#reading-pane");
      if (!newPane) {
        console.error("Failed to find #reading-pane in fetched content");
        return
      }

      readingPane.replaceChildren(newPane);
      window.scrollTo(0, 0);
      history.pushState({}, "", asset.path);
    }
    shadowRoot.addEventListener("pick", pickHandler);

    const layout = this;

    /**
     * @function
     * @type {(event: Event) => void}
     * @param {SidebarEvent} event
     */
    function sidebarHandler(event) {
      switch (event.detail) {
        case "open":
          layout.sidebar = true;
          break;
        case "close":
          layout.sidebar = false;
          break;
        case "toggle":
          layout.sidebar = !layout.sidebar;
          break;
      }
    }
    shadowRoot.addEventListener("sidebar", sidebarHandler);

    shadowRoot.addEventListener("search", () => {
      this.sidebar = true;
      search.focus();
      search.select();
    });
  }

  connectedCallback() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("sidebar") === "false") {
      this.sidebar = false;
    } else {
      this.sidebar = true;
    }
  }

  /**
   * @param {string} name 
   * @param {string} _oldValue
   * @param {string} newValue
   */
  attributeChangedCallback(name, _oldValue, newValue) {
    switch (name) {
      case "sidebar":
        this.#sidebar = newValue === "true";
        break;
    }
  }

  /** Set the visibility of the sidebar */
  set sidebar(newValue) {
    this.setAttribute("sidebar", newValue.toString());
  }

  /** Is the sidebar visible? */
  get sidebar() {
    return this.#sidebar;
  }

}

customElements.define("rpr-layout", Layout);
