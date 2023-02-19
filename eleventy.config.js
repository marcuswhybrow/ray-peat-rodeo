const getLibGenSearchURI = (query) => 'https://libgen.is/search.php?req=' + encodeURIComponent(query);
const getSciHubSearchURI = (query) => 'https://sci-hub.ru/' + encodeURIComponent(query);
const getGoogleSearchURI = (query) => 'https://www.google.com/search?q=' + encodeURIComponent(query);

const renderExternalLink = (href, content) => `<a href="${href}" target="_blank">${content}</a>`;

module.exports = function(eleventyConfig) {

  eleventyConfig.addPassthroughCopy({ "src/public": "/public" });

  // Shortcode appends a Library Genesis search link
  // TODO: Build a site-wide index of all ISBNs mentioned
  eleventyConfig.addPairedShortcode("publication", function(content, title, author, year) {
    return renderExternalLink(getLibGenSearchURI(`${title} ${author.split(' ').pop()}`), content);
  });

  // Shortcode appends Sci-Hub link
  // TODO: Build a site-wide index of all DOIs mentioned
  eleventyConfig.addPairedShortcode("doi", function(content, doi) {
    return renderExternalLink(getSciHubSearchURI(doi), content);
  });

  // Shortcode appends Google search for a person
  // TODO: Build a site-wide index of all names mentioned
  // TODO: Build a duplicate name detector to catch typos
  eleventyConfig.addPairedShortcode("person", function(content, name) {
    const derivedName = name || (content.endsWith("'s") ? content.substring(0, content.length - 2) : content);
    return renderExternalLink(getLibGenSearchURI(derivedName), content);
  });

  eleventyConfig.addShortcode("ray", () => `[<b>Ray</b>]`);
  eleventyConfig.addShortcode("host", () => `[<b>Host</b>]`);

  // TODO: Build a site-wide index of dates mentioned.
  eleventyConfig.addPairedShortcode("date", (content) => content);

  return {
    dir: {
      input: "src",
      output: "_site",
    }
  };
  
};