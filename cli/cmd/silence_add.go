package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/user"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/model"

	"github.com/prometheus/alertmanager/cli"
	"github.com/prometheus/alertmanager/types"
)

type addResponse struct {
	Status string `json:"status"`
	Data   struct {
		SilenceID string `json:"silenceId"`
	} `json:"data,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

func username() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

var (
	addCmd         = silenceCmd.Command("add", "Add a new alertmanager silence")
	author         = addCmd.Flag("author", "Username for CreatedBy field").Short('a').Default(username()).String()
	requireComment = addCmd.Flag("require-comment", "Require comment to be set").Hidden().Default("true").Bool()
	duration       = addCmd.Flag("duration", "Duration of silence").Short('d').Default("1h").String()
	addStart       = addCmd.Flag("start", "Set when the silence should start. RFC3339 format 2006-01-02T15:04:05Z07:00").String()
	addEnd         = addCmd.Flag("end", "Set when the silence should end (overwrites duration). RFC3339 format 2006-01-02T15:04:05Z07:00").String()
	comment        = addCmd.Flag("comment", "A comment to help describe the silence").Short('c').String()
	addArgs        = addCmd.Arg("matcher-groups", "Query filter").Strings()
)

func init() {
	addCmd.Action(add)
	longHelpText["silence add"] = `Add a new alertmanager silence

  Amtool uses a simplified prometheus syntax to represent silences. The
  non-option section of arguments constructs a list of "Matcher Groups"
  that will be used to create a number of silences. The following examples
  will attempt to show this behaviour in action:

  amtool silence add alertname=foo node=bar

	This statement will add a silence that matches alerts with the
	alertname=foo and node=bar label value pairs set.

  amtool silence add foo node=bar

	If alertname is ommited and the first argument does not contain a '=' or a
	'=~' then it will be assumed to be the value of the alertname pair.

  amtool silence add 'alertname=~foo.*'

	As well as direct equality, regex matching is also supported. The '=~' syntax
	(similar to prometheus) is used to represent a regex match. Regex matching
	can be used in combination with a direct match.
`
}

func add(element *kingpin.ParseElement, ctx *kingpin.ParseContext) error {
	var err error

	matchers, err := parseMatchers(*addArgs)
	if err != nil {
		return err
	}

	if len(matchers) < 1 {
		return fmt.Errorf("no matchers specified")
	}

	var endsAt time.Time
	if *addEnd != "" {
		endsAt, err = time.Parse(time.RFC3339, *addEnd)
		if err != nil {
			return err
		}
	} else {
		d, err := model.ParseDuration(*duration)
		if err != nil {
			return err
		}
		if d == 0 {
			return fmt.Errorf("silence duration must be greater than 0")
		}
		endsAt = time.Now().UTC().Add(time.Duration(d))
	}

	if *requireComment && *comment == "" {
		return errors.New("comment required by config")
	}

	var startsAt time.Time
	if *addStart != "" {
		startsAt, err = time.Parse(time.RFC3339, *addStart)
		if err != nil {
			return err
		}

	} else {
		startsAt = time.Now().UTC()
	}

	if startsAt.After(endsAt) {
		return errors.New("silence cannot start after it ends")
	}

	typeMatchers, err := TypeMatchers(matchers)
	if err != nil {
		return err
	}

	silence := types.Silence{
		Matchers:  typeMatchers,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		CreatedBy: *author,
		Comment:   *comment,
	}

	client, err := api.NewClient(api.Config{Address: (*alertmanagerUrl).String()})
	if err != nil {
		return err
	}
	silenceAPI := cli.NewSilenceAPI(client)
	silenceID, err := silenceAPI.Set(context.Background(), silence)
	if err != nil {
		return err
	}

	_, err = fmt.Println(silenceID)
	return err
}
