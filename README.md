# Weather CLI with Ollama AI

A command-line interface for getting weather information with natural language processing powered by Ollama AI.

## Features

- Natural language queries for weather information
- AI-powered weather descriptions
- Detailed weather data including:
  - Temperature and "feels like" temperature
  - Weather conditions
  - Wind speed, direction, and gusts
  - Humidity and pressure
  - Visibility and cloud cover
  - UV index
  - Precipitation chance
- Clean, formatted terminal output

## Prerequisites

- Go 1.21 or higher
- [Ollama](https://ollama.ai/) installed and running
- The phi model installed in Ollama

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/weather-cli
cd weather-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file:
```bash
cp .env.example .env
```

4. Update the .env file with your settings

5. Install the phi model in Ollama:
```bash
ollama run phi4
```

## Usage

1. Start the Ollama server:
```bash
ollama serve
```

2. Run the CLI:
```bash
go run main.go
```

3. Ask about the weather using natural language:
```
> What's the weather like in Miami?
> Is it raining in Seattle right now?
> How's the temperature in New York City?
```

4. Type 'quit' to exit

## Environment Variables

- `API_KEY`: Your API key for authentication
- `RATE_LIMIT`: Rate limit for requests (default: 1)
- `PORT`: Server port (default: 8080)
- `OLLAMA_URL`: Ollama API URL (default: http://localhost:11434/api)
- `OLLAMA_MODEL`: Ollama model to use (default: phi4)

## License

MIT License

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. 