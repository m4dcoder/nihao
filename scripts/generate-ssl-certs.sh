#!/bin/bash
set -e

function usage() {
    echo "Usage: $0 [-a app name] [-c common name] [-s SAN]... [-i IP]... [-d output directory]" >&2
}

APP_NAME=rbx
COMMON_NAME=""
OUTPUT_DIR=`pwd`

while getopts "a:c:s:i:d:" o; do
    case "${o}" in
        a)
            APP_NAME=${OPTARG}
            ;;
        c)
            COMMON_NAME=${OPTARG}
            ;;
        s)
            SANS+=("${OPTARG}")
            ;;
        i)
            IP_SANS+=("${OPTARG}")
            ;;
        d)
            OUTPUT_DIR=${OPTARG}
            mkdir -p ${OUTPUT_DIR} || true
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            usage
            exit 2
            ;;
        :)
            echo "Option -$OPTARG requires an argument." >&2
            usage
            exit 2
            ;;
    esac
done

if [[ -z "${COMMON_NAME}" ]]; then
    echo "ERROR: The common name for the cert is not provided."
    exit 1
fi

echo "COMMON NAME: ${COMMON_NAME}"
echo ""

echo "SUBJECT ALTERNATIVE NAMES:"
for SAN in "${SANS[@]}"; do echo "- ${SAN}"; done
for IP_SAN in "${IP_SANS[@]}"; do echo "- ${IP_SAN}"; done
echo ""

SSL_CONFIG_FILE="${OUTPUT_DIR}/${APP_NAME}.cnf"
CA_KEY_FILE="${OUTPUT_DIR}/cacert.key"
CA_CRT_FILE="${OUTPUT_DIR}/cacert.crt"
CA_PEM_FILE="${OUTPUT_DIR}/cacert.pem"
SSL_KEY_FILE="${OUTPUT_DIR}/${APP_NAME}.key"
SSL_CRT_FILE="${OUTPUT_DIR}/${APP_NAME}.crt"
SSL_PKCS12_CERT_FILE="${OUTPUT_DIR}/${APP_NAME}.pfx"
RSA_SIZE=2048

# CREATE THE SSL CONFIG FILE
SSL_CONFIG_CONTENT="[req]
distinguished_name = req_distinguished_name
[req_distinguished_name]
[v3_ca]
basicConstraints = critical, CA:TRUE
keyUsage = critical, digitalSignature, keyEncipherment, keyCertSign
[v3_req_server]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1=${COMMON_NAME}
"

SAN_IDX=2
for SAN in "${SANS[@]}"; do
    SSL_CONFIG_CONTENT+="DNS.${SAN_IDX}=${SAN}\n";
    SAN_IDX=$((SAN_IDX + 1))
done

IP_SAN_IDX=1
for IP_SAN in "${IP_SANS[@]}"; do
    SSL_CONFIG_CONTENT+="IP.${IP_SAN_IDX}=${IP_SAN}\n";
    IP_SAN_IDX=$((IP_SAN_IDX + 1))
done

echo -e "${SSL_CONFIG_CONTENT}" > "${SSL_CONFIG_FILE}"

# MANUALLY TOUCH THE RND FILE
touch ${HOME}/.rnd

# GENERATE THE CA CERT
echo "Generating the CA cert..."
openssl genrsa -out ${CA_KEY_FILE} ${RSA_SIZE}
chmod 600 ${CA_KEY_FILE}
openssl req -x509 -new -sha256 -nodes \
    -key ${CA_KEY_FILE} -days 365 -out ${CA_CRT_FILE} \
    -subj "/CN=TohaHeavyIndustries-CA" -extensions v3_ca \
    -config "${SSL_CONFIG_FILE}"
openssl x509 -in ${CA_CRT_FILE} -out ${CA_PEM_FILE} -outform PEM

# GENERATE THE SSL CERT
echo "Generating the SSL cert..."
openssl genrsa -out ${SSL_KEY_FILE} ${RSA_SIZE}
chmod 600 ${SSL_KEY_FILE}
openssl req -new -sha256 -key ${SSL_KEY_FILE} -subj "/CN=${COMMON_NAME}" | \
    openssl x509 -req -sha256 -CA ${CA_CRT_FILE} -CAkey ${CA_KEY_FILE} -CAcreateserial \
        -out ${SSL_CRT_FILE} -days 365 -extensions v3_req_server \
        -extfile ${SSL_CONFIG_FILE}
