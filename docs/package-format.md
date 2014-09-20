# Package Format

pm packages are distrubted as gzip compressed tar files with the following contents:

- metadata.json
- executable binary

## Naming Convention

```
${name}-${tag}-${platform}-${architecture}.tar.gz
```

### Example

```
hello-1.0.0-linux-amd64.tar.gz
```

## Binary File

```
hello
```

## Metadata File

metadata.json

```
{
  "name": "hello",
  "maintainer": "Kelsey Hightower <kelsey.hightower@gmail.com>",
  "sourceUrl": "https://github.com/kelseyhightower/hello",
  "description": "The hello app.",
  "tag": "1.0.0",
  "platform": "linux",
  "architecture": "amd64"
}
```

### Digest

hello-1.0.0-linux-amd64.sha256

```
9c1daf153bfb5534dd89b919ab42bbdb5afdf446c235c909b7b2366982c066c6 hello-1.0.0-linux-amd64.tar.gz
```

