package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"learn-go/ollama"
	"learn-go/weather"

	"github.com/joho/godotenv"
)

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if utf8.RuneCountInString(currentLine)+utf8.RuneCountInString(word)+1 <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func printWeather(w *tabwriter.Writer, weatherData *ollama.WeatherResponse) {
	fmt.Fprintf(w, "\n%s\n", strings.Repeat("─", 50))
	fmt.Fprintf(w, "🌡️  Weather Data\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 50))

	fmt.Fprintf(w, "Temperature:\t%.1f°F\n", weatherData.Weather.Temperature)
	if weatherData.Weather.FeelsLike > 0 {
		fmt.Fprintf(w, "Feels Like:\t%.1f°F\n", weatherData.Weather.FeelsLike)
	}
	fmt.Fprintf(w, "Conditions:\t%s\n", weatherData.Weather.Conditions)
	fmt.Fprintf(w, "Humidity:\t%d%%\n", weatherData.Weather.Humidity)

	if weatherData.Weather.WindSpeed > 0 {
		fmt.Fprintf(w, "Wind Speed:\t%.1f %s\n", weatherData.Weather.WindSpeed, weatherData.Weather.WindSpeedUnit)
	}
	if weatherData.Weather.WindDirection != "" {
		fmt.Fprintf(w, "Wind Direction:\t%s\n", weatherData.Weather.WindDirection)
	}
	if weatherData.Weather.WindGust > 0 {
		fmt.Fprintf(w, "Wind Gust:\t%.1f %s\n", weatherData.Weather.WindGust, weatherData.Weather.WindSpeedUnit)
	}

	if weatherData.Weather.Visibility > 0 {
		fmt.Fprintf(w, "Visibility:\t%.1f miles\n", weatherData.Weather.Visibility)
	}
	if weatherData.Weather.Pressure > 0 {
		fmt.Fprintf(w, "Pressure:\t%.1f mb\n", weatherData.Weather.Pressure)
	}
	if weatherData.Weather.DewPoint > 0 {
		fmt.Fprintf(w, "Dew Point:\t%.1f°F\n", weatherData.Weather.DewPoint)
	}
	if weatherData.Weather.UVIndex > 0 {
		fmt.Fprintf(w, "UV Index:\t%.1f\n", weatherData.Weather.UVIndex)
	}
	if weatherData.Weather.CloudCover > 0 {
		fmt.Fprintf(w, "Cloud Cover:\t%d%%\n", weatherData.Weather.CloudCover)
	}
	if weatherData.Weather.PrecipitationChance > 0 {
		fmt.Fprintf(w, "Precipitation:\t%d%%\n", weatherData.Weather.PrecipitationChance)
	}

	fmt.Fprintf(w, "Last Updated:\t%s\n", weatherData.Weather.Timestamp)
	w.Flush()

	fmt.Printf("\n%s\n", strings.Repeat("─", 50))
	fmt.Println("🤖 AI Description")
	fmt.Printf("%s\n", strings.Repeat("─", 50))
	fmt.Println(wrapText(weatherData.Description, 50))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	weatherSvc := weather.NewNWSService()
	client := ollama.NewClient(weatherSvc)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Printf("\n%s\n", strings.Repeat("═", 50))
	fmt.Println("🌤️  Weather Assistant")
	fmt.Printf("%s\n", strings.Repeat("─", 50))
	fmt.Println("Type 'quit' to exit")
	fmt.Println("Ask me about the weather anywhere in the US!")
	fmt.Println("Example: 'What's the weather like in Miami?'")
	fmt.Printf("%s\n", strings.Repeat("═", 50))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}

		if strings.ToLower(query) == "quit" {
			fmt.Println("\nGoodbye! 👋")
			break
		}

		weatherData, err := client.GetWeather(query, 3)
		if err != nil {
			fmt.Printf("\n❌ Error: %v\n", err)
			continue
		}

		printWeather(w, weatherData)
	}
}
