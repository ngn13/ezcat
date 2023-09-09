/*
 *
 *  ezcat | easy netcat revshell handler
 *  =====================================
 *  this program is licensed under GNU
 *  General Public License Version 2
 *  (GPLv2), please see LICENSE.txt
 *
 *  written by ngn - https://ngn13.fun
 *
*/

package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/cmds"
	"github.com/ngn13/ezcat/http"
	"github.com/ngn13/ezcat/term"
)

func completer(d prompt.Document) []prompt.Suggest {
  suggestions := []prompt.Suggest{}
  text := d.TextBeforeCursor()

  if text == "" {
    return suggestions
  }

  if text == "shell" || 
    text == "upload" || 
    text == "download" {
    for _, s := range bridge.SHELLS {
      if !s.Connected {
        continue
      }
      suggestions = append(suggestions, prompt.Suggest{Text: text+" "+s.ID})
    }
  }

  return suggestions
}

func main(){
  go http.Listen()
  cmds.LoadCmds()
  term.PrintBanner()
  term.Cyan("\nStarting HTTP beacon at %s", bridge.GetAddr())
  term.Newline()

  p := prompt.New(
    nil,
    completer,
    prompt.OptionPrefix("[ezcat]# "),
    prompt.OptionPrefixTextColor(prompt.Red),
    prompt.OptionSuggestionBGColor(prompt.DarkRed),
    prompt.OptionSuggestionTextColor(prompt.White),
    prompt.OptionDescriptionBGColor(prompt.DarkRed),
    prompt.OptionDescriptionTextColor(prompt.White),
    prompt.OptionSelectedSuggestionBGColor(prompt.Red),
    prompt.OptionSelectedDescriptionBGColor(prompt.Red),
    prompt.OptionScrollbarThumbColor(prompt.Black),
  )

  for {
    inp := p.Input()
    ran := false

    if inp == "exit" {
      break
    }

    args := strings.Split(inp, " ")
    cmd := args[0]
    args = append(args[:0], args[1:]...)

    if cmd == "" {
      continue
    }

    for _, c := range cmds.CMDS {
      if c.Name == cmd {
        c.Func(args)
        ran = true
      }
    }

    if ran {
      continue
    }

    term.Error("Command not found!")
  }
}
