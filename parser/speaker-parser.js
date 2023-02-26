const defaultSpeakerKey = (speakers = {}) => {
  const speakerEntries = Object.entries(speakers);
  switch (speakerEntries.length) {
    case 0:
      speakers.H = 'Host';
      return 'H';
    case 1:
      return speakerEntries[0][0];
    default:
      return null
  }
};

const SPEAKER_PARAGRAPH_REGEX = /(?<speakerDefinition>!(?<speakerKey>\S*))?\s*(?<paragraph>.*)$/;

const fillSpeakerKeys = (lines, data) => {
  const speakers = data.speakers || {};
  const otherSpeakers = [];
  let prevSpeakerKey = null;
  const results = [];
  lines.forEach(line => {
    const { speakerDefinition, speakerKey, paragraph } = SPEAKER_PARAGRAPH_REGEX.exec(line).groups;
    if (speakerDefinition) {
      let computedSpeakerKey = speakerKey || prevSpeakerKey || defaultSpeakerKey(speakers);
      if (computedSpeakerKey) {
        if (!otherSpeakers.includes(computedSpeakerKey))
          otherSpeakers.push(computedSpeakerKey);
        results.push({
          type: 'Other Speaker',
          speakerKey: computedSpeakerKey,
          speakerName: speakers[computedSpeakerKey],
          speakerNumber: otherSpeakers.indexOf(computedSpeakerKey),
          text: paragraph
        });
      } else {
        throw new Error(`${data.page.inputPath}: Cannot resolve speaker:\n${line}`);
      }
      prevSpeakerKey = computedSpeakerKey;
    } else {
      if (paragraph) {
        results.push({
          type: 'Ray Peat',
          text: paragraph
        });
      } else {
        results.push({
          type: 'Empty',
          text: line
        });
      }
    }
  });
  return results;
};

const groupLines = (lines) => {
  if (lines.length === 0) return lines;

  const results = [lines.shift()];
  lines.forEach(line => {
    const tailIndex = results.length - 1;
    const tail = results[tailIndex];
    if (
      line.type === 'Empty' ||
      (line.type === 'Ray Peat' && tail.type === 'Ray Peat') ||
      (line.type === 'Other Speaker' && tail.type === "Other Speaker" && line.speakerKey === tail.speakerKey)
    ) {
      results[tailIndex] = {...tail, text: `${tail.text}\n${line.text}`};
    } else {
      results.push(line);
    }
  });
  return results;
};

module.exports = (inputContent, data) => {
  const lines = fillSpeakerKeys(inputContent.split(/\r?\n/), data);
  const groupedLines = groupLines(lines);
  const sections = groupedLines.map(line => {
    const sectionName = line.type === 'Ray Peat'
      ? 'ray'
      : line.type === 'Other Speaker'
        ? `speaker ${line.speakerNumber} ${line.speakerName}`
        : null;
    if (sectionName) {
      return  `::: ${sectionName}\n${line.text}\n:::`;
    } else {
      return line.text;
    }
  });
  return sections.join('\n');
};