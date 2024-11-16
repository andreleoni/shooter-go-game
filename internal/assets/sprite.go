package assets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type StaticSprite struct {
	Filename string `json:"filename"`
	Image    *ebiten.Image
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

func (a *StaticSprite) Draw(
	screen *ebiten.Image,
	x, y float64,
	invertHorizontal bool) {
	op := &ebiten.DrawImageOptions{}
	if invertHorizontal {
		op.GeoM.Scale(-1, 1)
	}

	op.GeoM.Translate(x, y)

	screen.DrawImage(a.Image, op)
}

func (a *StaticSprite) DrawAngle(
	screen *ebiten.Image,
	x, y float64,
	angle float64) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(a.Image.Bounds().Dx())/2, -float64(a.Image.Bounds().Dy())/2) // Centralizar a origem

	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)

	screen.DrawImage(a.Image, op)
}

func (a *StaticSprite) DrawWithSize(
	screen *ebiten.Image,
	x, y, width, height float64,
	invertHorizontal bool) {

	op := &ebiten.DrawImageOptions{}
	if invertHorizontal {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(width, 0) // Ajustar a posição após inverter
	}

	// Redimensionar o sprite
	op.GeoM.Scale(width/float64(a.Image.Bounds().Dx()), height/float64(a.Image.Bounds().Dy()))
	op.GeoM.Translate(x, y)

	screen.DrawImage(a.Image, op)
}
