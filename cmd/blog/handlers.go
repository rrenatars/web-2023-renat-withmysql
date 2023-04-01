package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type indexPageData struct {
	FeaturedPosts []featuredPostData
	MostRecent    []mostRecentData
}

type featuredPostData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	PublishDate  string `db:"publish_date"`
	Author 		 string `db:"author"`
	AuthorAvatar string `db:"author_url"`
	Image		 string `db:"image_url"`
}

type mostRecentData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	PublishDate  string `db:"publish_date"`
	Author 		 string `db:"author"`
	AuthorAvatar string `db:"author_url"`
	Image		 string `db:"image_url"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		featuredPosts, err := featuredPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return 
		}

		recentPosts, err := mostRecent(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return 
		}

		ts, err := template.ParseFiles("pages/index.html") 
		if err != nil {
			http.Error(w, "Internal Server Error", 500) 
			log.Println(err)
			return 
		}

		data := indexPageData{
			FeaturedPosts: featuredPosts,
			MostRecent:    recentPosts,
		}

		err = ts.Execute(w, data) 
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/post.html") 
	if err != nil {
		http.Error(w, "Internal Server Error", 500) 
		log.Println(err.Error())                    
		return                                      
	}

	err = ts.Execute(w, nil) 
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func featuredPosts(db *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url
		FROM
			post
		WHERE featured = 1
	` 

	var posts []featuredPostData 

	err := db.Select(&posts, query) 
	if err != nil {                
		return nil, err
	}

	return posts, nil
}

func mostRecent(db *sqlx.DB) ([]mostRecentData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url
		FROM
			post
		WHERE featured = 0
	` 

	var posts []mostRecentData 

	err := db.Select(&posts, query) 
	if err != nil {                 
		return nil, err
	}

	return posts, nil
}
