import { LitElement, html, css } from "lit"

class Contribution extends LitElement {
  static properties = {
    name: { type: String, reflect: true },
    initials: { type: String, reflect: true },
    avatar: { type: String, reflect: true },
    filterable: { type: Boolean, reflect: true },
  };

  constructor() {
    super();
    this.initials = "";
    this.name = "";
    this.avatar = "";
    this.filterable = false;
  }

  static styles = css`
    :host(*) {
      font-family: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
      position: relative;
      margin-top: 4rem;
      display: block;
    }
    :host(:first) {
      margin-top: 0;
    }

    .avatar {
      border: 1px dashed var(--slate-300);
      width: 2rem;
      height: 2rem;
      border-radius: 9999px;
      display: inline-block;
      float: left;
      margin-right: 1rem;
      margin-bottom: 0;
      overflow: hidden;
      position: absolute;
      left: -4rem;
      top: -0rem;
      pointer-events: none;
    }
    .avatar .some {
      display: block;
      width: 9999px;
    }
    .avatar .none {
      display: table-cell;
      width: 2rem;
      height: 2rem;
      text-align: center;
      vertical-align: middle;
      color: var(--slate-400);
    }

    .avatar .some img {
      height: 2rem;
    }

    .avatar:has(.some) {
      border: 1px solid transparent;
    }
    @media (max-width: 1485px) {
      .avatar {
        left: 0;
        position: relative;
      }
    }

    .byline {
      font-size: var(--font-size-sm);
      line-height: var(--line-height-sm);
      margin-bottom: 1rem;
      display: block;
      letter-spacing: -0.025em;
      color: var(--slate-950);
    }

    .byline rpr-filter {
      border: 1px solid #E2E8F0;
      border-radius: 1rem;
      padding: 0.25rem 0.75rem;
    }

    .byline .no-pin {
      padding: 0.25rem 0;
      border: 1px solid transparent;
    }

    .content {
      display: block;
      color: var(--slate-700);
      font-size: var(--font-size-md);
      line-height: var(--line-height-md);
    }

    .content ::slotted(p) {
      margin-top: 0;
      margin-bottom: 1.5rem;
    }
    .content ::slotted(p:last-child) {
      margin-bottom: 0 !important;
    }
    .content > blockquote {
      padding-left: 1rem;
      font-size: 0.75rem;
      line-height: 1rem;
    }
  `;

  render() {
    return html`
      <div class="avatar">
        ${this.avatar
        ? html`<div class="some"><img src="${this.avatar}" alt="${this.name}" /></div>`
        : html`<div class="none">${this.initials}</div>`
      }
      </div>
      <div class="byline">
        ${this.filterable
        ? html`
          <rpr-filter 
            part="filter"
            exportparts="wrapper: filter-wrapper, key: filter-key, value: filter-value, checkbox: fitler-checkbox" 
            key="contributor" 
            hidekey 
            value="${this.name}"
          ></rpr-pin>`
        : html`<div class="no-pin">${this.name}</div>`
      }
      </div>
      <div class="content">
        <slot></slot>
      </div>
    `;
  }

}

customElements.define("rpr-contribution", Contribution);
