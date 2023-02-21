const SPEAKER_PARAGRAPH_REGEX = /(?<speakerDefinition>!(?<speaker>\S*))?\s*(?<paragraph>.*)$/;

module.exports = (inputContent) => {
const lines = inputContent.split(/\r?\n/).map(line => {
    const { speakerDefinition, speaker, paragraph } = SPEAKER_PARAGRAPH_REGEX.exec(line).groups;
    if (paragraph) {
    return { speaker: speakerDefinition ? speaker || "Host" : "Ray Peat", text: paragraph }
    }
    return { speaker: null, text: line };
});

let prevSpeaker = null;
const linesAndContainers = [];
lines.forEach(line => {
    if (line.speaker && line.speaker !== prevSpeaker) {
    if (prevSpeaker)
        linesAndContainers.push(`:::\n`);
    linesAndContainers.push(`::: speaker ${line.speaker === "Ray Peat" ? "ray" : "other"}\n`);
    prevSpeaker = line.speaker;
    }
    linesAndContainers.push(line.text);
});
linesAndContainers.push(':::');

return linesAndContainers.join('\n');
};