const interviewParser = require('./parser');
const slugify = require('slugify');

module.exports = function(eleventyConfig) {
  eleventyConfig.addExtension("md", {
    compile: (inputContent) => function(data) {
      if (data.page.inputPath.startsWith("./src/content/")) {
        return interviewParser(inputContent, data);
      } else {
        return this.defaultRenderer(data);
      }
    },
    compileOptions: {
      permalink: function(contents, inputPath) {
        return (data) => {
          if (data.page.inputPath.startsWith("./src/content/")) {
            return `${slugify(data.title, { lower: true })}/`;
          } else {
            return data.permalink;
          }
        }
      }
    }
  });
  eleventyConfig.addPassthroughCopy({ "src/public": "/public" });
  return {
    dir: {
      input: "src",
      output: "_site",
    }
  };
  
};