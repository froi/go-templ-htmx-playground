# Learning playground with Go, Templ, HTMX, and Tailwind

Snall playground repo to learn basic web dev with Go, Templ, and HTMX. This is not meant to be a fully functional app.

## What does it do

Serves a simple web app to rgister a user and login. Uses SQLite for it's DB.

## How to run this project

Start by cloning and installing the repo dependencies. You will need to have Go 1.23.4+ and Node 23.6.0+ installed.

```sh
git clone git@github.com:froi/go-templ-htmx-playground.git && \
cd go-templ-htmx-playground && \
go get ./... && \
npm install
```

After installing all dependencies you will need to run three command in separate terminal windows or TMUX panes

1. Generate and watch Templ templates
    ```sh
    make templ
    ```
2. Watch for CSS changes
    ```sh
    make tailwindcss
    ```
3. Run Air for auto-reloading
    ```sh
    make air
    ```

> [!IMPORTANT]
> Can this be more efficient? Probably, but as a simple POC / playground I'm not really going spend the time to do so.

With that you'll be able to visit localhost on port 8080.
