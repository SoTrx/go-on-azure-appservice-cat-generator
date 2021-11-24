# Dockerfile for an appservice go container
####################################################################################################
## Builder
####################################################################################################
FROM golang:alpine as builder
WORKDIR /app 
COPY . .
# Build the app, strip it (LDFLAGS) and optimize it with UPX
RUN apk add upx && \
    cd src && \ 
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" . && \
    upx --best --lzma cat-generator

####################################################################################################
## Final image
####################################################################################################
FROM alpine as release
# Copy the config files. Chmod are required is these can be copied from Windows
# which has no ACL
COPY --chmod=500 --chown=root:root init_container.sh /bin/
COPY --chmod=600 --chown=sshd:sshd sshd_config /etc/ssh/
COPY nginx.conf /etc/nginx/nginx.conf
# Add Runtime deps. Openrc/openssh for sshd, nginx as a reverse-proxy
# and su-exec to step down from root at runtime
ENV sysdirs="/bin   /etc    /lib    /sbin"
RUN apk add --no-cache openrc openssh nginx su-exec &&\
    # As weird at it seems, the root password MUST be "Docker!" to allow \
    # for SSH connections from the Azure portal
    echo "root:Docker!" | chpasswd &&\
    # The runtime user, having no home dir nor password
    adduser -HD -s /bin/ash appuser &&\
    # Generate SSH keys pairs. On Windows, the authorized cyphers may not work
    cd /etc/ssh && \
    ssh-keygen -A &&\
    cd - && \
    # Hardening part. As the root password is predefined, let's prevent the user
    # from root access some other ways \
    # First, remove all packages confs
    find $sysdirs -xdev -regex '.*apk.*' -exec rm -fr {} + && \
    # Next, ensure all system directories are owned by root and root only
    find $sysdirs -xdev -type d \
    -exec chown root:root {} \; \
    -exec chmod 0755 {} \; && \
    # Remove all SUID (files that can be exec with the ACL of another user)
    find $sysdirs -xdev -type f -a -perm +4000 -delete && \
    # Finally, remove all ACL-related programs
    find $sysdirs -xdev \( \
    -name hexdump -o \
    -name chgrp -o \
    -name chmod -o \
    -name chown -o \
    -name ln -o \
    -name od -o \
    -name strings -o \
    -name su \
    \) -delete

WORKDIR /app
# Copy the built app, only allowing our app user to execute it
COPY --from=builder --chmod=0500 --chown=appuser:appuser  /app/src/cat-generator ./

WORKDIR /home/site/wwwroot
EXPOSE 80 2222
ENTRYPOINT [ "/bin/init_container.sh" ]