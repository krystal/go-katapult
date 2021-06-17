# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

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
