class NewIssue extends HTMLElement {

  /** @type {HTMLAnchorElement} */
  #newIssueElement

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.innerHTML = `
      <style>
        :host([show="true"]) {
          display: block;
        }
        :host(:not([show="true"])) {
          display: none;
        }
        #new-issue {
          position: absolute;
          background: rgb(254, 240, 138);
          border-radius: 0.5rem;
          padding: 2rem;
          badding-bottom: 1.5rem;
          z-index: 40;
          transition: all;
          transition-duration: 100ms;
          box-shadow: 0 10px 15px -3px rgba(250, 204, 21, 0.5), 0 4px 6px -4px rgba(250, 204, 21, 0.5);
          text-decoration: none;
        }
        #new-issue:hover {
          box-shadow: 0 25px 50px -12px rgba(202, 138, 4, 0.6);
          transition: scale(1.1);
        }
        #new-issue:hover img {
          opacity: 0.75;
        }
        #new-issue:hover div {
          border: rgba(0,0,0,0.75);
          color: rgba(0,0,0,0.65);
        }

        #top {
          display: grid;
          grid-template-columns: repeat(2, minmax(0, 1fr));
          gap: 1rem;
        }
        #top #left {
          padding: 1rem;
          margin-bottom: 1rem;
          border: 0.125rem rgba(0,0,0,0.5) dashed;
          border-radius: 0.5rem;
        }
        #top #left img {
          display: block;
          width: 2rem;
          margin: 0 auto;
          opacity: 0.5;
        }
        #top #right {
          font-style: italic;
          width: 5rem;
          line-height: 1.375;
          color: rgba(0,0,0,0.6);
        }
        #bottom {
          font-weight: 700;
          text-transform: uppercase;
          align: center;
          color: rgba(0,0,0,0.5);
          letter-spacing: 0.05em;
          font-size: 0.875rem;
          line-height: 1.125rem;
        }
      </style>
      <a
        id="new-issue"
        href="https://github.com/marcuswhybrow/ray-peat-rodeo/issues/new"
        target="_blank"
        data-pagefind-ignore
      >
        <div id="top">
          <div id="left">
            <img src="/assets/images/plus-line-icon.svg" alt="Plus icon" />
          </div>
          <div id="right">Spotted an issue with this text?</div>
        </div>
        <div id="bottom">Email Ray Peat Rodeo</div>
      </a>
    `;

    this.#newIssueElement = /** @type {HTMLAnchorElement} */ (shadowRoot.querySelector("#new-issue"));

    "mouseup touchend".split(" ").forEach(e => window.addEventListener(e, () => {
      const selection = window.getSelection();
      if (!selection) {
        return;
      }

      const selectedText = selection.toString();
      const gap = 16;

      // The mouseup event may indicate either that a text selection had been
      // created or removed. When removed we must wait 1ms before checking to 
      // get the correct state.
      setTimeout(() => {
        if (!selection.isCollapsed) {
          this.style.display = "block";
          const selRect = selection.getRangeAt(0).getBoundingClientRect();
          const articleElement = document.querySelector("article");
          if (!articleElement) {
            return;
          }
          const articleRect = articleElement.getBoundingClientRect();

          this.#newIssueElement.style.top = (selRect.bottom + gap - articleRect.top) + "px";
          this.#newIssueElement.style.left = (selRect.left - articleRect.left) + "px";

          const titleElement = document.querySelector("h1");
          if (!titleElement) {
            return;
          }

          const assetTitle = titleElement.textContent;
          const assetLink = `https://raypeat.rodeo${window.location.pathname}`;
          const title = encodeURIComponent(`Issue with "${assetTitle}"`);
          const body = encodeURIComponent(`
<p>Hi, I've found an issue with this text from <a href="${assetLink}">${assetTitle}</a>:</p>
<blockquote>${selectedText}</blockquote>
<p>[consider describing the issue here]</p>
          `);
          this.#newIssueElement.href = `mailto:contact@raypeat.rodeo?subject=${title}&body=${body}`;
        } else {
          this.style.display = "none";
        }
      }, 1);
    }));

  }
}

customElements.define("rpr-new-issue", NewIssue);
