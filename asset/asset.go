package asset

//go:generate becky -wrap Blob style.css favicon.ico

type blob struct {
	asset
}

func Blob(a asset) blob {
	return blob{a}
}
