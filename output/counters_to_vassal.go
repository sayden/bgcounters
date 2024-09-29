package output

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
)

type gameModule struct {
	XMLName                    xml.Name `xml:"VASSAL.build.GameModule"`
	ModuleOther1               string   `xml:"ModuleOther1,attr"`
	ModuleOther2               string   `xml:"ModuleOther2,attr"`
	VassalVersion              string   `xml:"VassalVersion,attr"`
	Description                string   `xml:"description,attr"`
	Name                       string   `xml:"name,attr"`
	NextPieceSlotId            string   `xml:"nextPieceSlotId,attr"`
	Version                    string   `xml:"version,attr"`
	BasicCommandEncoder        capture  `xml:"VASSAL.build.module.BasicCommandEncoder"`
	Documentation              capture  `xml:"VASSAL.build.module.Documentation"`
	Chatter                    capture  `xml:"VASSAL.build.module.Chatter"`
	KeyNamer                   capture  `xml:"VASSAL.build.module.KeyNamer"`
	PieceWindow                PieceWindow
	DiceButton                 []DiceButton `xml:"VASSAL.build.module.DiceButton"`
	PlayerRoster               capture      `xml:"VASSAL.build.module.PlayerRoster"`
	GlobalOptions              capture      `xml:"VASSAL.build.module.GlobalOptions"`
	GamePieceDefinitions       capture      `xml:"VASSAL.build.module.gamepieceimage.GamePieceImageDefinitions"`
	GlobalProperties           capture      `xml:"VASSAL.build.module.properties.GlobalProperties"`
	GlobalTranslatableMessages capture      `xml:"VASSAL.build.module.properties.GlobalTranslatableMessages"`
	PrototypesContainer        capture      `xml:"VASSAL.build.module.PrototypesContainer"`
	Language                   capture      `xml:"VASSAL.i18n.Language"`
	Map                        capture      `xml:"VASSAL.build.module.Map"`
}

type capture struct {
	Raw string `xml:",innerxml"`
}
type DiceButton struct {
	Raw          string `xml:",innerxml"`
	AddToTotal   int    `xml:"addToTotal,attr"`
	CanDisable   bool   `xml:"canDisable,attr"`
	DisabledIcon string `xml:"disabledIcon,attr"`
	Hotkey       string `xml:"hotkey,attr"`
	Icon         string `xml:"icon,attr"`
	KeepCount    string `xml:"keepCount,attr"`
	KeepDice     string `xml:"keepDice,attr"`
	KeepOption   string `xml:"keepOption,attr"`
	LockAdd      string `xml:"lockAdd,attr"`
	LockDice     string `xml:"lockDice,attr"`
	LockPlus     string `xml:"lockPlus,attr"`
	LockSides    string `xml:"lockSides,attr"`
	NDice        string `xml:"nDice,attr"`
	NSides       string `xml:"nSides,attr"`
	Name         string `xml:"name,attr"`
	Plus         string `xml:"plus,attr"`
	Prompt       string `xml:"prompt,attr"`
	PropertyGate string `xml:"propertyGate,attr"`
	ReportFormat string `xml:"reportFormat,attr"`
	ReportTotal  string `xml:"reportTotal,attr"`
	SortDice     string `xml:"sortDice,attr"`
	Text         string `xml:"text,attr"`
	Tooltip      string `xml:"tooltip,attr"`
}

type PieceWindow struct {
	XMLName      xml.Name `xml:"VASSAL.build.module.PieceWindow"`
	DefaultWidth string   `xml:"defaultWidth,attr"`
	Hidden       string   `xml:"hidden,attr"`
	Hotkey       string   `xml:"hotkey,attr"`
	Icon         string   `xml:"icon,attr"`
	Scale        string   `xml:"scale,attr"`
	Text         string   `xml:"text,attr"`
	ToolTip      string   `xml:"tooltip,attr"`
	TabWidget    tabWidget
}

type tabWidget struct {
	XMLName    xml.Name     `xml:"VASSAL.build.widget.TabWidget"`
	EntryName  string       `xml:"entryName,attr"`
	ListWidget []listWidget `xml:"VASSAL.build.widget.ListWidget"`
}

type listWidget struct {
	XMLName   xml.Name `xml:"VASSAL.build.widget.ListWidget"`
	Divider   string   `xml:"divider,attr"`
	EntryName string   `xml:"entryName,attr"`
	Height    string   `xml:"height,attr"`
	Scale     string   `xml:"scale,attr"`
	Width     string   `xml:"width,attr"`
	PieceSlot []pieceSlot
}

type pieceSlot struct {
	XMLName   xml.Name `xml:"VASSAL.build.widget.PieceSlot"`
	EntryName string   `xml:"entryName,attr"`
	Gpid      string   `xml:"gpid,attr"`
	Height    int      `xml:"height,attr"`
	Width     int      `xml:"width,attr"`
	Data      string   `xml:",chardata"`
}

type templateData struct {
	Filename  string
	PieceName string
	Id        string
}

// GetVassalDataForCounters returns the Vassal module data for the counters
func GetVassalDataForCounters(t *counters.CounterTemplate, xmlFilepath string) ([]byte, error) {
	var g gameModule
	err := fsops.ReadMarkupFile(xmlFilepath, &g)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to decode content")
	}

	tw := tabWidget{
		EntryName:  "Forces",
		ListWidget: make([]listWidget, 0),
	}

	forces := make(map[string]listWidget)

	forces["Markers"] = listWidget{
		EntryName: "Markers",
		PieceSlot: make([]pieceSlot, 0),
		Scale:     "1.0",
		Height:    "215",
		Width:     "562",
		Divider:   "194",
	}

	xmlTemplateString := `+/null/prototype;BasicPrototype	piece;;;{{ .Filename }};{{ .PieceName}}/	null;0;0;{{ .Id }};0`
	xmlTemplate, err := template.New("xml").Parse(xmlTemplateString)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse template string")
	}

	// Read special files from images folder to statically load them into the module as pieces (terrain, -1, -2 markers, Disorganized, Spent, OOS, etc,)
	files, err := readFiles(counters.BASE_FOLDER + "/images")
	if err != nil {
		return nil, errors.Wrap(err, "could not read files")
	}

	gpid := 200
	id := 200

	for _, file := range files {
		buf := bytes.NewBufferString("")
		if err = xmlTemplate.ExecuteTemplate(
			buf, "xml", templateData{
				Filename:  file, //+1 because file number starts in 1 instead of 0 when they are generated
				PieceName: file,
				Id:        fmt.Sprintf("%d", id),
			},
		); err != nil {
			return nil, errors.Wrap(err, "error trying to write Vassal xml file using templates")
		}
		id++

		piece := pieceSlot{
			EntryName: file,
			Gpid:      fmt.Sprintf("%d", gpid),
			Height:    t.Height,
			Width:     t.Width,
			Data:      buf.String(),
		}

		temp := forces["Markers"]
		temp.PieceSlot = append(forces["Markers"].PieceSlot, piece)
		forces["Markers"] = temp

		gpid++
	}

	filenamesInUse := make(map[string]bool)

	for i, counter := range t.Counters {
		buf := bytes.NewBufferString("")
		if err = xmlTemplate.ExecuteTemplate(
			buf, "xml", templateData{
				Filename: counter.GetCounterFilename(
					t.PositionNumberForFilename,
					"",
					i+1,
					filenamesInUse,
				), //+1 because file number starts in 1 instead of 0 when they are generated
				PieceName: counter.GetCounterFilename(t.PositionNumberForFilename, "", -1, filenamesInUse),
				Id:        fmt.Sprintf("%d", id),
			},
		); err != nil {
			return nil, errors.Wrap(err, "error trying to write Vassal xml file using templates")
		}
		id++

		piece := pieceSlot{
			EntryName: counter.GetTextInPosition(t.PositionNumberForFilename),
			Gpid:      fmt.Sprintf("%d", gpid),
			Height:    t.Height,
			Width:     t.Width,
			Data:      buf.String(),
		}

		if _, ok := forces[counter.Extra.Side]; !ok {
			forces[counter.Extra.Side] = listWidget{
				EntryName: counter.Extra.Side,
				PieceSlot: make([]pieceSlot, 0),
				Scale:     "1.0",
				Height:    "215",
				Width:     "562",
				Divider:   "194",
			}
		}

		temp := forces[counter.Extra.Side]
		temp.PieceSlot = append(forces[counter.Extra.Side].PieceSlot, piece)
		forces[counter.Extra.Side] = temp

		gpid++
	}

	tw.ListWidget = append(tw.ListWidget, mapToArray[listWidget](forces)...)
	g.PieceWindow.TabWidget = tw

	byt, err := xml.MarshalIndent(g, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal the final game module data")
	}

	return byt, nil
}

func readFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			if file.Name()[0] == '_' {
				filenames = append(filenames, file.Name())
			}
		}
	}

	return filenames, nil
}

func mapToArray[T any](m map[string]T) []T {
	temp := make([]T, len(m))

	i := 0
	for _, item := range m {
		temp[i] = item

		i++
	}

	return temp
}
