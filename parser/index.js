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
      var args = tokens[idx].info.trim().match(/^speaker\s+(\d*?)\s+(.*)$/);
  
      if (tokens[idx].nesting === 1) {
        return `<div class="speaker speaker-other speaker-other-${args[1]}">\n<span class="speaker-name">${args[2]}:</span>`;
      } else {
        return `</div>\n`;
      }
    }
  })
  .use(require('markdown-it-container'), 'ray', {
    render: function (tokens, idx) {
      if (tokens[idx].nesting === 1) {
        return `<div class="speaker speaker-ray">\n<span class="speaker-name">Ray Peat:</span>`;
      } else {
        return `</div>\n`;
      }
    }
  });

module.exports = (inputContent, data) =>
  md.render(
    liquid.parseAndRenderSync(
      chatTagsParser(data).render(
        speakerParser(inputContent, data),
        data
      )
    )
  );