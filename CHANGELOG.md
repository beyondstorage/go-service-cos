# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v2.1.0] - 2021-05-24

### Added

- *: Implement GSP-87 Feature Gates (#27)
- storage: Add CreateDir (#28)

### Changed

- *: Implement GSP-97, GSP-109 and GSP-117 (#32)

### Upgraded

- build(deps): bump github.com/tencentyun/cos-go-sdk-v5 to 0.7.27 (#34)

## [v2.0.0] - 2021-05-24

### Added

- storage: Implement SSE support (#17)
- services: implement GSP-47 & GSP-51 (#21)
- storage: Implement multipart support (#23)

### Changed

- storage: Idempotent storager delete operation (#20)
- *: Implement GSP-73 Organization rename (#24)

## [v1.1.0] - 2021-04-24

### Added

- pair: Implement default pair support for service (#4)
- storage: Implement Create API (#13)
- *: Add UnimplementedStub (#15)
- tests: Introduce STORAGE_COS_INTEGRATION_TEST (#16)
- tests: Add docs for how to run tests 
- storage: Implement GSP-40 (#18)

### Changed

- docs: Migrate zulip to matrix
- build: Fix build scripts
- ci: Only run Integration Test while push to master

### Upgraded

- build(deps): bump github.com/tencentyun/cos-go-sdk-v5 from 0.7.19 to 0.7.24

## v1.0.0 - 2021-02-08

### Added

- Implement cos services.

[v2.1.0]: https://github.com/beyondstorage/go-service-cos/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-cos/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/beyondstorage/go-service-cos/compare/v1.0.0...v1.1.0
