# Changelog

## 0.3.0
- A route read by id next to its list is named `get` + the list's name: `client.API.AffiliatesGetAffiliators` (was `client.API.AffiliatesAffiliators2`), `client.API.AffiliatesGetPrograms` (was `client.API.AffiliatesPrograms2`), `client.API.BlogGetPublic` (was `client.API.BlogPublic`), `client.API.BlogGetPublic2` (was `client.API.BlogPublic2`), `client.API.InboxGetThreads` (was `client.API.InboxThreads2`), `client.API.MarketplaceGetCampaigns` (was `client.API.MarketplaceCampaigns2`), `client.API.MarketplaceGetCreators` (was `client.API.MarketplaceCreators2`). Each old name stays as a deprecated alias.
- Query fields the API refuses a request without are now required: `code` on GET /api/v1/creator-stats/connect/{platform}/callback, `state` on GET /api/v1/creator-stats/connect/{platform}/callback, `key` on GET /api/v1/uploads/avatar, `id` on GET /api/v1/uploads/deliverable, `key` on GET /api/v1/uploads/merchant-asset.

## 0.2.0
- `Client.API`: every feature route of the API, one method each, generated
  from the API spec (`api_generated.go`) and signed like every other call.

## 0.1.0
- Initial release. Module path is github.com/hachimi-cat/ripllo-go.
