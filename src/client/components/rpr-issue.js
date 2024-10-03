import { LitElement, html, css } from "lit"
import { unsafeHTML } from "lit/directives/unsafe-html.js";

class Issue extends LitElement {

  static properties = {
    issueId: { type: Number, reflect: true },
    issueTitle: { type: String, reflect: true },
  };

  static styles = css`
    #issue {
      z-index: 10;
      display: block;
      transition: all;
      transition-duration: 100ms;
      margin: 1rem;
      margin-right: -4rem;
      padding: 1rem;
      border-radius: 0.375rem;
      background: linear-gradient(170deg, var(--yellow-100) 30%, var(--amber-200) 100%);
      float: right;
      clear: right;
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      position: relative;
      text-decoration: none;
      vertical-align: top;
      width: 25%;
      color: var(--yellow-700);
    }
    #issue:hover {
      transform: translate(0, 0.25rem); 
      box-shadow: 0 25px 50px -12px rgba(202, 138, 4, 0.4);
      background: linear-gradient(135deg, var(--yellow-100) 70%, var(--amber-200) 100%);
    }
    @media (max-width: 1485px) {
      #issue {
        margin-right: 0;
      }
    }
    #heading {
      color: var(--yellow-700);
      font-weight: 700;
      margin-right: 0.125rem;
      letter-spacing: -0.05em;
    }
    #heading img {
      height: 1rem;
      width: 1rem;
      display: inline-block;
      position: relative;
      top: 2px;
      margin-right: 0.125rem;
      display: none;
    }
    #title {
    }
  `;

  constructor() {
    super();
    this.issueId = 0;
    this.issueTitle = "";
    this.issueUrl = "";
  }

  render() {
    this.id = `issue-${this.issueId}`;

    return html`
      <a id="issue" href="https://github.com/marcuswhybrow/ray-peat-rodeo/issues/${this.issueId}">
        <span id="heading"><img src="/static/images/github-mark.svg" /> #<span id="id">${this.issueId}</span></span>
        <span id="title">${unsafeHTML(this.issueTitle)}</span> →
      </a>
    `;
  }
}

customElements.define("rpr-issue", Issue);
