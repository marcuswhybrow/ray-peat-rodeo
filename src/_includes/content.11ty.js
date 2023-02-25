const { DateTime } = require('luxon');
const fs = require('fs');
const frontMatter = require('front-matter');
const chatTagsParser = require('../../parser/chat-tags-parser');

const countNodes = (nodes, counts) =>
  Object.entries(counts).reduce((results, [countName, [nodeType, nodeKeyMap]]) => ({
    ...results,
    [countName]: nodes
      .filter(node => node.type === nodeType)
      .map(nodeKeyMap)
      .reduce((countResults, key) => ({...countResults, [key]: (countResults[key] || 0) + 1}), {})
  }), {});

class Content {
  data(data) {
    return {
      layout: 'base.njk',
      permalink: data => `${data.page.fileSlug.toLowerCase()}/`,
      tags: 'content',
      eleventyComputed: {
        references: data => {
          const markdown = frontMatter(fs.readFileSync(data.page.inputPath, 'utf8'));
          const nodes = chatTagsParser(data).toTree(markdown.body);
          return countNodes(nodes, {
            people: ['person', node => node.value.indexKey],
            books: ['book', node => node.value.bibliographKey],
            externalLinks: ['external-link', node => node.value],
            sciencePapers: ['science-paper', node => node.value.doi]
          });
        }
      }
    }
  }

  render(data) {
    const transcriptionLink = data.transcription.source ? `<a href="${data.transcription.source}" target="_blank">Transcribed</a>` : `Transcribed`;
    const date = DateTime.fromJSDate(data.date);
    const transcriptionDate = DateTime.fromJSDate(data.transcription.date);
    return `
      <article class="interview">
        <div id="top-bar"><a href="/">Ray Peat Rodeo</a></div>
        <header>
          <h1>${data.title}</h1>
          <p>
            <a href="${data.source}" target="_blank">Convened</a> by ${data.series}, ${date.toFormat("LLLL dd, yyyy")}. ${transcriptionLink} by ${data.transcription.author}, ${transcriptionDate.toFormat("LLLL dd, yyyy")}.
            <br />
            <a href="${data.site.githubEdit}/${data.page.inputPath}" target="_blank">Edit this page on GiHub</a> to contribute corrections of any kind.
          </p>
        </header>
        <main>
          ${data.content}
        </main>
      </article>
    `;
  }
}

module.exports = Content;