#!/bin/bash
# 
# Eleven server init.
# 
# This is the first script to run during 
# the creation of the server (via cloud-init).
# 
# See: https://community.hetzner.com/tutorials/basic-cloud-config
# 
# In a nutshell, this script:
#   - create and configure the user "eleven" (notably the SSH access)
#   - configure and install the Eleven agent
#
# The next steps are assured by the Eleven agent via gRPC through SSH.
set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

log () {
  echo -e "${1}" >&2
}

log "\n\n"
log "---- Eleven server init (start) ----"
log "\n\n"

# Remove "debconf: unable to initialize frontend: Dialog" warnings
echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

# We use "jq" in our exit trap 
# and "curl" to download the Eleven agent
apt-get --assume-yes --quiet --quiet update
apt-get --assume-yes --quiet --quiet install jq curl

constructExitJSONResponse () {
  JSON_RESPONSE=$(jq --null-input \
  --arg exitCode "${1}" \
  --arg sshHostKeys "${2}" \
  --arg cloudInitLogs "${3}" \
  '{"exit_code": $exitCode, "ssh_host_keys": $sshHostKeys, "cloud_init_logs": $cloudInitLogs}')

  echo "${JSON_RESPONSE}"
}

ELEVEN_SSH_SERVER_HOST_KEY_FILE_PATH="/home/eleven/.ssh/eleven-ssh-server-host-key.pub"
ELEVEN_INIT_RESULTS_FILE_PATH="/tmp/eleven-init-results"

handleExit () {
  EXIT_CODE=$?
  CLOUD_INIT_LOGS="$(cat /var/log/cloud-init-output.log)"

  rm --force "${ELEVEN_INIT_RESULTS_FILE_PATH}"

  log "\n\n"
  if [[ "${EXIT_CODE}" != 0 ]]; then
    constructExitJSONResponse "${EXIT_CODE}" "" "${CLOUD_INIT_LOGS}" >> "${ELEVEN_INIT_RESULTS_FILE_PATH}"
    log "---- Eleven server init (failed) (exit code ${EXIT_CODE}) ----"
  else
    SSH_HOST_KEYS="$(cat "${ELEVEN_SSH_SERVER_HOST_KEY_FILE_PATH}")"
    constructExitJSONResponse "${EXIT_CODE}" "${SSH_HOST_KEYS}" "${CLOUD_INIT_LOGS}" >> "${ELEVEN_INIT_RESULTS_FILE_PATH}"
    
    log "---- Eleven server init (success) ----"
  fi
  log "\n\n"

  exit "${EXIT_CODE}"
}

trap 'handleExit' EXIT

# -- System configuration

# Looking up the server architecture to
# download the corresponding Eleven agent binary.
# See below.
SERVER_ARCH=""
case $(uname -m) in
  i386)       SERVER_ARCH="386" ;;
  i686)       SERVER_ARCH="386" ;;
  x86_64)     SERVER_ARCH="amd64" ;;
  arm)        dpkg --print-architecture | grep -q "arm64" && SERVER_ARCH="arm64" || SERVER_ARCH="armv6" ;;
  aarch64_be) SERVER_ARCH="arm64" ;;
  aarch64)    SERVER_ARCH="arm64" ;;
  armv8b)     SERVER_ARCH="arm64" ;;
  armv8l)     SERVER_ARCH="arm64" ;;
esac

# -- Creating the user "eleven"

log "Creating user \"eleven\""

groupadd --force eleven
id -u eleven >/dev/null 2>&1 || useradd --gid eleven --home /home/eleven --create-home --shell /bin/bash eleven

if [[ ! -f "/etc/sudoers.d/eleven" ]]; then
  echo "eleven ALL=(ALL) NOPASSWD:ALL" | tee /etc/sudoers.d/eleven > /dev/null
fi

# If the user "eleven" is updated after 
# the agent has started, we will be 
# forced to restart it to get the new permissions.
# To avoid that, we add the group "docker" here even 
# if Docker is not asked as a runtime.
groupadd --force docker
usermod --append --groups docker eleven

# -- Creating the Eleven configuration directory

log "Creating Eleven configuration directory"

# Needed by the Eleven agent 
# to store the gRPC server socket.
# These two variables are replaced at runtime.
mkdir --parents "${ELEVEN_CONFIG_DIR}"
mkdir --parents "${ELEVEN_AGENT_CONFIG_DIR}"

chown --recursive eleven:eleven "${ELEVEN_CONFIG_DIR}"

chmod 700 "${ELEVEN_CONFIG_DIR}"
chmod 700 "${ELEVEN_AGENT_CONFIG_DIR}"

# -- Configuring SSH access for the user "eleven"

log "Configuring SSH access for user \"eleven\""

# We want the user "eleven" to be able to 
# connect through SSH via the generated SSH key.
# See below.
SERVER_SSH_PUBLIC_KEY="$(cat /root/.ssh/authorized_keys)"

# Run as "eleven"
sudo --set-home --login --user eleven -- env \
	SERVER_SSH_PUBLIC_KEY="${SERVER_SSH_PUBLIC_KEY}" \
bash << 'EOF'

mkdir --parents .ssh
chmod 700 .ssh

if [[ ! -f ".ssh/eleven-ssh-server-host-key" ]]; then
  ssh-keygen -t ed25519 -f .ssh/eleven-ssh-server-host-key -q -N ""
fi

chmod 644 .ssh/eleven-ssh-server-host-key.pub
chmod 600 .ssh/eleven-ssh-server-host-key

if [[ ! -f ".ssh/authorized_keys" ]]; then
  echo "${SERVER_SSH_PUBLIC_KEY}" >> .ssh/authorized_keys
fi

chmod 600 .ssh/authorized_keys

EOF

# -- Installing the Eleven agent
#
# /!\ The SSH server host key ("eleven-ssh-server-host-key")
#     needs to be generated. See above.
#
# /!\ The Eleven configuration directory needs to be created 
#     for the agent to be able to create the gRPC server socket. 
#     See above.

log "Installing the Eleven agent"

ELEVEN_AGENT_VERSION="0.0.2"
ELEVEN_AGENT_TMP_ARCHIVE_PATH="/tmp/eleven-agent.tar.gz"
ELEVEN_AGENT_NAME="eleven-agent"
ELEVEN_AGENT_DIR="/usr/local/bin"
ELEVEN_AGENT_PATH="${ELEVEN_AGENT_DIR}/${ELEVEN_AGENT_NAME}"
ELEVEN_AGENT_FOREVER_PATH="${ELEVEN_AGENT_DIR}/forever"
ELEVEN_AGENT_SYSTEMD_SERVICE_NAME="eleven-agent.service"

if [[ ! -f "${ELEVEN_AGENT_PATH}" ]]; then
  # curl --fail --silent --show-error --location --header "Accept: application/octet-stream" https://api.github.com/repos/eleven-sh/agent/releases/assets/81893338 --output "${ELEVEN_AGENT_PATH}"
  rm --recursive --force "${ELEVEN_AGENT_TMP_ARCHIVE_PATH}"
  curl --fail --silent --show-error --location --header "Accept: application/octet-stream" "https://github.com/eleven-sh/agent/releases/download/v${ELEVEN_AGENT_VERSION}/agent_${ELEVEN_AGENT_VERSION}_linux_${SERVER_ARCH}.tar.gz" --output "${ELEVEN_AGENT_TMP_ARCHIVE_PATH}"
  tar --directory "${ELEVEN_AGENT_DIR}" --extract --file "${ELEVEN_AGENT_TMP_ARCHIVE_PATH}"
  rm --recursive --force "${ELEVEN_AGENT_TMP_ARCHIVE_PATH}"
fi

chmod +x "${ELEVEN_AGENT_PATH}"

if [[ ! -f "/etc/systemd/system/${ELEVEN_AGENT_SYSTEMD_SERVICE_NAME}" ]]; then
  tee /etc/systemd/system/"${ELEVEN_AGENT_SYSTEMD_SERVICE_NAME}" > /dev/null << EOF
[Unit]
Description=the agent used to establish connection with the Eleven CLI.

[Service]
Type=simple
ExecStart=${ELEVEN_AGENT_PATH}
WorkingDirectory=${ELEVEN_AGENT_DIR}
Restart=always
User=eleven
Group=eleven

[Install]
WantedBy=multi-user.target
EOF
fi

systemctl enable "${ELEVEN_AGENT_SYSTEMD_SERVICE_NAME}"
systemctl start "${ELEVEN_AGENT_SYSTEMD_SERVICE_NAME}"

if [[ ! -f "${ELEVEN_AGENT_FOREVER_PATH}" ]]; then
  tee "${ELEVEN_AGENT_FOREVER_PATH}" > /dev/null << EOF
#!/bin/bash
set -euo pipefail

${ELEVEN_AGENT_PATH} forever \$@
EOF
fi

chmod +x "${ELEVEN_AGENT_FOREVER_PATH}"
