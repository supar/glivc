.PHONY: all build clean

NAME=glivc
VERSION=0.1.1
SOURCES = $(shell ls -AB | grep -i 'makefile$$\|\.go$$')

# Debian build root
DEB_DIR = $(shell pwd)/build/debian
DEBIAN = $(DEB_DIR)/$(NAME)-$(VERSION)/debian
DEB_CONF = $(DEBIAN)/glivc.conf
DEB_SYS_CONF = $(DEBIAN)/glivc.service
DEB_SOURCE = $(DEB_DIR)/$(NAME)_$(VERSION).orig.tar.gz
DEB = $(DEB_DIR)/$(NAME)_$(VERSION)-$(DEBRELEASE)_$(DEBARCH).deb

default:
	@echo "Usage: make task"
	@echo "\t devdeps - install go packages for this project"
	@echo "\t build - build application"


clean:
	@test ! -e ./${NAME} || rm ./${NAME}

devdeps:
	go get -v github.com/astaxie/beego/logs
	go get -v github.com/go-martini/martini
	go get -v gopkg.in/libgit2/git2go.v22

build: clean
	go build -o ./${NAME} -v

deb: cleandeb $(DEB)

cleandeb:
	@rm -rf $(DEB_DIR)

$(DEBIAN): contribs/debian
	mkdir -p $(DEBIAN)
	cp -ad $</* $@/

$(DEB_SOURCE): $(SOURCES)
	mkdir -p $(@D)
	tar --transform "s,^,$(NAME)-$(VERSION)/src/$(NAME)/," -f $@ -cz $^

$(DEB_CONF): contribs/conf/glivc.conf
	mkdir -p $(@D)
	cp -ad $< $@

$(DEB_SYS_CONF): contribs/systemd/glivc.service
	mkdir -p $(@D)
	cp -ad $< $@

$(DEB): $(DEBIAN) $(DEB_SOURCE) $(DEB_CONF) $(DEB_SYS_CONF)
	cd $(DEB_DIR)/$(NAME)-$(VERSION) && \
	debuild -us -uc -b
