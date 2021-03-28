FROM alpine
COPY app /
ENTRYPOINT ["/app", "--logtostderr=true"]