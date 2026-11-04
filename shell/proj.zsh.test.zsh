#!/usr/bin/env zsh
# Tests for shell/proj.zsh.
#
# Floor: `zsh -n` must parse proj.zsh without error. Bonus: a shim-eval smoke
# that stubs the `proj` binary to write a known line to $PROJ_EXEC and asserts
# the shim would eval it (with a no-op eval override so the test stays inert).

emulate -L zsh
set -u

here="${0:A:h}"
target="$here/proj.zsh"
fail=0

# --- Floor: syntax check -----------------------------------------------------
if zsh -n "$target"; then
  print -r -- "ok: zsh -n proj.zsh"
else
  print -r -- "FAIL: zsh -n proj.zsh" >&2
  fail=1
fi

# --- Smoke: shim evals the $PROJ_EXEC line -----------------------------------
# Run in a subshell so stubs/overrides don't leak. The stub `proj` writes a
# sentinel to $PROJ_EXEC; `command` is stubbed so `command proj` hits it and the
# on-PATH check passes; `eval` is overridden to capture rather than execute.
smoke() {
  # Suppress the source-time auto-join hook: pretend we're already in tmux.
  TMUX=fake NO_AUTO_TMUX=1 source "$target"

  local captured=""
  captured_file="$(mktemp "${TMPDIR:-/tmp}/proj-smoke.XXXXXX")"

  # Stub `command`: `command -v proj` succeeds; `command proj ...` writes the
  # sentinel to $PROJ_EXEC (the binary's contract).
  command() {
    if [[ "$1" == "-v" ]]; then
      return 0
    fi
    if [[ "$1" == "proj" ]]; then
      print -r -- 'SENTINEL_ATTACH' >| "$PROJ_EXEC"
      return 0
    fi
    builtin command "$@"
  }

  # Override eval to capture the string the shim would have executed.
  eval() { print -r -- "$*" >| "$captured_file"; }

  proj somearg
  captured="$(<"$captured_file")"
  rm -f "$captured_file"

  [[ "$captured" == "SENTINEL_ATTACH" ]]
}

if ( smoke ); then
  print -r -- "ok: shim evals \$PROJ_EXEC line"
else
  print -r -- "FAIL: shim did not eval \$PROJ_EXEC line" >&2
  fail=1
fi

# --- Stale TMUX: __proj_clear_stale_tmux against a throwaway server ----------
# A terminal launched from inside tmux hands its shells a TMUX and TMUX_PANE
# they are not in. The check keeps them only when the pane's tty is the shell's
# own. Every case runs against a private server (-L, no config), never the
# user's, with the outer tmux and git environment cleared.
unset TMUX TMUX_PANE
for v in ${(k)parameters[(I)GIT_*]}; do unset $v; done

if ! builtin command -v tmux >/dev/null 2>&1; then
  print -r -- "skip: stale TMUX cases (tmux not on PATH)"
  exit $fail
fi

sock="proj-zsh-test-$$"
scratch="$(mktemp -d "${TMPDIR:-/tmp}/proj-zsh-test.XXXXXX")"
t() { tmux -L "$sock" -f /dev/null "$@"; }
socket=""
trap 't kill-server >/dev/null 2>&1; [[ -n "$socket" ]] && rm -f "$socket"; rm -rf "$scratch"' EXIT

# A session that only holds a pane open, so the server and pane stay up.
t new-session -d -s keep -x 80 -y 24 'sleep 3600'
socket="$(t display -p -t keep '#{socket_path}')"
pane="$(t display -p -t keep '#{pane_id}')"
server_tmux="$socket,$(t display -p -t keep '#{pid}'),0"

# stale <name> <TMUX> <TMUX_PANE> <want> [PATH]: in a subshell, source the
# wrapper outside tmux, set the pair, run the check and compare what is left
# ("unset" for an unset variable) with want.
stale() {
  local name=$1 want=$4 got
  got="$(
    unset TMUX TMUX_PANE
    source "$target"
    export TMUX="$2" TMUX_PANE="$3"
    [[ -n "${5:-}" ]] && PATH="$5"
    __proj_clear_stale_tmux
    print -r -- "${TMUX:-unset} ${TMUX_PANE:-unset}"
  )"
  if [[ "$got" == "$want" ]]; then
    print -r -- "ok: $name"
  else
    print -r -- "FAIL: $name: want '$want', got '$got'" >&2
    fail=1
  fi
}

# 1. A shell really in the pane keeps TMUX, also with stdin redirected (an
# editor reading the shell's environment runs `zsh -i -c … </dev/null`).
in_pane() {
  local name=$1 session=$2 redirect=$3
  local out="$scratch/$session"
  print -r -- "before=\$TMUX; source ${(q)target}; __proj_clear_stale_tmux; print -r -- \"\$before|\${TMUX:-unset}\" >| ${(q)out}" >| "$scratch/$session.zsh"
  t new-session -d -s "$session" -x 80 -y 24 "zsh -f ${(q)scratch}/$session.zsh $redirect; sleep 3600"
  for i in {1..50}; do [[ -s "$out" ]] && break; sleep 0.1; done
  if [[ -s "$out" ]]; then
    local res="$(<"$out")"
    if [[ "${res%%|*}" == "$socket,"* && "${res#*|}" == "${res%%|*}" ]]; then
      print -r -- "ok: $name"
    else
      print -r -- "FAIL: $name: got '$res'" >&2
      fail=1
    fi
  else
    print -r -- "FAIL: $name: the pane wrote nothing" >&2
    fail=1
  fi
}
in_pane "real pane keeps TMUX" inpane ""
in_pane "real pane with stdin redirected keeps TMUX" inpane-nostdin "</dev/null"

# 2–5. Outside the pane, every stale shape loses both.
stale "outside a real pane loses both" "$server_tmux" "$pane" "unset unset"
stale "empty TMUX_PANE loses both" "$server_tmux" "" "unset unset"
stale "missing pane loses both" "$server_tmux" "%99999" "unset unset"
stale "gone server loses both" "$scratch/no-such-socket,1,0" "$pane" "unset unset"

# 6. With tmux not on PATH yet, both are left alone.
stale "no tmux on PATH keeps both" "$server_tmux" "$pane" "$server_tmux $pane" /nonexistent

exit $fail
