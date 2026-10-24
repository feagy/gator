package rssfeed

import (
    "net/http"
    "html"
    "encoding/xml"
    "context"
    "io"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
	    return nil, err
    }
    req.Header.Set("User-Agent", "gator")

    client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

    dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeedResp := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeedResp)
	if err != nil {
		return nil, err
	}
    
    rssFeedResp.Channel.Title = html.UnescapeString(rssFeedResp.Channel.Title)
    rssFeedResp.Channel.Description = html.UnescapeString(rssFeedResp.Channel.Description)
    for i, item := range rssFeedResp.Channel.Item {
        rssFeedResp.Channel.Item[i].Title = html.UnescapeString(item.Title)
        rssFeedResp.Channel.Item[i].Description = html.UnescapeString(item.Description)
    }

    return &rssFeedResp, nil
}
