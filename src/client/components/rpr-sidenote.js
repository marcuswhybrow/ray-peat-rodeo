import { LitElement, html, css } from "lit";

class Sidenode extends LitElement {
  constructor() {
    super();
  }

  static styles = css`
    label {
      counter-increment: sidenote;
      font-family: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
    }
    label::after {
      content: counter(sidenote);
      top: -0.25rem;
      left: 0;
      vertical-align: baseline;
      font-size: 0.875rem;
      line-height: 1.25rem;
      position: relative;
      background: white;
      border-radius: 0.375rem;
      box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
      color: rgb(75, 85, 99);
      padding: 0.25rem 0.5rem;
    }
    #sidenote {
      z-index: 10;
      display: block;
      background: white;
      background: linear-gradient(170deg, var(--gray-50) 10%, var(--slate-50) 100%);
      box-shadow: 1px 1px 2px rgba(0,0,0,0.1);
      font-size: 0.875rem;
      line-height: 1.25rem;
      position: relative;
      line-height: 1.25rem;
      vertical-align: middle;
      padding: 1rem;
      margin: 1rem;
      margin-right: -4rem;
      float: right;
      clear: right;
      border-radius: 0.25rem;
      width: 25%;
    }
    @media (max-width: 1485px) {
      #sidenote {
        margin-right: 0;
      }
    }
    #sidenote::before {
      content: counter(sidenote) ".";
      float: left;
      margin-right: 0.25rem;
      color: rgb(107, 114, 128);
      font-size: 0.875rem;
      line-height: 1.25rem;
    }
    #sidenote ::slotted(img) {
      padding: 0.5rem 0;
      max-width: 100%;
    }
    #sidenote img:last-child {
      padding-bottom: 0;
    }
  `;

  render() {
    return html`
      <label for="sidenote"></label><span id="sidenote"><slot></slot></span>
    `;
  }

}

customElements.define("rpr-sidenote", Sidenode);
