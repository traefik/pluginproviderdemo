# Provider Plugin Demo

Example:

```yaml
entryPoints:
  web:
    address: :80

log:
  level: DEBUG

pilot:
  token: xxxxx

experimental:
  devPlugin:
    goPath: /plugins/go
    moduleName: github.com/traefik/pluginproviderdemo

providers:
  plugin:
    dev:
      pollInterval: 2s
```
