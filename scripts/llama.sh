#!/usr/bin/env bash

manllama() {
    local man_cmd="$1"
    local query="$2"
    local model="$3"

    if [[ -z "$query" || -z "$man_cmd" ]]; then
        echo "Usage: manllama <man-cmd> <query> <model>"
        echo "Model can be medium or high (default: medium)"
        echo "Medium model is mistral, High is glm-4.7-flash"
        return 1
    fi

    local model_name="gemma4"
    if [[ "$model" == "high" ]]; then
        model_name="qwen3.6"
    fi

    local documentation=$(eval "$man_cmd")
    if [[ -z "$documentation" ]]; then
        echo "Error: No documentation found for $man_cmd"
        return 1
    fi

    ollama run $model_name --think=false "You will now help mew answer a question using some potentially relevant documentation. Here is the documentation: <documentation>$documentation</documentation>. If the exact answer to the query is not in the documentation, DON'T ANSWER and simply say 'the documentation provided does not contain the information you are looking for'. If you provide any unverified information, ALWAYS mention what is unverified from the documentation. Now answer my question: <question>$query</question>. DON'T ANSWER UNVERIFIED INFORMATION."
}

shellama() {
	local query="$*"

    local prompt=$(cat <<EOF
Give me the shell command to run to execute the provided query. Only output the command, no introduction, no backticks, no code block, no markdown, just the command.

The current environment has access to all common linux commands. Prefer modern alternatives like rg, fd & sd. Other tools available are ffmpeg, fzf, docker, node, python, jq, duckdb, git and more.

Give the command for the user query now: <query>$query</query>.
EOF
)

	local cmd=$(ollama run gemma4 --think=false "$prompt")

    echo "$cmd"
    wl-copy "$cmd"

    print -z -- "$cmd"
}
