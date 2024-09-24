package counters

import (
	"github.com/fogleman/gg"
	"image"
)

func DrawGuides(cc Settings) (*image.Image, error) {
	var p Positioner
	width := float64(cc.Width)
	height := float64(cc.Height)

	temp := gg.NewContext(cc.Width, cc.Height)
	temp.SetRGBA(1, 0, 0, 0.6)

	//Vertical
	x, _, err := p.GetXYPosition(2, cc)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x, 0, x, height)
	temp.Stroke()
	x, _, err = p.GetXYPosition(4, cc)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(x, 0, x, height)
	temp.Stroke()

	//Horizontal
	_, y, err := p.GetXYPosition(16, cc)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(0, y, width, y)
	temp.Stroke()
	_, y, err = p.GetXYPosition(14, cc)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(0, y, width, y)
	temp.Stroke()

	//Horizontal center
	x, _, err = p.GetXYPosition(0, cc)
	if err != nil {
		return nil, err
	}

	temp.SetRGBA(0, 0, 1, 0.6)
	temp.DrawLine(x, 0, x, height)
	temp.Stroke()

	// Vertical center
	_, y, err = p.GetXYPosition(0, cc)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(0, y, width, y)
	temp.Stroke()

	temp.SetRGBA(0, 1, 1, 0.8)

	//Horizontal margins
	x1, y1, err := p.GetXYPosition(1, cc)
	if err != nil {
		return nil, err
	}
	x2, y2, err := p.GetXYPosition(5, cc)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()
	x1, y1, err = p.GetXYPosition(13, cc)
	if err != nil {
		return nil, err
	}

	x2, y2, err = p.GetXYPosition(9, cc)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	//Vertical margins
	x1, y1, err = p.GetXYPosition(1, cc)
	if err != nil {
		return nil, err
	}
	x2, y2, err = p.GetXYPosition(13, cc)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	x1, y1, err = p.GetXYPosition(5, cc)
	if err != nil {
		return nil, err
	}
	x2, y2, err = p.GetXYPosition(9, cc)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	temp.SetRGB(0, 0, 0)

	img := temp.Image()

	return &img, nil
}
