package assets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type StaticSprite struct {
	Filename string `json:"filename"`
	Image    *ebiten.Image
}

type DrawInput struct {
	Width        float64
	Height       float64
	ImageOptions *ebiten.DrawImageOptions
}

func NewStaticSprite() *StaticSprite {
	return &StaticSprite{}
}

func (a *StaticSprite) Load(spritePath string) error {
	// Carregar o tileset
	file, err := assets.Open(spritePath)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	a.Image = ebiten.NewImageFromImage(img)

	return nil
}

func (a *StaticSprite) Draw(screen *ebiten.Image, input DrawInput) {
	input.ImageOptions.GeoM.Scale(
		input.Width/float64(a.Image.Bounds().Dx()),
		input.Height/float64(a.Image.Bounds().Dy()),
	)

	screen.DrawImage(a.Image, input.ImageOptions)
}
