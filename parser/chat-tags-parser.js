const { Parser } = require('simple-text-parser');
const yaml = require("js-yaml");
const fs = require("fs");
const getLibGenSearchURL = (query) => 'https://libgen.is/search.php?req=' + encodeURIComponent(query);
const getSciHubSearchURL = (query) => 'https://sci-hub.ru/' + encodeURIComponent(query);
const getGoogleSearchURL = (query) => 'https://www.google.com/search?q=' + encodeURIComponent(query);

const people = yaml.load(fs.readFileSync('./src/_data/people.yml', 'utf8'));
const dois = yaml.load(fs.readFileSync('./src/_data/doi.yml', 'utf8'));

class Person {
  constructor(name) {
    this.name = name;
    const searchExternallyAs = people[this.name] ? people[this.name].searchExternallyAs || this.name : this.name;
    this.libGenURL = getLibGenSearchURL(searchExternallyAs);
    this.googleURL = getGoogleSearchURL(searchExternallyAs);
  }
};

class Book {
  constructor(isbn) {
    this.isbn = isbn;
    this.libGenURL = getLibGenSearchURL(this.isbn);
    this.googleURL = getGoogleSearchURL(this.isbn);
  }
}

class SciencePaper {
  constructor(doi) {
    this.doi = doi;
    this.sciHubURL = getSciHubSearchURL(this.doi);
    this.url = dois[this.doi] ? dois[this.doi].url || this.sciHubURL : this.sciHubURL;
  }
}

const computeNode = (pipeAndDisplayText, fallback, f) => {
  if (pipeAndDisplayText) {
    if (pipeAndDisplayText.length >= 1) return f(pipeAndDisplayText.substring(1) || fallback);
    return { type: "text", text: "" };
  }
  return f(fallback);
}

const chatTagsParser = new Parser();
chatTagsParser.addRule(/\[\[(\|.*)\]\]/gi, (tag, pipeAndDisplayText) => computeNode(pipeAndDisplayText, "", (displayText) => ({
  type: "internal-link-broken",
  text: displayText
})));
chatTagsParser.addRule(/\[\[(\d{13}|\d{10})(\|.*)?\]\]/gi, (tag, isbn, pipeAndDisplayText) => {
  const book = new Book(isbn);
  return computeNode(pipeAndDisplayText, isbn, (displayText) => ({
    type: "book",
    text: `<a href="${book.libGenURL}" target="_blank" class="isbn">${displayText}</a>`,
    value: book
  }));
});
chatTagsParser.addRule(/\[\[doi\:(.*?)(\|.*?)?\]\]/gi, (tag, doi, pipeAndDisplayText) => {
  const sciencePaper = new SciencePaper(doi);
  return computeNode(pipeAndDisplayText, sciencePaper.url, (displayText) => ({
    type: "science-paper",
    text: `<a href="${sciencePaper.url}" target="_blank" class="science-paper">${displayText}</a>`,
    value: sciencePaper
  }));
});
chatTagsParser.addRule(/\[\[(https?\:\/\/.*?)(\|.*)?\]\]/gi, (tag, url, pipeAndDisplayText) => {
  return computeNode(pipeAndDisplayText, url, (displayText) => ({
    type: "external-link",
    text: `<a href="${url}" target="_blank" class="external">${displayText}</a>`,
    value: url
  }));
});
chatTagsParser.addRule(/\[\[(([^\|\]]*?)(?:['’]s?)?)(\|(.*?))?\]\]/gi, (tag, nameAsWritten, nameWithoutPluralisation, pipeAndDisplayText) => {
  const person = new Person(nameWithoutPluralisation);
  return computeNode(pipeAndDisplayText, nameAsWritten, (displayText) => ({
    type: "person",
    text: `<a href="${person.libGenURL}" target="_blank" class="person">${displayText}</a>`,
    value: person
  }));
});

module.exports = chatTagsParser;