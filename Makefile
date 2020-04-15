NAME=chromium-downloader
BINDIR=bin
VERSION=$(shell git describe --tags || echo "unknown version")
BUILDTIME=$(shell date -u)
GOBUILD=CGO_ENABLED=0 go build -ldflags '-X "main.Version=$(VERSION)" \
		-X "main.Buildtime=$(BUILDTIME)" \
		-w -s'

PLATFORM_LIST = \
	darwin-amd64 \
	linux-amd64

WINDOWS_ARCH_LIST = \
	windows-amd64

all: linux-amd64 darwin-amd64 windows-amd64

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

gz_releases=$(addsuffix .tar, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(gz_releases): %.tar : %
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	tar -zcf $(BINDIR)/$(NAME)-$(VERSION)-$(basename $@).tar.gz $(BINDIR)/$(NAME)-$(basename $@) --remove-files

$(zip_releases): %.zip : %
	zip -m -j $(BINDIR)/$(NAME)-$(VERSION)-$(basename $@).zip $(BINDIR)/$(NAME)-$(basename $@).exe

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

releases: $(gz_releases) $(zip_releases)
clean:
	rm $(BINDIR)/*