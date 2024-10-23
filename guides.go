package counters

import (
	"image"

	"github.com/fogleman/gg"
)

func DrawGuides(s *Settings) (*image.Image, error) {
	var p Positioner
	width := float64(s.Width)
	height := float64(s.Height)

	temp := gg.NewContext(s.Width, s.Height)
	temp.SetRGBA(1, 0, 0, 0.6)

	//Vertical
	x, _, err := p.GetXYPosition(2, s)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x, 0, x, height)
	temp.Stroke()
	x, _, err = p.GetXYPosition(4, s)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(x, 0, x, height)
	temp.Stroke()

	//Horizontal
	_, y, err := p.GetXYPosition(16, s)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(0, y, width, y)
	temp.Stroke()
	_, y, err = p.GetXYPosition(14, s)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(0, y, width, y)
	temp.Stroke()

	//Horizontal center
	x, _, err = p.GetXYPosition(0, s)
	if err != nil {
		return nil, err
	}

	temp.SetRGBA(0, 0, 1, 0.6)
	temp.DrawLine(x, 0, x, height)
	temp.Stroke()

	// Vertical center
	_, y, err = p.GetXYPosition(0, s)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(0, y, width, y)
	temp.Stroke()

	temp.SetRGBA(0, 1, 1, 0.8)

	//Horizontal margins
	x1, y1, err := p.GetXYPosition(1, s)
	if err != nil {
		return nil, err
	}
	x2, y2, err := p.GetXYPosition(5, s)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()
	x1, y1, err = p.GetXYPosition(13, s)
	if err != nil {
		return nil, err
	}

	x2, y2, err = p.GetXYPosition(9, s)
	if err != nil {
		return nil, err
	}

	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	//Vertical margins
	x1, y1, err = p.GetXYPosition(1, s)
	if err != nil {
		return nil, err
	}
	x2, y2, err = p.GetXYPosition(13, s)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	x1, y1, err = p.GetXYPosition(5, s)
	if err != nil {
		return nil, err
	}
	x2, y2, err = p.GetXYPosition(9, s)
	if err != nil {
		return nil, err
	}
	temp.DrawLine(x1, y1, x2, y2)
	temp.Stroke()

	temp.SetRGB(0, 0, 0)

	img := temp.Image()

	return &img, nil
}
