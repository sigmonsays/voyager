// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package asset

var picture = Blob(asset.init(asset{Name: "picture.html", Content: "" +
	"<html>\n<head>\n   <title>Pictures in {{$.Path}}</title>\n   <link href=\"/s/style.css\" media=\"all\" rel=\"stylesheet\" />\n</head>\n<body>\n\n<h1>{{.Title}}</h1>\n\n<h3>{{$.LocalPath}}</h3>\n\n{{range $i, $f := .Files}}\n\n<div class=thumbnail><a href=\"{{$f.Url}}\"><img src=\"/image/?path={{$.LocalPath}}/{{$f.Name}}\"></img></a></div>\n\n{{end}}\n\n</body>\n" +
	""}))
