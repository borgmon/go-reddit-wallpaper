package main

type RedditPayload struct {
	Data PayloadData
}

type PayloadData struct {
	Children []PayloadDataChild
}

type PayloadDataChild struct {
	Data PayloadDataChildData
}

type PayloadDataChildData struct {
	Preview PayloadDataChildDataPreview
	Url     string
}

type PayloadDataChildDataPreview struct {
	Images []PayloadDataChildDataPreviewImage
}
type PayloadDataChildDataPreviewImage struct {
	Source PayloadDataChildDataPreviewImageSource
}

type PayloadDataChildDataPreviewImageSource struct {
	Width  int16
	Height int16
}
