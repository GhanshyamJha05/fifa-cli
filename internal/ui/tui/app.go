package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/styles"
)

type tab int

const (
	tabHome tab = iota
	tabTeams
	tabMatches
	tabStandings
	tabStats
)

type matchFilter string

const (
	filterAll      matchFilter = "all"
	filterToday    matchFilter = "today"
	filterLive     matchFilter = "live"
	filterUpcoming matchFilter = "upcoming"
	filterResults  matchFilter = "results"
)

type screen int

const (
	screenMain screen = iota
	screenDetail
	screenHelp
)

type errMsg error

type dataMsg struct {
	tab  tab
	data any
}

type detailMsg struct {
	content string
}

type teamItem struct {
	team domain.Team
}

func (i teamItem) Title() string       { return i.team.Name }
func (i teamItem) Description() string { return fmt.Sprintf("Group %s · %s", i.team.Group, i.team.Code) }
func (i teamItem) FilterValue() string { return i.team.Name }

type searchItem struct {
	result domain.SearchResult
}

func (i searchItem) Title() string { return i.result.Title }
func (i searchItem) Description() string {
	return fmt.Sprintf("%s — %s", i.result.Type, i.result.Subtitle)
}
func (i searchItem) FilterValue() string { return i.result.Title }

type keys struct {
	Up, Down, Left, Right key.Binding
	Enter, Back, Quit     key.Binding
	Help, Search, Filter  key.Binding
	TabNext, TabPrev      key.Binding
}

func defaultKeys() keys {
	return keys{
		Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "prev")),
		Right: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Filter: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "filter")),
		TabNext: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
		TabPrev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
	}
}

// Model is the interactive dashboard.
type Model struct {
	svc   *service.Service
	theme styles.Theme
	keys  keys

	width, height int
	activeTab     tab
	screen        screen
	matchFilter   matchFilter
	loading       bool
	err           error
	ready         bool

	spinner  spinner.Model
	viewport viewport.Model
	teamList list.Model
	searchList list.Model

	showSearch bool
	searchInput textinput.Model
	searchQuery string

	// cached data
	info          *domain.TournamentInfo
	todayMatches  []domain.Match
	allMatches    []domain.Match
	teams         []domain.Team
	standings     []domain.GroupStanding
	stats         *domain.TournamentStats
	searchResults []domain.SearchResult
	loaded        map[tab]bool

	detailContent string
}

// New creates the TUI model.
func New(svc *service.Service, themeName string) Model {
	theme := styles.Get(themeName)

	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = lipgloss.NewStyle().Foreground(theme.Primary)

	ti := textinput.New()
	ti.Placeholder = "Team, player, or match..."
	ti.CharLimit = 64
	ti.Prompt = "🔍 "

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(theme.Background).Background(theme.Primary).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(theme.Background).Background(theme.Primary)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(theme.Text)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(theme.Muted)

	teamList := list.New([]list.Item{}, delegate, 0, 0)
	teamList.Title = "Teams"
	teamList.SetShowStatusBar(false)
	teamList.SetFilteringEnabled(false)
	teamList.SetShowHelp(false)

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().Foreground(theme.Text)

	searchList := list.New([]list.Item{}, delegate, 0, 0)
	searchList.Title = "Search"
	searchList.SetShowStatusBar(false)
	searchList.SetFilteringEnabled(false)
	searchList.SetShowHelp(false)

	return Model{
		svc: svc, theme: theme, keys: defaultKeys(),
		spinner: sp, searchInput: ti,
		activeTab: tabHome, screen: screenMain, matchFilter: filterAll,
		loaded: make(map[tab]bool),
		teamList: teamList, searchList: searchList, viewport: vp,
	}
}

// Run launches the dashboard.
func Run(svc *service.Service, theme string) error {
	m := New(svc, theme)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m Model) Init() tea.Cmd {
	m.svc.RefreshCache(context.Background())
	return tea.Batch(m.spinner.Tick, m.loadTab(tabHome))
}

func (m Model) loadTab(t tab) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		switch t {
		case tabHome:
			dashboard, err := m.svc.LoadDashboard(ctx)
			if err != nil {
				return errMsg(err)
			}
			return dataMsg{tab: t, data: dashboard}
		case tabTeams:
			teams, err := m.svc.GetTeams(ctx)
			if err != nil {
				return errMsg(err)
			}
			return dataMsg{tab: t, data: teams}
		case tabMatches:
			matches, err := m.svc.GetMatches(ctx)
			if err != nil {
				return errMsg(err)
			}
			return dataMsg{tab: t, data: matches}
		case tabStandings:
			standings, err := m.svc.GetStandings(ctx)
			if err != nil {
				return errMsg(err)
			}
			return dataMsg{tab: t, data: standings}
		case tabStats:
			stats, err := m.svc.GetStats(ctx)
			if err != nil {
				return errMsg(err)
			}
			return dataMsg{tab: t, data: stats}
		}
		return nil
	}
}

func (m Model) loadTeamDetail(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		team, err := m.svc.GetTeam(ctx, name)
		if err != nil {
			return errMsg(err)
		}
		squad, _ := m.svc.GetSquad(ctx, name)
		form, _ := m.svc.GetTeamForm(ctx, name)
		return detailMsg{content: formatTeamDetail(team, squad, form, m.theme)}
	}
}

func (m Model) loadSearchDetail(r domain.SearchResult) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		switch r.Type {
		case "team":
			team, err := m.svc.GetTeam(ctx, r.Title)
			if err != nil {
				return errMsg(err)
			}
			squad, _ := m.svc.GetSquad(ctx, r.Title)
			form, _ := m.svc.GetTeamForm(ctx, r.Title)
			return detailMsg{content: formatTeamDetail(team, squad, form, m.theme)}
		case "player":
			p, err := m.svc.GetPlayer(ctx, r.Title)
			if err != nil {
				return errMsg(err)
			}
			return detailMsg{content: formatPlayer(p, m.theme)}
		default:
			return detailMsg{content: m.theme.Dim.Render("  Match details — use Matches tab.")}
		}
	}
}

func (m Model) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		results, err := m.svc.Search(ctx, query)
		if err != nil {
			return errMsg(err)
		}
		return dataMsg{tab: tab(-1), data: results}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.layout()
		return m, nil

	case tea.KeyMsg:
		if m.showSearch {
			return m.handleSearchKey(msg)
		}
		return m.handleKey(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case errMsg:
		m.loading = false
		m.err = error(msg)
		return m, nil

	case detailMsg:
		m.loading = false
		m.screen = screenDetail
		m.detailContent = msg.content
		header := m.theme.Dim.Render("  Press Esc to go back") + "\n\n"
		m.viewport.SetContent(header + msg.content)
		m.viewport.GotoTop()
		m.layout()
		return m, nil

	case dataMsg:
		m.loading = false
		m.err = nil

		if msg.tab == tab(-1) {
			m.searchResults = msg.data.([]domain.SearchResult)
			items := make([]list.Item, len(m.searchResults))
			for i, r := range m.searchResults {
				items[i] = searchItem{result: r}
			}
			m.searchList.SetItems(items)
			m.showSearch = true
			m.screen = screenMain
			m.layout()
			m.refreshContent()
			return m, nil
		}

		m.loaded[msg.tab] = true
		switch msg.tab {
		case tabHome:
			d := msg.data.(*service.DashboardData)
			m.info = d.Info
			m.todayMatches = d.TodayMatches
			m.teams = d.Teams
			m.standings = d.Standings
			m.stats = d.Stats
		case tabTeams:
			m.teams = msg.data.([]domain.Team)
			items := make([]list.Item, len(m.teams))
			for i, t := range m.teams {
				items[i] = teamItem{team: t}
			}
			m.teamList.SetItems(items)
		case tabMatches:
			m.allMatches = msg.data.([]domain.Match)
		case tabStandings:
			m.standings = msg.data.([]domain.GroupStanding)
		case tabStats:
			m.stats = msg.data.(*domain.TournamentStats)
		}
		m.refreshContent()
		m.layout()
		return m, nil
	}

	// Delegate to sub-components
	if m.activeTab == tabTeams && m.screen == screenMain && !m.showSearch {
		var cmd tea.Cmd
		m.teamList, cmd = m.teamList.Update(msg)
		m.refreshContent()
		cmds = append(cmds, cmd)
	}
	if m.showSearch && len(m.searchResults) > 0 {
		var cmd tea.Cmd
		m.searchList, cmd = m.searchList.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.screen == screenDetail || m.activeTab != tabTeams || m.showSearch {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab != tabTeams {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.showSearch = false
		m.searchInput.SetValue("")
		m.searchResults = nil
		m.screen = screenMain
		m.refreshContent()
		return m, nil
	case key.Matches(msg, m.keys.Enter):
		query := strings.TrimSpace(m.searchInput.Value())
		if query == "" {
			return m, nil
		}
		m.searchQuery = query
		m.loading = true
		m.showSearch = false
		return m, m.doSearch(query)
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Back):
		if m.screen == screenDetail {
			m.screen = screenMain
			m.detailContent = ""
			m.refreshContent()
			m.layout()
			return m, nil
		}
		if m.showSearch {
			m.showSearch = false
			m.searchResults = nil
			m.refreshContent()
			return m, nil
		}
		if m.screen == screenHelp {
			m.screen = screenMain
			return m, nil
		}
		return m, nil

	case key.Matches(msg, m.keys.Help):
		if m.screen == screenHelp {
			m.screen = screenMain
		} else {
			m.screen = screenHelp
		}
		return m, nil

	case key.Matches(msg, m.keys.Search):
		m.showSearch = true
		m.searchResults = nil
		m.searchQuery = ""
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		return m, textinput.Blink

	case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down):
		if m.screen == screenDetail || m.activeTab != tabTeams || m.showSearch {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
		return m, nil

	case key.Matches(msg, m.keys.TabNext):
		return m.switchTab(1)
	case key.Matches(msg, m.keys.TabPrev):
		return m.switchTab(-1)
	case key.Matches(msg, m.keys.Left):
		return m.switchTab(-1)
	case key.Matches(msg, m.keys.Right):
		return m.switchTab(1)

	case key.Matches(msg, m.keys.Filter):
		if m.activeTab == tabMatches && m.screen == screenMain {
			m.matchFilter = m.matchFilter.next()
			m.refreshContent()
		}
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		if m.showSearch && len(m.searchResults) > 0 {
			if item, ok := m.searchList.SelectedItem().(searchItem); ok {
				m.loading = true
				return m, m.loadSearchDetail(item.result)
			}
		}
		if m.activeTab == tabTeams && m.screen == screenMain {
			if item, ok := m.teamList.SelectedItem().(teamItem); ok {
				m.loading = true
				return m, m.loadTeamDetail(item.team.Name)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) switchTab(dir int) (Model, tea.Cmd) {
	if m.screen == screenDetail {
		m.screen = screenMain
		m.detailContent = ""
	}
	tabs := []tab{tabHome, tabTeams, tabMatches, tabStandings, tabStats}
	idx := 0
	for i, t := range tabs {
		if t == m.activeTab {
			idx = i
			break
		}
	}
	idx = (idx + dir + len(tabs)) % len(tabs)
	m.activeTab = tabs[idx]
	m.screen = screenMain
	m.showSearch = false

	if !m.loaded[m.activeTab] {
		m.loading = true
		return *m, m.loadTab(m.activeTab)
	}
	m.refreshContent()
	return *m, nil
}

func (f matchFilter) next() matchFilter {
	order := []matchFilter{filterAll, filterToday, filterLive, filterUpcoming, filterResults}
	for i, v := range order {
		if v == f {
			return order[(i+1)%len(order)]
		}
	}
	return filterAll
}

func (f matchFilter) label() string {
	switch f {
	case filterToday:
		return "Today"
	case filterLive:
		return "Live"
	case filterUpcoming:
		return "Upcoming"
	case filterResults:
		return "Results"
	default:
		return "All"
	}
}

func (m *Model) layout() {
	headerH := 5
	footerH := 2
	contentH := m.height - headerH - footerH
	if contentH < 8 {
		contentH = 8
	}

	if m.activeTab == tabTeams && m.screen == screenMain && !m.showSearch {
		listW := m.width / 3
		if listW < 28 {
			listW = 28
		}
		m.teamList.SetSize(listW-4, contentH-4)
		m.viewport.Width = m.width - listW - 6
		m.viewport.Height = contentH - 2
	} else {
		m.viewport.Width = m.width - 4
		m.viewport.Height = contentH - 2
	}

	if m.showSearch && len(m.searchResults) > 0 {
		m.searchList.SetSize(m.width/3, contentH-4)
	}
}

func (m *Model) refreshContent() {
	if m.screen == screenDetail {
		return
	}
	var content string
	switch m.activeTab {
	case tabHome:
		content = formatHome(m.info, m.todayMatches, m.theme)
	case tabTeams:
		if len(m.teams) > 0 {
			if item, ok := m.teamList.SelectedItem().(teamItem); ok {
				content = m.theme.Dim.Render(fmt.Sprintf("\n  Select a team and press Enter.\n\n  Highlighted: %s (Group %s)\n  Press ↑↓ to browse, Enter for full squad.",
					item.team.Name, item.team.Group))
			}
		} else {
			content = m.theme.Dim.Render("  Loading teams...")
		}
	case tabMatches:
		content = formatMatches(m.filteredMatches(), string(m.matchFilter), m.theme)
	case tabStandings:
		content = formatStandings(m.standings, m.theme)
	case tabStats:
		if m.stats != nil {
			content = formatStats(m.stats, m.theme)
		} else {
			content = m.theme.Dim.Render("  No stats available.")
		}
	}
	if m.showSearch && len(m.searchResults) > 0 {
		content = formatSearchResults(m.searchResults, m.searchQuery, m.theme)
	}
	m.viewport.SetContent(content)
	m.viewport.GotoTop()
}

func (m Model) filteredMatches() []domain.Match {
	now := time.Now()
	var out []domain.Match
	for _, match := range m.allMatches {
		switch m.matchFilter {
		case filterToday:
			if sameDay(match.Date, now) {
				out = append(out, match)
			}
		case filterLive:
			if match.Status == domain.StatusLive {
				out = append(out, match)
			}
		case filterUpcoming:
			if match.Status == domain.StatusScheduled && match.Date.After(now) {
				out = append(out, match)
			}
		case filterResults:
			if match.Status == domain.StatusFinished {
				out = append(out, match)
			}
		default:
			out = append(out, match)
		}
	}
	return out
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func (m Model) View() string {
	if !m.ready {
		return m.theme.Dim.Render("  Initializing...")
	}

	var sections []string
	sections = append(sections, m.renderHeader())
	sections = append(sections, m.renderBody())
	sections = append(sections, m.renderFooter())
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	title := m.theme.AppTitle.Width(m.width).Render("⚽ FIFA World Cup 2026  ·  USA · CAN · MEX")

	tabLabels := []struct {
		t     tab
		label string
	}{
		{tabHome, "Home"}, {tabTeams, "Teams"}, {tabMatches, "Matches"},
		{tabStandings, "Standings"}, {tabStats, "Stats"},
	}
	var tabs []string
	for _, tl := range tabLabels {
		label := tl.label
		if tl.t == m.activeTab && m.screen != screenHelp {
			tabs = append(tabs, m.theme.TabActive.Render(" "+label+" "))
		} else {
			tabs = append(tabs, m.theme.Tab.Render(" "+label+" "))
		}
	}
	tabBar := m.theme.TabBar.Width(m.width).Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	return title + "\n" + tabBar
}

func (m Model) renderBody() string {
	if m.screen == screenHelp {
		return m.theme.Content.Width(m.width).Height(m.viewport.Height).Render(helpText())
	}

	if m.showSearch && len(m.searchResults) == 0 && m.searchQuery == "" {
		box := m.theme.Content.Width(m.width - 4).Render(
			m.theme.Header.Render("Search") + "\n\n" + m.searchInput.View() + "\n\n" +
				m.theme.Dim.Render("  Type a name and press Enter · Esc to cancel"),
		)
		return box
	}

	if m.loading {
		return m.theme.Content.Width(m.width).Render(
			fmt.Sprintf("\n\n  %s  Loading...\n", m.spinner.View()),
		)
	}

	if m.err != nil {
		errLine := m.theme.Loss.Render(fmt.Sprintf("  ⚠ %v", m.err)) + "\n" +
			m.theme.Dim.Render("  Press Esc to go back · check API key in config.yaml")
		if m.screen == screenDetail || m.activeTab != tabTeams {
			return m.theme.Content.Width(m.width).Render("\n" + errLine)
		}
	}

	// Teams: split list + preview
	if m.activeTab == tabTeams && m.screen == screenMain && !m.showSearch {
		listPanel := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.theme.Border).
			Width(m.teamList.Width() + 4).
			Height(m.teamList.Height() + 2).
			Render(m.teamList.View())

		preview := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.theme.Border).
			Width(m.viewport.Width + 2).
			Height(m.viewport.Height + 2).
			Render(m.viewport.View())

		return lipgloss.JoinHorizontal(lipgloss.Top, "  "+listPanel, "  "+preview)
	}

	// Search results: split list + hint
	if m.showSearch && len(m.searchResults) > 0 && m.screen == screenMain {
		listPanel := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.theme.Border).
			Width(m.searchList.Width() + 4).
			Height(m.searchList.Height() + 2).
			Render(m.searchList.View())

		hint := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.theme.Border).
			Width(m.viewport.Width + 2).
			Height(m.viewport.Height + 2).
			Render(m.viewport.View())

		return lipgloss.JoinHorizontal(lipgloss.Top, "  "+listPanel, "  "+hint)
	}

	// Detail or full-width content
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Width(m.viewport.Width + 2).
		Height(m.viewport.Height + 2)

	if m.screen == screenDetail {
		return "  " + border.Render(m.viewport.View())
	}

	content := m.viewport.View()

	// Match filter chips
	extra := ""
	if m.activeTab == tabMatches && m.screen == screenMain {
		chips := []matchFilter{filterAll, filterToday, filterLive, filterUpcoming, filterResults}
		var parts []string
		for _, c := range chips {
			if c == m.matchFilter {
				parts = append(parts, m.theme.ChipActive.Render(c.label()))
			} else {
				parts = append(parts, m.theme.Chip.Render(c.label()))
			}
		}
		extra = "\n  " + m.theme.Dim.Render("Filter (f): ") + strings.Join(parts, " ") + "\n\n"
		inner := extra + content
		return "  " + border.Render(inner)
	}

	return "  " + border.Render(content)
}

func (m Model) renderFooter() string {
	if m.screen == screenHelp {
		return m.theme.Footer.Width(m.width).Render("  ? close help  ·  q quit")
	}
	if m.showSearch && m.searchInput.Focused() {
		return m.theme.Footer.Width(m.width).Render("  enter search  ·  esc cancel")
	}
	if m.screen == screenDetail {
		return m.theme.Footer.Width(m.width).Render("  esc back  ·  ←→ switch tab  ·  / search  ·  q quit")
	}
	if m.activeTab == tabMatches {
		return m.theme.Footer.Width(m.width).Render("  f filter matches  ·  ←→ tabs  ·  / search  ·  ? help  ·  q quit")
	}
	if m.activeTab == tabTeams {
		return m.theme.Footer.Width(m.width).Render("  ↑↓ browse  ·  enter view squad  ·  ←→ tabs  ·  / search  ·  q quit")
	}
	return m.theme.Footer.Width(m.width).Render("  ←→ switch tab  ·  / search  ·  ? help  ·  q quit")
}

func helpText() string {
	return `
  FIFA CLI — Quick Guide

  NAVIGATION
    ← / → or Tab     Switch between Home, Teams, Matches, Standings, Stats
    ↑ / ↓            Browse lists (Teams, Search results)
    Enter            Open team squad or search result
    Esc              Go back from detail view
    /                Search teams and players
    f                Cycle match filters (Matches tab)
    ?                Toggle this help
    q                Quit

  TIPS
    • Teams tab shows a list on the left — pick one and press Enter
    • Matches tab lets you filter: All, Today, Live, Upcoming, Results
    • Data is cached locally for faster browsing
`
}
