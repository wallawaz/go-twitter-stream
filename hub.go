package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"math/rand"
	"time"
)

// hub maintains the set of active clients and inbound tweets channel
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound tweets from twitterStream
	// stream *twitter.Stream
	twitterClient *twitter.Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// incoming tweets from demux
	hubTweets chan *twitter.Tweet

	// incoming votes from clients
	votes chan []byte

	// currentVotes Tally
	currentVotes map[string]int
}

func newHub(c *twitter.Client) *Hub {
	startingVotes := map[string]int{
		"POO": 0,
	}
	return &Hub{
		twitterClient: c,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		hubTweets:     make(chan *twitter.Tweet),
		votes:         make(chan []byte),
		currentVotes:  startingVotes,
	}
}

func (h *Hub) handle() {

	// filter tweets
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"Trump", "@realDonaldTrump"},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}
	stream, err := h.twitterClient.Streams.Filter(filterParams)
	if err != nil {
		panic(err)
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		rand := r1.Intn(10)
		if rand <= 5 {
			h.hubTweets <- tweet
			time.Sleep(time.Duration(1) * time.Second)

		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	go demux.HandleChan(stream.Messages)
}

func (h *Hub) run() {

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Client Connected:", client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				fmt.Println("Client Disconnected:", client)
				//clients not sending yet
				//close(client.send)
			}
		case tweet := <-h.hubTweets:
			for client := range h.clients {
				select {
				case client.recvTweets <- tweet:
				default:
					close(client.recvTweets)
					delete(h.clients, client)
				}
			}
		case vote := <-h.votes:
			voteStr := string(vote[:])
			if voteStr == "1" {
				h.currentVotes["POO"] += 1
			}
			for client := range h.clients {
				select {
				case client.recvVotes <- h.currentVotes:
				default:
					close(client.recvVotes)
				}
			}
		}
	}
}
