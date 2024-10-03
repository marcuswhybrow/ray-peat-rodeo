/**
  * @param {Section} section
  * @returns {Section}
  */
export function lastSection(section) {
  if (section.subsections.length === 0) return section;
  return lastSection(section.subsections[section.subsections.length - 1]);
}

/**
 * @param {Asset} parsed
 * @returns {Section|undefined}
 */
export function findLastSection(parsed) {
  if (parsed.sections.length === 0) return undefined;
  let section = parsed.sections[parsed.sections.length - 1];
  while (section.subsections.length > 0) {
    section = section.subsections[section.subsections.length - 1];
  }
  return section;
}

/**
 * @param {Asset} parsed
 * @param {number} depth
 * @return {Section[]|undefined}
 */
export function findSectionParent(parsed, depth) {
  let sections = parsed.sections;
  for (let currentDepth = 2; currentDepth < depth; currentDepth++) {
    if (sections.length === 0) return undefined;
    sections = sections[sections.length - 1].subsections;
  }
  return sections;
}
