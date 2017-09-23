FROM scratch
COPY iugaserver iugaserver
ENTRYPOINT [ "/iugaserver" ]
