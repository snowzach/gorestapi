# ui

## Installation

Install the application dependencies by running:

```sh
yarn
```

## Development

Start the application in development mode by running:

```sh
yarn dev
```

## Production

Build the application in production mode by running:

```sh
yarn build
```

## DataProvider

The included data provider use [FakeREST](https://github.com/marmelab/fakerest) to simulate a backend.
You'll find a `data.json` file in the `src` directory that includes some fake data for testing purposes.

It includes two resources, posts and comments.
Posts have the following properties: `id`, `title` and `content`.
Comments have the following properties: `id`, `post_id` and `content`.

