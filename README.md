api-score
=========

api-score is a project that creates a "p value" for a business, where the p
value represents the likely hood that a business will be able to repay it's
loans.

# Usage

To run api-score, you'll need to have [Golang](https://golang.org/dl/)
installed.

Setup your twitter api config by editing the config.sh.example file (and moving
it to config.sh if you'd like), and then running `source config.sh`.

Compile the program using `go build`, and run it by specifying the required
flags, e.g.

```
./api-score --twuser stripe --business Stripe --owner Stripe
```

# Design Decisions

TODO
