tour-of-go-web-crawler
======================

Personal answer for the last excersize, Web Crawler, on a Tour of Go

I enjoyed working on channels of golang thinking about answers to [the last excersize of a Tour of Go](http://tour.golang.org/#73)

- The idea of handing the map over a channel is from [A Tour of Go #69 Exercise: Web Crawler | Sonia Codes](http://soniacodes.wordpress.com/2011/10/09/a-tour-of-go-69-exercise-web-crawler/)
- The final code prints messages sent through a channel so multiple lines are seperated while `Crawl()` go routines are running concurrently
- The main function waits for all the `Crawl()` to finish using `sync.WaitGroup`. The point here was that I had to pass reference to the WaitGroup rather than the value
