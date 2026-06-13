package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/render"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/styles"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/tui"
)

var (
	svc    *service.Service
	theme  styles.Theme
	logger *slog.Logger
)

// Execute runs the root command.
func Execute() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	level := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	svc, err = service.New(cfg, logger)
	if err != nil {
		return err
	}
	defer svc.Close()

	theme = styles.Get(cfg.Theme)

	root := &cobra.Command{
		Use:   "fifa",
		Short: "FIFA World Cup 2026 terminal dashboard",
		Long:  "A premium terminal application for the 2026 FIFA World Cup.\nRun without arguments for interactive mode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run(svc, cfg.Theme)
		},
	}

	root.AddCommand(
		teamsCmd(), teamCmd(), squadCmd(), playerCmd(),
		matchesCmd(), nextCmd(), resultsCmd(),
		standingsCmd(), bracketCmd(), statsCmd(), searchCmd(),
		h2hCmd(), exportCmd(), serveCmd(),
	)

	return root.Execute()
}

func ctx() context.Context {
	return context.Background()
}

func printOut(s string) {
	fmt.Println(s)
}

func teamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "teams",
		Short: "List all World Cup teams",
		RunE: func(cmd *cobra.Command, args []string) error {
			teams, err := svc.GetTeams(ctx())
			if err != nil {
				return err
			}
			printOut(render.Teams(teams, theme))
			return nil
		},
	}
}

func teamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "team [name]",
		Short: "Show team details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			team, err := svc.GetTeam(ctx(), args[0])
			if err != nil {
				return err
			}
			form, _ := svc.GetTeamForm(ctx(), args[0])
			printOut(render.TeamDetail(team, form, theme))
			return nil
		},
	}
}

func squadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "squad [team]",
		Short: "Show team squad",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			squad, err := svc.GetSquad(ctx(), args[0])
			if err != nil {
				return err
			}
			printOut(render.Squad(squad, args[0], theme))
			return nil
		},
	}
}

func playerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "player [name]",
		Short: "Show player profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			player, err := svc.GetPlayer(ctx(), args[0])
			if err != nil {
				return err
			}
			printOut(render.PlayerProfile(player, theme))
			return nil
		},
	}
}

func matchesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matches",
		Short: "View match fixtures and results",
		RunE: func(cmd *cobra.Command, args []string) error {
			matches, err := svc.GetMatches(ctx())
			if err != nil {
				return err
			}
			printOut(render.Matches(matches, "⚽ All Matches", theme))
			return nil
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "today",
		Short: "Show today's matches",
		RunE: func(cmd *cobra.Command, args []string) error {
			matches, err := svc.GetMatchesToday(ctx())
			if err != nil {
				return err
			}
			printOut(render.Matches(matches, "📅 Today's Matches", theme))
			return nil
		},
	})
	return cmd
}

func nextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Show upcoming fixtures",
		RunE: func(cmd *cobra.Command, args []string) error {
			matches, err := svc.GetUpcoming(ctx(), 10)
			if err != nil {
				return err
			}
			printOut(render.Matches(matches, "🔜 Upcoming Fixtures", theme))
			return nil
		},
	}
}

func resultsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "results",
		Short: "Show completed match results",
		RunE: func(cmd *cobra.Command, args []string) error {
			matches, err := svc.GetResults(ctx())
			if err != nil {
				return err
			}
			printOut(render.Matches(matches, "✅ Results", theme))
			return nil
		},
	}
}

func standingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "standings",
		Short: "Show group standings",
		RunE: func(cmd *cobra.Command, args []string) error {
			standings, err := svc.GetStandings(ctx())
			if err != nil {
				return err
			}
			printOut(render.Standings(standings, theme))
			return nil
		},
	}
}

func bracketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bracket",
		Short: "Show knockout bracket",
		RunE: func(cmd *cobra.Command, args []string) error {
			bracket, err := svc.GetBracket(ctx())
			if err != nil {
				return err
			}
			printOut(render.Bracket(bracket, theme))
			return nil
		},
	}
}

func statsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show tournament statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			stats, err := svc.GetStats(ctx())
			if err != nil {
				return err
			}
			printOut(render.Stats(stats, theme))
			return nil
		},
	}
}

func searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Search teams, players, and matches",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := svc.Search(ctx(), args[0])
			if err != nil {
				return err
			}
			printOut(render.SearchResults(results, args[0], theme))
			return nil
		},
	}
}

func h2hCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "h2h [team-a] [team-b]",
		Short: "Head-to-head comparison between two teams",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h2h, err := svc.GetHeadToHead(ctx(), args[0], args[1])
			if err != nil {
				return err
			}
			printOut(render.HeadToHead(h2h, theme))
			return nil
		},
	}
}

func exportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export data to JSON or CSV",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "fixtures [file.json]",
		Short: "Export all fixtures to JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportJSON(args[0], func() (any, error) {
				return svc.GetMatches(ctx())
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "results [file.json]",
		Short: "Export results to JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportJSON(args[0], func() (any, error) {
				return svc.GetResults(ctx())
			})
		},
	})
	return cmd
}
