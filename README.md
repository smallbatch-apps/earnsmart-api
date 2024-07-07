# earnsmart-api

API for EarnSmart, a financial services investment platform powered by Golang, Tigerbeetle and PostgreSQL.

## About Tigerbeetle

Tigerbeetle is a high performance, ultra-high availability database solution for financial transactions. Tigerbeetle handles all user account management, balances, and bookkeeping for EarnSmart.

## Current feature set

At present features are being aggressively added, and nothing is fully completed, though price listings should be correct. With that said the worker to update prices is not yet implemented, and the API is not yet fully functional.

This list will be updated as features are added.

- [x] logging middleware
- [x] admin middleware
- [x] user middleware
- [x] CORS and heading middleware
- [x] User login
- [x] User logout
- [x] User session
- [ ] User password reset
- [x] Price listings
- [x] Model creation and database structure
- [x] Database migration
- [x] Fund routing
- [ ] Fund service functionality
- [ ] Account generation and searching
- [ ] Transfer/transaction model implementation
- [ ] Worker to update prices
- [ ] Worker to handle payouts on investments
- [ ] Worker to handle matured divestment
- [ ] User registration
- [ ] Initial database seeding
- [ ] User account management
- [ ] User investment and divestment routes
- [ ] RFQ routing and functionality
- [ ] Swap routing and functionality
