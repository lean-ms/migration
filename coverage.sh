export CC_TEST_REPORTER_URL=https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64
export CC_TEST_REPORTER_ID=584dbd5c892f225f70261e9af11ab6dd5e363bb3a3b1d02109993300faee8ecc
curl -L $CC_TEST_REPORTER_URL > /tmp/cc-test-reporter
chmod +x /tmp/cc-test-reporter
/tmp/cc-test-reporter before-build
go test -coverprofile c.out ./...
/tmp/cc-test-reporter after-build --exit-code $?