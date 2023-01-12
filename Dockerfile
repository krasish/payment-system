FROM golang:1.19.5 as builder

ENV BASE_APP_DIR /go/src/github.com/krasish/payment-system
WORKDIR ${BASE_APP_DIR}

# Download dependencies
COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

# Copy files
COPY . .

# Build app
RUN go build -v -o main ./cmd/
RUN mkdir /app && mv ./main /app/main \
    && mv ./internal/views /app/view-templates

FROM golang:1.19.5
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/view-templates /app/view-templates

ENV APP_VIEW_TEMPLATES_PATH /app/view-templates

CMD ["/app/main"]