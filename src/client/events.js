export class ResultClickEvent extends Event {
  /** @type {string} */
  slug

  /** @type {string|null} */
  section

  /**
  * @param {string} slug 
  * @param {string|null} section 
  */
  constructor(slug, section = null) {
    super("result-click", { bubbles: true, composed: true });
    this.slug = slug;
    this.section = section;
  }
}

export class SearchChangeEvent extends Event {
  /** @type {string} */
  query

  /**
  * @param {string} query
  */
  constructor(query) {
    super("search-change", { bubbles: true, composed: true });
    this.query = query;
  }
}
