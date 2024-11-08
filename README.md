# earnsmart-api

API for EarnSmart, a financial services investment platform powered by Golang, Tigerbeetle and PostgreSQL. It is intended to form a base for a fully featured application that adds additional business rules and requirements onto the existing scaffold. The application as developed is by-design minimal in requirements. Questions such as fee schedules, KYC, and other business rules are not yet implemented, because they are highly dependent on the specific use-case and business model.

Near infinite additional functionality can be added to the base scaffold, and the application is designed to allow for such additions.

## About Tigerbeetle

Tigerbeetle is a high performance, ultra-high availability database solution for financial transactions. Tigerbeetle handles all user account management, balances, and book-keeping for EarnSmart.

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
- [x] Activity logging
- [x] Model creation and database structure
- [x] Database migration
- [x] Fund routing
- [x] Fund service functionality
- [x] Account generation and searching
- [x] Transfer/transaction model implementation
- [x] Worker to update prices
- [ ] Worker to handle payouts on investments
- [ ] Worker to handle matured divestment
- [x] User registration
- [x] Initial database seeding
- [ ] User account management
- [x] User investment and divestment routes
- [x] RFQ routing and functionality
- [x] Swap routing and functionality
- [ ] Role Based Access Control - RBAC
- [ ] Crypto management and wallet generation
