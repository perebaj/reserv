As we don't have a swagger, I'm providing some examples of how to use the API.

# Properties

```bash
    curl -i -X GET "https://reserv-production.up.railway.app/properties?id=b18c5092-13fd-4aeb-bd38-122744d3f865"
    curl -i -X POST "https://reserv-production.up.railway.app/properties" -d '{"title": "My Property", "description": "My Description", "price_per_night_cents": 1000, "currency": "USD", "host_id": "b18c5092-13fd-4aeb-bd38-122744d3f865"}'
    curl -i -X PUT "https://reserv-production.up.railway.app/properties?id=b18c5092-13fd-4aeb-bd38-122744d3f865" -d '{"title": "My Property", "description": "My Description", "price_per_night_cents": 1000, "currency": "USD"}'
    curl -i -X DELETE "https://reserv-production.up.railway.app/properties?id=b18c5092-13fd-4aeb-bd38-122744d3f865"
```
