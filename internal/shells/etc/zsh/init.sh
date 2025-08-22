#!/bin/sh

if [ "$KX_DID_SOURCE" ]; then
	return 0
fi

precmd() {
	eval "$PROMPT_COMMAND"
}

export PROMPT_COMMAND="$KX_BIN checktime || exit"

export PS1="
%{%F{#111111}%}\$($KX_BIN prompt)%{%f%}${PS1}"

export KX_DID_SOURCE=true

echo "---START---"