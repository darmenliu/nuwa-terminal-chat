# Release Notes

## What's Changed

[Changes will be automatically inserted here by the release workflow]

## Installation

You can download the pre-built binary for your platform from the assets below, or install from source:

```bash
go install github.com/darmenliu/nuwa-terminal-chat@latest
```

## Environment Setup

Set the following environment variables:

```bash
# For Gemini (default)
export LLM_BACKEND=gemini
export LLM_MODEL_NAME=gemini-1.5-pro
export LLM_API_KEY=<your-api-key>
export LLM_TEMPERATURE=0.8

# For Ollama
export LLM_BACKEND=ollama
export LLM_MODEL_NAME=llama2
export LLM_API_KEY=apikey
export LLM_TEMPERATURE=0.8
export OLLAMA_SERVER_URL=http://localhost:8000

# For Groq
export LLM_BACKEND=groq
export LLM_MODEL_NAME=llama3-8b-8192
export LLM_API_KEY=<your-groq-api-key>
export LLM_TEMPERATURE=0.8
```

## Supported Platforms

- Linux (amd64, arm64)

## Known Issues

Please report any issues you encounter at our [GitHub Issues](https://github.com/darmenliu/nuwa-terminal-chat/issues) page. 