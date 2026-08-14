# Changelog

## [4.0.0](https://github.com/Liphium/magic/compare/v3.2.0...v4.0.0) (2026-08-14)


### ⚠ BREAKING CHANGES

* **runner:** Reuse containers when possible

### Features

* Add SeaweedFS driver to Magic ([7c69d54](https://github.com/Liphium/magic/commit/7c69d547a22608181419a84e48dced49f21463a6))
* **examples:** Add tests and scripts to file sharing example ([4c9aafe](https://github.com/Liphium/magic/commit/4c9aafee8087dd3917569dc1fd96c2f953a6a9f1))
* Implement base of SeaweedFS driver ([3738df5](https://github.com/Liphium/magic/commit/3738df5fbe345450cada67b2db604e2bcb0413b9))
* Release SeaweedFS driver ([f76fd9d](https://github.com/Liphium/magic/commit/f76fd9dff2f1ac375deb504da1efb6d67c93a96e))
* **runner:** Reuse containers when possible ([2d2cff2](https://github.com/Liphium/magic/commit/2d2cff2e6d434da98a6cc73a8f6df37ac65bdef0))
* Upgrade to v4 + dependency upgrade ([648cba4](https://github.com/Liphium/magic/commit/648cba427910ccdac7f197fb9994587eb07a5e46))

## [3.2.0](https://github.com/Liphium/magic/compare/v3.1.0...v3.2.0) (2026-07-17)


### Features

* redis driver ([2a8ebdb](https://github.com/Liphium/magic/commit/2a8ebdbe7716b7935bfa10ed2020d0c6062cb708))

## [3.1.0](https://github.com/Liphium/magic/compare/v3.0.0...v3.1.0) (2026-07-14)


### Features

* add option to run magic without docker ([39bdd53](https://github.com/Liphium/magic/commit/39bdd53cf54d0a6d65b7992a648845a27bdcd98b))
* Add PostgreSQL 18 driver + driver refactor ([9872fae](https://github.com/Liphium/magic/commit/9872fae22cb60577570da655616cb852be42f7d0))
* Add PostgreSQL 18 driver + upgrade real project example ([5f00260](https://github.com/Liphium/magic/commit/5f00260cbe3253a918e388131829f420015ef927))
* add release please + dependency updates ([627c615](https://github.com/Liphium/magic/commit/627c61551bbc3dfbfc7e0e21bc0f7e644f37458f))
* Basic idea for the new structure ([e5568c3](https://github.com/Liphium/magic/commit/e5568c3621bf72d3ed1419aa645b4bb92d8dba05))
* Better test runner + README for real-project example ([77145ce](https://github.com/Liphium/magic/commit/77145ce9969deb1733d097841934fd240d6a2bc1))
* **database:** Add MAGIC_POSTGRES_IMAGE to set the postgres image ([04e11c1](https://github.com/Liphium/magic/commit/04e11c11a2075e4dcf3a01c32ff99978e02f5d1a))
* **database:** Make the postgres image declarable using MAGIC_POSTGRES_IMAGE ([8c55995](https://github.com/Liphium/magic/commit/8c55995a8f401bf390f1a0d72477515bf90e5c47))
* Detect broken container images before collisions ([204688a](https://github.com/Liphium/magic/commit/204688a532ea5ecf2fd4e25671c2df7bfc4c5cce))
* Draft for new number field ([e541082](https://github.com/Liphium/magic/commit/e54108284b1e3dd62a70d85bc70099e8c8a4515d))
* Finished test runner + real-project example ([4ea2d60](https://github.com/Liphium/magic/commit/4ea2d60680fe2a362194ae05a2ada3f6582c8ad5))
* First attempt at new scripting engine ([2b3c108](https://github.com/Liphium/magic/commit/2b3c1081c53584f3a2d083741e8fb0a078a07484))
* First draft implementation of new main functions ([3a55b26](https://github.com/Liphium/magic/commit/3a55b26dc0b0bf1e7260c3c5d6dbe75b06fd3b13))
* First example + New execution engine ([3eeaa1c](https://github.com/Liphium/magic/commit/3eeaa1c83f99cb0009ef5715e579a82059a17ba4))
* Integrated services properly ([23260a8](https://github.com/Liphium/magic/commit/23260a8c3caee0e49f90fd68381fce763f25306f))
* Magic v2 ([4755dc1](https://github.com/Liphium/magic/commit/4755dc1a78ecd1f6392f44c51974a1e3da67077b))
* Magic v2 ([bb4f14f](https://github.com/Liphium/magic/commit/bb4f14f79c6e75af38e209c991dc50075ec5c0f2))
* Make sure major version upgrades are detected and starting of ([863ca22](https://github.com/Liphium/magic/commit/863ca224e1e8d8f5ce65399f98f18a8b220ed808))
* Migrate some PostgreSQL driver logic ([876aabb](https://github.com/Liphium/magic/commit/876aabbf4e94d0a25e715966e68454a2472cbe25))
* More progress towards a service layer ([1bdde47](https://github.com/Liphium/magic/commit/1bdde47643049a3824c70af165cf468d77016fa1))
* More scripting progress ([39411d4](https://github.com/Liphium/magic/commit/39411d477caa9ab071fbd1544da32fcf703ccb73))
* Parameters and return types for scripts are now less restrictive ([5d146d1](https://github.com/Liphium/magic/commit/5d146d15523d7b32aafd5f8a10e183121f41f66b))
* Parameters and return types for scripts are now less restrictive ([dde7902](https://github.com/Liphium/magic/commit/dde7902b24554007391dd18a51fa8d34f4cbc897))
* Postgres Driver and ServiceDriver interface finished ([ccacb05](https://github.com/Liphium/magic/commit/ccacb05f31897b9f3d6a3c376fc5e62c6333f7f6))
* Real project example and fixes ([956c7ae](https://github.com/Liphium/magic/commit/956c7ae0ffa2cad7f1b71cf2d15433c1cb4c578a))
* real-project example tests + lock files work properly now ([8e10003](https://github.com/Liphium/magic/commit/8e10003e2779dc3e239c76cb91a90ffcad4eaf39))
* Start of real-project example ([634fd57](https://github.com/Liphium/magic/commit/634fd579810ccacced783254b3a790a9dd9e4aaf))
* Start work on database drivers ([6847c48](https://github.com/Liphium/magic/commit/6847c48c66e8386b588eb4c24b956be8adba7a1d))
* Update to latest version of packages ([d17c76f](https://github.com/Liphium/magic/commit/d17c76feb9ed0a0ac4453434bf98c25d763fb01d))
* Update to latest version of packages ([0a6d2c4](https://github.com/Liphium/magic/commit/0a6d2c4abeaa7d097bbc3158a47df5137c2e6575))


### Bug Fixes

* Better message for when Docker is not running ([dd3cba3](https://github.com/Liphium/magic/commit/dd3cba32ee767e5146f33cbeb6375c9764652e21))
* **deploy:** Make sure volumes are properly deleted ([e59585f](https://github.com/Liphium/magic/commit/e59585f32c7ec01025e053173eae91ea2d42aa29))
* **deploy:** Make sure volumes are properly deleted and containers found again ([e4df491](https://github.com/Liphium/magic/commit/e4df4912cc1a8c5af32fb00c294d60914d562ce0))
* Fix invalid line break ([d92c111](https://github.com/Liphium/magic/commit/d92c111e57cfcf427bbed356b34e24cc1e66b065))
* Make sure error messages are properly done ([255b080](https://github.com/Liphium/magic/commit/255b08079f8a86663e5e92f0a5dc9f07d801d933))
* Make sure Magic finds the project directory ([d7c4549](https://github.com/Liphium/magic/commit/d7c45492307722062da3976e27f2c04447028b88))
* Make sure Magic finds the project directory ([19a3a0c](https://github.com/Liphium/magic/commit/19a3a0c743592aa71aea260b780b744ce348a735))
* Make sure port allocation is skipped properly ([bd76642](https://github.com/Liphium/magic/commit/bd76642bd184861da70243e645cf009e05527a44))
* Make sure previous versions are retracted ([6bbe0fe](https://github.com/Liphium/magic/commit/6bbe0fe5edc89a3829b2534644003489fcb647fc))
* Make sure previous versions are retracted ([fa96bef](https://github.com/Liphium/magic/commit/fa96bef9fdcff338f972f8c6f018abb0c8345317))
* Make sure service drivers are persisted properly ([199a72b](https://github.com/Liphium/magic/commit/199a72b2ce87f91193f0fffbf85c23c76059f3dd))
* Make sure version checks are skipped when no version is found ([6a9f8f7](https://github.com/Liphium/magic/commit/6a9f8f7a994e8ea46c91e9b0404f341f7f5ecf84))
* Only use postgres connection in drivers ([1d6abc8](https://github.com/Liphium/magic/commit/1d6abc88df395a1a061beffab0f527db73ce7ce8))
* Properly check container mounts ([50a4f84](https://github.com/Liphium/magic/commit/50a4f8404f22ca4453e12a535e3fc4ece8090d52))
* Properly unlock factory ([e902f9e](https://github.com/Liphium/magic/commit/e902f9e6bd8edcf9614ca410b690ee578aecdebe))
* Remove integration package + Fix real-project and bugs with driver ([38c325f](https://github.com/Liphium/magic/commit/38c325fd0acd5823e76c1510af264ad03afae3b3))
* typo ([269227a](https://github.com/Liphium/magic/commit/269227ad6d87aa3d8cb3aaaae8e21d1b357397ff))
