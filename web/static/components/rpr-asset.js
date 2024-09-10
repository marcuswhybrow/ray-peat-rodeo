window.customElements.define("rpr-asset", class extends HTMLElement {
  #weight = null;
  #excerpts = [];
  #elements = {};
  #active = false;
  #path = null;
  #date = null;
  #title = null;
  #kind = null;
  #series = null;
  #issues = [];
  #showIssues = null;

  static observedAttributes = [
    "date", "title", "kind", "active", "has-matches", "path", "series", "issues",
    "show-issues"
  ];

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
      <style>
        :host(*) {
          font-size: 0.75rem;
          line-height: 1rem;
          color: #64748B;
          border-radius: 0.25rem;
          overflow: hidden;
          display: inline-flex;
          flex-direction: column;
          cursor: pointer;
          font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
          letter-spacing: 0.025em;
          margin-bottom: 1rem;
        }
        :host([active="true"]) #title,
        :host([active="true"]) #stats {
          color: #292524 !important;
          text-shadow: 0 0 1px #292524;
        }
        :host([score]) #stats {
          display: none;
        }
        :host([score="0"]) #stats,
        :host(:not([score])) #stats {
          display: block;
        }
        .top {
          display: flex;
        }
        #title {
          flex: 1 1 auto;
          display: flex;
          font-size: 1rem;
          line-height: 1.5rem;
          letter-spacing: 0.025em;
          margin-bottom: 0.5rem;
        }
        #results {
          margin-top: 0.25rem;
        }
        #results:empty {
          margin-top: 0;
        }
        #results > div {
          margin-top: 0.5rem;
        }
        #results > div:first-child {
          margin-top: 0;
        }
        
        .issues {
          display: flex;
          flex-direction: column;
          gap: 0.25rem;
          margin: 0;
          padding: 0.25rem 0 0 0rem;
        }
        :host(:not([show-issues="true"])) .issues {
          display: none;
        }
        .issues td {
          vertical-align: top;
        }
        .issues td:first-child {
          padding-right: 0.5rem;
        }
        .issues a {
          color: #444;
        }
        .issues a:hover {
          color: #000;
        }
      </style>
      <span id="title">${this.title}</span>
      <div id="stats">
        <span id="date">${this.date}</span>
        <span id="series">${this.series}</span>
      </div>
      <div id="results">
      </div>
      <table class="issues"></table>
    `;

    this.addEventListener("click", event => {
      this.dispatchEvent(new CustomEvent("pick", {
        bubbles: true,
        detail: {
          asset: this, 
          issue: null,
        }
      }));
    });
    this.addEventListener("keyup", () => {
      if (event.key === "Enter") this.click();
    });
  }

  connectedCallback() {
  }

  disonnectedCallback() {
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "title":
        this.#title = newValue;
        this.shadowRoot.querySelector("#title").textContent = newValue;
        break;
      case "date":
        this.#date = newValue;
        this.shadowRoot.querySelector("#date").textContent = newValue;
        break;
      case "active":
        this.#active = newValue === "true";
        break;
      case "path":
        this.#path = newValue;
        break;
      case "series":
        this.#series = newValue;
        this.shadowRoot.querySelector("#series").textContent = newValue;
        break;
      case "issues":
        this.#issues = JSON.parse(newValue);
        const elements = [];
        for (const issue of this.issues) {
          const a = document.createElement("a");
          a.href = this.path + "#issue-" + issue.Id;
          a.textContent = issue.Title;
          a.addEventListener("click", e => {
            e.preventDefault();
            e.stopPropagation();
            this.dispatchEvent(new CustomEvent("pick", {
              bubbles: true,
              detail: {
                asset: this,
                issue: issue.Id,
              },
            }));
          });

          const left = document.createElement("td");
          left.textContent = `#${issue.Id}`;

          const right = document.createElement("td");
          right.replaceChildren(a);

          const tr = document.createElement("tr");
          tr.replaceChildren(left, right);
          elements.push(tr);
        }
        this.shadowRoot.querySelector(".issues").replaceChildren(...elements);
        break;
      case "show-issues":
        this.#showIssues = newValue === "true";
        break;
    }
  };

  get excerpts() {
    return this.#excerpts;
  }

  get active() {
    return this.#active;
  }

  set active(newValue) {
    this.setAttribute("active", newValue);
  }

  get date() {
    return this.#date;
  }

  set date(date) {
    this.setAttribute("date", date);
  }

  get title() {
    return this.#title;
  }

  set title(title) {
    this.setAttribute("title", title);
  }

  get kind() {
    return this.#kind;
  }

  set kind(kind) {
    this.setAttribute("kind", kind);
  }

  get link() {
    return this.getAttribute("link") || null;
  }

  get path() {
    return this.#path;
  }

  set path(newValue) {
    this.setAttribute("path", newValue);
  }

  get series() {
    return this.#series;
  }

  set series(newValue) {
    this.setAttribute("series", newValue);
  }

  get issues() {
    return this.#issues;
  }

  set issues(newValue) {
    this.setAttribute("issues", JSON.stringify(newValue));
  }

  get showIssues() {
    return this.#showIssues;
  }

  set showIssues(newValue) {
    this.setAttribute("show-issues", newValue);
  }

  deriveActive() {
    const currentPath = window.location.pathname.replace(/\/+$/, "");
    if (this.path === currentPath) {
      this.active = true;
      return true;
    }
    return false;
  }

  replaceResults(...pagefindSubResult) {
    const elements = [];
    for (const psr of pagefindSubResult) {
      const result = document.createElement("div");
      result.innerHTML = psr.excerpt;
      elements.push(result);
    }
    this.shadowRoot.querySelector("#results").replaceChildren(...elements);
  }
});
