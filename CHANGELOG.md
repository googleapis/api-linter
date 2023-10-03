# Changelog

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
