const hljs = require('highlight.js/lib/core');

hljs.registerLanguage('python', require('highlight.js/lib/languages/python'))

function highlight() {
	hljs.highlightAll();
	console.log("hello");
}
highlight()
