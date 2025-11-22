package db

import (
	"context"
	"math/rand"
	"fmt"
	"log"
	"social/internal/store"
)

var usernames = []string{
	"alice",
	"bob",
	"carol",
	"dave",
	"erin",
	"frank",
	"grace",
	"heidi",
	"ivan",
	"judy",
	"mallory",
	"oscar",
	"peggy",
	"trent",
	"victor",
	"wendy",
	"yvonne",
	"zack",
	"nina",
	"oliver",
	"paula",
	"quinn",
	"rachel",
	"steve",
	"tina",
	"ursula",
	"vince",
	"walter",
	"xavier",
	"yuri",
	"zoe",
	"harry",
	"isabel",
	"jake",
	"karen",
	"liam",
	"mia",
	"noah",
	"olga",
	"peter",
	"queen",
	"roger",
	"sara",
	"tom",
	"uma",
	"vera",
	"will",
	"xena",
	"yara",
	"zeke",
}

var titles = []string{

	"Getting Started with Go",
	"Understanding Goroutines",
	"Mastering Concurrency",
	"Building REST APIs",
	"Error Handling Best Practices",
	"Testing in Go",
	"Working with JSON",
	"Structs and Interfaces",
	"Go Modules Explained",
	"Optimizing Performance",
	"Logging and Monitoring",
	"Database Access with Go",
	"Writing Clean Go Code",
	"Microservices in Go",
	"Deploying Go Applications",
	"Go for CLI Tools",
	"Configuration Management",
	"Handling Files in Go",
	"Networking with net/http",
	"Tips for Go Beginners",
}

var contents = []string{
	"Hello world!",
	"Just deployed a new feature 🚀",
	"Working on some Go code today.",
	"Any recommendations for a good book?",
	"Coffee break time ☕",
	"Debugging is like solving a mystery.",
	"Unit tests are finally passing!",
	"Refactoring old code feels so good.",
	"Taking a quick walk outside.",
	"Just smashed a nasty bug 💥",
	"Learning concurrency with goroutines.",
	"What’s your favorite programming language?",
	"Trying out a new Dark Mode theme.",
	"Reading about system design patterns.",
	"Music + coding = perfect combo.",
	"Who else loves clean code?",
	"Documenting APIs this morning.",
	"Pair programming session went great.",
	"Time to review some pull requests.",
	"Calling it a day. Good work everyone!",
}

var tags = []string{
	"golang",
	"backend",
	"api",
	"concurrency",
	"microservices",
	"testing",
	"performance",
	"database",
	"cli",
	"http",
	"json",
	"logging",
	"devops",
	"cloud",
	"docker",
	"kubernetes",
	"best-practices",
	"clean-code",
	"tutorial",
	"beginner",
}
var comments = []string{
	"This is super helpful, thanks!",
	"I’ve been stuck on this for hours, glad I found this.",
	"Can you share a full example?",
	"Nice explanation 👍",
	"This didn’t work for me, any ideas why?",
	"Amazing, it works perfectly now!",
	"Is this still the recommended approach?",
	"Which Go version are you using for this?",
	"Can you compare this with another method?",
	"Short and clear, love it.",
	"I’m new to Go, this was easy to follow.",
	"What about error handling in this case?",
	"Does this scale well under heavy load?",
	"Thanks, learned something new today!",
	"Could you add tests for this?",
	"Any performance benchmarks available?",
	"This fixed my production issue 🙏",
	"Great writeup, keep them coming!",
	"I had the same problem, this solved it.",
	"Subscribing for more content like this.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding completed")

}


func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			UserName: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}

	}
	return posts
}


func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
    cms := make([]*store.Comment, num)

    for i := 0; i < num; i++ {
        post := posts[rand.Intn(len(posts))]
        user := users[rand.Intn(len(users))]

        cms[i] = &store.Comment{
            PostID:  post.ID,
            UserID:  user.ID,
            Content: comments[rand.Intn(len(comments))], 
        }
    }
    return cms
}
