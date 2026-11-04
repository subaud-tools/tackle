# proj — thin shim over the tackle `proj` binary.
#
# The binary renders its TUI on the terminal and, for the one case it cannot
# handle itself (attaching a bare shell to a session on another server), writes
# an `exec tmux … attach` line to $PROJ_EXEC. We eval that line here so the
# attach replaces THIS shell. The in-tmux switch/detach cases run in-process in
# the binary and leave $PROJ_EXEC empty.
proj() {
  command -v proj >/dev/null 2>&1 || { echo "proj: binary not on PATH" >&2; return 127; }
  local ef; ef="$(mktemp "${TMPDIR:-/tmp}/proj-exec.XXXXXX")" || return 1
  PROJ_EXEC="$ef" command proj "$@"
  local rc=$?
  if [[ -s "$ef" ]]; then
    local cmd; cmd="$(<"$ef")"; rm -f "$ef"
    [[ -n "$cmd" ]] && eval "$cmd"
    return
  fi
  rm -f "$ef"
  return $rc
}

# A terminal app launched from inside tmux passes TMUX and TMUX_PANE on to every
# shell it opens, though those shells are not in tmux. A shell really in a pane
# sits on that pane's terminal, so keep the pair only when $TMUX_PANE's tty on
# the server in $TMUX is this shell's own; otherwise unset both, so proj and
# auto-join treat the shell as outside tmux. No-op when tmux is not on PATH yet
# (06-proj.zsh calls __proj_autojoin again once it is).
__proj_clear_stale_tmux() {
  [[ -z "${TMUX:-}" ]] && return 0
  command -v tmux >/dev/null 2>&1 || return 0
  local pt=""
  [[ -n "${TMUX_PANE:-}" ]] &&
    pt="$(tmux -S "${TMUX%%,*}" display -p -t "$TMUX_PANE" '#{pane_tty}' 2>/dev/null)"
  [[ -n "$pt" && "$pt" == "$(tty 2>/dev/null)" ]] && return 0
  unset TMUX TMUX_PANE
}

# Auto-join: a fresh interactive shell inside a project dir (and not already in
# tmux) opens that project's proj view — so ⌘T in a project's Ghostty window
# lands on proj. Opt-outs preserved from the zsh reference
# (~/dotfiles/config/zsh/06-tmux-autojoin.zsh):
#   * a stale TMUX is cleared first (__proj_clear_stale_tmux)
#   * interactive shells only
#   * skip if already inside tmux ($TMUX)
#   * skip per-shell via NO_AUTO_TMUX
#   * skip globally via ~/.no-auto-tmux
#   * silent no-op if tmux or proj is absent, or $PWD is outside every root
# AUTO_CLAUDE=1 keeps its meaning: the project view is opened with claude
# preselected as the agent for any new work.
__proj_autojoin() {
  __proj_clear_stale_tmux
  [[ $- != *i* ]] && return 0                 # interactive only
  [[ -n "${TMUX:-}" ]] && return 0            # already in tmux
  [[ -n "${NO_AUTO_TMUX:-}" ]] && return 0    # per-shell opt-out
  [[ -f "$HOME/.no-auto-tmux" ]] && return 0  # global opt-out
  command -v tmux >/dev/null 2>&1 || return 0
  command -v proj >/dev/null 2>&1 || return 0
  # Is $PWD inside a known project? Ask the binary (exit 0 + non-empty name).
  local name; name="$(command proj __autojoin-project 2>/dev/null)"
  [[ -z "$name" ]] && return 0
  if [[ -n "${AUTO_CLAUDE:-}" ]]; then
    proj --claude "$name"
  else
    proj "$name"
  fi
}

# Fire once at source time — we want this immediately, not on the next precmd.
__proj_autojoin
