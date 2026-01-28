FROM node:20

ARG TZ
ENV TZ="$TZ"

ARG CLAUDE_CODE_VERSION=latest

# Install basic development tools and iptables/ipset
RUN apt-get update && apt-get install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  dnsmasq \
  aggregate \
  jq \
  nano \
  vim \
  wget \
  postgresql-client \
  # Playwright/Chromium dependencies for agent-browser
  libnss3 \
  libnspr4 \
  libatk1.0-0 \
  libatk-bridge2.0-0 \
  libcups2 \
  libxkbcommon0 \
  libatspi2.0-0 \
  libxcomposite1 \
  libxdamage1 \
  libxfixes3 \
  libxrandr2 \
  libgbm1 \
  libasound2 \
  libpango-1.0-0 \
  libcairo2 \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install Go
ARG GO_VERSION=1.24.1
RUN ARCH=$(dpkg --print-architecture) && \
  if [ "$ARCH" = "amd64" ]; then GO_ARCH="amd64"; \
  elif [ "$ARCH" = "arm64" ]; then GO_ARCH="arm64"; \
  else GO_ARCH="$ARCH"; fi && \
  wget "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" && \
  tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" && \
  rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"

ENV PATH=$PATH:/usr/local/go/bin:/home/node/go/bin
ENV GOPATH=/home/node/go

# Install golangci-lint
ARG GOLANGCI_LINT_VERSION=1.64.5
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${ARCH}.tar.gz" && \
  tar -xzf "golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${ARCH}.tar.gz" && \
  mv "golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${ARCH}/golangci-lint" /usr/local/bin/ && \
  rm -rf "golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${ARCH}" "golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${ARCH}.tar.gz"

# Install golang-migrate for database migrations
ARG GOLANG_MIGRATE_VERSION=4.18.1
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/golang-migrate/migrate/releases/download/v${GOLANG_MIGRATE_VERSION}/migrate.linux-${ARCH}.tar.gz" && \
  tar -xzf "migrate.linux-${ARCH}.tar.gz" && \
  mv migrate /usr/local/bin/ && \
  rm "migrate.linux-${ARCH}.tar.gz"

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.claude && \
  chown -R node:node /workspace /home/node/.claude

WORKDIR /workspace

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Install beans (issue tracker)
ARG BEANS_VERSION=0.3.4
RUN ARCH=$(dpkg --print-architecture) && \
  if [ "$ARCH" = "amd64" ]; then BEANS_ARCH="x86_64"; \
  elif [ "$ARCH" = "arm64" ]; then BEANS_ARCH="arm64"; \
  else BEANS_ARCH="$ARCH"; fi && \
  wget "https://github.com/hmans/beans/releases/download/v${BEANS_VERSION}/beans_Linux_${BEANS_ARCH}.tar.gz" && \
  tar -xzf "beans_Linux_${BEANS_ARCH}.tar.gz" && \
  mv beans /usr/local/bin/ && \
  rm "beans_Linux_${BEANS_ARCH}.tar.gz"

# Enable Corepack for pnpm support (uses packageManager field in package.json)
RUN corepack enable

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Set the default editor and visual
ENV EDITOR=nano
ENV VISUAL=nano

# Default powerline10k theme
ARG ZSH_IN_DOCKER_VERSION=1.2.0
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install Claude
RUN npm install -g @anthropic-ai/claude-code@${CLAUDE_CODE_VERSION}

# Install agent-browser for browser automation
RUN npm install -g agent-browser && \
  npx playwright install chromium


# Copy and set up firewall scripts
COPY init-firewall.sh /usr/local/bin/
COPY dynamic-domains.conf /usr/local/bin/
COPY warmup-dns.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh /usr/local/bin/warmup-dns.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
