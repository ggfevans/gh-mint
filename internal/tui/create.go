package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ggfevans/gh-mint/internal/config"
)

type CreateModel struct {
	form        *huh.Form
	styles      *Styles
	cfg         *config.Config
}

func NewCreateModel(styles *Styles, cfg *config.Config) CreateModel {
	m := CreateModel{styles: styles, cfg: cfg}

	profileOpts := make([]huh.Option[string], 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		profileOpts = append(profileOpts, huh.NewOption(name, name))
	}

	defaultProfile := cfg.DefaultProfile
	defaultVisibility := "public"

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Repository name").
				Validate(func(s string) error {
					return config.ValidateRepoName(s)
				}),

			huh.NewInput().
				Key("description").
				Title("Description"),

			huh.NewSelect[string]().
				Key("visibility").
				Title("Visibility").
				Options(
					huh.NewOption("Public", "public"),
					huh.NewOption("Private", "private"),
				).
				Value(&defaultVisibility),

			huh.NewSelect[string]().
				Key("profile").
				Title("Profile").
				Options(profileOpts...).
				Value(&defaultProfile),

			huh.NewConfirm().
				Key("confirm").
				Title("Create repository?").
				Affirmative("Create").
				Negative("Back"),
		),
	).WithShowHelp(true)

	return m
}

func (m CreateModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m CreateModel) Update(msg tea.Msg) (CreateModel, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m CreateModel) View() string {
	return m.form.View()
}

func (m CreateModel) Done() bool {
	return m.form.State == huh.StateCompleted
}

func (m CreateModel) Confirmed() bool {
	return m.form.GetBool("confirm")
}

func (m CreateModel) IsPublic() bool {
	return m.form.GetString("visibility") == "public"
}

func (m CreateModel) Profile() config.Profile {
	return m.cfg.Profiles[m.form.GetString("profile")]
}
