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
	X            float64
	Y            float64
	ImageOptions *ebiten.DrawImageOptions
	Angle        *float64
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
	if input.ImageOptions == nil {
		input.ImageOptions = &ebiten.DrawImageOptions{}
	}

	if input.Angle != nil {
		input.ImageOptions.GeoM.Translate(-float64(a.Image.Bounds().Dx())/2, -float64(a.Image.Bounds().Dy())/2)

		input.ImageOptions.GeoM.Rotate(*input.Angle)

		input.ImageOptions.GeoM.Translate(float64(a.Image.Bounds().Dx())/2, float64(a.Image.Bounds().Dy())/2)
	}

	input.ImageOptions.GeoM.Scale(
		input.Width/float64(a.Image.Bounds().Dx()),
		input.Height/float64(a.Image.Bounds().Dy()),
	)

	input.ImageOptions.GeoM.Translate(input.X, input.Y)

	screen.DrawImage(a.Image, input.ImageOptions)
}
