package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/internal/store"
	"math/rand"
)

var usernames = []string{
	"shadowfox21",
	"pixelstorm",
	"neonbyte",
	"quantumjay",
	"crypticnova",
	"lunar_echo",
	"bytebandit",
	"silentvector",
	"glitchrider",
	"orbitzero",
	"codephantom",
	"velvetlogic",
	"zerogravity",
	"hexhunter",
	"midnightapi",
	"datawanderer",
	"syntaxsage",
	"cloudstriker",
	"binaryblaze",
	"alpha_kernel",
	"nexusflare",
	"voidnavigator",
	"debugdragon",
	"echo_circuit",
	"fusioncoder",
	"omega_patch",
	"pixelpirate",
	"stackshifter",
	"cybertrail",
	"infinite_loop",
	"astralbyte",
	"codecascade",
	"vortexrunner",
	"digitalnomicon",
	"hyperthreader",
	"nebulacoder",
	"phantomstack",
	"runtime_rebel",
	"circuitdream",
	"logicforge",
	"cosmicbuffer",
	"packetpilot",
	"bytealchemy",
	"devshifter",
	"quantumcache",
	"syntaxshaper",
	"kernelknight",
	"datadrifter",
	"bitvoyager",
	"looplegend",
	"zerodayz",
	"scriptwarden",
	"fluxcoder",
	"arrayassassin",
	"binarynomad",
	"hashharvester",
	"stackoracle",
	"cybersprint",
	"modulomancer",
	"threadtitan",
	"compileking",
	"memorymystic",
	"debugdynamo",
	"codevoyant",
	"bytevanguard",
	"asyncarcher",
	"logicwanderer",
	"virtualviper",
	"packetwhisper",
	"functionfury",
	"stackraider",
	"nullnavigator",
	"devoverdrive",
	"scriptphantom",
	"bitstormer",
	"quantumforge",
	"codeharbinger",
	"dataknighthood",
	"pixeloverlord",
	"recursionhero",
	"fractalcoder",
	"bytehorizon",
	"signalrunner",
	"dataphantom",
	"codelattice",
	"orbitcoder",
	"hashvoyager",
	"bitforge",
	"threadpulse",
	"logicnova",
	"compilecraft",
	"devnebula",
	"stackpulse",
	"binaryhorizon",
	"kernelpulse",
	"scriptforge",
	"packetnova",
	"debugpulse",
	"asyncvoyager",
	"functionnova",
}

var titles = []string{
	"To Kill a Mockingbird",
	"1984",
	"Pride and Prejudice",
	"The Great Gatsby",
	"The Catcher in the Rye",
	"Jane Eyre",
	"Wuthering Heights",
	"Moby Dick",
	"The Lord of the Rings",
	"The Hunger Games",
	"Gone with the Wind",
	"The Adventures of Huckleberry Finn",
	"The Picture of Dorian Gray",
	"Frankenstein",
	"Dracula",
	"The Count of Monte Cristo",
	"The Scarlet Letter",
	"War and Peace",
	"Moby Dick",
	"The Age of Innocence",
	"The Color Purple"}
var content = []string{
	"Explaining the darkness of the human psyche",
	"Romance and love in Victorian England",
	"Mysteries and adventures on the high seas",
	"Epic quests and legendary heroes",
	"Survival and rebellion in a dystopian world",
	"Historical events through a fictional lens",
	"Coming-of-age stories that resonate with readers",
	"The power of revenge and redemption",
	"Exploring the unknown through science fiction",
	"Vampires, ghosts, and supernatural creatures",
	"A tale of obsession and corruption",
	"Unraveling the mysteries of crime and detection",
	"A glimpse into the past's social commentary",
	"Inspiring stories that transcend time and cultures",
}
var tags = []string{
	"go",
	"api",
	"backend",
	"database",
	"sql",
	"grpc",
	"rest",
	"json",
	"authentication",
	"authorization",
	"middleware",
	"microservices",
	"docker",
	"kubernetes",
	"testing",
	"logging",
	"monitoring",
	"caching",
	"performance",
	"security",
}
var comments = []string{
	"This looks great!",
	"Needs a bit more testing.",
	"I think there’s a bug here.",
	"Can we optimize this later?",
	"Nice work overall.",
	"I’m not sure this is the best approach.",
	"Let’s revisit this section.",
	"This part is confusing.",
	"Works as expected.",
	"Consider adding more logging.",
	"Edge cases might break this.",
	"Looks clean and readable.",
	"Maybe refactor this function.",
	"Good use of interfaces.",
	"This could be simplified.",
	"I like this solution.",
	"Performance might be an issue here.",
	"Add some comments for clarity.",
	"Test coverage is missing here.",
	"Ship it 🚀",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)

	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error creating user", user, err)
			return
		}

	}
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post", post, err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating post", comment, err)
			return
		}
	}
	log.Println("Seeding Complete")
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)

	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: content[rand.Intn(len(content))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, n)
	for i := 0; i < n; i++ {
		cms[i] = &store.Comment{
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
