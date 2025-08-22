#!/bin/sh

if [ "$KX_DID_SOURCE" ]; then
	return 0
fi

export PROMPT_COMMAND="$KX_BIN checktime || exit"

export PS1="\$($KX_BIN prompt)
${PS1}"

export KX_DID_SOURCE=true

echo "---START---"