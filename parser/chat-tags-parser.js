const { Parser } = require('simple-text-parser');
const getLibGenSearchURI = (query) => 'https://libgen.is/search.php?req=' + encodeURIComponent(query);
const getSciHubSearchURI = (query) => 'https://sci-hub.ru/' + encodeURIComponent(query);
const getGoogleSearchURI = (query) => 'https://www.google.com/search?q=' + encodeURIComponent(query);

const chatTagsParser = new Parser();
chatTagsParser.addRule(/\[\[(\d{13}|\d{10})\|(.*?)\]\]/gi, (tag, isbn, displayText) => ({
  type: "isbn",
  text: `<a href="${getLibGenSearchURI(isbn)}" target="_blank" class="isbn">${displayText}</a>`,
  value: isbn
}));
chatTagsParser.addRule(/\[\[([^\|\]]*)(?:\|(.*?))?\]\]/gi, (tag, name, displayText) => {
  const sanitisedName = name.endsWith("'s") ? name.substring(0, name.length - 2) : name;
  return {
    type: "person",
    text: `<a href="${getLibGenSearchURI(sanitisedName)}" target="_blank" class="person">${displayText || name}</a>`,
    value: sanitisedName
  }
});

module.exports = chatTagsParser;