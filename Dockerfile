# The 'runtime' container only contains ssl cert chain, and is otherwise an empty base image.
# If you find this too bare, you can instead switch to using a normal base image.
FROM scratch


# Add ssl certs
ADD ./build/ca-certificates.crt /etc/ssl/certs/

# Add binary
ADD ./build/tmp/github.com/underarmour/libra /bin/libra

# Expose service app ports
EXPOSE 8080

# Start the service app. Note we have to use the array style because this container does not include /bin/sh
ENTRYPOINT ["/bin/libra"]

ADD ./example /etc/libra