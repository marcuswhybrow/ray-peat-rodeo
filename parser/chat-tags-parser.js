const { Parser } = require('simple-text-parser');
const yaml = require("js-yaml");
const fs = require("fs");
const getLibGenSearchURI = (query) => 'https://libgen.is/search.php?req=' + encodeURIComponent(query);
const getSciHubSearchURI = (query) => 'https://sci-hub.ru/' + encodeURIComponent(query);
const getGoogleSearchURI = (query) => 'https://www.google.com/search?q=' + encodeURIComponent(query);

const people = yaml.load(fs.readFileSync('./src/_data/people.yml', 'utf8'));

class Person {
  constructor(name) {
    this.name = name;
    const searchExternallyAs = people[this.name] ? people[this.name].searchExternallyAs || this.name : this.name;
    this.libGenURI = getLibGenSearchURI(searchExternallyAs);
    this.googleURI = getGoogleSearchURI(searchExternallyAs);
  }
};

class Book {
  constructor(isbn) {
    this.isbn = isbn;
    this.libGenURI = getLibGenSearchURI(this.isbn);
    this.googleURI = getGoogleSearchURI(this.isbn);
  }
}

const chatTagsParser = new Parser();
chatTagsParser.addRule(/\[\[(\d{13}|\d{10})\|(.*?)\]\]/gi, (tag, isbn, displayText) => {
  const book = new Book(isbn);
  return {
    type: "book",
    text: `<a href="${book.libGenURI}" target="_blank" class="isbn">${displayText}</a>`,
    value: book
  }
});
chatTagsParser.addRule(/\[\[(([^\|\]]*?)(?:['’]s?)?)(?:\|(.*?))?\]\]/gi, (tag, nameAsWritten, nameWithoutPluralisation, displayText) => {
  const person = new Person(nameWithoutPluralisation);
  return {
    type: "person",
    text: `<a href="${person.libGenURI}" target="_blank" class="person">${displayText || nameAsWritten}</a>`,
    value: person
  }
});

module.exports = chatTagsParser;