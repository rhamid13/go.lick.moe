FROM golang

# copy source files
ADD . /go/src/gitlab.com/drylemon/go.lick.moe

# working directory
WORKDIR /go/src/gitlab.com/drylemon/go.lick.moe

# install dependencies
RUN go get github.com/gorilla/mux
RUN go get github.com/mattn/go-sqlite3

# compile, drop it in /go/bin
RUN go build -o /go/bin/lick.moe

# run lick.moe when the container starts
ENTRYPOINT /go/bin/lick.moe

# document that the service listens on 8000
EXPOSE 8000
