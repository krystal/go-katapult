# Changelog

## [0.2.1](https://github.com/krystal/go-katapult/compare/v0.2.0...v0.2.1) (2023-03-21)


### Features

* **file_storage_volume:** add concrete type and constants for state ([56162b4](https://github.com/krystal/go-katapult/commit/56162b42cb8f515a177ad3e88cc1abecd510d91d))
* **file_storage_volumes:** add *TrashObject return value to Delete method ([232b1d4](https://github.com/krystal/go-katapult/commit/232b1d4a817a671c26aa5a065a5bc426fb4fe30c))
* **file_storage_volumes:** add all expected CRUD+List operations ([91acddf](https://github.com/krystal/go-katapult/commit/91acddf742bfbeafd28e02395dde292176b73c66))


### Bug Fixes

* **file_storage_volume:** enabling updating associations to empty list ([7869f80](https://github.com/krystal/go-katapult/commit/7869f802e3349b6451fe4d1edf8010a0a404003a))

## [0.2.0](https://github.com/krystal/go-katapult/compare/v0.1.7...v0.2.0) (2023-03-07)


### ⚠ BREAKING CHANGES

* **core:** Changed return type of SecurityGroupRules.List().
* **core:** Changed return type of LoadBalancerRules.List().
* **deps:** Minimum required Go version changed from 1.15 to 1.18.
* **dns_zones:** DNS Zone related methods have been changed due to upstream deprecation and changes in Katapult's API.
* **errors:** katapult.ErrScopeNotGrantedError is now named katapult.ErrScopeNotGranted

### Features

* **core/SecurityGroupRule:** add SecurityGroup ref field ([467011a](https://github.com/krystal/go-katapult/commit/467011a68124cb0b35ac21ab60b0c78defd08355))
* **dns_zones:** implement new DNS Zone endpoints ([fea6d83](https://github.com/krystal/go-katapult/commit/fea6d83abe18519ff28a7464671c6e65930a52a1))
* **errors:** re-generate error types from bundled API schema ([1363d69](https://github.com/krystal/go-katapult/commit/1363d6989c2b30fe128050be939278ea8e135c5b))


### Bug Fixes

* **buildspec:** add tests to verify XML entity substitution is not supported ([68d7e0d](https://github.com/krystal/go-katapult/commit/68d7e0d98f747a36fea9efa2984c4698ac4de1d7))
* **codegen:** improve error variable naming ([f59cb42](https://github.com/krystal/go-katapult/commit/f59cb425d94da7b6e3682e7ade2fab14450e5808))
* **core:** change return type of LoadBalancerRules.List() method ([7689e27](https://github.com/krystal/go-katapult/commit/7689e278b25222a0e42dbe289fbd9a1ef33f3040))
* **core:** change return type of SecurityGroupRules.List() method ([731e89b](https://github.com/krystal/go-katapult/commit/731e89b2cc6731637c7904b9296eb046a82ba793))


### Miscellaneous Chores

* **deps:** minimum required Go version changed from 1.15 to 1.18 ([a09c267](https://github.com/krystal/go-katapult/commit/a09c2673127ce9e9e71c379e361f90006622789a))

### [0.1.7](https://github.com/krystal/go-katapult/compare/v0.1.6...v0.1.7) (2021-08-19)


### Features

* add request option support ([03da1ac](https://github.com/krystal/go-katapult/commit/03da1ace0f7b9f292161bd9a076f6905301dc218))
* add request options support ([b022bf2](https://github.com/krystal/go-katapult/commit/b022bf2b68d4e0f78e646ff1df482cf9b483d5ee))

### [0.1.6](https://github.com/krystal/go-katapult/compare/v0.1.5...v0.1.6) (2021-08-18)


### Features

* **errors:** update generated errors based on latest Katapult API schema ([462b2c4](https://github.com/krystal/go-katapult/commit/462b2c4212af5e7c7a11c8fa35f4b9594e7d583e))
* **ssh_keys:** add support for managing organization SSH keys ([#119](https://github.com/krystal/go-katapult/issues/119)) ([782b3dd](https://github.com/krystal/go-katapult/commit/782b3dd6c06ac1f0bfb51e486eb0a7ab306d0ee2))
* **tags:** Add tag management support ([#118](https://github.com/krystal/go-katapult/issues/118)) ([0a78954](https://github.com/krystal/go-katapult/commit/0a78954f5f5eaeed6b7601e5b7c3755b1779670b))


### Bug Fixes

* **codegen:** fix issue caused by a recent change to Katapult's API Schema ([7120dd7](https://github.com/krystal/go-katapult/commit/7120dd7533c6f9dff283de161bfacea6a416cae0))

### [0.1.5](https://github.com/krystal/go-katapult/compare/v0.1.4...v0.1.5) (2021-06-17)


### Features

* **apischema:** add package to parse Katapult API JSON schema ([932de00](https://github.com/krystal/go-katapult/commit/932de00ad64c3d7c633a3a1b912974885b5207fd))
* **errors:** add custom code generator tool for generating error structs ([b481bdf](https://github.com/krystal/go-katapult/commit/b481bdf5c3b7a4cb857e5c928f572873708547ec))
* **errors:** generate error structs from Katapult API schema ([903a4b8](https://github.com/krystal/go-katapult/commit/903a4b851aa990caa20fbb3fd2bf516e1d5b171d))
* **security_groups:** add support for security group rules ([#112](https://github.com/krystal/go-katapult/issues/112)) ([0580d7a](https://github.com/krystal/go-katapult/commit/0580d7a9491ea823f9c1ef5db1567cc003359c69))
* add support for katapult security groups ([#103](https://github.com/krystal/go-katapult/issues/103)) ([e5b1fb4](https://github.com/krystal/go-katapult/commit/e5b1fb4da06c3d89e4d9d228ffb5cbc0d09d2daf))

### [0.1.4](https://github.com/krystal/go-katapult/compare/v0.1.3...v0.1.4) (2021-05-31)


### Bug Fixes

* **load_balancer:** enable removing all certs from a rule ([788df99](https://github.com/krystal/go-katapult/commit/788df995b96f88b1c7be46bd781c82af6bac8901))
* **load_balancer:** use pointer to arguments struct for the sake of consistency ([54b0009](https://github.com/krystal/go-katapult/commit/54b000943c0296a760545a903c78e88187c866d2))

### [0.1.3](https://github.com/krystal/go-katapult/compare/v0.1.2...v0.1.3) (2021-05-28)


### Bug Fixes

* **load_balancer:** use CertificateRef when creating/updating rules ([a05dbe8](https://github.com/krystal/go-katapult/commit/a05dbe8c0ac09410176eed532e7ecab26c759a66))

### [0.1.2](https://github.com/krystal/go-katapult/compare/v0.1.1...v0.1.2) (2021-05-27)


### Features

* **client:** add WithHTTPClient option ([dfa0e09](https://github.com/krystal/go-katapult/commit/dfa0e0990d1cf2f356c98dec5ea20f9279dc2909))

### [0.1.1](https://github.com/krystal/go-katapult/compare/v0.1.0...v0.1.1) (2021-05-27)


### Features

* **data_center:** add DefaultNetwork method ([5cc2716](https://github.com/krystal/go-katapult/commit/5cc2716b063e05ab7920deeb5f2919e9cf6ae630))
