# Generate tarball with new build of osx_vpn_manager
#
# NOTE: OSX only
VERSION=$$(cat main.go | grep -i "cliVersion =" | awk {'print$$3'} | tr -d '"')


all: clean build compress report

clean:
	@rm -f /tmp/osx_vpn_manager-*.tar.gz
	@rm -f vpn

build:
	@echo Building osx_vpn_manager version $(VERSION)
	@go build -o vpn

compress:
	@tar czf /tmp/osx_vpn_manager-$(VERSION).tar.gz ./vpn

report:
	@rm -f vpn
	@shasum -a 256 /tmp/osx_vpn_manager-$(VERSION).tar.gz

.PHONY: all clean build

