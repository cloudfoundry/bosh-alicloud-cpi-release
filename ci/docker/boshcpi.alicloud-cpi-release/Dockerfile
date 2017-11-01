FROM ubuntu:15.04

# Update base image
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8
RUN locale-gen en_US.UTF-8; dpkg-reconfigure locales
RUN apt-get update; apt-get -y upgrade; apt-get clean

# Install dependencies
RUN apt-get -y install build-essential git-core curl tar make wget jq sudo; apt-get clean
RUN apt-get install -y sqlite3 libsqlite3-dev; apt-get clean
RUN apt-get install -y mysql-client libmysqlclient-dev; apt-get clean
RUN apt-get install -y libpq-dev; apt-get clean
RUN apt-get install -y python python-pip libyaml-dev libpython-dev; apt-get clean

RUN apt-get install -y --no-install-recommends \
    ca-certificates \
    openssh-client \
    realpath \
    apt-get clean

# Install NodeJs
RUN curl -sL https://deb.nodesource.com/setup_4.x | bash -; apt-get install -y nodejs; apt-get clean; rm -f /usr/bin/node; ln -s /usr/bin/nodejs /usr/bin/node

# Install Golang
ENV GOLANG_VERSION 1.8.1
RUN curl -sSL https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV GOROOT /usr/local/go
ENV PATH $PATH:$GOROOT/bin

# Import Postmodern PGP key
RUN wget -nv https://raw.github.com/postmodern/postmodern.github.io/master/postmodern.asc && \
    gpg --import postmodern.asc && \
    gpg --fingerprint 0xB9515E77 | grep 'Key fingerprint = 04B2 F3EA 6541 40BC C7DA  1B57 54C3 D9E9 B951 5E77' && \
    if [ "$?" != "0" ]; then echo "Invalid PGP key!"; exit 1; fi

# chruby
RUN /bin/bash -l -c " \
    curl https://codeload.github.com/postmodern/chruby/tar.gz/v0.3.9 | tar -xz -C /tmp/ && \
    pushd /tmp/chruby-0.3.9 && \
    ./scripts/setup.sh && \
    popd && \
    rm -rf /tmp/chruby-0.3.9 \
"

# ruby-install
RUN /bin/bash -l -c " \
    curl https://codeload.github.com/postmodern/ruby-install/tar.gz/v0.5.0 | tar -xz -C /tmp/ && \
    pushd /tmp/ruby-install-0.5.0 && \
    make install && \
    popd && \
    rm -rf /tmp/ruby-install-0.5.0 \
"

# ruby
ENV RUBY_VERSION 2.2.4
RUN ruby-install ruby ${RUBY_VERSION}

# Bundler and BOSH CLI
RUN /bin/bash -l -c "                            \
  source /etc/profile.d/chruby.sh ;              \
  chruby ${RUBY_VERSION} ;                       \
  gem install bundler bosh_cli --no-ri --no-rdoc \
"

# receipt generator
RUN cd /tmp && \
    curl -o certify-artifacts -L https://s3.amazonaws.com/bosh-certification-generator-releases/certify-artifacts-linux-amd64 && \
    chmod +x certify-artifacts && \
    mv certify-artifacts /bin/certify-artifacts