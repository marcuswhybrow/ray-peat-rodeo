window.customElements.define("rpr-new-issue", class extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.shadowRoot.innerHTML = `
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

    "mouseup touchend".split(" ").forEach(e => window.addEventListener(e, event => {
      const newIssue = this.shadowRoot.querySelector("#new-issue");
      const selection = window.getSelection();
      const selectedText = selection.toString();
      const gap = 16;

      // The mouseup event may indicate either that a text selection had been
      // created or removed. When removed we must wait 1ms before checking to 
      // get the correct state.
      setTimeout(() => {
        if (!selection.isCollapsed) {
          this.style.display = "block";
          const selRect = selection.getRangeAt(0).getBoundingClientRect();
          const articleRect = document.querySelector("article").getBoundingClientRect();
          newIssue.style.top = (selRect.bottom + gap - articleRect.top) + "px";
          newIssue.style.left = (selRect.left - articleRect.left) + "px";

          const assetTitle = document.querySelector("h1").textContent;
          const assetLink = `https://raypeat.rodeo${window.location.pathname}`;
          const title = encodeURIComponent(`Issue with "${assetTitle}"`);
          const body = encodeURIComponent(`
<p>Hi, I've found an issue with this text from <a href="${assetLink}">${assetTitle}</a>:</p>
<blockquote>${selectedText}</blockquote>
<p>[consider describing the issue here]</p>
          `);
          newIssue.href = `mailto:contact@raypeat.rodeo?subject=${title}&body=${body}`;
        } else {
          this.style.display = "none";
        }
      }, 1);
    }));

  }

  connectedCallback() {

  }

  disconnectedCallback() {

  }

  attributeChangedCallback(name, oldValue, newValue) {
  }
});
