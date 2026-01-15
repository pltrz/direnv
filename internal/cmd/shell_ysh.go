package cmd

import "fmt"

type ysh struct{}

var Ysh Shell = ysh{}

const yshHook = `
proc _direnv_hook {
  var previous_exit_status = _status
  var vars = $(
    "{{.SelfPath}}" export ysh
  )
  trap --add SIGINT {
    true
  }
  eval $vars
  hash -r
  trap --remove SIGINT
  return $previous_exit_status
}

var prompt_command = getVar('PROMPT_COMMAND')
if (prompt_command === null) {
  setvar prompt_command = ''
}
if (not (';' ++ prompt_command ++ ';').contains(';_direnv_hook;')) {
  if (prompt_command !== '') {
    setglobal PROMPT_COMMAND = "_direnv_hook;" ++ prompt_command
  } else {
    setglobal PROMPT_COMMAND = "_direnv_hook"
  }
}
`

func (sh ysh) Hook() (string, error) {
	return yshHook, nil
}

func (sh ysh) Export(e ShellExport) (string, error) {
	var out string
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out, nil
}

func (sh ysh) Dump(env Env) (string, error) {
	var out string
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out, nil
}

func (sh ysh) export(key, value string) string {
	return "setglobal ENV." + key + " = " + sh.escape(value) + ";"
}

func (sh ysh) unset(key string) string {
	return "setglobal ENV." + key + " = null;"
}

func (sh ysh) escape(str string) string {
	return YshEscape(str)
}

// YshEscape returns a single-line b'' J8 string with only ASCII bytes.
func YshEscape(str string) string {
	out := "b'"
	for _, b := range []byte(str) {
		switch b {
		case '\b':
			out += `\b`
		case '\f':
			out += `\f`
		case '\n':
			out += `\n`
		case '\r':
			out += `\r`
		case '\t':
			out += `\t`
		case '\\':
			out += `\\`
		case '\'':
			out += `\'`
		default:
			if b < 0x20 || b == 0x7f || b >= 0x80 {
				out += fmt.Sprintf("\\y%02x", b)
			} else {
				out += string([]byte{b})
			}
		}
	}
	out += "'"
	return out
}
