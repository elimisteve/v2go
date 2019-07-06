## v2go

v2go translates [V](https://vlang.io/) code to Go.

**STATUS:** Active Development | Alpha Quality


### Getting Started

Fetch, build, and install it:

    go get github.com/elimisteve/v2go
    go get golang.org/x/tools/cmd/goimports

Try it out:

    cd $(go env GOPATH)/src/github.com/elimisteve/v2go
    cat ./test_v_files/hello_world.v
    v2go run ./test_v_files/hello_world.v

Or, for a more interesting example:

    cd $(go env GOPATH)/src/github.com/elimisteve/v2go
    cat ./test_v_files/links_scraper2.v
    v2go run ./test_v_files/links_scraper2.v

`v2go run ./path/to/mycode.v` translates the given `.v` file(s) into
Go, then tells Go to build and run them -- all with one command! :tada:

To just translate your `.v` files _without_ also running the
translated Go files afterward:

    v2go translate ./test_v_files/*.v

Now you can run them as one normally would:

    go run ./test_v_files/hello_world.go
    go run ./test_v_files/hello_world_interpolated.go
    go run ./test_v_files/links_scraper2.go


For a bit more info, run

    v2go help

Now go forth and write useful, world-changing software!


### Why v2go?

As someone who thinks a lot about the future of computer programming,
and as someone who's been coding for 10 years,
**I consider V to be by far the best-designed programming language of all time**.

Why?  Because V gives us the simplicity, compilation speed, and
concurrency model of Go; the safety and runtime speed of Rust; and the
syntactic convenience of a scripting language like Python.

The only problem is that V is very new and does not have much code
that V programmers can use; V was open-sourced just last month, and
[the C-to-V translator c2v](https://github.com/vlang/c2v) has not
quite been released yet.

Thus, this short-term solution of v2go, which translates V code into
Go _so that all existing Go code can be used **today**_ while enjoying
the benefits of V.

**v2go enables us to use the brilliant simplicity of V's syntax while importing arbitrary Go code _today_.**

v2go does not yet support translating imported V code (e.g., from its
standand library), so I recommend simply importing Go code directly;
v2go will not rewrite these import paths.


### Near-future Work

- [ ] Consider first compiling the given `.v` files with the V compiler so as to benefit from the additional safely checks that V has and Go doesn't (e.g., enforcing immutability and ensuring that all errors are checked)
- [ ] Make v2go much more robust by making it a legit compiler with a lexer, parser, AST walking, etc


### Less-near-future Work

Probably not to be done by me:

- [ ] Write **go2v**, a Go-to-V translator/transpiler/compiler

Why?  So that the V standard library and more of its runtime can be
auto-generated rather than _manually_ written (how old-fashioned that
would be!)


### Hacking on v2go

Get started in the same way as stated above:

    go get github.com/elimisteve/v2go
    go get golang.org/x/tools/cmd/goimports
    cd $(go env GOPATH)/src/github.com/elimisteve/v2go

Make sure all tests pass:

    go test ./translate

Hack on. :fist: Consider starting by [creating an issue](https://github.com/elimisteve/v2go/issues)
with a proposed new feature, or by submitting a pull request.


### More About V

See <https://vlang.io/>, especially <https://vlang.io/docs>.  For a
comparison of V to other languages, see <https://vlang.io/compare>.
