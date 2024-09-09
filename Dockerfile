FROM golang:1.22-alpine AS fetch-stage
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# 
# RUN templ generate
# RUN go build -o /server .

# Generate
FROM ghcr.io/a-h/templ AS generate-stage
WORKDIR /appstage
COPY --chown=65532:65532 . .
RUN ["templ", "generate"]

FROM golang:1.22-alpine AS build-stage
WORKDIR /appbuild
COPY --chown=user:user --from=generate-stage /appstage .
# RUN apk --no-cache add sqlite
# # CMD ["sqlite3" "-init database.sql books.db"]
# RUN sqlite3 -init database.sql books.db	 ""
RUN go build -o ./server


FROM gcr.io/distroless/static-debian11 AS deploy-stage
# FROM golang:1.22 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /appbuild/server /app/server
ENV PORT=3000
EXPOSE ${PORT}
CMD ["/app/server"]