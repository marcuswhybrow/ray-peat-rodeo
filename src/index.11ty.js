const chatTagsParser = require('../parser/chat-tags-parser');
const { DateTime } = require('luxon');

class Index {
  data() {
    return {
      title: "Ray Peat Rodeo",
      layout: "base.njk"
    }
  }
  render(data) {
    const people = {};
    data.collections.all.forEach(item => {
      if (item.data.references) {
        if (item.data.references.people) {
          Object.entries(item.data.references.people).forEach(([name, count]) => {
            if (!people[name]) {
              people[name] = [count, [item]];
            } else {
              people[name] = [people[name][0] + count, [...people[name][1], item]];
            }
          })
        }
      }
    });
    const peopleSorted = Object.entries(people).sort((a, b) => a[0].localeCompare(b[0]));
    const peopleSortedByMentionCount = Object.entries(people).sort((a, b) => b[1][0] - a[1][0]);
    const topTenMentions = peopleSortedByMentionCount
      .slice(0, 10)
      .map(([name, [count]]) => `<span><span class="name">${name}</span> x${count}</span>`)
      .join(' — ');

    const peopleMentioned = peopleSorted.map(([name, [count, pages]]) => {
      const pageLinks = pages.map(item => `<a href="${item.url}">${item.data.title}</a>`).join(', ');
      return `
        <dt><span class="name">${name}</span> <sup>x${count}</sup></dt>
        <dd>${pageLinks}</dd>.
      `;
    }).join('\n');

    const transcripts = data.collections.content.map(item => {
      return `<li>${DateTime.fromJSDate(item.date).toFormat("yyyy LL dd")} <a href="${item.url}">${item.data.title}</a></li>`
    }).join('\n');
      
    return `
      <h1>Ray Peat Rodeo</h1>
      <p><em>Roll-up for a round-up of <a href="${data.site.github}" target="_blank">open-source</a> Ray Peat transcripts.</em></p>
      <form name="contact" netlify>
          <p>Suggest a <input type="text" name="url" placeholder="URL" /> I should <button type="submit">Transcribe</button></p>
      </form>
      
      <section id="transcripts">
        <h2>Transcriptions</h2>
        <ul>
          ${transcripts}
        </ul>
      </section>

      <section id="references">
        <h2>People Discussed By Ray Peat</h2>
        <p>
          TOP 10 — ${topTenMentions}
        </p>
        <br/>
        <dl>
          ${peopleMentioned}
        </dl>
      </section>
    `;
  }
};

module.exports = Index;