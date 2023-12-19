ARG BASE_IMAGE_VERSION="37"

FROM fedora:${BASE_IMAGE_VERSION}

ARG OCM_CLI_VERSION
ARG BACKPLANE_CLI_VERSION

RUN dnf update -y && \
    dnf install -y procps \
    wget \
    golang \
    sudo \
    jq \
    python3 \
    python-pip \
    git && \
    yum install -y net-tools \
    make && \
    curl -Lo /usr/bin/ocm https://github.com/openshift-online/ocm-cli/releases/download/v${OCM_CLI_VERSION}/ocm-linux-amd64 && \
    chmod +x /usr/bin/ocm && \
    wget https://github.com/openshift/backplane-cli/releases/download/v${BACKPLANE_CLI_VERSION}/ocm-backplane_${BACKPLANE_CLI_VERSION}_Linux_x86_64.tar.gz && \
    tar -xvf ocm-backplane_${BACKPLANE_CLI_VERSION}_Linux_x86_64.tar.gz && \
    mv $PWD/ocm-backplane /usr/bin/ocm-backplane && \
    chmod +x /usr/bin/ocm-backplane && \
    wget https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz && \
    tar -xvf openshift-client-linux.tar.gz && \
    mv $PWD/oc /usr/bin/oc && \
    mv $PWD/kubectl /usr/bin/kubectl

ARG BUILD_SHA=
ENV BUILD_SHA=${BUILD_SHA}

RUN mkdir -p /hc
WORKDIR /hc
ADD ./terminal ./terminal
ADD ./hc .
RUN chmod +x ./hc && cp ./hc /usr/bin/hc

