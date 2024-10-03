/**
 * Is any part of [element] visible in the viewport. [margin] expands the area
 * considered visible in pixel units.
 *
 * @param {Element} element 
 * @param {number} margin
 * @returns {Boolean}
 */
export function isVisible(element, margin = 0) {
  const rect = element.getBoundingClientRect();

  const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
  const viewportwidth = window.innerWidth || document.documentElement.clientWidth;

  // Visual Reference:
  // - https://developer.mozilla.org/en-US/docs/Web/API/Element/getBoundingClientRect#return_value

  if (rect.bottom < margin) return false;
  if (rect.top > viewportHeight + margin) return false;
  if (rect.left > viewportwidth + margin) return false;
  if (rect.right < margin) return false;
  return true;
}


/**
 * @param {Result} a 
 * @param {Result} b
 * @returns {number}
 */
export function resultsByScore(a, b) {
  return (b.score || 0) - (a.score || 0);
}

/**
 * @param {Result[]} localResults
 * @param {PagefindResult[]} pagefindResults 
 * @param {string[]} pagefindResuldIdList
 * @returns {Promise<Result[]>} The new results shape
 */
export async function loadResultData(localResults, pagefindResults, pagefindResuldIdList) {
  const results = structuredClone(localResults);

  const promises = pagefindResuldIdList.map(async id => {
    const localIndex = results.findIndex(result => result.pagefindResultId === id);
    const localResult = results[localIndex];
    if (!localResult) return;
    if (localResult.loaded) return;

    const pagefindIndex = pagefindResults.findIndex(pagefindResult => {
      return pagefindResult.id === id;
    });

    if (pagefindIndex === -1) return;

    /** @type {PagefindResult} */
    const pagefindResult = pagefindResults[pagefindIndex];


    const data = await pagefindResult.data();

    const preludeSubResult = data?.sub_results.find(s => !s.hasOwnProperty("anchor"));
    localResult.excerpt = preludeSubResult?.excerpt || "";

    walkSections(localResult.sections, section => {
      const subResult = data?.sub_results.find(s => (s.anchor?.id || "") === section.id);
      section.excerpt = subResult?.excerpt || "";
    });
    localResult.loaded = true;
  });

  await Promise.all([promises]);
  return results;
};

/**
  * @param {Section[]} sections 
  * @param {(section:Section) => void} callback
  */
function walkSections(sections, callback) {
  if (!sections) return;
  for (const section of sections) {
    callback(section);
    walkSections(section.subsections, callback);
  }
}

/**
  * @param {Filter[]} filterArr
  * @returns {Object.<String, String[]>}
  */
export function pagefindFilters(filterArr) {
  const filters = {};
  filterArr.forEach(filter => {
    if (filters.hasOwnProperty(filter.key)) {
      filters[filter.key].push(filter.value);
    } else {
      filters[filter.key] = [filter.value];
    }
  });
  return filters;
}

/**
  * @param {Section[]} sections
  * @param {string} sectionId
  * @returns {Section|undefined}
  */
export function findSection(sections, sectionId) {
  if (sectionId === "") return undefined;
  for (const section of sections) {
    if (section.id === sectionId) return section;
    return findSection(section.subsections, sectionId);
  }
}
