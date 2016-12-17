var mimeTypes = require("./mimetypes.json")

var allMimeTypes = {}

for (var i = 0; i < mimeTypes.length; i++) {
	for (var k in mimeTypes[i]) {
		var val = mimeTypes[i][k];
		var constName = val.charAt(0).toUpperCase() + val.slice(1).replace(/\//g, "_").replace(/\./g, "_").replace(/\+/g, "_").replace(/-/g, "_")

		var underScoreIndex = constName.indexOf("_");

		while (underScoreIndex != -1 && underScoreIndex < (constName.length - 1)) {
			var orig = "_" + constName.charAt(underScoreIndex + 1)
			var replacement = constName.charAt(underScoreIndex + 1).toUpperCase()
			constName = constName.replace(orig, replacement)

			underScoreIndex = constName.indexOf("_");
		}
		allMimeTypes[constName] = val
	}
}

for (var constName in allMimeTypes) {
	console.log("// " + constName + " is the MIME type for '" + allMimeTypes[constName] + "'")
	console.log(constName + " MIMEType = \"" + allMimeTypes[constName] + "\"")
}