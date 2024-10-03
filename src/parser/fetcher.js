/**
 * @typedef {object} FetcherResponse
 * @param {string} url
 * @param {Promise<Response>} response
 */

/**
 * @typedef {Object.<string, Object.<string, string>>} FetcherCache 
 */

export class Fetcher {
  /** @type {FetcherResponse[]} */
  #responses = []

  /** @type [url:string, key:string, value:string][] */
  #knownValues = []

  /** @type [url:string, key:string, value:string][] */
  #unknownValues = []

  /** @type [url:string, key:string][] */
  #seen = []

  /** @param {FetcherCache} cache */
  constructor(cache) {
    for (const [url, keys] of Object.entries(cache)) {
      for (const [key, value] of Object.entries(keys)) {
        this.#knownValues.push([url, key, value]);
      }
    }
  }

  /**
  * @param {string} url 
  * @param {string} key 
  * @param {(response:Response) => Promise<string>} handler 
  * @returns Promise<string>
  */
  async fetch(url, key, handler) {
    this.#seen.push([url, key]);

    let entry = this.#knownValues.find(e => e[0] === url && e[1] === key);
    if (entry) {
      return entry[2];
    }

    /** @type {FetcherResponse} */
    let response = this.#responses.find(r => r.url === url);

    if (!response) {
      response = { url, response: fetch(url) };
      this.#responses.push(response);
    }

    const value = await handler(await response.response);
    this.#unknownValues.push([url, key, value]);
    return value;
  }
}
