FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-jenkins"]
COPY baton-jenkins /