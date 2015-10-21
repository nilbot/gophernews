package gophernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// create data structures

// Client type
type Client struct {
	BaseURI    string
	Version    string
	Suffix     string
	HTTPClient *http.Client
}

// Initializes and returns an API client
func NewClient(HTTPClient *http.Client) *Client {
	if HTTPClient == nil {
		HTTPClient = http.DefaultClient
	}
	var c Client
	c.BaseURI = "https://hacker-news.firebaseio.com/"
	c.Version = "v0"
	c.Suffix = ".json"
	c.HTTPClient = HTTPClient
	return &c
}

// GetStory makes an API request and puts response into a Story struct
func (c *Client) GetStory(id int) (Story, error) {
	item, err := c.getItem(id)

	if err != nil {
		return Story{}, err
	}

	if item.Type() != "story" {
		emptyStory := Story{}
		return emptyStory, fmt.Errorf("Called GetStory on ID #%v which is not a _story_. "+
			"Item is of type _%v_.", id, item.Type)
	}
	story := item.ToStory()
	return story, nil
}

// GetComment makes an API request and puts response into a Comment struct
func (c *Client) GetComment(id int) (Comment, error) {
	item, err := c.getItem(id)

	if err != nil {
		return Comment{}, err
	}

	if item.Type() != "comment" {
		emptyComment := Comment{}
		return emptyComment, fmt.Errorf("Called GetComment on ID #%v which is not a _comment_. "+
			"Item is of type _%v_.", id, item.Type)
	}
	comment := item.ToComment()
	return comment, nil
}

// GetPoll makes an API request and puts response into a Poll struct
func (c *Client) GetPoll(id int) (Poll, error) {
	item, err := c.getItem(id)

	if err != nil {
		return Poll{}, err
	}

	if item.Type() != "poll" {
		emptyPoll := Poll{}
		return emptyPoll, fmt.Errorf("Called GetPoll on ID #%v which is not a _poll_. "+
			"Item is of type _%v_.", id, item.Type)
	}
	poll := item.ToPoll()
	return poll, nil
}

// GetPart makes an API request and puts response into a Part struct
func (c *Client) GetPart(id int) (Part, error) {
	item, err := c.getItem(id)

	if err != nil {
		return Part{}, err
	}

	if item.Type() != "pollopt" {
		emptyPart := Part{}
		return emptyPart, fmt.Errorf("Called GetPart on ID #%v which is not a _part_. "+
			"Item is of type _%v_.", id, item.Type)
	}
	part := item.ToPart()
	return part, nil
}

// GetUser makes an API request and puts response into a User struct
func (c *Client) GetUser(id string) (User, error) {
	// TODO - refactor URL call into separate method
	url := c.BaseURI + c.Version + "/user/" + id + c.Suffix

	var u User

	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return u, err
	}

	// TODO - other checking around errors (wrong type, nonexistent user, etc.)
	return u, nil
}

// getItem makes an API request and puts response into a item struct
// items are then converted into Stories, Comments, Polls, and Parts (of polls)
func (c *Client) getItem(id int) (item, error) {
	url := c.BaseURI + c.Version + "/item/" + strconv.Itoa(id) + c.Suffix

	var i map[string]interface{}

	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return i, err
	}

	if string(body) == "404 page not found" {
		return i, fmt.Errorf("404 page not found")
	}

	err = json.Unmarshal(body, &i)

	return i, err
}

// GetTopStories makes an API request on top stories and fill the array of id
func (c *Client) GetTopStories() ([]int, error) {
	url := c.BaseURI + c.Version + "/topstories" + c.Suffix

	body, err := c.MakeHTTPRequest(url)

	var top100 []int

	err = json.Unmarshal(body, &top100)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return top100, nil
}

// GetMaxItem makes an API request and return the Item with max ID
func (c *Client) GetMaxItem() (Item, error) {
	url := c.BaseURI + c.Version + "/maxitem" + c.Suffix

	body, err := c.MakeHTTPRequest(url)

	var maxItemID int

	err = json.Unmarshal(body, &maxItemID)
	if err != nil {
		return item{}, err
	}

	maxItem, err := c.getItem(maxItemID)

	return maxItem, err
}

// GetChanges makes an API request and return a Changes type
func (c *Client) GetChanges() (Changes, error) {
	url := c.BaseURI + c.Version + "/updates" + c.Suffix

	body, err := c.MakeHTTPRequest(url)

	var changes Changes

	err = json.Unmarshal(body, &changes)

	return changes, err
}

// MakeHTTPRequest wraps a http.Get and return the byte slice
func (c *Client) MakeHTTPRequest(url string) ([]byte, error) {
	response, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf(http.StatusText(http.StatusNotFound))
	}
	return body, nil
}

// Convert an item to a Story
func (i item) ToStory() Story {
	var s Story
	s.By = i.By()
	s.ID = i.ID()
	s.Kids = i.Kids()
	s.Score = i.Score()
	s.Time = i.Time()
	s.Title = i.Title()
	s.Type = i.Type()
	s.URL = i.URL()
	return s
}

// Convert an item to a Comment
func (i item) ToComment() Comment {
	var c Comment
	c.By = i.By()
	c.ID = i.ID()
	c.Kids = i.Kids()
	c.Parent = i.Parent()
	c.Text = i.Text()
	c.Time = i.Time()
	c.Type = i.Type()
	return c
}

// Convert an item to a Poll
func (i item) ToPoll() Poll {
	var p Poll
	p.By = i.By()
	p.ID = i.ID()
	p.Kids = i.Kids()
	p.Parts = i.Parts()
	p.Score = i.Score()
	p.Text = i.Text()
	p.Time = i.Time()
	p.Title = i.Title()
	p.Type = i.Type()
	return p
}

// Convert an item to a Part
func (i item) ToPart() Part {
	var p Part
	p.By = i.By()
	p.ID = i.ID()
	p.Parent = i.Parent()
	p.Score = i.Score()
	p.Text = i.Text()
	p.Time = i.Time()
	p.Type = i.Type()
	return p
}

func main() {
	client := NewClient(nil)

	// README
	s, err := client.GetStory(8412605) //=> Actual Story
	// c, err := client.GetComment(2921983) //=> Actual Comment
	// p, err := client.GetPoll(126809) //=> Actual Poll
	// pp, err := client.GetPart(160705) //=> Actual Part of Poll
	// u, err := client.GetUser("pg") //=> User

	if err != nil {
		panic(err)
	} else {
		// fmt.Println(u.About, "\n", u.Created, "\n", u.Karma)
		fmt.Println(s.By, "\n", s.Title, "\n", s.Score)
	}

	// write accessors to get stories, comments, polls, parts, and users
	// write accessors for special cases (top stories, updates, etc.)
	// write special accessors for stories, comments, etc. to get objects instead of
	// IDs (ints) of parents, children, etc.
}
