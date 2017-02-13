api-score
=========

api-score is a project that creates a "P score" for a business, where the score
represents the likelihood that a business will be able to repay it's loans. 1
meaning 100% likely and 0 meaning 0% likely.

# Usage

To run api-score, you'll need to have [Golang](https://golang.org/dl/)
installed.

Setup your Twitter API config by editing the `config.sh.example` file (and
moving it to `config.sh` if you'd like), and then running `source config.sh`.

Compile the program using `go build`, and run it by specifying the required
flags, e.g.

```
./api-score --twuser stripe --business Stripe --owner Stripe
```

The `--twuser` flag is required, and is the twitter user name of the company in
question. The `--business` and `--owner` flags are optional.

You may also specify a `--verbose` flag to show additional information in the
output.

# P Score Rationale

Admittedly I am not very familiar with how to evaluate a business's
credit-worthiness. However, I do believe that a business which is more
consistent is a good thing. So I structured my use of the Twitter API to
calculate a general level of consistent tweeting, and used that to calculate a
"P score".

### Tweet Counts

A company that consistently tweets a slightly larger amount over time is
considered healthy. If their number of tweets drops off greatly, that's likely
a bad sign. Likewise, if they start tweeting a much greater amount, that could
be a sign of desperation or lack of customer traction. Therefore, comparing the
most recent 60 days of tweeting with the last 180 days will give a sufficient
time to measure an increase or decrease in number of tweets. An increase of 5%
to 10% is considered healthy, and everything else is penalized depending on how
far away from this the percentage increase is.

### Tweets vs Replies

A company is likely to have both original tweets and replies to customers'
tweets. A company that is only replying to tweets might not be sufficiently
advertising to or educating customers. A ratio of 1 original to 100 replies was
used as the lower bounds for this. A significant, but more minor penalty of 0.3
is applied is this case.

Likewise, if a company is tweeting many more times than they're replying, they
might not have sufficient customer interaction or engagement. A ratio of 1
original tweet to 1 reply tweet is used as the upper bound for this. Because
this is likely a larger indication of a problem than replying to customers, a
larger penalty of 0.5 is applied.

# Design Decisions

### API client library `anaconda`

A large part of api-score is working with the Twitter API. Though it would be
possible to implement an API client of my own, that is out of scope (and likely
unnecessary) for this program. Therefore, I searched the Twitter documentation,
as well as Github, for a good library to use.
[anaconda](https://github.com/ChimeraCoder/anaconda) was the recommended
library, and it didn't have any outstanding issues relating to getting user
timelines, so it seemed like a good choice.

An additional reason for the choice of `anaconda` is the library's use of Go
standard library structures for HTTP requests. Though none of my tests ended up
using this, a test I would write given more time would allow me to use the
default HttpClient in Go to fake sending messages and give fake results. That
way I could test more complicated code that otherwise would require use of the
Twitter API.

### Language

I used Go to write this program, both for my own comfort with the language, as
well as Go's powerful standard library that handles API work fairly seamlessly.
The fact that Go has a good standard library for testing is a nice bonus.

I did rely on an external library for the Twitter specific auth, as well as
implementation of the individual endpoints, so as not to duplicate work that had
already been done. Though I did read much of the library code when architecting
my solution, to ensure it was doing what I expected.

Another reason I used Go is that it is a fairly simple and readable language. It
has many of the same keywords as other languages (if statements, for loops,
etc), and uses them heavily over many other abstractions. It is likely that
those that are not fluent in Go would still be able to make some sense of it,
particularly when a program stays away from many Go specific abstractions
(goroutines, channels, etc.), as this program does.

### Configuration

Go has an excellent standard library package called
[flag](https://golang.org/pkg/flag/) which makes it easy to setup configuration
from the command line, while combining that with optional environment variables.

I decided to separate out this section into its own file for clarity, as well as
include a sample `config.sh` file for easier project setup.

### Files

This particular project has 3 main file groupings. A `main.go` file where the
program is executed. A `config.go` and `config.sh.example` file where
configuration and setup happens. And a `twitter.go` (and `twitter_test.go`) file
where the Twitter specific API work is done.

These files are broken up according to function, with the extensibility kept in
mind. If I were to add another API, for example Yelp, I would add an additional
file (likely called `yelp.go`) to handle all of the Yelp work. Configuration for
Yelp would be added to the `config.go` file, and the final calculation would be
added to `main.go`, using a similar structure to the way Twitter was calculated
and weighted for the final score.

# TODO (e.g. further work given more time)

With more time, additional testing would be done. This includes both more test
cases, as well as testing certain functions which were only manually tested for
the sake of time.

With more time, additional APIs would be added, and weighted into the final
calculation of the P score.

Better documentation, particularly of the functions and what they are doing,
would also be included.
