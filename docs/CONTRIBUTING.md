# Notes for contributors

There are multiple ways to contrubute. Bug reports, documentation, examples, tests, new features

## Base requirements for all code contributions
* Make sure the code builds without warnings with `go build ./...`
* Make sure it is properly formatted with `gofmt -s .`
* Make sure you checked the code with `golint ./...`. You can install golint with `go get -u golang.org/x/lint/golint`

## Bug reports
If you have found a bug then the best you can do is to submit a PR with a test that fails.
If you are not able to do that then open an issue and provide minimal reproducible example of the code
and describe differences between expected behaviour and observed behaviour

## Documentation
If you see how documentation can be improved, feel free to open a PR.

## Examples
If you have an interesting use case that you feel could be a good illustration of package usage, please
submit a PR with implementation in "examples" directory. Create a new directory for your new example.

## Tests
Help increase code coverage or spot problems in the current tests by submitting a PR

## New features
Before submitting a PR with a new feature please open an issue to discuss it.
