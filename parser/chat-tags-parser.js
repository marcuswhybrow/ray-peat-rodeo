const { Parser } = require('simple-text-parser');
const yaml = require("js-yaml");
const fs = require("fs");
const getLibGenSearchURL = (query) => 'https://libgen.is/search.php?req=' + encodeURIComponent(query);
const getSciHubSearchURL = (query) => 'https://sci-hub.ru/' + encodeURIComponent(query);
const getGoogleSearchURL = (query) => 'https://www.google.com/search?q=' + encodeURIComponent(query);

const people = yaml.load(fs.readFileSync('./src/_data/people.yml', 'utf8'));
const dois = yaml.load(fs.readFileSync('./src/_data/doi.yml', 'utf8'));
const books = yaml.load(fs.readFileSync('./src/_data/books.yml', 'utf8'));

class Person {
  constructor(name) {
    this.name = name;
    this.searchExternallyAs = people[this.name] ? people[this.name].searchExternallyAs || this.name : this.name;
    this.libGenURL = getLibGenSearchURL(this.searchExternallyAs);
    this.googleURL = getGoogleSearchURL(this.searchExternallyAs);
  }
};

class Book {
  constructor(title, primaryAuthorFullName) {
    this.title = title,
    this.author = primaryAuthorFullName;
    const person = new Person(this.author);
    const query = `${this.title} ${person.searchExternallyAs}`;
    this.libGenURL = getLibGenSearchURL(query);
    this.googleURL = getGoogleSearchURL(query);

    const key = `${title} -by- ${primaryAuthorFullName}`;
    console.log(key, books[key]);
    const book = books[key];
    if (book) {
      this.url = book.openAsURL || this.libGenURL;
      this.linkTitle = book.openAsURLMessage;
    } else {
      this.url = this.libGenURL;
    }
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

module.exports = (data) => {
  const chatTagsParser = new Parser();
  chatTagsParser.addRule(/\[\[(\|.*?)\]\]/gi, (tag, pipeAndDisplayText) => computeNode(pipeAndDisplayText, "", (displayText) => ({
    type: "internal-link-broken",
    text: displayText
  })));
  chatTagsParser.addRule(/\[\[([^\]\|]*?)\s*-by-\s*([^\]\|]*?)(\|[^\]]*?)?\]\]/gi, (tag, bookTitle, primaryAuthorFullName, pipeAndDisplayText) => {
    const book = new Book(bookTitle, primaryAuthorFullName);
    return computeNode(pipeAndDisplayText, bookTitle, (displayText) => ({
      type: "book",
      text: `<a href="${book.url}" target="_blank" class="book" title="${book.linkTitle}">${displayText}</a>`,
      value: book
    }));
  });
  chatTagsParser.addRule(/\[\[doi\:([^\]\|]*?)(\|[^\]]*?)?\]\]/gi, (tag, doi, pipeAndDisplayText) => {
    const sciencePaper = new SciencePaper(doi);
    return computeNode(pipeAndDisplayText, sciencePaper.url, (displayText) => ({
      type: "science-paper",
      text: `<a href="${sciencePaper.url}" target="_blank" class="science-paper">${displayText}</a>`,
      value: sciencePaper
    }));
  });
  chatTagsParser.addRule(/\[\[(https?\:\/\/[^\]\|]*?)(\|[^\]]*?)?\]\]/gi, (tag, url, pipeAndDisplayText) => {
    return computeNode(pipeAndDisplayText, url, (displayText) => ({
      type: "external-link",
      text: `<a href="${url}" target="_blank" class="external">${displayText}</a>`,
      value: url
    }));
  });
  chatTagsParser.addRule(/\[\[(([^\|\]]*?)(?:['’]s?)?)(\|([^\]]*?))?\]\]/gi, (tag, nameAsWritten, nameWithoutPluralisation, pipeAndDisplayText) => {
    const person = new Person(nameWithoutPluralisation);
    return computeNode(pipeAndDisplayText, nameAsWritten, (displayText) => ({
      type: "person",
      text: `<a href="${person.libGenURL}" target="_blank" class="person">${displayText}</a>`,
      value: person
    }));
  });
  chatTagsParser.addRule(/\[(\d+)\:(\d+)(?:\:(\d+))?\]/gi, (tag, hoursStr, minutesStr, secondsStr) => {
    if (!secondsStr) {
      // Handle case: two inputs (00:00). Not three (00:00:00)
      secondsStr = minutesStr;
      minutesStr = hoursStr;
      hoursStr = '00';
    }

    // Format numbers
    hoursStr = hoursStr.padStart(2, '0');
    hoursInt = parseInt(hoursStr);
    minutesStr = minutesStr.padStart(2, '0');
    secondsStr = secondsStr.padStart(2, '0');

    const youTubeFormat = hoursInt ? `${hoursStr}h${minutesStr}m${secondsStr}s` : `${minutesStr}m${secondsStr}s`;
    const timecode = {
      youTubeFormat,
      localFormat: hoursInt ? `${hoursStr}:${minutesStr}:${secondsStr}` : `${minutesStr}:${secondsStr}`,
      originalURL: data.source,

      // TODO: Don't assume source URL is YouTube
      // TODO: Don't assume source URL has no hash
      url: `${data.source}#t=${youTubeFormat}`
    };

    return ({
      type: "timecode",
      text: `<a href="${timecode.url}" target="_blank" class="timecode">${timecode.localFormat}</a>`,
      value: timecode
    })
  })
  return chatTagsParser;
};