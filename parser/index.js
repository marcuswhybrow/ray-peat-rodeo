const chatTagsParser = require('./chat-tags-parser');
const speakerParser = require('./speaker-parser');
const { Liquid } = require('liquidjs');


const liquid = new Liquid({
  extname: '.liquid',
  dynamicPartials: false,
  strictFilters: false,
  root: ['_includes']
});

const md = require('markdown-it')({
  html: true,
  linkify: true,
  typographer: true,
})
  .use(require('markdown-it-footnote'))
  .use(require('markdown-it-container'), 'speaker', {
    render: function (tokens, idx) {
      var speakerClass = tokens[idx].info.trim().match(/^speaker\s+(.*)$/);
  
      if (tokens[idx].nesting === 1) {
        return `<div class="speaker speaker-${speakerClass[1]}">\n`;
      } else {
        return `</div>\n`;
      }
    }
  });

module.exports = (inputContent, data) =>
  md.render(
    liquid.parseAndRenderSync(
      speakerParser(
        chatTagsParser(data).render(inputContent)
      ),
      data
    )
  );