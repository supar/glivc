.PHONY: all build clean

NAME=glivc

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
	go build -o ./${NAME} -v $(LDFLAGS)
