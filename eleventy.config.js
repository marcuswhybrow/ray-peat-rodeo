const interviewParser = require('./parser');
const yaml = require("js-yaml");

module.exports = function(eleventyConfig) {
  eleventyConfig.addDataExtension("yml", contents => yaml.load(contents));
  eleventyConfig.addExtension("md", {
    compile: inputContent =>
      data => data.page.inputPath.startsWith("./src/content/") ?
        interviewParser(inputContent, data) :
        this.defaultRenderer(data)
  });
  eleventyConfig.addPassthroughCopy({ "src/public": "/public" });
  return {
    dir: {
      input: "src",
      output: "_site",
    }
  };
  
};