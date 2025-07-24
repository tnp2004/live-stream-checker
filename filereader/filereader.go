package filereader

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/tnp2004/live-stream-checker/models"
)

const (
	CHANNEL_LIST_FILE_NAME = "channel_list.csv"
	PLATFORM_INDEX         = 0
	LINK_INDEX             = 1
)

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
			Platform: record[PLATFORM_INDEX],
			Link:     record[LINK_INDEX],
		}
		channelList = append(channelList, ch)
	}

	return channelList
}
