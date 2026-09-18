# Changelog

## [1.4.3](https://github.com/Wassergeist/parse-dmarc/compare/v1.7.0...v1.4.3) (2026-09-18)


### Features

* add CI and goreleaser with optimized distroless docker ([66c6828](https://github.com/Wassergeist/parse-dmarc/commit/66c682807d5ed12328349e2713c9d814d6337e78))
* add DigitalOcean Droplet and Dokploy deployment options ([#67](https://github.com/Wassergeist/parse-dmarc/issues/67)) ([21d494f](https://github.com/Wassergeist/parse-dmarc/commit/21d494f60d7b7450d29cfd9a50fa25a914b100e8))
* add MCP server integration to project ([#46](https://github.com/Wassergeist/parse-dmarc/issues/46)) ([807b8d6](https://github.com/Wassergeist/parse-dmarc/commit/807b8d677c81aef6ac39605bffb7125680841d4a))
* add production-ready Grafana dashboard ([6f857fe](https://github.com/Wassergeist/parse-dmarc/commit/6f857fe660144093b8766f529c15b17a21edc3b2))
* add production-ready Prometheus metrics ([#39](https://github.com/Wassergeist/parse-dmarc/issues/39)) ([17f6968](https://github.com/Wassergeist/parse-dmarc/commit/17f69681b23eee92ff1bc908fde7da2de11b8a52))
* add sort and refresh to UI ([9d6853e](https://github.com/Wassergeist/parse-dmarc/commit/9d6853e95319529c40958ae60d95d05d8cb5a675))
* allow changing api endpoint from the settings modal ([#76](https://github.com/Wassergeist/parse-dmarc/issues/76)) ([e6cb48b](https://github.com/Wassergeist/parse-dmarc/commit/e6cb48be538e5c3e997d86dda9250068a2fb377e))
* allow customizing SEEN & move-folder behavior after processing ([80346df](https://github.com/Wassergeist/parse-dmarc/commit/80346dfddc56844dd59c9474fb788ca1d0d9312c)), closes [#118](https://github.com/Wassergeist/parse-dmarc/issues/118)
* **CI:** add man page and shell completion to brew ([0e45ed2](https://github.com/Wassergeist/parse-dmarc/commit/0e45ed2a10d26a06fadc6a2fdfe3b81ad5a01197))
* **CI:** add mcp registry publishing to goreleaser ([a0d6f88](https://github.com/Wassergeist/parse-dmarc/commit/a0d6f88eb54df004c3b35ef7e256bb8d7be36cc9))
* **CI:** add prettier job ([39af52d](https://github.com/Wassergeist/parse-dmarc/commit/39af52d186fcf75dba2f8ac05ef313808f1e01b4))
* **CI:** build docker via goreleaser ([8555265](https://github.com/Wassergeist/parse-dmarc/commit/855526505aa48c9ab1cdc946f10370a875ee47d8))
* **CI:** install pandoc and use official goreleaser action ([0d0a307](https://github.com/Wassergeist/parse-dmarc/commit/0d0a3077395a9879a9d036166b4afbab43ee164e))
* clean up the cli and optimize the docker ([db23ba4](https://github.com/Wassergeist/parse-dmarc/commit/db23ba4ff2fd80ad1c016daa73a57519998dd271))
* **docs:** add 1-click deployment at the top of README ([#41](https://github.com/Wassergeist/parse-dmarc/issues/41)) ([50d2ff2](https://github.com/Wassergeist/parse-dmarc/commit/50d2ff2c54a506801e802daa6b5d375b79d1281d))
* **docs:** add nerdy badges to README ([#23](https://github.com/Wassergeist/parse-dmarc/issues/23)) ([5026a45](https://github.com/Wassergeist/parse-dmarc/commit/5026a45da22a9fa0a3c831945da4afaa05c2dd7c))
* **docs:** add northflank deployment button ([81215a5](https://github.com/Wassergeist/parse-dmarc/commit/81215a5b5e8cd87aa5f601f22fb460b4929d8d15))
* **docs:** add the self-hosted options to deployments ([42ab966](https://github.com/Wassergeist/parse-dmarc/commit/42ab966218b9e875ed438a9cb4a4b4ac3adb9b33))
* **docs:** add zeabur deployment template ([4abaa36](https://github.com/Wassergeist/parse-dmarc/commit/4abaa36d9b077f71a3e6da841af056dffbe046ba))
* **docs:** update readme to the latest changes ([c20e0cf](https://github.com/Wassergeist/parse-dmarc/commit/c20e0cf837b91b05e61153e4e55e7315fb09d396))
* ensure db path exists and add demo screenshot ([42d6d14](https://github.com/Wassergeist/parse-dmarc/commit/42d6d14fd96f2818830c73ad4e23b8f97719c497))
* **frontend:** add DMARC DNS record generator ([34491ba](https://github.com/Wassergeist/parse-dmarc/commit/34491ba244f9af48a25e3a0d6b09979661678c0d))
* **frontend:** implement dark mode with theme toggle ([#35](https://github.com/Wassergeist/parse-dmarc/issues/35)) ([569fd36](https://github.com/Wassergeist/parse-dmarc/commit/569fd3649931f94fa4f39845d0f141afacd6fbc0))
* **imap:** support internal CA, skip-verify and STARTTLS ([#194](https://github.com/Wassergeist/parse-dmarc/issues/194)) ([98abb52](https://github.com/Wassergeist/parse-dmarc/commit/98abb526d856a2c560fabc090a383a206fcec2cc)), closes [#193](https://github.com/Wassergeist/parse-dmarc/issues/193)
* Implement DMARC report parser with Vue.js dashboard ([#1](https://github.com/Wassergeist/parse-dmarc/issues/1)) ([23b8ac0](https://github.com/Wassergeist/parse-dmarc/commit/23b8ac0cc63a81aad7dfa2385ed24f6d74915775))
* log attachments rejected as non-DMARC ([#208](https://github.com/Wassergeist/parse-dmarc/issues/208)) ([31b6a23](https://github.com/Wassergeist/parse-dmarc/commit/31b6a230fe56c749296dcac7b487675651b852a6))
* **mcp:** add OAuth2 authentication for MCP HTTP server ([361f078](https://github.com/Wassergeist/parse-dmarc/commit/361f0782e2f20b7160873a085e600adc68988432))
* **UI:** revamp the dashboard for actionable insights ([#68](https://github.com/Wassergeist/parse-dmarc/issues/68)) ([dbe073d](https://github.com/Wassergeist/parse-dmarc/commit/dbe073d35a7f349770e9fab6df4923e221c350f2))
* update the UI footer with OS friendly text ([5008fdf](https://github.com/Wassergeist/parse-dmarc/commit/5008fdf5e59e5b4af1fc12ccafc644d27ccdbdfc))
* update UI footer & readme and optimize vite ([9578d05](https://github.com/Wassergeist/parse-dmarc/commit/9578d0530adab10334ad1a8a20c563abdc589854))


### Bug Fixes

* accept DMARC reports sent as inline MIME parts ([#207](https://github.com/Wassergeist/parse-dmarc/issues/207)) ([20b2be0](https://github.com/Wassergeist/parse-dmarc/commit/20b2be03d8a43197eabde63eff452e276eaf2b63))
* accept fetch-interval without prefix ([06fd520](https://github.com/Wassergeist/parse-dmarc/commit/06fd520d1a816795d8e35f6a7c9b1b5eed5c096c)), closes [#162](https://github.com/Wassergeist/parse-dmarc/issues/162)
* address unhandled charset ([#87](https://github.com/Wassergeist/parse-dmarc/issues/87)) ([99c245a](https://github.com/Wassergeist/parse-dmarc/commit/99c245a670ce1cb44f5d08b24b5d0c99c502ec69))
* **CI:** add current dir to docker build context ([254cc6a](https://github.com/Wassergeist/parse-dmarc/commit/254cc6a6b0b88afc67af8a1cdc4a41d1b71ff1a6))
* **CI:** add the brew token ([4081621](https://github.com/Wassergeist/parse-dmarc/commit/4081621ae9d1771185712512a912aa350a2e0305))
* **CI:** build on major version as well ([fdff51d](https://github.com/Wassergeist/parse-dmarc/commit/fdff51d2c77d55e890e5b6d9c12ce316ddec2542))
* **CI:** change mcp auth type to gh oidc ([1132b50](https://github.com/Wassergeist/parse-dmarc/commit/1132b5038c16fda2c83fe90a83f9f6e802c06af4))
* **CI:** disable attestation digest for docker build ([95d6665](https://github.com/Wassergeist/parse-dmarc/commit/95d666543f59b3a890ea8921275f840eebb32f8e))
* **CI:** disable publishing to mcp registry ([c2e3969](https://github.com/Wassergeist/parse-dmarc/commit/c2e3969e90554479921437be4df0d011b4214a89))
* **CI:** disable sbom on docker ([36e5cc4](https://github.com/Wassergeist/parse-dmarc/commit/36e5cc4d829efa633fc75879e0a132d076866318))
* **CI:** disable sbom via config ([6b0f3be](https://github.com/Wassergeist/parse-dmarc/commit/6b0f3be7d79cafe1eeb0930d2f30bf02f53c56af))
* **CI:** ensure a fake dist exists before golangci-lint ([4e01983](https://github.com/Wassergeist/parse-dmarc/commit/4e01983f13a16ab0a8ea7cb880e43b39c285ba1f))
* **CI:** make linter happy ([cd04dfc](https://github.com/Wassergeist/parse-dmarc/commit/cd04dfc8beae30b0405d38834680011cd7d14fe5))
* **CI:** move bundled dist to server directory ([59dabf1](https://github.com/Wassergeist/parse-dmarc/commit/59dabf1becfc22fddda4bfdab6dd22188396b4c3))
* **CI:** name docker repo explicitly ([0ebbf87](https://github.com/Wassergeist/parse-dmarc/commit/0ebbf87a6aad51ad735434982734105d79a97893))
* **CI:** only pin main with latest ([f120b2e](https://github.com/Wassergeist/parse-dmarc/commit/f120b2e7b7985adbaa815b6104f8aaee44e8e9ce))
* **CI:** perform keyless docker sign ([1f322d3](https://github.com/Wassergeist/parse-dmarc/commit/1f322d33999d3586880d82c5ce3570b60157cd21))
* **CI:** publish the image to the new repo on main ([f89a9ec](https://github.com/Wassergeist/parse-dmarc/commit/f89a9ece752e3de0bf5c42f3583f41209cf46d55))
* **CI:** publish to casks for a change ([0e8346b](https://github.com/Wassergeist/parse-dmarc/commit/0e8346b3e1c3dcc557445370f59ea7c93588c1fd))
* **CI:** publish to ghcr.io as well ([9c0bf28](https://github.com/Wassergeist/parse-dmarc/commit/9c0bf284a08c50487bb6e8fe475b0248a836cedc))
* **CI:** publish to new and old docker hub repo ([e7eda52](https://github.com/Wassergeist/parse-dmarc/commit/e7eda52d98e7d3d7d127f199f0f06376007e5cb4))
* **CI:** remove annotations from docker build ([b5d0cd6](https://github.com/Wassergeist/parse-dmarc/commit/b5d0cd650ba0dd780ea95ec35c82940620aad6bb))
* **CI:** remove changelog from release note ([127cb80](https://github.com/Wassergeist/parse-dmarc/commit/127cb800f3f285b2d5899561fa390e2894bd52de))
* **CI:** remove docker repo from old ghcr.io ([d3f24c0](https://github.com/Wassergeist/parse-dmarc/commit/d3f24c00961876ac7bd8a976a16ada168a59be28))
* **CI:** remove incomplete mcp docker from goreleaser ([356bb7b](https://github.com/Wassergeist/parse-dmarc/commit/356bb7b8ead6c3b783811178e49e10145a876d29))
* **CI:** remove mcp entry from goreleaser altogether ([1d769a4](https://github.com/Wassergeist/parse-dmarc/commit/1d769a45d05ea5afc5340e87fb60cec3fe3514ce))
* **CI:** reverse the conditional for build-dev job ([ce8b417](https://github.com/Wassergeist/parse-dmarc/commit/ce8b417c65c62563886787df3871200679c8de66))
* **CI:** reverse the digest conditional for provenance ([9c0e165](https://github.com/Wassergeist/parse-dmarc/commit/9c0e165ec1cfeb58c7c535ec02bd0c81e60ba501))
* **CI:** set up buildx action for multi platform build ([a36e41c](https://github.com/Wassergeist/parse-dmarc/commit/a36e41c02e6af54297275042e269c66b85545a78))
* **CI:** update dockerignore after main.go change ([1241b2f](https://github.com/Wassergeist/parse-dmarc/commit/1241b2f22b8138348b7bf087ba00659e0b1f6fe8))
* **CI:** update goreleaser after moving FE to root ([357d37e](https://github.com/Wassergeist/parse-dmarc/commit/357d37ed69d9cd7400c05c2f7ca2667072f4d46d))
* **CI:** update mcp registry version identifier ([a135890](https://github.com/Wassergeist/parse-dmarc/commit/a135890747ef5fc07e43a7c9ed16e4dc21ac592b))
* **CI:** update the digest name for attestation ([eb692a8](https://github.com/Wassergeist/parse-dmarc/commit/eb692a81fb08205990e639ab6c54c8998c705cd0))
* **CI:** update version of the binary ([0a8f023](https://github.com/Wassergeist/parse-dmarc/commit/0a8f0230b30229349db270aa065e3332da63990a))
* **CI:** use go version file when setting up go ([40fbcb1](https://github.com/Wassergeist/parse-dmarc/commit/40fbcb1dfbe962c7d48ef2c821ee32d997e7c78c))
* **CI:** use oxc the default vite minifier ([2cc92e1](https://github.com/Wassergeist/parse-dmarc/commit/2cc92e1d85734bc1cc34b91b209eb66b13336b5d))
* **dashboard:** judge health over delivered mail, in the backend ([#191](https://github.com/Wassergeist/parse-dmarc/issues/191)) ([53d8944](https://github.com/Wassergeist/parse-dmarc/commit/53d8944fdc95ed4199e3ad47505e461e74851d5a))
* **dashboard:** show null reports as EMPTY instead of FAIL ([#187](https://github.com/Wassergeist/parse-dmarc/issues/187)) ([dcf9d98](https://github.com/Wassergeist/parse-dmarc/commit/dcf9d98cff2ebda51c6157bcb603f0d9550941e2))
* **deps:** update module github.com/caarlos0/env/v11 to v11.4.0 ([#104](https://github.com/Wassergeist/parse-dmarc/issues/104)) ([e14e6d2](https://github.com/Wassergeist/parse-dmarc/commit/e14e6d2d7e8ecbfd8d7e5564c600b165b4135946))
* **deps:** update module github.com/caarlos0/env/v11 to v11.4.1 ([#155](https://github.com/Wassergeist/parse-dmarc/issues/155)) ([31398d3](https://github.com/Wassergeist/parse-dmarc/commit/31398d31413b2c204fb4fd0946bc3d087a490b77))
* **deps:** update module github.com/coreos/go-oidc/v3 to v3.17.0 ([#64](https://github.com/Wassergeist/parse-dmarc/issues/64)) ([306583c](https://github.com/Wassergeist/parse-dmarc/commit/306583cc37ba2c4fc7e7fffd542cdbab3edd119d))
* **deps:** update module github.com/goccy/go-json to v0.10.6 ([#119](https://github.com/Wassergeist/parse-dmarc/issues/119)) ([7ad51d3](https://github.com/Wassergeist/parse-dmarc/commit/7ad51d3a0b0395d99b542293ae5dd2a0b1e7ad32))
* **deps:** update module github.com/mattn/go-sqlite3 to v1.14.34 ([#85](https://github.com/Wassergeist/parse-dmarc/issues/85)) ([5336d3a](https://github.com/Wassergeist/parse-dmarc/commit/5336d3a64e647c32f46ddeb7d08c76fabd2a1200))
* **deps:** update module github.com/mattn/go-sqlite3 to v1.14.38 ([#129](https://github.com/Wassergeist/parse-dmarc/issues/129)) ([40fb4d4](https://github.com/Wassergeist/parse-dmarc/commit/40fb4d4a9f2048af86f3cf77be12df5a16cf134b))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1 ([#49](https://github.com/Wassergeist/parse-dmarc/issues/49)) ([58c0e79](https://github.com/Wassergeist/parse-dmarc/commit/58c0e79699392fe6661a8e6cce3092c4dd0d598f))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.3.1 ([#102](https://github.com/Wassergeist/parse-dmarc/issues/102)) ([68a80a4](https://github.com/Wassergeist/parse-dmarc/commit/68a80a4ab077af0ce4f45046211aaa77288de42b))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.4.0 ([#107](https://github.com/Wassergeist/parse-dmarc/issues/107)) ([39f1fcb](https://github.com/Wassergeist/parse-dmarc/commit/39f1fcb833098df1e2b3ab9d89c47350475055e6))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.4.1 ([#127](https://github.com/Wassergeist/parse-dmarc/issues/127)) ([ce63d91](https://github.com/Wassergeist/parse-dmarc/commit/ce63d9146a8865e1f1bf5869b6502b6f62c30dc5))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.6.1 ([#151](https://github.com/Wassergeist/parse-dmarc/issues/151)) ([eb1fb95](https://github.com/Wassergeist/parse-dmarc/commit/eb1fb953b6e91199a7de4df959802a2dffdeb9ef))
* **deps:** update module github.com/urfave/cli/v3 to v3.6.2 ([#25](https://github.com/Wassergeist/parse-dmarc/issues/25)) ([955e4e0](https://github.com/Wassergeist/parse-dmarc/commit/955e4e04d4d1d40617721ac4ec4aa70eaabad4ca))
* **deps:** update module github.com/urfave/cli/v3 to v3.7.0 ([#110](https://github.com/Wassergeist/parse-dmarc/issues/110)) ([8ce7090](https://github.com/Wassergeist/parse-dmarc/commit/8ce7090ce717ce25b563d2532b6a1a4abe7555b0))
* **deps:** update module github.com/urfave/cli/v3 to v3.8.0 ([#135](https://github.com/Wassergeist/parse-dmarc/issues/135)) ([dc06765](https://github.com/Wassergeist/parse-dmarc/commit/dc06765f9d65c0dd67a2debcb12ed3499813ed05))
* **deps:** update module modernc.org/sqlite to v1.45.0 ([#27](https://github.com/Wassergeist/parse-dmarc/issues/27)) ([2daa3b4](https://github.com/Wassergeist/parse-dmarc/commit/2daa3b4177bd5298e0f04eea8405fede27446ac2))
* **deps:** update module modernc.org/sqlite to v1.46.1 ([#101](https://github.com/Wassergeist/parse-dmarc/issues/101)) ([087c4e7](https://github.com/Wassergeist/parse-dmarc/commit/087c4e79f1fe433f82e2245a654c4e09701b69de))
* **dev:** bring the frontend to the root and simplify docker ([edc64e7](https://github.com/Wassergeist/parse-dmarc/commit/edc64e7417c910b7200bb1d727f0ba200c1a787d))
* **dev:** minify html with vite plugin ([10eb4f7](https://github.com/Wassergeist/parse-dmarc/commit/10eb4f73ce936cc7def91d5aa633f2398b19fa33))
* **dev:** remove extra compose file ([eb71b09](https://github.com/Wassergeist/parse-dmarc/commit/eb71b097005bb9d024454ee6d84c2ce5fbb59d7b))
* **dev:** remove user from compose ([e9fbf73](https://github.com/Wassergeist/parse-dmarc/commit/e9fbf73b72a2a9958d62ed6f6ffea4bad2f2771f))
* do not overwrite release notes from goreleaser ([43bec10](https://github.com/Wassergeist/parse-dmarc/commit/43bec1038b8a4ee46c6a898c89d918596bab0fdc))
* **docker:** create DB parent dirs ([e8d6497](https://github.com/Wassergeist/parse-dmarc/commit/e8d6497227119e4e7e0d348f10386e4b657b630d)), closes [#138](https://github.com/Wassergeist/parse-dmarc/issues/138)
* **docker:** make /data writable under non-root runtimes ([#171](https://github.com/Wassergeist/parse-dmarc/issues/171)) ([90a36d3](https://github.com/Wassergeist/parse-dmarc/commit/90a36d3359edf9dd7e2529b412261c3ecbfe67f7))
* **docker:** publish to the new repository ([c29118f](https://github.com/Wassergeist/parse-dmarc/commit/c29118feff7ce3b2cb4c045b6dabb2e621b7c050)), closes [#201](https://github.com/Wassergeist/parse-dmarc/issues/201)
* **docs:** add the missing png ([c9937ce](https://github.com/Wassergeist/parse-dmarc/commit/c9937cea0649a48ed92c973f223a25f8ef112f0f))
* **docs:** clean up the providers ([6879ca7](https://github.com/Wassergeist/parse-dmarc/commit/6879ca70d2a9e5d1adc7f7732d2a85869aecdd2e))
* **docs:** pin to full version for now ([2060e45](https://github.com/Wassergeist/parse-dmarc/commit/2060e451a3b7322d933ad4b8f389f28474f85fb1))
* **docs:** reference dockerhub image ([b2d2b06](https://github.com/Wassergeist/parse-dmarc/commit/b2d2b0627dbf423d760b5ee9598060334b13f290))
* **docs:** remove report card ([6fc3871](https://github.com/Wassergeist/parse-dmarc/commit/6fc3871905dfb970558dfe28e30830aa72548c13))
* **docs:** simplify README for first-time viewers ([#15](https://github.com/Wassergeist/parse-dmarc/issues/15)) ([7341b7b](https://github.com/Wassergeist/parse-dmarc/commit/7341b7bb022e5e96798e6d8d3523a377e282cf6f))
* **docs:** update brew installation ref ([d62606e](https://github.com/Wassergeist/parse-dmarc/commit/d62606e99b80c5eddfe1b2d33c5aaed2631ee0aa))
* **docs:** update docker-run snippet to the new image repo ([38d8e59](https://github.com/Wassergeist/parse-dmarc/commit/38d8e59914b3cdef25165b667e8414cd5a71fa5d))
* **docs:** update imap-host env var ([#116](https://github.com/Wassergeist/parse-dmarc/issues/116)) ([6d5aa4f](https://github.com/Wassergeist/parse-dmarc/commit/6d5aa4f32a13956723138ed786f599ad81b229c8))
* **docs:** update quickstart command with volume ([617c258](https://github.com/Wassergeist/parse-dmarc/commit/617c2584402cb309e1a652630f5f75de92df0b2e))
* **goreleaser:** take build arg for copying binary in dockerfile ([5fb9504](https://github.com/Wassergeist/parse-dmarc/commit/5fb9504f1d7a946edf0c7c0a62eed13c1abbeb8d))
* **imap:** unwrap nested message/rfc822 report attachments ([c82f118](https://github.com/Wassergeist/parse-dmarc/commit/c82f118c375481956cccd352ec00c3275aa01e0a))
* **ingest:** bound report decompression to stop a gzip/zip bomb ([#189](https://github.com/Wassergeist/parse-dmarc/issues/189)) ([f50ee41](https://github.com/Wassergeist/parse-dmarc/commit/f50ee4129f8f56fd8a1647f9859978f73e778d41))
* make github star non-intrusive ([831047b](https://github.com/Wassergeist/parse-dmarc/commit/831047b88394cd468ff57bc08c53907880e0feee))
* make linter happy ([900e314](https://github.com/Wassergeist/parse-dmarc/commit/900e314e9c32909475db562c5b6d4b038496550e))
* make server context aware and do not fail UI on empty result ([588c466](https://github.com/Wassergeist/parse-dmarc/commit/588c4660534d4111cf24738a952cbad56ebbf345))
* **mcp:** handle json error ([d17a04f](https://github.com/Wassergeist/parse-dmarc/commit/d17a04f23994911118f30522ea3ca909fb203efd))
* move cli init to cmd/ and grab values using Destination syntax ([e04af1a](https://github.com/Wassergeist/parse-dmarc/commit/e04af1aefd1c33c74780afd2b4e00805a951c1b0))
* move main.go to the root ([0dd5787](https://github.com/Wassergeist/parse-dmarc/commit/0dd57878ec5f7d6707faca0189834a6d91c266e5))
* **precommit:** run prettier on changed files only ([729f64d](https://github.com/Wassergeist/parse-dmarc/commit/729f64db6044da83bd7bfc356b926a8e05d062cf))
* propagate errors with context ([85014c8](https://github.com/Wassergeist/parse-dmarc/commit/85014c8aafaf03a6d6b512a7be7eaccab971a4b0))
* publish to grafana community dashboards instead ([86f40b0](https://github.com/Wassergeist/parse-dmarc/commit/86f40b052607861de3d1a843f9e9c9e08ba57481))
* remove broken dns generator ([df39352](https://github.com/Wassergeist/parse-dmarc/commit/df393527710b487bc0bdbe7d67d9a3f909623b92))
* stay compliant with brew core version print contraint ([35f8261](https://github.com/Wassergeist/parse-dmarc/commit/35f826169b3864811cc7684e7c87190518c373a5))
* switch env and json config parsing order ([#83](https://github.com/Wassergeist/parse-dmarc/issues/83)) ([12afcd2](https://github.com/Wassergeist/parse-dmarc/commit/12afcd2e295c31b60a4696ba5c80edac1ae424e2)), closes [#84](https://github.com/Wassergeist/parse-dmarc/issues/84)
* **UI:** handle no data in compliance score ([#90](https://github.com/Wassergeist/parse-dmarc/issues/90)) ([eced435](https://github.com/Wassergeist/parse-dmarc/commit/eced435ab7b08a73829d4b7b6e09ed31d324834b))
* update docker image reference to dockerhub ([6afcb3a](https://github.com/Wassergeist/parse-dmarc/commit/6afcb3abe19beda674d19ab7911019c384ba54b1))
* update the UI on reports API parsing ([b85b59e](https://github.com/Wassergeist/parse-dmarc/commit/b85b59ef577ffeb6b4dacd1188e6e55d549b1b5a))
* use app for static config ([3820862](https://github.com/Wassergeist/parse-dmarc/commit/3820862f8fffe26fd3dad7c004d1c0b04fc8be9e))
* use PEEK when mark-seen is set to false ([#125](https://github.com/Wassergeist/parse-dmarc/issues/125)) ([e54e269](https://github.com/Wassergeist/parse-dmarc/commit/e54e26900faed10f916afed47d23b3a2a1f65467))
* use png for favicon ([4704ad7](https://github.com/Wassergeist/parse-dmarc/commit/4704ad7c23fbd31913adf677e50f001393a1adda))


### Build & Dependencies

* **deps:** pin dependencies ([#158](https://github.com/Wassergeist/parse-dmarc/issues/158)) ([7f6e974](https://github.com/Wassergeist/parse-dmarc/commit/7f6e974626f2906216172e7020f42a5a0b1c67ee))
* **deps:** update actions/checkout action to v7 ([#165](https://github.com/Wassergeist/parse-dmarc/issues/165)) ([cd5f10f](https://github.com/Wassergeist/parse-dmarc/commit/cd5f10f8937ca82991f2d921f9e82c1c7c3f5460))
* **deps:** update actions/setup-go action to v7 ([#175](https://github.com/Wassergeist/parse-dmarc/issues/175)) ([99e4a48](https://github.com/Wassergeist/parse-dmarc/commit/99e4a48f9e04e854c09b2ecd906ed543131c6432))
* **deps:** update golang docker tag to v1.27 ([#184](https://github.com/Wassergeist/parse-dmarc/issues/184)) ([afdb202](https://github.com/Wassergeist/parse-dmarc/commit/afdb202e842fcddaf88d6f54dcbeb779dd4ac9d8))
* **deps:** update module github.com/coreos/go-oidc/v3 to v3.19.0 ([#150](https://github.com/Wassergeist/parse-dmarc/issues/150)) ([a718263](https://github.com/Wassergeist/parse-dmarc/commit/a7182637435b878fbf70711926871d83367e51ca))
* **deps:** update module github.com/coreos/go-oidc/v3 to v3.20.0 ([#170](https://github.com/Wassergeist/parse-dmarc/issues/170)) ([19e0261](https://github.com/Wassergeist/parse-dmarc/commit/19e026181157b375a80606b6688c6b0a8b9d29e5))
* **deps:** update module github.com/coreos/go-oidc/v3 to v3.21.0 ([#188](https://github.com/Wassergeist/parse-dmarc/issues/188)) ([675492d](https://github.com/Wassergeist/parse-dmarc/commit/675492dfb372adcb03b5b5b56db9c367e85ddcaa))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.7.0 ([#178](https://github.com/Wassergeist/parse-dmarc/issues/178)) ([f10b023](https://github.com/Wassergeist/parse-dmarc/commit/f10b0233751ebdaa1e9d94a28906fb69ac417672))
* **deps:** update module github.com/prometheus/client_golang to v1.24.0 ([#173](https://github.com/Wassergeist/parse-dmarc/issues/173)) ([175d56d](https://github.com/Wassergeist/parse-dmarc/commit/175d56d3b2d352fde15eb491df7f4d0fec952818))
* **deps:** update module github.com/urfave/cli/v3 to v3.10.0 ([#160](https://github.com/Wassergeist/parse-dmarc/issues/160)) ([0d822fe](https://github.com/Wassergeist/parse-dmarc/commit/0d822fefcc3da78cb1aed015e775fc815e0efc50))
* **deps:** update module github.com/urfave/cli/v3 to v3.11.0 ([#181](https://github.com/Wassergeist/parse-dmarc/issues/181)) ([e5d3fef](https://github.com/Wassergeist/parse-dmarc/commit/e5d3fefa373a85687bae2e13421e04308e0a7f44))
* **deps:** update module modernc.org/sqlite to v1.52.0 ([#161](https://github.com/Wassergeist/parse-dmarc/issues/161)) ([68a9241](https://github.com/Wassergeist/parse-dmarc/commit/68a9241a8509d3ccc095fa2067f1db16d8792647))
* **deps:** update module modernc.org/sqlite to v1.54.0 ([#163](https://github.com/Wassergeist/parse-dmarc/issues/163)) ([31006bf](https://github.com/Wassergeist/parse-dmarc/commit/31006bfcb3ec668c668d03199bd7c157c25f4c1f))
* **deps:** update module modernc.org/sqlite to v1.58.0 ([#179](https://github.com/Wassergeist/parse-dmarc/issues/179)) ([21daa5b](https://github.com/Wassergeist/parse-dmarc/commit/21daa5b8c797968a5d3c9d51533cc41138082843))
* **deps:** update non-major dependencies ([#159](https://github.com/Wassergeist/parse-dmarc/issues/159)) ([c1db71f](https://github.com/Wassergeist/parse-dmarc/commit/c1db71f6da10a0bab1087cbeaaba7af7cc49f9fe))
* **deps:** update non-major dependencies ([#164](https://github.com/Wassergeist/parse-dmarc/issues/164)) ([7bac6aa](https://github.com/Wassergeist/parse-dmarc/commit/7bac6aaa98778bff48b1f547e3ddac56bc747d35))
* **deps:** update non-major dependencies ([#174](https://github.com/Wassergeist/parse-dmarc/issues/174)) ([4bae98f](https://github.com/Wassergeist/parse-dmarc/commit/4bae98f85398c2d160401c64c4f37a0bdab5f349))
* **deps:** update sigstore/cosign-installer action to v4 ([#167](https://github.com/Wassergeist/parse-dmarc/issues/167)) ([cbdda89](https://github.com/Wassergeist/parse-dmarc/commit/cbdda8962697194726f28b8c7a8d63f3f5f248d6))
* **deps:** upgrade JS dependencies (ky v2, pinia v4, sharp 0.35, vite 8.2.2) ([#192](https://github.com/Wassergeist/parse-dmarc/issues/192)) ([d80f128](https://github.com/Wassergeist/parse-dmarc/commit/d80f1289638e230e2193ab8e1f3b349341679411))


### Chores

* add brew publish to goreleaser ([49d4057](https://github.com/Wassergeist/parse-dmarc/commit/49d40576c9c89107adcf0b5850a4dec02ae4b460))
* add missing contributor from last release ([ab4fc45](https://github.com/Wassergeist/parse-dmarc/commit/ab4fc45a9fbc31104e4b0f8b540184a238ac8a84))
* **CI:** build docker edge on push to main ([dd9c090](https://github.com/Wassergeist/parse-dmarc/commit/dd9c090c98b70d5f15273de74b9c6e0507310407))
* **CI:** build goreleaser in dry-run mode on main ([26b1f1a](https://github.com/Wassergeist/parse-dmarc/commit/26b1f1aad5fcf82002c82aa31f58808ce6244786))
* **CI:** enrich release changelog ([5f9f9e8](https://github.com/Wassergeist/parse-dmarc/commit/5f9f9e85e69a0816176df62f2619e6291168d12b))
* **CI:** update checkout to v6 ([9bb92d7](https://github.com/Wassergeist/parse-dmarc/commit/9bb92d727658bd4f442597e7aaf28ae971e7a454))
* **CI:** use separate dockerfile for goreleaser ([a2ce77c](https://github.com/Wassergeist/parse-dmarc/commit/a2ce77c495002059e373a3dcdd97c26c9118c7bf))
* **config:** migrate config renovate.json ([#66](https://github.com/Wassergeist/parse-dmarc/issues/66)) ([09871ee](https://github.com/Wassergeist/parse-dmarc/commit/09871ee40afcdbc7c5f2a5bc86d536f998bc6372))
* **deps:** lock file maintenance ([#168](https://github.com/Wassergeist/parse-dmarc/issues/168)) ([32b0edf](https://github.com/Wassergeist/parse-dmarc/commit/32b0edfc42e572436599129ee7ed55029874dbad))
* **deps:** lock file maintenance ([#177](https://github.com/Wassergeist/parse-dmarc/issues/177)) ([b3da7f2](https://github.com/Wassergeist/parse-dmarc/commit/b3da7f2d067e413323580ef20f6eb4b17b99a64e))
* **deps:** lock file maintenance ([#198](https://github.com/Wassergeist/parse-dmarc/issues/198)) ([a4d5d96](https://github.com/Wassergeist/parse-dmarc/commit/a4d5d96a7b4e69392babefbe0541a25e4001a640))
* **deps:** update actions/attest-build-provenance action to v4 ([#106](https://github.com/Wassergeist/parse-dmarc/issues/106)) ([9301b0f](https://github.com/Wassergeist/parse-dmarc/commit/9301b0f2507284960e58182f629b239f9f734b4a))
* **deps:** update dependency @vitejs/plugin-vue to v6.0.4 ([b7eb415](https://github.com/Wassergeist/parse-dmarc/commit/b7eb41518fef6b72471e389b0da9814abb6a37c8))
* **deps:** update dependency ky to v1.14.2 ([#81](https://github.com/Wassergeist/parse-dmarc/issues/81)) ([97f402c](https://github.com/Wassergeist/parse-dmarc/commit/97f402cedd7f9db7419eb9b9da43d16f78691975))
* **deps:** update dependency ky to v1.14.3 ([#96](https://github.com/Wassergeist/parse-dmarc/issues/96)) ([2666c4c](https://github.com/Wassergeist/parse-dmarc/commit/2666c4c496259a82dff2359d799e5a1109a072bb))
* **deps:** update dependency lightningcss to v1.31.1 ([#98](https://github.com/Wassergeist/parse-dmarc/issues/98)) ([357add7](https://github.com/Wassergeist/parse-dmarc/commit/357add70be8d5deb2156e862f02d52e73f9838bd))
* **deps:** update dependency lightningcss to v1.32.0 ([#114](https://github.com/Wassergeist/parse-dmarc/issues/114)) ([120608e](https://github.com/Wassergeist/parse-dmarc/commit/120608eeb990fcfecc14665c62cfcfdff5551f48))
* **deps:** update dependency vite to v7.2.2 ([#14](https://github.com/Wassergeist/parse-dmarc/issues/14)) ([e13aa2c](https://github.com/Wassergeist/parse-dmarc/commit/e13aa2cb072e4e9fe3d03f1f45b77663e811bcaf))
* **deps:** update dependency vite to v7.3.1 ([#28](https://github.com/Wassergeist/parse-dmarc/issues/28)) ([f10c5c4](https://github.com/Wassergeist/parse-dmarc/commit/f10c5c47edbabab4608186e72759ab6fce6da450))
* **deps:** update dependency vite to v8 ([#123](https://github.com/Wassergeist/parse-dmarc/issues/123)) ([e288c1c](https://github.com/Wassergeist/parse-dmarc/commit/e288c1c0af03e4d7ee50fe5a5ad7bafdb63558d7))
* **deps:** update dependency vite to v8.0.1 ([#132](https://github.com/Wassergeist/parse-dmarc/issues/132)) ([07ecd5b](https://github.com/Wassergeist/parse-dmarc/commit/07ecd5b13a11e43e480522bd3ec1d69eec51e166))
* **deps:** update dependency vite to v8.0.3 ([#137](https://github.com/Wassergeist/parse-dmarc/issues/137)) ([9c011cd](https://github.com/Wassergeist/parse-dmarc/commit/9c011cd2963e436141124b9539aa8f6b18c94772))
* **deps:** update dependency vite-plugin-compression2 to v2.4.0 ([#42](https://github.com/Wassergeist/parse-dmarc/issues/42)) ([d03427b](https://github.com/Wassergeist/parse-dmarc/commit/d03427b58f65dbe4983530f5b311b0c76b7e752d))
* **deps:** update dependency vite-plugin-compression2 to v2.5.0 ([#109](https://github.com/Wassergeist/parse-dmarc/issues/109)) ([c565f70](https://github.com/Wassergeist/parse-dmarc/commit/c565f706984dac64e38b25552235b3969ea91a0f))
* **deps:** update dependency vite-plugin-compression2 to v2.5.1 ([#117](https://github.com/Wassergeist/parse-dmarc/issues/117)) ([9d410ef](https://github.com/Wassergeist/parse-dmarc/commit/9d410ef8e7c1234fa418a1a11fcecf07f3c8fafc))
* **deps:** update dependency vite-plugin-compression2 to v2.5.3 - abandoned ([#128](https://github.com/Wassergeist/parse-dmarc/issues/128)) ([7d81cec](https://github.com/Wassergeist/parse-dmarc/commit/7d81cecbdb57e14874b3e035d8dbb750ec4c3724))
* **deps:** update dependency vue to v3.5.24 ([#11](https://github.com/Wassergeist/parse-dmarc/issues/11)) ([126431e](https://github.com/Wassergeist/parse-dmarc/commit/126431ebd05aef47932ec31d8a8d4145ec70a167))
* **deps:** update dependency vue to v3.5.28 ([#31](https://github.com/Wassergeist/parse-dmarc/issues/31)) ([0a80797](https://github.com/Wassergeist/parse-dmarc/commit/0a80797ad6c28cad75e554f059bdc5878e61c17e))
* **deps:** update dependency vue to v3.5.29 ([#105](https://github.com/Wassergeist/parse-dmarc/issues/105)) ([c515dcb](https://github.com/Wassergeist/parse-dmarc/commit/c515dcbea08f1c79ef949aea67044cc708d7930d))
* **deps:** update dependency vue to v3.5.30 ([#115](https://github.com/Wassergeist/parse-dmarc/issues/115)) ([7b71009](https://github.com/Wassergeist/parse-dmarc/commit/7b7100946744414f512a9a0a416bb34d1ead0f50))
* **deps:** update dependency vue to v3.5.31 ([1a2ff72](https://github.com/Wassergeist/parse-dmarc/commit/1a2ff722303e59b023af82547bd5d78ee03582b2))
* **deps:** update dependency zod to v4.3.6 ([#97](https://github.com/Wassergeist/parse-dmarc/issues/97)) ([8855827](https://github.com/Wassergeist/parse-dmarc/commit/8855827a8e1ebcf853ca18d26098f85c3d40abd3))
* **deps:** update docker/login-action action to v4 ([#112](https://github.com/Wassergeist/parse-dmarc/issues/112)) ([7d3b905](https://github.com/Wassergeist/parse-dmarc/commit/7d3b905b7322bd700b75b7ead8ce8e0fd6dade31))
* **deps:** update docker/setup-buildx-action action to v4 ([#113](https://github.com/Wassergeist/parse-dmarc/issues/113)) ([a83de1a](https://github.com/Wassergeist/parse-dmarc/commit/a83de1a0652611f15d2305ecbaab6c06e50d7632))
* **deps:** update docker/setup-qemu-action action to v4 ([#111](https://github.com/Wassergeist/parse-dmarc/issues/111)) ([7b4b802](https://github.com/Wassergeist/parse-dmarc/commit/7b4b802596a1c901d09afd131aa68af51b57bad9))
* **deps:** update frontend libs to latest ([9c41a85](https://github.com/Wassergeist/parse-dmarc/commit/9c41a85001e73a8e16a9eb5b226396eb838f141f))
* **deps:** update go libs to latest ([455bc5b](https://github.com/Wassergeist/parse-dmarc/commit/455bc5be5e68b26963f209842887578e8810ecfa))
* **deps:** update go.mod ([1dad5ff](https://github.com/Wassergeist/parse-dmarc/commit/1dad5ff217aca6a786546339bcf5d11222ee0901))
* **deps:** update golang docker tag to v1.26 ([#99](https://github.com/Wassergeist/parse-dmarc/issues/99)) ([e062abd](https://github.com/Wassergeist/parse-dmarc/commit/e062abde8aa1d1ebb758318c061bc7f268fdd8fd))
* **deps:** update golangci/golangci-lint-action action to v9 ([#22](https://github.com/Wassergeist/parse-dmarc/issues/22)) ([dbc73b3](https://github.com/Wassergeist/parse-dmarc/commit/dbc73b38626e34a863c78a5615daf0f95f0db82a))
* **deps:** update googleapis/release-please-action action to v5 ([#152](https://github.com/Wassergeist/parse-dmarc/issues/152)) ([90ea76d](https://github.com/Wassergeist/parse-dmarc/commit/90ea76ddab01aafa7bfb3f378dd7a57134aa903c))
* **deps:** update goreleaser/goreleaser-action action to v7 ([#103](https://github.com/Wassergeist/parse-dmarc/issues/103)) ([075c548](https://github.com/Wassergeist/parse-dmarc/commit/075c548ef436ad3ca21233cf6b58362f6daf9faa))
* **deps:** update js dependencies to latest ([2aff67f](https://github.com/Wassergeist/parse-dmarc/commit/2aff67f61eddfcd8a0903822594b029fe3f98941))
* **deps:** update UI libs ([40ef5a4](https://github.com/Wassergeist/parse-dmarc/commit/40ef5a477b0e74b4fcb67102bd9410ea638ae7a1))
* **deps:** update UI libs, again ([cb9b4c5](https://github.com/Wassergeist/parse-dmarc/commit/cb9b4c5775788c32718463b3f50a6f63a3b9497f))
* **deps:** upgrade go libs ([c391ac6](https://github.com/Wassergeist/parse-dmarc/commit/c391ac6a6c2b3241bff4bab9007fa5fce1ddc3c5))
* **deps:** upgrade vite compression ([4928961](https://github.com/Wassergeist/parse-dmarc/commit/49289618e4546fc10b7e0f0295689cba7baa9d93))
* **dev:** add flake nixos ([6d29961](https://github.com/Wassergeist/parse-dmarc/commit/6d29961e1ddb4f7f9c9e4a8fb52f44f5be25179b))
* **dev:** do not track local files ([8946d5b](https://github.com/Wassergeist/parse-dmarc/commit/8946d5b673b7a392a164004809457eb34109ded2))
* **dev:** generate png for images ([2f096a1](https://github.com/Wassergeist/parse-dmarc/commit/2f096a1109e5f82c03f8c1757473cbe0c92df52a))
* **dev:** replace fetch with ky http client ([#75](https://github.com/Wassergeist/parse-dmarc/issues/75)) ([115c91c](https://github.com/Wassergeist/parse-dmarc/commit/115c91cc98aebc70528eaafa0a673832a8ad737e))
* **dev:** update nix formatter ([417a3ab](https://github.com/Wassergeist/parse-dmarc/commit/417a3abde508fc0ad3f10b85416a91286f16c9ed))
* **dev:** update pre-commit hooks ([51b83c6](https://github.com/Wassergeist/parse-dmarc/commit/51b83c6a8a6da2e0a12bcdc1a9ca076c2ca557c0))
* do not track envrc ([6ec3175](https://github.com/Wassergeist/parse-dmarc/commit/6ec3175e9364b30bfd87ce87d260da3f09533f39))
* **docs:** add claude.md project guide ([#32](https://github.com/Wassergeist/parse-dmarc/issues/32)) ([21f2fb2](https://github.com/Wassergeist/parse-dmarc/commit/21f2fb2568bda5186cacaf552963ae81b1cf9aca))
* **docs:** add social preview to the top of readme ([1871652](https://github.com/Wassergeist/parse-dmarc/commit/1871652c6485b4b7847ec54038b74287d2a18745))
* **docs:** add some badges to readme ([8ac0762](https://github.com/Wassergeist/parse-dmarc/commit/8ac07621b3181a3deacaf5eb2ef15bc5e13ee83f))
* **docs:** publish brew to the new org-wide repository ([432f010](https://github.com/Wassergeist/parse-dmarc/commit/432f010856838969b99a078339d72830421a9678))
* **docs:** remove one cmd installation from brew ([9e4864a](https://github.com/Wassergeist/parse-dmarc/commit/9e4864ad0dd142806224c76f240b7fe2d10ac9bb))
* **docs:** rewrite README around Parse DMARC, split metrics and MCP docs, credit DMARCguard ([#196](https://github.com/Wassergeist/parse-dmarc/issues/196)) ([aed41c1](https://github.com/Wassergeist/parse-dmarc/commit/aed41c1f1a7bd37e13d1276bbc8b3db5d0a2299d))
* **docs:** simplify quickstart deployment ([cb5cd9b](https://github.com/Wassergeist/parse-dmarc/commit/cb5cd9b0f792bb9b83376c47a08fd1c73124b66d))
* go mod tidy ([8a9cd68](https://github.com/Wassergeist/parse-dmarc/commit/8a9cd68905c82a52da0251d65013246867487063))
* **license:** add contact info ([4013b0b](https://github.com/Wassergeist/parse-dmarc/commit/4013b0b93887ec4372caa744f7080a737bf90cc7))
* **local:** add dovecot for local development ([dbb76b2](https://github.com/Wassergeist/parse-dmarc/commit/dbb76b2ba1848cc9a159319735341710bda8bcd3))
* **main:** release 1.0.0 ([#4](https://github.com/Wassergeist/parse-dmarc/issues/4)) ([43204f0](https://github.com/Wassergeist/parse-dmarc/commit/43204f0fd7a4a42fc9077717cc8738372984889c))
* **main:** release 1.0.1 ([#16](https://github.com/Wassergeist/parse-dmarc/issues/16)) ([925f279](https://github.com/Wassergeist/parse-dmarc/commit/925f2799c9cf85b7a8cce71aa6b0d77cd05ed66b))
* **main:** release 1.0.2 ([#21](https://github.com/Wassergeist/parse-dmarc/issues/21)) ([844e246](https://github.com/Wassergeist/parse-dmarc/commit/844e2461889dbfdcf8933437e486de6dd46f2590))
* **main:** release 1.1.0 ([#24](https://github.com/Wassergeist/parse-dmarc/issues/24)) ([c2e1f23](https://github.com/Wassergeist/parse-dmarc/commit/c2e1f239751585a4a7ad59e51d0a1051e9d26c4b))
* **main:** release 1.2.0 ([#36](https://github.com/Wassergeist/parse-dmarc/issues/36)) ([d481792](https://github.com/Wassergeist/parse-dmarc/commit/d4817923b496da418bfc6d1cba8d8ac87357762a))
* **main:** release 1.2.1 ([#43](https://github.com/Wassergeist/parse-dmarc/issues/43)) ([dfe37c9](https://github.com/Wassergeist/parse-dmarc/commit/dfe37c9766fb7911948da18c5a14b867c138067e))
* **main:** release 1.2.2 ([#44](https://github.com/Wassergeist/parse-dmarc/issues/44)) ([1a39541](https://github.com/Wassergeist/parse-dmarc/commit/1a39541d2651c8ee4b819ecdacb77c23610af8de))
* **main:** release 1.3.0 ([#47](https://github.com/Wassergeist/parse-dmarc/issues/47)) ([221d2b3](https://github.com/Wassergeist/parse-dmarc/commit/221d2b362dbb1245c6ff3517a26545f98bbc3de9))
* **main:** release 1.3.1 ([#53](https://github.com/Wassergeist/parse-dmarc/issues/53)) ([6115466](https://github.com/Wassergeist/parse-dmarc/commit/61154668a20bc526fbe65e6a445d050185ea9017))
* **main:** release 1.3.10 ([#62](https://github.com/Wassergeist/parse-dmarc/issues/62)) ([075b27a](https://github.com/Wassergeist/parse-dmarc/commit/075b27a9d2eea2be4c7bbb68678290cd09ec824e))
* **main:** release 1.3.2 ([#54](https://github.com/Wassergeist/parse-dmarc/issues/54)) ([335f593](https://github.com/Wassergeist/parse-dmarc/commit/335f5931a9e05210390f43d21bcf798dd3ef0549))
* **main:** release 1.3.3 ([#55](https://github.com/Wassergeist/parse-dmarc/issues/55)) ([2e919c7](https://github.com/Wassergeist/parse-dmarc/commit/2e919c76fb446ac078c83b0def529b0269dca08e))
* **main:** release 1.3.4 ([#56](https://github.com/Wassergeist/parse-dmarc/issues/56)) ([6e0062a](https://github.com/Wassergeist/parse-dmarc/commit/6e0062a20e2b51066578fb86a37d8d82536c1092))
* **main:** release 1.3.5 ([#57](https://github.com/Wassergeist/parse-dmarc/issues/57)) ([ac749db](https://github.com/Wassergeist/parse-dmarc/commit/ac749dbf306335e4dcb7828171af2a81dc1c8c3b))
* **main:** release 1.3.6 ([#58](https://github.com/Wassergeist/parse-dmarc/issues/58)) ([85b66ed](https://github.com/Wassergeist/parse-dmarc/commit/85b66ed1aeffa6507282287ccace13d3898cb1d2))
* **main:** release 1.3.7 ([#59](https://github.com/Wassergeist/parse-dmarc/issues/59)) ([395f524](https://github.com/Wassergeist/parse-dmarc/commit/395f524ac88029fff72aa7f26f2fc7956bb98ace))
* **main:** release 1.3.8 ([#60](https://github.com/Wassergeist/parse-dmarc/issues/60)) ([a80a8ce](https://github.com/Wassergeist/parse-dmarc/commit/a80a8ce1b9b3ad72a61b3492ac7515a96fd2b91c))
* **main:** release 1.3.9 ([#61](https://github.com/Wassergeist/parse-dmarc/issues/61)) ([1accb98](https://github.com/Wassergeist/parse-dmarc/commit/1accb98614e91f6e011c4a28dc80901ec5053c4f))
* **main:** release 1.4.0 ([#63](https://github.com/Wassergeist/parse-dmarc/issues/63)) ([656b909](https://github.com/Wassergeist/parse-dmarc/commit/656b909f9b89bd8d69a7f74f0020d184302b2f5a))
* **main:** release 1.4.1 ([#69](https://github.com/Wassergeist/parse-dmarc/issues/69)) ([021ba02](https://github.com/Wassergeist/parse-dmarc/commit/021ba02f13941add258fb29a198b1bc31588bca1))
* **main:** release 1.4.2 ([#74](https://github.com/Wassergeist/parse-dmarc/issues/74)) ([cdf2502](https://github.com/Wassergeist/parse-dmarc/commit/cdf25025aa3df65178bdf07f2da27fd8ce4b395a))
* **main:** release 1.4.3 ([#77](https://github.com/Wassergeist/parse-dmarc/issues/77)) ([8a15e34](https://github.com/Wassergeist/parse-dmarc/commit/8a15e34e0d65fd6202963999104c040127580c6d))
* **main:** release 1.4.4 ([d06ccd7](https://github.com/Wassergeist/parse-dmarc/commit/d06ccd7ff5224d810a9efd5b9098569d5a3a8b7d))
* **main:** release 1.4.5 ([#88](https://github.com/Wassergeist/parse-dmarc/issues/88)) ([21b86a4](https://github.com/Wassergeist/parse-dmarc/commit/21b86a4a0573d4fc9987a6dcb2896dd065a388f4))
* **main:** release 1.4.6 ([#91](https://github.com/Wassergeist/parse-dmarc/issues/91)) ([fac0e0f](https://github.com/Wassergeist/parse-dmarc/commit/fac0e0f0415ca96bf5aaa90a4dba88b9ae3d8ec8))
* **main:** release 1.4.7 ([#92](https://github.com/Wassergeist/parse-dmarc/issues/92)) ([3afe5b2](https://github.com/Wassergeist/parse-dmarc/commit/3afe5b2593b4628d9f022ccdedb70f3e314810da))
* **main:** release 1.4.8 ([#93](https://github.com/Wassergeist/parse-dmarc/issues/93)) ([5e8d293](https://github.com/Wassergeist/parse-dmarc/commit/5e8d293f2402e38eb5b8833583c6c7f1e67d7298))
* **main:** release 1.5.0 ([#95](https://github.com/Wassergeist/parse-dmarc/issues/95)) ([bdc5377](https://github.com/Wassergeist/parse-dmarc/commit/bdc53770eebec88d0d9fdf565b2cc8885d39da97))
* **main:** release 1.5.1 ([#124](https://github.com/Wassergeist/parse-dmarc/issues/124)) ([dd76450](https://github.com/Wassergeist/parse-dmarc/commit/dd764509de03eb50d179c84446e49f08013e6223))
* **main:** release 1.5.2 ([#126](https://github.com/Wassergeist/parse-dmarc/issues/126)) ([57055c3](https://github.com/Wassergeist/parse-dmarc/commit/57055c3a32e09658ea1dc0a6ee33d764e26386d5))
* **main:** release 1.5.3 ([#136](https://github.com/Wassergeist/parse-dmarc/issues/136)) ([d909ce4](https://github.com/Wassergeist/parse-dmarc/commit/d909ce41740b3053fc0f643701659027951c4e8a))
* **main:** release 1.5.4 ([#140](https://github.com/Wassergeist/parse-dmarc/issues/140)) ([2a8d6df](https://github.com/Wassergeist/parse-dmarc/commit/2a8d6dfe688567d4056be7cc4a12157ee08423ce))
* **main:** release 1.5.5 ([#172](https://github.com/Wassergeist/parse-dmarc/issues/172)) ([e23f266](https://github.com/Wassergeist/parse-dmarc/commit/e23f26636724010cfa1afd5765b8ddb77b0ac109))
* **main:** release 1.5.6 ([#183](https://github.com/Wassergeist/parse-dmarc/issues/183)) ([8285cbd](https://github.com/Wassergeist/parse-dmarc/commit/8285cbdbd00a1d86b07fb96a7c2452c5bcd50048))
* **main:** release 1.5.7 ([#190](https://github.com/Wassergeist/parse-dmarc/issues/190)) ([04101e9](https://github.com/Wassergeist/parse-dmarc/commit/04101e9c64db412ff5fd71d60d09308d8e0420bc))
* **main:** release 1.6.0 ([#195](https://github.com/Wassergeist/parse-dmarc/issues/195)) ([0d95078](https://github.com/Wassergeist/parse-dmarc/commit/0d950780d662074e49ef0c7374dc2ebdec0f79c6))
* **main:** release 1.6.1 ([#197](https://github.com/Wassergeist/parse-dmarc/issues/197)) ([1f1ef4d](https://github.com/Wassergeist/parse-dmarc/commit/1f1ef4d9beeb36ef2b60e65f0b4d70d6091dab17))
* **main:** release 1.7.0 ([#202](https://github.com/Wassergeist/parse-dmarc/issues/202)) ([ae659f7](https://github.com/Wassergeist/parse-dmarc/commit/ae659f762861dee5b486c2952f1f38f0afebc808))
* migrate from stdlib logs to zerolog project-wide ([b8685f3](https://github.com/Wassergeist/parse-dmarc/commit/b8685f3278979622740fd0d340c694b4f3284bcc))
* move screenshots to their own dir ([1f36c8a](https://github.com/Wassergeist/parse-dmarc/commit/1f36c8a9bf966c51d6d17df43b6639515620bb40))
* prepare next release 1.4.3 ([141e2ed](https://github.com/Wassergeist/parse-dmarc/commit/141e2ed8c97a34e3d24d359645c0aa67237a47eb))
* remove noisy outputs from .goreleaser.yml ([0807b96](https://github.com/Wassergeist/parse-dmarc/commit/0807b9647d8b27ba600bf6a3ca4aa24d4aa5cadb))
* **renovate:** extend the default base ([3071e76](https://github.com/Wassergeist/parse-dmarc/commit/3071e7654736fb49e73a95239d4d4441ad9d0825))
* run go mod tidy and fix linting errors ([e51baae](https://github.com/Wassergeist/parse-dmarc/commit/e51baae966db15724cdf2dc15216158750eae9b0))
* tidy up goreleaser config ([02d8789](https://github.com/Wassergeist/parse-dmarc/commit/02d8789e007be307122dfdecbcd7b5f4dc77e705))
* update claude.md with frontend structure and features ([#38](https://github.com/Wassergeist/parse-dmarc/issues/38)) ([faa99a8](https://github.com/Wassergeist/parse-dmarc/commit/faa99a8ccdfb95529884697c65ce214187b277f8))
* update CLAUDE.md with mcp, metrics and deployments ([#71](https://github.com/Wassergeist/parse-dmarc/issues/71)) ([3cd362d](https://github.com/Wassergeist/parse-dmarc/commit/3cd362d61748c51a7bf7f679dfbcb6ba336cc9d3))
* update generated PNGs ([e8a0bc9](https://github.com/Wassergeist/parse-dmarc/commit/e8a0bc94311c80a50065856286a2d93f9dc27f72))
* update package-lock.json ([157b0b0](https://github.com/Wassergeist/parse-dmarc/commit/157b0b0edeae8885771e196be53d2ac7082f66d4))


### Documentation

* add comprehensive product roadmap ([#33](https://github.com/Wassergeist/parse-dmarc/issues/33)) ([d13aa0d](https://github.com/Wassergeist/parse-dmarc/commit/d13aa0d5af81138744910aa49764ceb347b7e31c))
* add Homebrew installation instructions to README ([#45](https://github.com/Wassergeist/parse-dmarc/issues/45)) ([7e93fb9](https://github.com/Wassergeist/parse-dmarc/commit/7e93fb980a0f55e8bfd146ed0a3babdbffb17c72))
* update CLAUDE.md with current codebase state ([#79](https://github.com/Wassergeist/parse-dmarc/issues/79)) ([534e640](https://github.com/Wassergeist/parse-dmarc/commit/534e640c7906827ddd39807874fad606137752c7))
* update claude.md with project improvements ([bfda097](https://github.com/Wassergeist/parse-dmarc/commit/bfda09706ba3e8d9668ae4b78d8784a17926795d))
* update README and ROADMAP for DNS Record Generator ([4103451](https://github.com/Wassergeist/parse-dmarc/commit/4103451273d93b000962242c7bf90acc9012f2d8))

## [1.7.0](https://github.com/dmarcguardhq/parse-dmarc/compare/v1.6.1...v1.7.0) (2026-09-18)


### Features

* log attachments rejected as non-DMARC ([#208](https://github.com/dmarcguardhq/parse-dmarc/issues/208)) ([31b6a23](https://github.com/dmarcguardhq/parse-dmarc/commit/31b6a230fe56c749296dcac7b487675651b852a6))


### Bug Fixes

* accept DMARC reports sent as inline MIME parts ([#207](https://github.com/dmarcguardhq/parse-dmarc/issues/207)) ([20b2be0](https://github.com/dmarcguardhq/parse-dmarc/commit/20b2be03d8a43197eabde63eff452e276eaf2b63))
* **docs:** update brew installation ref ([d62606e](https://github.com/dmarcguardhq/parse-dmarc/commit/d62606e99b80c5eddfe1b2d33c5aaed2631ee0aa))


### Chores

* **docs:** publish brew to the new org-wide repository ([432f010](https://github.com/dmarcguardhq/parse-dmarc/commit/432f010856838969b99a078339d72830421a9678))

## [1.6.1](https://github.com/dmarcguardhq/parse-dmarc/compare/v1.6.0...v1.6.1) (2026-09-15)


### Bug Fixes

* **CI:** publish the image to the new repo on main ([f89a9ec](https://github.com/dmarcguardhq/parse-dmarc/commit/f89a9ece752e3de0bf5c42f3583f41209cf46d55))
* **docker:** publish to the new repository ([c29118f](https://github.com/dmarcguardhq/parse-dmarc/commit/c29118feff7ce3b2cb4c045b6dabb2e621b7c050)), closes [#201](https://github.com/dmarcguardhq/parse-dmarc/issues/201)
* **docs:** remove report card ([6fc3871](https://github.com/dmarcguardhq/parse-dmarc/commit/6fc3871905dfb970558dfe28e30830aa72548c13))
* **docs:** update docker-run snippet to the new image repo ([38d8e59](https://github.com/dmarcguardhq/parse-dmarc/commit/38d8e59914b3cdef25165b667e8414cd5a71fa5d))


### Chores

* **deps:** lock file maintenance ([#198](https://github.com/dmarcguardhq/parse-dmarc/issues/198)) ([a4d5d96](https://github.com/dmarcguardhq/parse-dmarc/commit/a4d5d96a7b4e69392babefbe0541a25e4001a640))
* **docs:** rewrite README around Parse DMARC, split metrics and MCP docs, credit DMARCguard ([#196](https://github.com/dmarcguardhq/parse-dmarc/issues/196)) ([aed41c1](https://github.com/dmarcguardhq/parse-dmarc/commit/aed41c1f1a7bd37e13d1276bbc8b3db5d0a2299d))

## [1.6.0](https://github.com/dmarcguardhq/dmarcguard/compare/v1.5.7...v1.6.0) (2026-09-11)


### Features

* **imap:** support internal CA, skip-verify and STARTTLS ([#194](https://github.com/dmarcguardhq/dmarcguard/issues/194)) ([98abb52](https://github.com/dmarcguardhq/dmarcguard/commit/98abb526d856a2c560fabc090a383a206fcec2cc)), closes [#193](https://github.com/dmarcguardhq/dmarcguard/issues/193)

## [1.5.7](https://github.com/dmarcguardhq/dmarcguard/compare/v1.5.6...v1.5.7) (2026-09-05)


### Bug Fixes

* **dashboard:** judge health over delivered mail, in the backend ([#191](https://github.com/dmarcguardhq/dmarcguard/issues/191)) ([53d8944](https://github.com/dmarcguardhq/dmarcguard/commit/53d8944fdc95ed4199e3ad47505e461e74851d5a))
* **dashboard:** show null reports as EMPTY instead of FAIL ([#187](https://github.com/dmarcguardhq/dmarcguard/issues/187)) ([dcf9d98](https://github.com/dmarcguardhq/dmarcguard/commit/dcf9d98cff2ebda51c6157bcb603f0d9550941e2))
* **ingest:** bound report decompression to stop a gzip/zip bomb ([#189](https://github.com/dmarcguardhq/dmarcguard/issues/189)) ([f50ee41](https://github.com/dmarcguardhq/dmarcguard/commit/f50ee4129f8f56fd8a1647f9859978f73e778d41))
* **precommit:** run prettier on changed files only ([729f64d](https://github.com/dmarcguardhq/dmarcguard/commit/729f64db6044da83bd7bfc356b926a8e05d062cf))


### Build & Dependencies

* **deps:** update actions/setup-go action to v7 ([#175](https://github.com/dmarcguardhq/dmarcguard/issues/175)) ([99e4a48](https://github.com/dmarcguardhq/dmarcguard/commit/99e4a48f9e04e854c09b2ecd906ed543131c6432))
* **deps:** update golang docker tag to v1.27 ([#184](https://github.com/dmarcguardhq/dmarcguard/issues/184)) ([afdb202](https://github.com/dmarcguardhq/dmarcguard/commit/afdb202e842fcddaf88d6f54dcbeb779dd4ac9d8))
* **deps:** update module github.com/coreos/go-oidc/v3 to v3.21.0 ([#188](https://github.com/dmarcguardhq/dmarcguard/issues/188)) ([675492d](https://github.com/dmarcguardhq/dmarcguard/commit/675492dfb372adcb03b5b5b56db9c367e85ddcaa))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.7.0 ([#178](https://github.com/dmarcguardhq/dmarcguard/issues/178)) ([f10b023](https://github.com/dmarcguardhq/dmarcguard/commit/f10b0233751ebdaa1e9d94a28906fb69ac417672))
* **deps:** update module github.com/urfave/cli/v3 to v3.11.0 ([#181](https://github.com/dmarcguardhq/dmarcguard/issues/181)) ([e5d3fef](https://github.com/dmarcguardhq/dmarcguard/commit/e5d3fefa373a85687bae2e13421e04308e0a7f44))
* **deps:** update module modernc.org/sqlite to v1.58.0 ([#179](https://github.com/dmarcguardhq/dmarcguard/issues/179)) ([21daa5b](https://github.com/dmarcguardhq/dmarcguard/commit/21daa5b8c797968a5d3c9d51533cc41138082843))
* **deps:** update non-major dependencies ([#174](https://github.com/dmarcguardhq/dmarcguard/issues/174)) ([4bae98f](https://github.com/dmarcguardhq/dmarcguard/commit/4bae98f85398c2d160401c64c4f37a0bdab5f349))
* **deps:** upgrade JS dependencies (ky v2, pinia v4, sharp 0.35, vite 8.2.2) ([#192](https://github.com/dmarcguardhq/dmarcguard/issues/192)) ([d80f128](https://github.com/dmarcguardhq/dmarcguard/commit/d80f1289638e230e2193ab8e1f3b349341679411))


### Chores

* **CI:** enrich release changelog ([5f9f9e8](https://github.com/dmarcguardhq/dmarcguard/commit/5f9f9e85e69a0816176df62f2619e6291168d12b))
* **deps:** lock file maintenance ([#177](https://github.com/dmarcguardhq/dmarcguard/issues/177)) ([b3da7f2](https://github.com/dmarcguardhq/dmarcguard/commit/b3da7f2d067e413323580ef20f6eb4b17b99a64e))

## [1.5.6](https://github.com/dmarcguardhq/dmarcguard/compare/v1.5.5...v1.5.6) (2026-08-18)


### Bug Fixes

* **imap:** unwrap nested message/rfc822 report attachments ([c82f118](https://github.com/dmarcguardhq/dmarcguard/commit/c82f118c375481956cccd352ec00c3275aa01e0a))

## [1.5.5](https://github.com/dmarcguardhq/dmarcguard/compare/v1.5.4...v1.5.5) (2026-07-22)


### Bug Fixes

* **docker:** make /data writable under non-root runtimes ([#171](https://github.com/dmarcguardhq/dmarcguard/issues/171)) ([90a36d3](https://github.com/dmarcguardhq/dmarcguard/commit/90a36d3359edf9dd7e2529b412261c3ecbfe67f7))

## [1.5.4](https://github.com/dmarcguardhq/dmarcguard/compare/v1.5.3...v1.5.4) (2026-06-20)


### Bug Fixes

* accept fetch-interval without prefix ([06fd520](https://github.com/dmarcguardhq/dmarcguard/commit/06fd520d1a816795d8e35f6a7c9b1b5eed5c096c)), closes [#162](https://github.com/dmarcguardhq/dmarcguard/issues/162)
* **CI:** name docker repo explicitly ([0ebbf87](https://github.com/dmarcguardhq/dmarcguard/commit/0ebbf87a6aad51ad735434982734105d79a97893))
* **CI:** publish to ghcr.io as well ([9c0bf28](https://github.com/dmarcguardhq/dmarcguard/commit/9c0bf284a08c50487bb6e8fe475b0248a836cedc))
* **CI:** publish to new and old docker hub repo ([e7eda52](https://github.com/dmarcguardhq/dmarcguard/commit/e7eda52d98e7d3d7d127f199f0f06376007e5cb4))
* **CI:** remove docker repo from old ghcr.io ([d3f24c0](https://github.com/dmarcguardhq/dmarcguard/commit/d3f24c00961876ac7bd8a976a16ada168a59be28))
* **deps:** update module github.com/caarlos0/env/v11 to v11.4.1 ([#155](https://github.com/dmarcguardhq/dmarcguard/issues/155)) ([31398d3](https://github.com/dmarcguardhq/dmarcguard/commit/31398d31413b2c204fb4fd0946bc3d087a490b77))
* **deps:** update module github.com/mattn/go-sqlite3 to v1.14.38 ([#129](https://github.com/dmarcguardhq/dmarcguard/issues/129)) ([40fb4d4](https://github.com/dmarcguardhq/dmarcguard/commit/40fb4d4a9f2048af86f3cf77be12df5a16cf134b))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.6.1 ([#151](https://github.com/dmarcguardhq/dmarcguard/issues/151)) ([eb1fb95](https://github.com/dmarcguardhq/dmarcguard/commit/eb1fb953b6e91199a7de4df959802a2dffdeb9ef))
* move cli init to cmd/ and grab values using Destination syntax ([e04af1a](https://github.com/dmarcguardhq/dmarcguard/commit/e04af1aefd1c33c74780afd2b4e00805a951c1b0))

## [1.5.3](https://github.com/meysam81/parse-dmarc/compare/v1.5.2...v1.5.3) (2026-03-28)


### Bug Fixes

* **deps:** update module github.com/urfave/cli/v3 to v3.8.0 ([#135](https://github.com/meysam81/parse-dmarc/issues/135)) ([dc06765](https://github.com/meysam81/parse-dmarc/commit/dc06765f9d65c0dd67a2debcb12ed3499813ed05))
* **docker:** create DB parent dirs ([e8d6497](https://github.com/meysam81/parse-dmarc/commit/e8d6497227119e4e7e0d348f10386e4b657b630d)), closes [#138](https://github.com/meysam81/parse-dmarc/issues/138)

## [1.5.2](https://github.com/meysam81/parse-dmarc/compare/v1.5.1...v1.5.2) (2026-03-13)


### Bug Fixes

* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.4.1 ([#127](https://github.com/meysam81/parse-dmarc/issues/127)) ([ce63d91](https://github.com/meysam81/parse-dmarc/commit/ce63d9146a8865e1f1bf5869b6502b6f62c30dc5))
* use PEEK when mark-seen is set to false ([#125](https://github.com/meysam81/parse-dmarc/issues/125)) ([e54e269](https://github.com/meysam81/parse-dmarc/commit/e54e26900faed10f916afed47d23b3a2a1f65467))

## [1.5.1](https://github.com/meysam81/parse-dmarc/compare/v1.5.0...v1.5.1) (2026-03-13)


### Bug Fixes

* **CI:** remove mcp entry from goreleaser altogether ([1d769a4](https://github.com/meysam81/parse-dmarc/commit/1d769a45d05ea5afc5340e87fb60cec3fe3514ce))

## [1.5.0](https://github.com/meysam81/parse-dmarc/compare/v1.4.8...v1.5.0) (2026-03-13)


### Features

* allow customizing SEEN & move-folder behavior after processing ([80346df](https://github.com/meysam81/parse-dmarc/commit/80346dfddc56844dd59c9474fb788ca1d0d9312c)), closes [#118](https://github.com/meysam81/parse-dmarc/issues/118)


### Bug Fixes

* **CI:** disable publishing to mcp registry ([c2e3969](https://github.com/meysam81/parse-dmarc/commit/c2e3969e90554479921437be4df0d011b4214a89))
* **CI:** use oxc the default vite minifier ([2cc92e1](https://github.com/meysam81/parse-dmarc/commit/2cc92e1d85734bc1cc34b91b209eb66b13336b5d))
* **deps:** update module github.com/caarlos0/env/v11 to v11.4.0 ([#104](https://github.com/meysam81/parse-dmarc/issues/104)) ([e14e6d2](https://github.com/meysam81/parse-dmarc/commit/e14e6d2d7e8ecbfd8d7e5564c600b165b4135946))
* **deps:** update module github.com/goccy/go-json to v0.10.6 ([#119](https://github.com/meysam81/parse-dmarc/issues/119)) ([7ad51d3](https://github.com/meysam81/parse-dmarc/commit/7ad51d3a0b0395d99b542293ae5dd2a0b1e7ad32))
* **deps:** update module github.com/mattn/go-sqlite3 to v1.14.34 ([#85](https://github.com/meysam81/parse-dmarc/issues/85)) ([5336d3a](https://github.com/meysam81/parse-dmarc/commit/5336d3a64e647c32f46ddeb7d08c76fabd2a1200))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1 ([#49](https://github.com/meysam81/parse-dmarc/issues/49)) ([58c0e79](https://github.com/meysam81/parse-dmarc/commit/58c0e79699392fe6661a8e6cce3092c4dd0d598f))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.3.1 ([#102](https://github.com/meysam81/parse-dmarc/issues/102)) ([68a80a4](https://github.com/meysam81/parse-dmarc/commit/68a80a4ab077af0ce4f45046211aaa77288de42b))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v1.4.0 ([#107](https://github.com/meysam81/parse-dmarc/issues/107)) ([39f1fcb](https://github.com/meysam81/parse-dmarc/commit/39f1fcb833098df1e2b3ab9d89c47350475055e6))
* **deps:** update module github.com/urfave/cli/v3 to v3.6.2 ([#25](https://github.com/meysam81/parse-dmarc/issues/25)) ([955e4e0](https://github.com/meysam81/parse-dmarc/commit/955e4e04d4d1d40617721ac4ec4aa70eaabad4ca))
* **deps:** update module github.com/urfave/cli/v3 to v3.7.0 ([#110](https://github.com/meysam81/parse-dmarc/issues/110)) ([8ce7090](https://github.com/meysam81/parse-dmarc/commit/8ce7090ce717ce25b563d2532b6a1a4abe7555b0))
* **deps:** update module modernc.org/sqlite to v1.45.0 ([#27](https://github.com/meysam81/parse-dmarc/issues/27)) ([2daa3b4](https://github.com/meysam81/parse-dmarc/commit/2daa3b4177bd5298e0f04eea8405fede27446ac2))
* **deps:** update module modernc.org/sqlite to v1.46.1 ([#101](https://github.com/meysam81/parse-dmarc/issues/101)) ([087c4e7](https://github.com/meysam81/parse-dmarc/commit/087c4e79f1fe433f82e2245a654c4e09701b69de))
* **docs:** update imap-host env var ([#116](https://github.com/meysam81/parse-dmarc/issues/116)) ([6d5aa4f](https://github.com/meysam81/parse-dmarc/commit/6d5aa4f32a13956723138ed786f599ad81b229c8))
* **goreleaser:** take build arg for copying binary in dockerfile ([5fb9504](https://github.com/meysam81/parse-dmarc/commit/5fb9504f1d7a946edf0c7c0a62eed13c1abbeb8d))

## [1.4.8](https://github.com/meysam81/parse-dmarc/compare/v1.4.7...v1.4.8) (2026-02-11)


### Bug Fixes

* **CI:** remove incomplete mcp docker from goreleaser ([356bb7b](https://github.com/meysam81/parse-dmarc/commit/356bb7b8ead6c3b783811178e49e10145a876d29))
* do not overwrite release notes from goreleaser ([43bec10](https://github.com/meysam81/parse-dmarc/commit/43bec1038b8a4ee46c6a898c89d918596bab0fdc))
* **mcp:** handle json error ([d17a04f](https://github.com/meysam81/parse-dmarc/commit/d17a04f23994911118f30522ea3ca909fb203efd))
* propagate errors with context ([85014c8](https://github.com/meysam81/parse-dmarc/commit/85014c8aafaf03a6d6b512a7be7eaccab971a4b0))
* stay compliant with brew core version print contraint ([35f8261](https://github.com/meysam81/parse-dmarc/commit/35f826169b3864811cc7684e7c87190518c373a5))

## [1.4.7](https://github.com/meysam81/parse-dmarc/compare/v1.4.6...v1.4.7) (2026-02-01)


### Bug Fixes

* **CI:** update mcp registry version identifier ([a135890](https://github.com/meysam81/parse-dmarc/commit/a135890747ef5fc07e43a7c9ed16e4dc21ac592b))

## [1.4.6](https://github.com/meysam81/parse-dmarc/compare/v1.4.5...v1.4.6) (2026-02-01)


### Bug Fixes

* make linter happy ([900e314](https://github.com/meysam81/parse-dmarc/commit/900e314e9c32909475db562c5b6d4b038496550e))
* **UI:** handle no data in compliance score ([#90](https://github.com/meysam81/parse-dmarc/issues/90)) ([eced435](https://github.com/meysam81/parse-dmarc/commit/eced435ab7b08a73829d4b7b6e09ed31d324834b))

## [1.4.5](https://github.com/meysam81/parse-dmarc/compare/v1.4.4...v1.4.5) (2026-01-25)


### Bug Fixes

* address unhandled charset ([#87](https://github.com/meysam81/parse-dmarc/issues/87)) ([99c245a](https://github.com/meysam81/parse-dmarc/commit/99c245a670ce1cb44f5d08b24b5d0c99c502ec69))

## [1.4.4](https://github.com/meysam81/parse-dmarc/compare/v1.4.3...v1.4.4) (2026-01-17)


### Bug Fixes

* **CI:** remove changelog from release note ([127cb80](https://github.com/meysam81/parse-dmarc/commit/127cb800f3f285b2d5899561fa390e2894bd52de))
* switch env and json config parsing order ([#83](https://github.com/meysam81/parse-dmarc/issues/83)) ([12afcd2](https://github.com/meysam81/parse-dmarc/commit/12afcd2e295c31b60a4696ba5c80edac1ae424e2)), closes [#84](https://github.com/meysam81/parse-dmarc/issues/84)

## [1.4.3](https://github.com/meysam81/parse-dmarc/compare/v1.4.2...v1.4.3) (2025-12-21)


### Features

* allow changing api endpoint from the settings modal ([#76](https://github.com/meysam81/parse-dmarc/issues/76)) ([e6cb48b](https://github.com/meysam81/parse-dmarc/commit/e6cb48be538e5c3e997d86dda9250068a2fb377e))


### Miscellaneous Chores

* prepare next release 1.4.3 ([141e2ed](https://github.com/meysam81/parse-dmarc/commit/141e2ed8c97a34e3d24d359645c0aa67237a47eb))

## [1.4.2](https://github.com/meysam81/parse-dmarc/compare/v1.4.1...v1.4.2) (2025-12-20)


### Bug Fixes

* **CI:** change mcp auth type to gh oidc ([1132b50](https://github.com/meysam81/parse-dmarc/commit/1132b5038c16fda2c83fe90a83f9f6e802c06af4))
* **dev:** minify html with vite plugin ([10eb4f7](https://github.com/meysam81/parse-dmarc/commit/10eb4f73ce936cc7def91d5aa633f2398b19fa33))

## [1.4.1](https://github.com/meysam81/parse-dmarc/compare/v1.4.0...v1.4.1) (2025-12-18)


### Features

* add DigitalOcean Droplet and Dokploy deployment options ([#67](https://github.com/meysam81/parse-dmarc/issues/67)) ([21d494f](https://github.com/meysam81/parse-dmarc/commit/21d494f60d7b7450d29cfd9a50fa25a914b100e8))
* **CI:** add mcp registry publishing to goreleaser ([a0d6f88](https://github.com/meysam81/parse-dmarc/commit/a0d6f88eb54df004c3b35ef7e256bb8d7be36cc9))


### Bug Fixes

* **docs:** add the missing png ([c9937ce](https://github.com/meysam81/parse-dmarc/commit/c9937cea0649a48ed92c973f223a25f8ef112f0f))
* use png for favicon ([4704ad7](https://github.com/meysam81/parse-dmarc/commit/4704ad7c23fbd31913adf677e50f001393a1adda))

## [1.4.0](https://github.com/meysam81/parse-dmarc/compare/v1.3.10...v1.4.0) (2025-12-17)


### Features

* **docs:** add 1-click deployment at the top of README ([#41](https://github.com/meysam81/parse-dmarc/issues/41)) ([50d2ff2](https://github.com/meysam81/parse-dmarc/commit/50d2ff2c54a506801e802daa6b5d375b79d1281d))
* **docs:** add northflank deployment button ([81215a5](https://github.com/meysam81/parse-dmarc/commit/81215a5b5e8cd87aa5f601f22fb460b4929d8d15))
* **docs:** add the self-hosted options to deployments ([42ab966](https://github.com/meysam81/parse-dmarc/commit/42ab966218b9e875ed438a9cb4a4b4ac3adb9b33))
* **docs:** add zeabur deployment template ([4abaa36](https://github.com/meysam81/parse-dmarc/commit/4abaa36d9b077f71a3e6da841af056dffbe046ba))
* **UI:** revamp the dashboard for actionable insights ([#68](https://github.com/meysam81/parse-dmarc/issues/68)) ([dbe073d](https://github.com/meysam81/parse-dmarc/commit/dbe073d35a7f349770e9fab6df4923e221c350f2))


### Bug Fixes

* **deps:** update module github.com/coreos/go-oidc/v3 to v3.17.0 ([#64](https://github.com/meysam81/parse-dmarc/issues/64)) ([306583c](https://github.com/meysam81/parse-dmarc/commit/306583cc37ba2c4fc7e7fffd542cdbab3edd119d))
* **docs:** clean up the providers ([6879ca7](https://github.com/meysam81/parse-dmarc/commit/6879ca70d2a9e5d1adc7f7732d2a85869aecdd2e))
* **docs:** reference dockerhub image ([b2d2b06](https://github.com/meysam81/parse-dmarc/commit/b2d2b0627dbf423d760b5ee9598060334b13f290))
* update docker image reference to dockerhub ([6afcb3a](https://github.com/meysam81/parse-dmarc/commit/6afcb3abe19beda674d19ab7911019c384ba54b1))

## [1.3.10](https://github.com/meysam81/parse-dmarc/compare/v1.3.9...v1.3.10) (2025-12-13)


### Bug Fixes

* **CI:** reverse the digest conditional for provenance ([9c0e165](https://github.com/meysam81/parse-dmarc/commit/9c0e165ec1cfeb58c7c535ec02bd0c81e60ba501))

## [1.3.9](https://github.com/meysam81/parse-dmarc/compare/v1.3.8...v1.3.9) (2025-12-13)


### Bug Fixes

* **CI:** perform keyless docker sign ([1f322d3](https://github.com/meysam81/parse-dmarc/commit/1f322d33999d3586880d82c5ce3570b60157cd21))
* **CI:** update the digest name for attestation ([eb692a8](https://github.com/meysam81/parse-dmarc/commit/eb692a81fb08205990e639ab6c54c8998c705cd0))

## [1.3.8](https://github.com/meysam81/parse-dmarc/compare/v1.3.7...v1.3.8) (2025-12-13)


### Bug Fixes

* **CI:** set up buildx action for multi platform build ([a36e41c](https://github.com/meysam81/parse-dmarc/commit/a36e41c02e6af54297275042e269c66b85545a78))

## [1.3.7](https://github.com/meysam81/parse-dmarc/compare/v1.3.6...v1.3.7) (2025-12-13)


### Bug Fixes

* **CI:** remove annotations from docker build ([b5d0cd6](https://github.com/meysam81/parse-dmarc/commit/b5d0cd650ba0dd780ea95ec35c82940620aad6bb))

## [1.3.6](https://github.com/meysam81/parse-dmarc/compare/v1.3.5...v1.3.6) (2025-12-13)


### Bug Fixes

* **CI:** disable sbom via config ([6b0f3be](https://github.com/meysam81/parse-dmarc/commit/6b0f3be7d79cafe1eeb0930d2f30bf02f53c56af))

## [1.3.5](https://github.com/meysam81/parse-dmarc/compare/v1.3.4...v1.3.5) (2025-12-13)


### Bug Fixes

* **CI:** disable attestation digest for docker build ([95d6665](https://github.com/meysam81/parse-dmarc/commit/95d666543f59b3a890ea8921275f840eebb32f8e))

## [1.3.4](https://github.com/meysam81/parse-dmarc/compare/v1.3.3...v1.3.4) (2025-12-13)


### Bug Fixes

* **CI:** disable sbom on docker ([36e5cc4](https://github.com/meysam81/parse-dmarc/commit/36e5cc4d829efa633fc75879e0a132d076866318))

## [1.3.3](https://github.com/meysam81/parse-dmarc/compare/v1.3.2...v1.3.3) (2025-12-13)


### Bug Fixes

* **CI:** add current dir to docker build context ([254cc6a](https://github.com/meysam81/parse-dmarc/commit/254cc6a6b0b88afc67af8a1cdc4a41d1b71ff1a6))

## [1.3.2](https://github.com/meysam81/parse-dmarc/compare/v1.3.1...v1.3.2) (2025-12-13)


### Features

* **CI:** build docker via goreleaser ([8555265](https://github.com/meysam81/parse-dmarc/commit/855526505aa48c9ab1cdc946f10370a875ee47d8))

## [1.3.1](https://github.com/meysam81/parse-dmarc/compare/v1.3.0...v1.3.1) (2025-12-13)


### Bug Fixes

* **CI:** publish to casks for a change ([0e8346b](https://github.com/meysam81/parse-dmarc/commit/0e8346b3e1c3dcc557445370f59ea7c93588c1fd))

## [1.3.0](https://github.com/meysam81/parse-dmarc/compare/v1.2.2...v1.3.0) (2025-12-12)


### Features

* add MCP server integration to project ([#46](https://github.com/meysam81/parse-dmarc/issues/46)) ([807b8d6](https://github.com/meysam81/parse-dmarc/commit/807b8d677c81aef6ac39605bffb7125680841d4a))
* **CI:** add man page and shell completion to brew ([0e45ed2](https://github.com/meysam81/parse-dmarc/commit/0e45ed2a10d26a06fadc6a2fdfe3b81ad5a01197))
* **CI:** add prettier job ([39af52d](https://github.com/meysam81/parse-dmarc/commit/39af52d186fcf75dba2f8ac05ef313808f1e01b4))
* **CI:** install pandoc and use official goreleaser action ([0d0a307](https://github.com/meysam81/parse-dmarc/commit/0d0a3077395a9879a9d036166b4afbab43ee164e))
* **mcp:** add OAuth2 authentication for MCP HTTP server ([361f078](https://github.com/meysam81/parse-dmarc/commit/361f0782e2f20b7160873a085e600adc68988432))

## [1.2.2](https://github.com/meysam81/parse-dmarc/compare/v1.2.1...v1.2.2) (2025-12-05)


### Bug Fixes

* **CI:** add the brew token ([4081621](https://github.com/meysam81/parse-dmarc/commit/4081621ae9d1771185712512a912aa350a2e0305))

## [1.2.1](https://github.com/meysam81/parse-dmarc/compare/v1.2.0...v1.2.1) (2025-12-05)


### Bug Fixes

* **CI:** move bundled dist to server directory ([59dabf1](https://github.com/meysam81/parse-dmarc/commit/59dabf1becfc22fddda4bfdab6dd22188396b4c3))
* **docs:** pin to full version for now ([2060e45](https://github.com/meysam81/parse-dmarc/commit/2060e451a3b7322d933ad4b8f389f28474f85fb1))
* **docs:** update quickstart command with volume ([617c258](https://github.com/meysam81/parse-dmarc/commit/617c2584402cb309e1a652630f5f75de92df0b2e))

## [1.2.0](https://github.com/meysam81/parse-dmarc/compare/v1.1.0...v1.2.0) (2025-12-05)


### Features

* add production-ready Grafana dashboard ([6f857fe](https://github.com/meysam81/parse-dmarc/commit/6f857fe660144093b8766f529c15b17a21edc3b2))
* add production-ready Prometheus metrics ([#39](https://github.com/meysam81/parse-dmarc/issues/39)) ([17f6968](https://github.com/meysam81/parse-dmarc/commit/17f69681b23eee92ff1bc908fde7da2de11b8a52))
* **frontend:** add DMARC DNS record generator ([34491ba](https://github.com/meysam81/parse-dmarc/commit/34491ba244f9af48a25e3a0d6b09979661678c0d))
* **frontend:** implement dark mode with theme toggle ([#35](https://github.com/meysam81/parse-dmarc/issues/35)) ([569fd36](https://github.com/meysam81/parse-dmarc/commit/569fd3649931f94fa4f39845d0f141afacd6fbc0))


### Bug Fixes

* **CI:** build on major version as well ([fdff51d](https://github.com/meysam81/parse-dmarc/commit/fdff51d2c77d55e890e5b6d9c12ce316ddec2542))
* **CI:** only pin main with latest ([f120b2e](https://github.com/meysam81/parse-dmarc/commit/f120b2e7b7985adbaa815b6104f8aaee44e8e9ce))
* **CI:** reverse the conditional for build-dev job ([ce8b417](https://github.com/meysam81/parse-dmarc/commit/ce8b417c65c62563886787df3871200679c8de66))
* **CI:** update dockerignore after main.go change ([1241b2f](https://github.com/meysam81/parse-dmarc/commit/1241b2f22b8138348b7bf087ba00659e0b1f6fe8))
* **CI:** update version of the binary ([0a8f023](https://github.com/meysam81/parse-dmarc/commit/0a8f0230b30229349db270aa065e3332da63990a))
* **CI:** use go version file when setting up go ([40fbcb1](https://github.com/meysam81/parse-dmarc/commit/40fbcb1dfbe962c7d48ef2c821ee32d997e7c78c))
* **dev:** remove extra compose file ([eb71b09](https://github.com/meysam81/parse-dmarc/commit/eb71b097005bb9d024454ee6d84c2ce5fbb59d7b))
* **dev:** remove user from compose ([e9fbf73](https://github.com/meysam81/parse-dmarc/commit/e9fbf73b72a2a9958d62ed6f6ffea4bad2f2771f))
* make github star non-intrusive ([831047b](https://github.com/meysam81/parse-dmarc/commit/831047b88394cd468ff57bc08c53907880e0feee))
* move main.go to the root ([0dd5787](https://github.com/meysam81/parse-dmarc/commit/0dd57878ec5f7d6707faca0189834a6d91c266e5))
* publish to grafana community dashboards instead ([86f40b0](https://github.com/meysam81/parse-dmarc/commit/86f40b052607861de3d1a843f9e9c9e08ba57481))
* remove broken dns generator ([df39352](https://github.com/meysam81/parse-dmarc/commit/df393527710b487bc0bdbe7d67d9a3f909623b92))

## [1.1.0](https://github.com/meysam81/parse-dmarc/compare/v1.0.2...v1.1.0) (2025-11-08)


### Features

* **docs:** add nerdy badges to README ([#23](https://github.com/meysam81/parse-dmarc/issues/23)) ([5026a45](https://github.com/meysam81/parse-dmarc/commit/5026a45da22a9fa0a3c831945da4afaa05c2dd7c))


### Bug Fixes

* **CI:** update goreleaser after moving FE to root ([357d37e](https://github.com/meysam81/parse-dmarc/commit/357d37ed69d9cd7400c05c2f7ca2667072f4d46d))

## [1.0.2](https://github.com/meysam81/parse-dmarc/compare/v1.0.1...v1.0.2) (2025-11-07)


### Bug Fixes

* **dev:** bring the frontend to the root and simplify docker ([edc64e7](https://github.com/meysam81/parse-dmarc/commit/edc64e7417c910b7200bb1d727f0ba200c1a787d))

## [1.0.1](https://github.com/meysam81/parse-dmarc/compare/v1.0.0...v1.0.1) (2025-11-07)


### Bug Fixes

* **docs:** simplify README for first-time viewers ([#15](https://github.com/meysam81/parse-dmarc/issues/15)) ([7341b7b](https://github.com/meysam81/parse-dmarc/commit/7341b7bb022e5e96798e6d8d3523a377e282cf6f))

## 1.0.0 (2025-11-05)


### Features

* add CI and goreleaser with optimized distroless docker ([66c6828](https://github.com/meysam81/parse-dmarc/commit/66c682807d5ed12328349e2713c9d814d6337e78))
* add sort and refresh to UI ([9d6853e](https://github.com/meysam81/parse-dmarc/commit/9d6853e95319529c40958ae60d95d05d8cb5a675))
* clean up the cli and optimize the docker ([db23ba4](https://github.com/meysam81/parse-dmarc/commit/db23ba4ff2fd80ad1c016daa73a57519998dd271))
* **docs:** update readme to the latest changes ([c20e0cf](https://github.com/meysam81/parse-dmarc/commit/c20e0cf837b91b05e61153e4e55e7315fb09d396))
* ensure db path exists and add demo screenshot ([42d6d14](https://github.com/meysam81/parse-dmarc/commit/42d6d14fd96f2818830c73ad4e23b8f97719c497))
* Implement DMARC report parser with Vue.js dashboard ([#1](https://github.com/meysam81/parse-dmarc/issues/1)) ([23b8ac0](https://github.com/meysam81/parse-dmarc/commit/23b8ac0cc63a81aad7dfa2385ed24f6d74915775))
* update the UI footer with OS friendly text ([5008fdf](https://github.com/meysam81/parse-dmarc/commit/5008fdf5e59e5b4af1fc12ccafc644d27ccdbdfc))
* update UI footer & readme and optimize vite ([9578d05](https://github.com/meysam81/parse-dmarc/commit/9578d0530adab10334ad1a8a20c563abdc589854))


### Bug Fixes

* **CI:** ensure a fake dist exists before golangci-lint ([4e01983](https://github.com/meysam81/parse-dmarc/commit/4e01983f13a16ab0a8ea7cb880e43b39c285ba1f))
* **CI:** make linter happy ([cd04dfc](https://github.com/meysam81/parse-dmarc/commit/cd04dfc8beae30b0405d38834680011cd7d14fe5))
* make server context aware and do not fail UI on empty result ([588c466](https://github.com/meysam81/parse-dmarc/commit/588c4660534d4111cf24738a952cbad56ebbf345))
* update the UI on reports API parsing ([b85b59e](https://github.com/meysam81/parse-dmarc/commit/b85b59ef577ffeb6b4dacd1188e6e55d549b1b5a))
* use app for static config ([3820862](https://github.com/meysam81/parse-dmarc/commit/3820862f8fffe26fd3dad7c004d1c0b04fc8be9e))
