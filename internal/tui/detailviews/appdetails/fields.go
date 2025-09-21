package appdetails

import (
	"fmt"
	"slices"
	"strings"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"
	"swimpeek/pkg/laneclient"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type applicationFieldList struct {
	analyzer    *analyzer.Analyzer
	innerFrame  *app.Frame
	outerFrame  *app.Frame
	fieldTable  *table.Model
	app         *graph.Node
	appResource *laneclient.Application
}

// newApplicationFieldList creates a new list view for displaying the fields associated with an application.
func newApplicationFieldList(analyzer *analyzer.Analyzer, innerFrame *app.Frame, outerFrame *app.Frame, appResource *laneclient.Application, app *graph.Node) tea.Model {
	return &applicationFieldList{
		analyzer:    analyzer,
		innerFrame:  innerFrame,
		outerFrame:  outerFrame,
		fieldTable:  createFieldTable(appResource, innerFrame),
		app:         app,
		appResource: appResource,
	}
}

func (m *applicationFieldList) openFieldAccess() tea.Msg {
	if len(m.appResource.Fields) == 0 {
		return nil
	}

	selectedRow := m.fieldTable.SelectedRow()
	fieldID := selectedRow[0]

	var selectedField *laneclient.ApplicationField
	for _, f := range m.appResource.Fields {
		if f.Id == fieldID {
			selectedField = &f
			break
		}
	}
	if selectedField == nil {
		return nil
	}

	accessEvents := m.analyzer.ApplicationFieldModifiedBy(m.app, selectedField)
	return app.CmdPushView(newAccessListView(m.analyzer, m.outerFrame, accessEvents, selectedField, m.app))
}

func (m *applicationFieldList) Init() tea.Cmd {
	return nil
}

func (m *applicationFieldList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavUp:
			m.fieldTable.MoveUp(1)
		case app.NavDown:
			m.fieldTable.MoveDown(1)
		case app.NavPageUp:
			m.fieldTable.MoveUp(min(len(m.fieldTable.Rows()), m.fieldTable.Height()) / 3)
		case app.NavPageDown:
			m.fieldTable.MoveDown(min(len(m.fieldTable.Rows()), m.fieldTable.Height()) / 3)
		case app.NavHome:
			m.fieldTable.GotoTop()
		case app.NavEnd:
			m.fieldTable.GotoBottom()
		case app.NavSelect:
			return m, m.openFieldAccess
		}
	}

	return m, nil
}

func (m *applicationFieldList) View() string {
	title := styles.TitleStyle.Render(m.app.Meta.Label+fmt.Sprintf(" - %d Fields", len(m.appResource.Fields))) + "\n"

	m.fieldTable.SetHeight(m.innerFrame.Height - lipgloss.Height(title))
	m.fieldTable.SetWidth(m.innerFrame.Width)

	content := m.renderFieldTable()

	return lipgloss.JoinVertical(lipgloss.Left, title, content)

}

func (m *applicationFieldList) renderFieldTable() string {
	if len(m.fieldTable.Rows()) == 0 {
		return styles.ResDescriptionStyle.Render("No fields defined for this application.")
	}

	return m.fieldTable.View()
}

// createFieldTable creates a table model for displaying the application fields.
func createFieldTable(app *laneclient.Application, frame *app.Frame) *table.Model {
	rows := make([]table.Row, len(app.Fields))

	for idx, field := range app.Fields {

		// Format field flags
		flags := make([]string, 0, 2)
		if field.Required {
			flags = append(flags, "REQ")
		}
		if field.ReadOnly {
			flags = append(flags, "RO")
		}

		// Append input type if available
		fieldType := strings.TrimSuffix(strings.TrimPrefix(field.Type, "Core.Models.Fields."), ", Core")
		if field.InputType != "" {
			fieldType = fmt.Sprintf("%s (%s)", fieldType, field.InputType)
		}

		row := table.Row{
			field.Id,
			field.Key,
			field.Name,
			fieldType,
			strings.Join(flags, ", "),
		}
		rows[idx] = row
	}

	// Sort rows by key
	slices.SortFunc(rows, func(a, b table.Row) int { return strings.Compare(a[1], b[1]) })

	columns := []table.Column{
		{Title: "ID"},
		{Title: "Key"},
		{Title: "Name"},
		{Title: "Type"},
		{Title: "Flags"},
	}

	for idx, w := range getColumnWidths(rows, len(columns)) {
		columns[idx].Width = max(w, len(columns[idx].Title)) + 2
	}

	t := table.New(
		table.WithFocused(true),
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(frame.Height-2),
		table.WithStyles(styles.TableStyle),
	)

	return &t
}

// getColumnWidths calculates the maximum width for each column based on the content of the rows.
func getColumnWidths(rows []table.Row, colCount int) []int {
	colWidths := make([]int, colCount)

	for _, row := range rows {
		for colIdx, cell := range row {
			colWidths[colIdx] = max(colWidths[colIdx], len(cell))
		}
	}

	return colWidths
}
