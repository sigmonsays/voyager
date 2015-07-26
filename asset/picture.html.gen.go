// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package asset

var picture = Blob(asset.init(asset{Name: "picture.html", Content: "" +
	"<html>\n<head>\n   <title>Pictures in {{$.Path}}</title>\n   <link href=\"/s/style.css\" media=\"all\" rel=\"stylesheet\" />\n</head>\n<body>\n<div class=page>\n\n{{range $i, $b := .Breadcrum}}\n<a href=\"{{$b.Url}}\">{{$b.Name}}</a> \n{{end}}\n\n<h2>{{.Title}} - {{$.LocalPath}}</h2>\n\n{{if ..Directories}}\n<div class=directories>\n   {{range $i, $f := .Directories}}\n   <a href=\"{{$f.Url}}\">{{$f.Basename}}</a> &nbsp;\n   {{end}}\n   <p>\n</div>\n{{end}}\n\n<div class=container>\n   {{range $i, $f := .Files}}\n   <div class=thumbnail>\n      <a href=\"{{$f.Url}}\"><img src=\"/image/?path={{$.LocalPath}}/{{$f.Name}}\"></img></a>\n   </div>\n   {{end}}\n</div>\n\n\n</div>\n</body>\n</html>\n" +
	""}))
