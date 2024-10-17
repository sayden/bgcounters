package output

import (
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

const (
	padding  = 2
	maxWidth = 80
)

type progressBar struct {
	total    float64
	percent  float64
	progress progress.Model
}

func (m *progressBar) Init() tea.Cmd {
	return nil
}

func (m *progressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case int:
		m.percent = float64(msg) / m.total
		if m.percent >= 1.0 {
			m.percent = 1.0
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width >= maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	default:
		return m, tea.Quit
	}
}

func (m *progressBar) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n"
}

type GlobalState struct {
	counterPos     int
	row            int
	columns        int
	fileNumber     int
	filenamesInUse map[string]bool
	template       *counters.CounterTemplate
	canvas         *gg.Context
}

func newGlobalState(template *counters.CounterTemplate, canvas *gg.Context) *GlobalState {
	return &GlobalState{
		filenamesInUse: make(map[string]bool),
		fileNumber:     1,
		canvas:         canvas,
		template:       template,
	}
}

func CountersToPNG(template *counters.CounterTemplate) error {
	var canvas *gg.Context
	if template.Mode == counters.TEMPLATE_MODE_TILES {
		// Create canvas
		canvas = gg.NewContext(template.Columns*template.Width, template.Rows*template.Height)
		if err := canvas.LoadFontFace(counters.DEFAULT_FONT_PATH, template.FontHeight); err != nil {
			log.WithField("font_path", counters.DEFAULT_FONT_PATH).Error(err)
			return errors.Wrap(err, "could not load font")
		}
	}

	var (
		total = 0

		counterPos = 0
		row        = 0
		// fileNumber = 1
	)

	gs := newGlobalState(template, canvas)

	// Progress bar
	for _, c := range template.Counters {
		if !c.Skip {
			total += c.Multiplier
		}
	}
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	program := tea.NewProgram(&progressBar{progress: prog, total: float64(total)})
	go program.Run()
	defer program.Quit()

	//Iterate rows and columns, painting a counter on each
iteration:
	for {
		// for columns := 0; columns < template.Columns; columns++ {
		for gs.columns = 0; gs.columns < template.Columns; gs.columns++ {
			if counterPos >= len(template.Counters) {
				break iteration
			}

			counter := template.Counters[counterPos]
			if counter.Skip {
				counterPos++
				continue
			}

			counterCanvas, err := GetCounterCanvas(&counter, template)
			if err != nil {
				return err
			}

			err = writeCounterToFile(counterCanvas, counter, gs)
			if err != nil {
				return err
			}

			// Update the progress bar
			program.Send(counterPos)
			counterPos++
		}

		row++
	}

	// Save result on a file
	if template.Mode == counters.TEMPLATE_MODE_TILES {
		return canvas.SavePNG(template.OutputFolder)
	}

	program.Send(100)
	program.Wait()

	return nil
}

type imageWriter struct {
	dc       *gg.Context
	template *counters.CounterTemplate
}

func writeCounterToFile(dc *gg.Context, counter counters.Counter, gs *GlobalState) error {
	iw := imageWriter{
		dc:       dc,
		template: gs.template,
	}
	var suffix string
	if counter.Back != nil {
		suffix = "-front"
	}

	switch gs.template.Mode {
	case counters.TEMPLATE_MODE_TEMPLATE:
		err := iw.createFile(&counter, gs, suffix)
		if err != nil {
			return errors.Wrap(err, "error trying to write counter file")
		}

		if counter.Back != nil {
			suffix = "-back"
			if err = iw.createFile(counter.Back, gs, suffix); err != nil {
				return errors.Wrap(err, "error trying to write counter file")
			}
			gs.fileNumber++
		}
	case counters.TEMPLATE_MODE_TILES:
		if counter.Back != nil {
			iw.createTiledFile(counter.Back, gs)
		}
		iw.createTiledFile(&counter, gs)
	default:
		return errors.New("unknown template mode or template mode is empty")
	}

	return nil
}

// createFile creates a file with the counter image. Filenumber is the filename, a pointer is passed to be able to use
// the multiplier to create more than one file with the same counter
func (iw *imageWriter) createFile(counter *counters.Counter, gs *GlobalState, suffix string) error {
	_ = os.MkdirAll(iw.template.OutputFolder, 0750)

	if counter.Multiplier == 0 {
		counter.Multiplier = 1
	}

	// Use sequencing of numbers or a position in the counter texts to name files
	for i := 0; i < counter.Multiplier; i++ {
		counterFilename := counter.GetCounterFilename(iw.template.PositionNumberForFilename, suffix, gs.fileNumber, gs.filenamesInUse)
		filepath := path.Join(iw.template.OutputFolder, counterFilename)
		if counter.Skip {
			continue
		}

		log.Debug("Saving file: ", filepath)
		if err := iw.dc.SavePNG(filepath); err != nil {
			log.WithField("file", filepath).Error("could not save PNG file")
			return err
		}
		gs.fileNumber++
	}

	return nil
}

// createTiledFile creates a file with the counter image. Filenumber is the filename, a row and column pointer
// is passed to be able to use the multiplier to create more than one counter in the same sheet
func (iw *imageWriter) createTiledFile(counter *counters.Counter, gs *GlobalState) {
	for i := 0; i < counter.Multiplier; i++ {
		gs.canvas.DrawImage(iw.dc.Image(), gs.columns*counter.Width, gs.row*counter.Height)
		gs.columns++
		if gs.columns == iw.template.Columns {
			gs.columns = 0
			gs.row++
		}
	}
	gs.columns--
}

// GetCanvasForCounter creates a canvas for a counter definition.
func GetCanvasForCounter(cc *counters.Counter, t *counters.CounterTemplate) (*gg.Context, error) {
	dc := gg.NewContext(cc.Width, cc.Height)
	if err := dc.LoadFontFace(cc.FontPath, cc.FontHeight); err != nil {
		log.WithFields(log.Fields{"font": "'" + cc.FontPath + "'", "height": cc.FontHeight}).Error(err)
		return nil, err
	}

	dc.Push()
	dc.SetColor(cc.BgColor)
	dc.DrawRectangle(0, 0, float64(cc.Width), float64(cc.Height))
	dc.Fill()
	dc.Pop()

	if cc.FontColorS != "" {
		counters.GetValidColorForString(cc.FontColorS, t.BgColor)
	}

	return dc, nil
}
