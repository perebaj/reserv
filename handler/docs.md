As we don't have a swagger, I'm providing some examples of how to use the API.

# Properties

```bash
    curl -i -X GET "https://reserv-production.up.railway.app/properties?id=af0c35de-e513-4b67-b113-e68be16fba24"
    curl -i -X POST "https://reserv-production.up.railway.app/properties" -d '{"title": "My Property", "description": "My Description", "price_per_night_cents": 1000, "currency": "USD", "host_id": "af0c35de-e513-4b67-b113-e68be16fba24"}'
    curl -i -X PUT "https://reserv-production.up.railway.app/properties?id=af0c35de-e513-4b67-b113-e68be16fba24" -d '{"title": "My Property", "description": "My Description", "price_per_night_cents": 1000, "currency": "USD"}'
    curl -i -X DELETE "https://reserv-production.up.railway.app/properties?id=af0c35de-e513-4b67-b113-e68be16fba24"
```
