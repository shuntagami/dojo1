package typing_game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GameConfig struct {
	Contents      []string
	ContentsCount int
	Score         int
}

func Initialize(rootDIR string) (*GameConfig, error) {
	file, err := os.Open(filepath.Join(rootDIR, "words", "list.en.txt"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	client := &GameConfig{}
	for scanner.Scan() {
		client.ContentsCount++
		client.Contents = append(client.Contents, scanner.Text())
	}
	return client, nil
}

func (c *GameConfig) Play() error {
	rand.Seed(time.Now().UnixNano())
	expected := c.Contents[rand.Intn(c.ContentsCount)]
	fmt.Printf("Type: %s\n", expected)
	reader := bufio.NewReader(os.Stdin)

	actual, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	actual = strings.TrimSuffix(actual, "\n")
	if expected == actual {
		fmt.Println("OK")
		c.Score++
	} else {
		fmt.Println("NG")
	}
	return nil
}
