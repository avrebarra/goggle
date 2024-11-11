# 🏊‍♂️ goggle

Sample golang service.

## Setting Up

**Dev Dependency and requirements:**

- Go 1.19 or later
- Ruby 2.x (to run scripts, the one included in mac is enough)
- NodeJS & npm (to develop ui)
- [mitranim/gow](https://github.com/mitranim/gow) (to run watch mode)
- [matryer/moq](https://github.com/matryer/moq) (for generationg interface mocks)

**How to setup the project:**

1. Clone the repo.
1. Reinit git [like this](https://stackoverflow.com/a/66689268/1914707).
1. Replace all `goggle` string into your desired service name. _If using VSCode you can press Cmd+Shift+F and replace all `goggle` to `nüservice`._
1. Read and remove components and pkgs and lines thats not used.
1. Update this README.md as you see fit.
1. Git commit, `day one start making impacts~ 🍃`.

**Running the project:**

```sh
$ make build # to build
$ make watch # run go app in development mode
```

See Makefile for available commands and script.

## Philosophies and Rules

- Keep things flat. _Resist the temptation to create subdirectory. Usually most of the reason to create directory is categorization, but categorization actually can be solved with just good namings._
- If it's hard to understand, it is a valid concern! Start by the assumptions that the code is not good enough.
- Keep things simple and easy. _Don't worry, hardships will come from every direction sooner or later, no need to set it up earlier._
- Try to not make people think when reading a code.
- Fast understandability over fast performance. _Golang is fast enough. And lately clouds are cheap enough to let us spin multiple instances for a service. So yeah, push for common sense but not too far to sacrifice simplicity._
- Code economy: If it's easy, it's copypasteable. If it actually copypasteable, then it's perfect.
