<h1 align="center">
  <br>
  <a href="http://www.amitmerchant.com/electron-markdownify"><img src="https://img001.prntscr.com/file/img001/MvCSWOg2SWSs2uPk8E7C3w.png" alt="Markdownify" width="400"></a>
  <br>
  Autocomplete
  <br>
</h1>

<h4 align="center">Toy project made with <a href="https://go.dev/" target="_blank">Golang</a>.</h4>

<p align="center">
  <a>
    <img src="https://img.shields.io/github/go-mod/go-version/matheusgb/autocomplete" alt="go version">
  </a>
</p>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#documentation">Documentation</a>
</p>

## Key Features

* Saves and returns data from Elasticsearch.
* List of most saved words in Elasticsearch.
* Uses WebSocket to list autocomplete suggestions from Elasticsearch.
* List of most frequently saved words in Elasticsearch using WebSocket.

## How To Use

Initialize Elasticsearch in docker with:

```
docker-compose up
```

If you have [air](https://github.com/air-verse/air) installed, you can run the backend with the command in root folder:

```
air
```

Also, you can run with:

```
go run main.go
```

You can run frontend with [LiveServer vscode extension](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer).

## Documentation

If you want Elasticsearch with some data, use the `GET - /populate` route.
