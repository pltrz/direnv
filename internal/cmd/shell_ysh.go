package cmd

import (
	"bytes"
	"encoding/json"
)

type ysh struct{}

// Ysh adds support for the Oils YSH shell.
var Ysh Shell = ysh{}

const yshHook = `
proc _direnv_hook {
  var payload = $("{{.SelfPath}}" export ysh)
  if (payload !== '') {
    var diff = {}
    write -- $payload | json read (&diff)
    var props = propView(ENV)
    for key, value in (diff) {
      if (value === null) {
        call props->erase(key)
      } else {
        setvar props[key] = value
      }
    }
  }
}

var prompt_command = getVar('PROMPT_COMMAND')
if (prompt_command === null) {
  setvar prompt_command = ''
}
if (not (';' ++ prompt_command ++ ';').contains(';_direnv_hook;')) {
  if (prompt_command !== '') {
    setglobal PROMPT_COMMAND = '_direnv_hook; ' ++ prompt_command
  } else {
    setglobal PROMPT_COMMAND = '_direnv_hook'
  }
}
`

func (sh ysh) Hook() (string, error) {
	return yshHook, nil
}

func (sh ysh) Export(e ShellExport) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(e)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (sh ysh) Dump(env Env) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(env)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var (
	_ Shell = (*ysh)(nil)
)
