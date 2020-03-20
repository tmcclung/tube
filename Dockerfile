# Build
FROM prologic/go-builder:latest AS build

# Runtime
FROM alpine

COPY --from=build /src/tube /tube

ENTRYPOINT ["/tube"]
CMD [""]
