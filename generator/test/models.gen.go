// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package test

type Book interface {
	IsBook()
	GetTitle() string
	GetAuthor() *Author
}

type MediaItem interface {
	IsMediaItem()
}

type AudioClip struct {
	Duration int    `json:"duration"`
	Typename string `json:"__typename"`
}

func (AudioClip) IsMediaItem() {}

type Author struct {
	Name string `json:"name"`
}

type Library struct {
	Books []Book `json:"books"`
}

type TextBook struct {
	Title                 string      `json:"title"`
	Author                *Author     `json:"author"`
	SupplementaryMaterial []MediaItem `json:"supplementaryMaterial"`
	Typename              string      `json:"__typename"`
}

func (TextBook) IsBook()                 {}
func (this TextBook) GetTitle() string   { return this.Title }
func (this TextBook) GetAuthor() *Author { return this.Author }

type VideoClip struct {
	PreviewURL string `json:"previewURL"`
	Typename   string `json:"__typename"`
}

func (VideoClip) IsMediaItem() {}
