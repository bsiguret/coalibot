package FortyTwo

import (
	"fmt"
	"strings"
	"time"

	"github.com/genesixx/coalibot/Struct"
	"github.com/genesixx/coalibot/Utils"
	"github.com/nlopes/slack"
	"gitlab.com/clafoutis/api42"
)

func Where(option string, event *Struct.Message) bool {
	params := api42.NewParameter()
	if len(strings.Split(option, " ")) == 4 && (strings.IndexAny(option, "le branle couille") != -1 || strings.IndexAny(option, "la branle couille") != -1) {
		y, m, d := time.Now().Date()
		rangeBegin := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
		rangeEnd := rangeBegin.AddDate(0, 0, -7)
		user := strings.Split(option, " ")[3]
		logtime := Utils.IntraLogtime(user, rangeEnd, rangeBegin, event.FortyTwo)
		if logtime.Hours() >= 35 {
			Utils.PostMsg(event, slack.MsgOptionText("*"+user+"* is not a branle couille", false))
			return true
		}
		data, err := event.FortyTwo.GetUserLocations(user, params)
		if err != nil {
			Utils.PostMsg(event, slack.MsgOptionText("login invalide", false))
			return false
		}

		if len(data) == 0 || data[0].EndAt != nil {
			var diff = time.Now().Sub(*data[0].EndAt)

			Utils.PostMsg(event, slack.MsgOptionText("*"+user+"* est hors-ligne depuis *"+Utils.FmtDuration(diff)+"*", false))
		} else {
			Utils.PostMsg(event, slack.MsgOptionText("*"+user+"* est à la place *"+data[0].Host+"*", false))
		}
		return true
	}
	user, error := Utils.GetLogin(option, event)
	if error == true {
		return false
	}

	if user[0] == '!' || user[0] == '?' {
		return false
	}
	if user == "dieu" {
		user = "elebouch"
	} else if user == "manager" {
		user = "vtennero"
	}
	if user == "guardians" || user == "gardiens" {

		var guardians = []string{
			"dcirlig",
			"vtennero",
			"elebouch",
			"fbabin",
			"tbailly-",
			"mmerabet",
			"aledru",
			"dlavaury",
			"jauplat",
			"jraymond",
		}
		var str string
		for i := 0; i < len(guardians); i++ {
			data, err := event.FortyTwo.GetUserLocations(guardians[i], params)
			if err != nil {
				fmt.Println(err)
				str += "login invalide\n"
				return false
			}
			if len(data) == 0 || data[0].EndAt != nil {
				var diff = time.Now().Sub(*data[0].EndAt)

				str += "*" + guardians[i] + "* est hors-ligne depuis " + Utils.PrettyDurationPrinting(diff) + "\n"
			} else {
				str += "*" + guardians[i] + "* est à la place *" + data[0].Host + "*\n"
			}
		}
		Utils.PostMsg(event, slack.MsgOptionText(str, false))

		return true
	}
	data, err := event.FortyTwo.GetUserLocations(user, params)
	if err != nil {
		Utils.PostMsg(event, slack.MsgOptionText("login invalide", false))
		return false
	}

	if len(data) == 0 || data[0].EndAt != nil {
		var diff = time.Now().Sub(*data[0].EndAt)

		Utils.PostMsg(event, slack.MsgOptionText("*"+user+"* est hors-ligne depuis *"+Utils.PrettyDurationPrinting(diff)+"*", false))
	} else {
		Utils.PostMsg(event, slack.MsgOptionText("*"+user+"* est à la place *"+data[0].Host+"*", false))
	}
	return true
}
