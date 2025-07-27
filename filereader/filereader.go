package filereader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/tnp2004/live-stream-checker/models"
)

const (
	NAME_INDEX = iota
	PLATFORM_INDEX
	LINK_INDEX
)

const CHANNEL_LIST_FILE_NAME = "channel_list.csv"
const UNCHECKED_STATUS = "unchecked"

func ReadChannelList() []*models.Channel {
	file, err := os.Open(CHANNEL_LIST_FILE_NAME)
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	// trim header
	records = records[1:]

	channelList := make([]*models.Channel, 0, len(records))
	for _, record := range records {
		ch := &models.Channel{
			Name:     record[NAME_INDEX],
			Platform: record[PLATFORM_INDEX],
			Link:     record[LINK_INDEX],
			Status:   UNCHECKED_STATUS,
		}
		channelList = append(channelList, ch)
	}

	return channelList
}

func AddChannel(name, url string) error {
	file, err := os.OpenFile(CHANNEL_LIST_FILE_NAME, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	defer file.Close()

	regex := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(youtube|twitch)\.(com|tv)`)

	match := regex.FindStringSubmatch(url)
	if len(match) < 1 {
		log.Println("Error: invalid platform")
		return fmt.Errorf("invalid platform")
	}

	row := []string{name, match[1], url}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(row)

	return nil
}
