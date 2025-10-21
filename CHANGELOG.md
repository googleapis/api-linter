# Changelog

## [1.72.0](https://github.com/googleapis/api-linter/compare/v1.71.0...v1.72.0) (2025-10-16)


### Features

* upgrade to Go 1.24 ([#1532](https://github.com/googleapis/api-linter/issues/1532)) ([8448403](https://github.com/googleapis/api-linter/commit/844840381c253125c4d68a165b4000c7f68be411))


### Bug Fixes

* **AIP-140:** restrict `uri` naming suggestions to fields with URI in comments ([#1541](https://github.com/googleapis/api-linter/issues/1541)) ([7dfafc4](https://github.com/googleapis/api-linter/commit/7dfafc4dc1a2a0cf6009f69760f790cec8b59b45))


### Documentation

* update http-uri-resource.md ([#1543](https://github.com/googleapis/api-linter/issues/1543)) ([be4bb0f](https://github.com/googleapis/api-linter/commit/be4bb0f858310dc2ce9ef21201a12867ae2f0a07))

## [1.71.0](https://github.com/googleapis/api-linter/compare/v1.70.2...v1.71.0) (2025-08-26)


### Features

* **AIP-160:** Validate filter field name. ([#1523](https://github.com/googleapis/api-linter/issues/1523)) ([73e4f83](https://github.com/googleapis/api-linter/commit/73e4f83da2399d9067df2c15310864969a33509f))
* **AIP-160:** Validate filter field type ([#1524](https://github.com/googleapis/api-linter/issues/1524)) ([f182a25](https://github.com/googleapis/api-linter/commit/f182a25e6ba6ff5c85a7d3300f6d94f19b36d95b))


### Bug Fixes

* **AIP-133:** skip method sig for non-standard create ([#1521](https://github.com/googleapis/api-linter/issues/1521)) ([e9306c6](https://github.com/googleapis/api-linter/commit/e9306c6f9bd46508fcaefefd3266bdd935c6f2a4))

## [1.70.2](https://github.com/googleapis/api-linter/compare/v1.70.1...v1.70.2) (2025-08-06)


### Bug Fixes

* **AIP-123:** handle errant slash in pattern parsing ([#1517](https://github.com/googleapis/api-linter/issues/1517)) ([40a69bd](https://github.com/googleapis/api-linter/commit/40a69bd75b7eddfa5c16e58aa0c73e441342512e)), refs [#1514](https://github.com/googleapis/api-linter/issues/1514)


### Documentation

* **config:** fix incorrect rule flags ([#1516](https://github.com/googleapis/api-linter/issues/1516)) ([2454606](https://github.com/googleapis/api-linter/commit/2454606c87d4f5c9a12647b13c0f7e8fe945219f))

## [1.70.1](https://github.com/googleapis/api-linter/compare/v1.70.0...v1.70.1) (2025-07-21)


### Bug Fixes

* **AIP-191:** fix php ruby casing strcase regression ([#1510](https://github.com/googleapis/api-linter/issues/1510)) ([6bb2d95](https://github.com/googleapis/api-linter/commit/6bb2d9519051bd75c6a668444eb312e4272ac278))

## [1.70.0](https://github.com/googleapis/api-linter/compare/v1.69.2...v1.70.0) (2025-07-16)


### Features

* **AIP-142:** add relative time segments comment rule ([5fab299](https://github.com/googleapis/api-linter/commit/5fab2997d1f81bf60b55c6ae6e3e3605212c9870))
* **AIP-142:** add time_offset type rule ([#1506](https://github.com/googleapis/api-linter/issues/1506)) ([945cff3](https://github.com/googleapis/api-linter/commit/945cff325fe38d4f8e3a7a620059a1e6b04c5107))
* **integration-tests:** add test harness for cli invocation ([#1493](https://github.com/googleapis/api-linter/issues/1493)) ([35be28f](https://github.com/googleapis/api-linter/commit/35be28f922581bff1f8ad6d2bb9289f6c57e2657))


### Bug Fixes

* **AIP-136:** support response msg lint with resource singular aligned field name ([#1499](https://github.com/googleapis/api-linter/issues/1499)) ([8dec010](https://github.com/googleapis/api-linter/commit/8dec01076c4bbeb0506c39610e325fc25bbda6ca))
* **AIP-158:** clarify pluralized response field finding ([#1498](https://github.com/googleapis/api-linter/issues/1498)) ([f0b7895](https://github.com/googleapis/api-linter/commit/f0b7895da8cd4b437ac0c3a9be2ac442560eeda8))
* **AIP-203:** field-behavior-required ignore imported request types ([#1504](https://github.com/googleapis/api-linter/issues/1504)) ([bb82f00](https://github.com/googleapis/api-linter/commit/bb82f006b37c85cf255ba8bddb1bf34a07993596)), refs [#1503](https://github.com/googleapis/api-linter/issues/1503)
* **cli:** unexpected lint warning when providing multiple files ([#1496](https://github.com/googleapis/api-linter/issues/1496)) ([7ecaa42](https://github.com/googleapis/api-linter/commit/7ecaa4200da7b5cbcbf1c273fc77d524f346ae1c)), refs [#1465](https://github.com/googleapis/api-linter/issues/1465)


### Documentation

* add comments to lint.Config and configuration page ([#1505](https://github.com/googleapis/api-linter/issues/1505)) ([39d0376](https://github.com/googleapis/api-linter/commit/39d0376281fb03f57b24efe0d82cb842e7316615))

## [1.69.2](https://github.com/googleapis/api-linter/compare/v1.69.1...v1.69.2) (2025-02-20)


### Bug Fixes

* **AIP-133/AIP-134:** handle qualified lro response type name comparison ([#1475](https://github.com/googleapis/api-linter/issues/1475)) ([5e8fe24](https://github.com/googleapis/api-linter/commit/5e8fe2442327ab2a3f1833ff77824723d8331e82))
* **cli:** only call ResolveFilenames with ProtoImportPaths if specified ([#1478](https://github.com/googleapis/api-linter/issues/1478)) ([6a0ddc6](https://github.com/googleapis/api-linter/commit/6a0ddc6d441083d60e7a6e1e35cb0f18f562021e))

## [1.69.1](https://github.com/googleapis/api-linter/compare/v1.69.0...v1.69.1) (2025-02-14)


### Bug Fixes

* **cli:** resolve against cwd separately ([#1474](https://github.com/googleapis/api-linter/issues/1474)) ([6206451](https://github.com/googleapis/api-linter/commit/620645169d3e717fb24651b6cffce3a4aa85b837))


### Documentation

* **AIP-215:** fix incorrect heading for `foreign-type-reference` ([#1472](https://github.com/googleapis/api-linter/issues/1472)) ([cd0f8a1](https://github.com/googleapis/api-linter/commit/cd0f8a1accaf504572248c3a3c2a13eec39e0dd2))

## [1.69.0](https://github.com/googleapis/api-linter/compare/v1.68.0...v1.69.0) (2025-02-11)


### Features

* **AIP-215:** augment foreign type checking ([#1467](https://github.com/googleapis/api-linter/issues/1467)) ([6c514fb](https://github.com/googleapis/api-linter/commit/6c514fb12f5839bb3dbf27742ca62af36466c6cf))


### Bug Fixes

* **cli:** exclude cwd from input path resolution ([#1466](https://github.com/googleapis/api-linter/issues/1466)) ([a14ed3d](https://github.com/googleapis/api-linter/commit/a14ed3de28a0d20ee82b9692d5d290d3732e690d))
* **rules:** fix HasParent check and utilify it ([#1468](https://github.com/googleapis/api-linter/issues/1468)) ([6ac3b57](https://github.com/googleapis/api-linter/commit/6ac3b57ca3bccd15806a714239a89e751ac42428))

## [1.68.0](https://github.com/googleapis/api-linter/compare/v1.67.6...v1.68.0) (2025-01-14)


### Features

* **AIP-123:** resource type name matches message ([#1452](https://github.com/googleapis/api-linter/issues/1452)) ([8f3e2ac](https://github.com/googleapis/api-linter/commit/8f3e2ac2ecdba4798a18cd8c3962ded1b0f86b6c))
* **deps:** update to min version to Go 1.22 ([#1457](https://github.com/googleapis/api-linter/issues/1457)) ([f34f16b](https://github.com/googleapis/api-linter/commit/f34f16b865968ce58dc8140d1bb4943ae984b0f4))


### Bug Fixes

* **AIP-126:** Allow prefixed UNKNOWN value ([#1455](https://github.com/googleapis/api-linter/issues/1455)) ([9353565](https://github.com/googleapis/api-linter/commit/93535656cef91b214b6b40edc5a0eac51db89134))


### Documentation

* **AIP-148:** linkify UUID4 format in rule ([#1456](https://github.com/googleapis/api-linter/issues/1456)) ([5eb475f](https://github.com/googleapis/api-linter/commit/5eb475fb13945253c20a29b24952354cb37e3a71))

## [1.67.6](https://github.com/googleapis/api-linter/compare/v1.67.5...v1.67.6) (2024-11-08)


### Bug Fixes

* **AIP-136:** allow google.api.HttpBody as body field ([#1444](https://github.com/googleapis/api-linter/issues/1444)) ([5327865](https://github.com/googleapis/api-linter/commit/5327865093c518404f621c1a7da2f81dd23997bc))

## [1.67.5](https://github.com/googleapis/api-linter/compare/v1.67.4...v1.67.5) (2024-11-05)


### Bug Fixes

* **deps:** upgrade webrick dep ([#1441](https://github.com/googleapis/api-linter/issues/1441)) ([30b0a84](https://github.com/googleapis/api-linter/commit/30b0a84f70c67e55b101abea3c11fc34f8ada01a))

## [1.67.4](https://github.com/googleapis/api-linter/compare/v1.67.3...v1.67.4) (2024-10-22)


### Bug Fixes

* **AIP-133:** use resource field in signature suggestion ([#1439](https://github.com/googleapis/api-linter/issues/1439)) ([20c96b6](https://github.com/googleapis/api-linter/commit/20c96b624560f7646342068ccc45984b114849fa))

## [1.67.3](https://github.com/googleapis/api-linter/compare/v1.67.2...v1.67.3) (2024-09-23)


### Bug Fixes

* **AIP-132:** refine List request response regex ([#1420](https://github.com/googleapis/api-linter/issues/1420)) ([5cc4d27](https://github.com/googleapis/api-linter/commit/5cc4d279c9cfc80545a9d2447b4fe13a8032b2aa))
* **AIP-136:** ignore revision methods ([#1432](https://github.com/googleapis/api-linter/issues/1432)) ([a6ba5f3](https://github.com/googleapis/api-linter/commit/a6ba5f32458cefc42b662019d97199d0a8e86551))
* **AIP-162:** handle LRO in response name rules ([#1431](https://github.com/googleapis/api-linter/issues/1431)) ([8bca1dd](https://github.com/googleapis/api-linter/commit/8bca1dd68ccf00c39a06da9862ac8c599029297e))
* **internal/utils:** refine Get, Create, Update, Delete request regex ([#1422](https://github.com/googleapis/api-linter/issues/1422)) ([487328c](https://github.com/googleapis/api-linter/commit/487328ca8708521562be2921d3c4f2aabaf8a5ae))
* **locations:** make source info access concurrent safe ([#1433](https://github.com/googleapis/api-linter/issues/1433)) ([223aa5b](https://github.com/googleapis/api-linter/commit/223aa5bb6cf4f24193ad6c6037e1b88160474f2e))


### Documentation

* **AIP-132:** fix incorrect field for AIP-217 ([#1423](https://github.com/googleapis/api-linter/issues/1423)) ([6a52a68](https://github.com/googleapis/api-linter/commit/6a52a6845bf8f240a4d9f9a305a26609a2699c17))
* **AIP-134:** change mandated to recommended ([#1426](https://github.com/googleapis/api-linter/issues/1426)) ([338b6a9](https://github.com/googleapis/api-linter/commit/338b6a95906b61ba5a83805bce92919dd53725dc))

## [1.67.2](https://github.com/googleapis/api-linter/compare/v1.67.1...v1.67.2) (2024-08-14)


### Bug Fixes

* **AIP-123:** multiword singleton reduction ([#1417](https://github.com/googleapis/api-linter/issues/1417)) ([7868552](https://github.com/googleapis/api-linter/commit/7868552ff7b27c2fa0f2ff9be3a538763f0450c5))
* **AIP-135:** allow required etag in Delete ([#1414](https://github.com/googleapis/api-linter/issues/1414)) ([aa9587b](https://github.com/googleapis/api-linter/commit/aa9587bc7184a78109f138c809baa00018ea75e9)), refs [#1395](https://github.com/googleapis/api-linter/issues/1395)
* **AIP-235:** allow hosting allow_missing ([#1416](https://github.com/googleapis/api-linter/issues/1416)) ([6bfbcdf](https://github.com/googleapis/api-linter/commit/6bfbcdfa8858ccdba98760d76e2d2a757855cc7b)), refs [#1404](https://github.com/googleapis/api-linter/issues/1404)
* exit rule if response type cannot be resolved ([#1415](https://github.com/googleapis/api-linter/issues/1415)) ([6874dab](https://github.com/googleapis/api-linter/commit/6874dabb4f0d3503f267bb0ab970d62785d12727)), refs [#1399](https://github.com/googleapis/api-linter/issues/1399)


### Documentation

* **AIP-143:** fix rule name used for implementation link ([#1411](https://github.com/googleapis/api-linter/issues/1411)) ([f9cf2eb](https://github.com/googleapis/api-linter/commit/f9cf2ebc9589abfce88317b1e3318a9e1547b41a))

## [1.67.1](https://github.com/googleapis/api-linter/compare/v1.67.0...v1.67.1) (2024-07-31)


### Bug Fixes

* **AIP-123:** skip resource-pattern-plural when there is no plural ([#1409](https://github.com/googleapis/api-linter/issues/1409)) ([93a601d](https://github.com/googleapis/api-linter/commit/93a601d92adbeb0c17fa8724212ee344f934a4aa))

## [1.67.0](https://github.com/googleapis/api-linter/compare/v1.66.2...v1.67.0) (2024-07-26)


### Features

* **AIP-123:** enforce singular as ID segment ([#1403](https://github.com/googleapis/api-linter/issues/1403)) ([088ec19](https://github.com/googleapis/api-linter/commit/088ec196da93a9338b2abf60469cb55ecec5c34d))
* **AIP-123:** resource collection matches plural ([#1408](https://github.com/googleapis/api-linter/issues/1408)) ([9025d3d](https://github.com/googleapis/api-linter/commit/9025d3d674df9f918483decc0f559f168737ee69))
* **AIP-134:** update_mask must be OPTIONAL ([#1394](https://github.com/googleapis/api-linter/issues/1394)) ([9fc0d05](https://github.com/googleapis/api-linter/commit/9fc0d05f3d89905ea7d3b22c9b44fbfa79edac07))
* **AIP-217:** add various rules for return_partial_success support ([#1406](https://github.com/googleapis/api-linter/issues/1406)) ([cf4ba12](https://github.com/googleapis/api-linter/commit/cf4ba1284bbab151275e7dedf291ffea0e488b1c))


### Bug Fixes

* docs typo ([#1401](https://github.com/googleapis/api-linter/issues/1401)) ([4acf04c](https://github.com/googleapis/api-linter/commit/4acf04c6829ffe7f57cf2997c4f9ccc956de9274))

## [1.66.2](https://github.com/googleapis/api-linter/compare/v1.66.1...v1.66.2) (2024-06-04)


### Bug Fixes

* **aip-130:** identify standard and custom methods ([#1396](https://github.com/googleapis/api-linter/issues/1396)) ([be41d72](https://github.com/googleapis/api-linter/commit/be41d72e50032b45f4263779e638fc8ec0ff9013))

## [1.66.1](https://github.com/googleapis/api-linter/compare/v1.66.0...v1.66.1) (2024-05-23)


### Bug Fixes

* **AIP-136:** handle LRO response names ([#1391](https://github.com/googleapis/api-linter/issues/1391)) ([ec79f53](https://github.com/googleapis/api-linter/commit/ec79f5392829fc58a44f577dce55a936ea112988))

## [1.66.0](https://github.com/googleapis/api-linter/compare/v1.65.2...v1.66.0) (2024-05-17)


### Features

* **AIP-136:** request message name ([#1386](https://github.com/googleapis/api-linter/issues/1386)) ([46a6e43](https://github.com/googleapis/api-linter/commit/46a6e43d1d6bb6a6e79131866f16f0b1dfd2e326))
* **aip-136:** response message name ([#1387](https://github.com/googleapis/api-linter/issues/1387)) ([9e43e3f](https://github.com/googleapis/api-linter/commit/9e43e3f1c98dfe716d4c8ede6fc213239425c6ef))

## [1.65.2](https://github.com/googleapis/api-linter/compare/v1.65.1...v1.65.2) (2024-04-22)


### Bug Fixes

* **AIP-203:** resource etag should not have behavior ([#1376](https://github.com/googleapis/api-linter/issues/1376)) ([1c0f838](https://github.com/googleapis/api-linter/commit/1c0f838941e064caa0eda046a5b4f1c2b7fb2788))

## [1.65.1](https://github.com/googleapis/api-linter/compare/v1.65.0...v1.65.1) (2024-04-15)


### Bug Fixes

* **ci:** build binaries with CGO_ENABLED=0 ([#1371](https://github.com/googleapis/api-linter/issues/1371)) ([f776eca](https://github.com/googleapis/api-linter/commit/f776ecaa0fbab579726120383139c13e9f57e479))
* ignore proto3_optional fields in oneof check ([#1370](https://github.com/googleapis/api-linter/issues/1370)) ([0d6e074](https://github.com/googleapis/api-linter/commit/0d6e0742603c377e86e6abbc557d6ff20e142a46)), refs [#1323](https://github.com/googleapis/api-linter/issues/1323)

## [1.65.0](https://github.com/googleapis/api-linter/compare/v1.64.0...v1.65.0) (2024-04-12)


### Features

* **AIP-133:** ignore create methods with invalid LRO response types ([#1366](https://github.com/googleapis/api-linter/issues/1366)) ([22d015a](https://github.com/googleapis/api-linter/commit/22d015afc1067f8895a2603ae859d11d33f06a36))

## [1.64.0](https://github.com/googleapis/api-linter/compare/v1.63.6...v1.64.0) (2024-03-04)


### Features

* remove legacy list revisions rules ([#1348](https://github.com/googleapis/api-linter/issues/1348)) ([2bc5c57](https://github.com/googleapis/api-linter/commit/2bc5c574eb2e33aee2df502bd3b70454dbfae542))


### Documentation

* **AIP-191:** remove ambiguity in java_outer_classname ([#1345](https://github.com/googleapis/api-linter/issues/1345)) ([1d8d76d](https://github.com/googleapis/api-linter/commit/1d8d76d561e5042735c63fa23ec26c7520d11498))

## [1.63.6](https://github.com/googleapis/api-linter/compare/v1.63.5...v1.63.6) (2024-02-20)


### Bug Fixes

* **AIP-4232:** support nested field of required field ([#1339](https://github.com/googleapis/api-linter/issues/1339)) ([e86a159](https://github.com/googleapis/api-linter/commit/e86a159cfecf8e19bff7d869e3c0bca9c140ce08))

## [1.63.5](https://github.com/googleapis/api-linter/compare/v1.63.4...v1.63.5) (2024-02-16)


### Bug Fixes

* **AIP-4232:** correct repeated field check ([#1337](https://github.com/googleapis/api-linter/issues/1337)) ([b383639](https://github.com/googleapis/api-linter/commit/b383639288fb14c776ad644368bf22d62c83e3b7))

## [1.63.4](https://github.com/googleapis/api-linter/compare/v1.63.3...v1.63.4) (2024-02-16)


### Bug Fixes

* **AIP 133-135:** fix false positive in 133-135 ([#1335](https://github.com/googleapis/api-linter/issues/1335)) ([d79af9c](https://github.com/googleapis/api-linter/commit/d79af9cc85959ce2f22d2a12f1d2fbfca0fd2e7b))

## [1.63.3](https://github.com/googleapis/api-linter/compare/v1.63.2...v1.63.3) (2024-01-25)


### Bug Fixes

* **AIP-123:** allow name in nested messages ([#1325](https://github.com/googleapis/api-linter/issues/1325)) ([16316a5](https://github.com/googleapis/api-linter/commit/16316a5405bd967e926a1482f3bd1e85e1c45eed))

## [1.63.2](https://github.com/googleapis/api-linter/compare/v1.63.1...v1.63.2) (2024-01-22)


### Bug Fixes

* tweak cli integration test ([#1320](https://github.com/googleapis/api-linter/issues/1320)) ([931ab2d](https://github.com/googleapis/api-linter/commit/931ab2dec5005d7c4fcc7b656bcd4141c55daeaa))

## [1.63.1](https://github.com/googleapis/api-linter/compare/v1.63.0...v1.63.1) (2024-01-08)


### Documentation

* **AIP-135:** fix title heading ([#1314](https://github.com/googleapis/api-linter/issues/1314)) ([963c7d8](https://github.com/googleapis/api-linter/commit/963c7d8ac016d4feec7e4b4d552dfb08ff421cfe))

## [1.63.0](https://github.com/googleapis/api-linter/compare/v1.62.0...v1.63.0) (2024-01-08)


### Features

* allow `request_id` for standard Get and List ([#1312](https://github.com/googleapis/api-linter/issues/1312)) ([b41be6e](https://github.com/googleapis/api-linter/commit/b41be6ea41dfc2fb230f5b9f5aa5de4e5d276849))

## [1.62.0](https://github.com/googleapis/api-linter/compare/v1.61.0...v1.62.0) (2024-01-02)


### Features

* **AIP-122:** Flag self-link fields in resources ([#1294](https://github.com/googleapis/api-linter/issues/1294)) ([d7228d3](https://github.com/googleapis/api-linter/commit/d7228d329ed90ced353dd6a9022d19570069adab))

## [1.61.0](https://github.com/googleapis/api-linter/compare/v1.60.0...v1.61.0) (2024-01-02)


### Features

* undelete should not be required for declarative-friendly interfaces ([#1304](https://github.com/googleapis/api-linter/issues/1304)) ([b40c90d](https://github.com/googleapis/api-linter/commit/b40c90d9b1a30d50c08e1373dd9c7b468dd651c2))


### Documentation

* update release docs ([#1301](https://github.com/googleapis/api-linter/issues/1301)) ([3cfd411](https://github.com/googleapis/api-linter/commit/3cfd4111355af9ac581a5bc45e860d8592869418))

## [1.60.0](https://github.com/googleapis/api-linter/compare/v1.59.2...v1.60.0) (2023-12-15)


### Features

* **AIP-142:** warn on creation_time ([#1291](https://github.com/googleapis/api-linter/issues/1291)) ([ebf2763](https://github.com/googleapis/api-linter/commit/ebf27633aed7afc0679664fab0b8493110a5462f))
* require golang 1.20 or later ([#1299](https://github.com/googleapis/api-linter/issues/1299)) ([6864876](https://github.com/googleapis/api-linter/commit/6864876c07c8f2adfd3e81bd651edbfdaa621a79))


### Bug Fixes

* **AIP-123:** disallow spaces in pattern ([#1293](https://github.com/googleapis/api-linter/issues/1293)) ([4d6e057](https://github.com/googleapis/api-linter/commit/4d6e0574c9bc8537968cc4f301e5fe2e4b121618))


### Documentation

* update help message in usage section ([#1288](https://github.com/googleapis/api-linter/issues/1288)) ([eb09eb6](https://github.com/googleapis/api-linter/commit/eb09eb6d8e2600431a326b3ab7b332054e5cf10b))

## [1.59.2](https://github.com/googleapis/api-linter/compare/v1.59.1...v1.59.2) (2023-11-13)


### Bug Fixes

* **AIP-133:** lint http collection ID for lookalikes ([#1286](https://github.com/googleapis/api-linter/issues/1286)) ([3397db6](https://github.com/googleapis/api-linter/commit/3397db63db4adab4f80f022bf247019483473644))

## [1.59.1](https://github.com/googleapis/api-linter/compare/v1.59.0...v1.59.1) (2023-11-01)


### Bug Fixes

* **AIP-203:** skip identifier check if  missing name field ([#1282](https://github.com/googleapis/api-linter/issues/1282)) ([2050717](https://github.com/googleapis/api-linter/commit/2050717c5f965a7374956f87b35ee048d1f2f53a))

## [1.59.0](https://github.com/googleapis/api-linter/compare/v1.58.1...v1.59.0) (2023-10-18)


### Features

* **AIP-148:** add uid-format rule ([#1270](https://github.com/googleapis/api-linter/issues/1270)) ([79cd201](https://github.com/googleapis/api-linter/commit/79cd20109925c348e7a898211d5d8a8533f5a0c3))
* **AIP-148:** lint IP Address field format ([#1271](https://github.com/googleapis/api-linter/issues/1271)) ([cb30ca8](https://github.com/googleapis/api-linter/commit/cb30ca877e249f9a0492a5b95742ed8f1a4f092b))
* **AIP-155:** enforce format uuid4 guidance ([#1272](https://github.com/googleapis/api-linter/issues/1272)) ([2f2e34b](https://github.com/googleapis/api-linter/commit/2f2e34b24c8ac967094418654ebbffcecbd2d04d))

## [1.58.1](https://github.com/googleapis/api-linter/compare/v1.58.0...v1.58.1) (2023-10-05)


### Bug Fixes

* **AIP-203:** remove resource name identifier suggestion ([#1263](https://github.com/googleapis/api-linter/issues/1263)) ([687fe7f](https://github.com/googleapis/api-linter/commit/687fe7f087018f140df677c8cb9836da0bf37b93))

## [1.58.0](https://github.com/googleapis/api-linter/compare/v1.57.1...v1.58.0) (2023-10-03)


### Features

* **AIP-202:** add string-only-format rule ([#1261](https://github.com/googleapis/api-linter/issues/1261)) ([b67339e](https://github.com/googleapis/api-linter/commit/b67339ed8cfb49fcafd0f47f4585a636f8da8600))

## [1.57.1](https://github.com/googleapis/api-linter/compare/v1.57.0...v1.57.1) (2023-09-12)


### Bug Fixes

* **AIP-121:** ignore standard method lookalikes ([#1250](https://github.com/googleapis/api-linter/issues/1250)) ([34aa3c8](https://github.com/googleapis/api-linter/commit/34aa3c8ef5bd58879d2982542e5b786abedd971b))
* **AIP-121:** ignore streaming lookalikes ([#1252](https://github.com/googleapis/api-linter/issues/1252)) ([e723668](https://github.com/googleapis/api-linter/commit/e7236687b64c1e465ce03048c8b67fe09db1091b))

## [1.57.0](https://github.com/googleapis/api-linter/compare/v1.56.1...v1.57.0) (2023-09-11)


### Features

* **AIP-121:** lint for mutable reference cycles ([#1238](https://github.com/googleapis/api-linter/issues/1238)) ([a3895eb](https://github.com/googleapis/api-linter/commit/a3895eba02890c72318bad0726b0599f5b37b261)), refs [#1109](https://github.com/googleapis/api-linter/issues/1109)
* **AIP-122:** disallow embedded resource ([#1245](https://github.com/googleapis/api-linter/issues/1245)) ([e3bbdb3](https://github.com/googleapis/api-linter/commit/e3bbdb313f7507139d7ec91995200b9bc543d7ae))
* **AIP-203:** add resource name IDENTIFIER enforcement ([#1241](https://github.com/googleapis/api-linter/issues/1241)) ([7d9daab](https://github.com/googleapis/api-linter/commit/7d9daab01d4da23d90e42fc78673d27086289cf4))
* **AIP-203:** disallow IDENTIFIER on non-name ([#1244](https://github.com/googleapis/api-linter/issues/1244)) ([1022df2](https://github.com/googleapis/api-linter/commit/1022df2e3160df675ca6dfd460cc6e77dfe5e954))


### Documentation

* Add contributing section to README ([#1242](https://github.com/googleapis/api-linter/issues/1242)) ([45f8534](https://github.com/googleapis/api-linter/commit/45f853426f1d737258df2345e14e60ff20e50645))

## [1.56.1](https://github.com/googleapis/api-linter/compare/v1.56.0...v1.56.1) (2023-08-14)


### Bug Fixes

* **AIP-121:** strict check of resource type ([#1235](https://github.com/googleapis/api-linter/issues/1235)) ([3764f3c](https://github.com/googleapis/api-linter/commit/3764f3c9eef1caae3a806c0de175f03e059fcf74))
* **AIP-192:** enforce deprecated comment on all non-file descriptors ([#1233](https://github.com/googleapis/api-linter/issues/1233)) ([efaced9](https://github.com/googleapis/api-linter/commit/efaced966aa2b7712259cae71dfd539b79bb01ed))

## [1.56.0](https://github.com/googleapis/api-linter/compare/v1.55.2...v1.56.0) (2023-08-11)


### Features

* **AIP-121:** enforce list requirement ([#1225](https://github.com/googleapis/api-linter/issues/1225)) ([7ad11d0](https://github.com/googleapis/api-linter/commit/7ad11d0add4228236bae63f8b3e87812428758e8))
* **AIP-156:** allow singleton list ([#1230](https://github.com/googleapis/api-linter/issues/1230)) ([ccb38cf](https://github.com/googleapis/api-linter/commit/ccb38cfed985c24dad9055704103583a2d8ff326))

## [1.55.2](https://github.com/googleapis/api-linter/compare/v1.55.1...v1.55.2) (2023-08-02)


### Bug Fixes

* **ci:** fix release asset job ([#1222](https://github.com/googleapis/api-linter/issues/1222)) ([af6066d](https://github.com/googleapis/api-linter/commit/af6066d83b00bacae8c3f3535729110580bc8f75))

## [1.55.1](https://github.com/googleapis/api-linter/compare/v1.55.0...v1.55.1) (2023-08-02)


### Bug Fixes

* **locations:** add arbitrary field option helper ([#1213](https://github.com/googleapis/api-linter/issues/1213)) ([b8a0992](https://github.com/googleapis/api-linter/commit/b8a09921324769f882d14efb95014cadc81b8644))
* refactor Standard Method message helpers into utils ([#1212](https://github.com/googleapis/api-linter/issues/1212)) ([c6b5d10](https://github.com/googleapis/api-linter/commit/c6b5d10eadf72b71437d5639d3ad17d07af22082))
