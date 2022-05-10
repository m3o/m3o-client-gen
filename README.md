# m3o-client-gen

A collection of tools we use to generate clients and examples of our M3O services and it consists of:

## m3oGen

To run the code generation, from the repo root issue:

```sh
go install ./m3oGen
```

The general flow is that protos get turned to an openapi json and this generator takes both files (JSON and proto) and generates clients and examples for the specified target like go, typescript, Dart and shell clients.

To generate Go clients localy, clone the micro/services repo and run this command from the root.

```sh
m3oGen go
```

similarly, to generate typescript, dart or shell:

```sh
m3oGen ts
```

```sh
m3oGen dart
```

```sh
m3oGen shell
```

## release-note

The purpose of this program is to fetch the latest commit metadata (sha, html_url and message) from the micro/services repo and output a release note that has the following format.

[9ae89b](https://github.com/micro/services/commit/9ae89b537680a949b4442c5f9f393bf845fb7fa4) Wordle API (#417)