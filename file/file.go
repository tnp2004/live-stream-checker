package file

import (
	"encoding/csv"
	"log"
	"os"
)

const (
	CHANNEL_LIST_FILE_NAME = "channel_list.csv"
	PLATFORM_INDEX         = 0
	LINK_INDEX             = 1
)

type ChannelList struct {
	Platform string
	Link     string
}

func ReadChannelList() []*ChannelList {
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

	channelList := make([]*ChannelList, 0, len(records))
	for _, record := range records {
		ch := &ChannelList{
			Platform: record[PLATFORM_INDEX],
			Link:     record[LINK_INDEX],
		}
		channelList = append(channelList, ch)
	}

	return channelList
}
