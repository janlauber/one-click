# Filters

## Timestamp

### Suffix

pattern: \d{10}-slim$
policy: timestamp

### Prefix

pattern: slim-\d{10}$
policy: timestamp

## Semver

### Plain

pattern: .*
policy: semver:1.2.x

### Prefix

pattern: test-.*
policy: semver:1.2.x

### Suffix

pattern: .*-slim
policy: semver:1.2.x

### Between

pattern: .*
policy: semver:>=1.24.x, <1.25.x

### Only numbers

```json
{
    "pattern": "^\\d+\\.\\d+\\.\\d+$",
    "policy": "semver:>=1.24.x, <1.25.x",
    "registry": "https://registry.hub.docker.com/v2/repositories",
    "repository": "library/nginx"
}
```
