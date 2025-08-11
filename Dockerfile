# Dockerfile para cross-compilation Go ARM64 com SDL2 TrimUI Smart Pro
FROM golang:1.24-bullseye

ENV DEBIAN_FRONTEND=noninteractive

# Instalar dependências básicas para cross-compilation ARM64
RUN dpkg --add-architecture arm64 && \
    apt-get update && \
    apt-get install -y \
    wget \
    build-essential \
    ca-certificates \
    pkg-config \
    tar \
    linux-libc-dev-arm64-cross \
    libc6-arm64-cross \
    libc6-dev-arm64-cross \
    binutils-aarch64-linux-gnu \
    gcc-aarch64-linux-gnu \
    g++-aarch64-linux-gnu \
    libsdl2-dev:arm64 \
    libsdl2-gfx-dev:arm64 \
    libsdl2-image-dev:arm64 \
    libsdl2-ttf-dev:arm64 \
    libfreetype6-dev:arm64 \
    && rm -rf /var/lib/apt/lists/*

# Baixar toolchain TrimUI como referência (opcional)
RUN wget https://github.com/trimui/toolchain_sdk_smartpro/releases/download/20231018/aarch64-linux-gnu-7.5.0-linaro.tgz && \
    tar -C /usr/local -xzf aarch64-linux-gnu-7.5.0-linaro.tgz && \
    rm aarch64-linux-gnu-7.5.0-linaro.tgz

# Baixar SDK adicional do TrimUI (apenas para referência)
RUN wget https://github.com/trimui/toolchain_sdk_smartpro/releases/download/20231018/SDK_usr_tg5040_a133p.tgz && \
    tar -xzf SDK_usr_tg5040_a133p.tgz -C /tmp && \
    rm -rf SDK_usr_tg5040_a133p.tgz /tmp/usr

# Configurar variáveis de ambiente simplificadas para cross-compilation
ENV PKG_CONFIG_PATH="/usr/lib/aarch64-linux-gnu/pkgconfig"
ENV CC="aarch64-linux-gnu-gcc"
ENV CXX="aarch64-linux-gnu-g++"
ENV AR="aarch64-linux-gnu-ar"
ENV STRIP="aarch64-linux-gnu-strip"
ENV GOOS="linux"
ENV GOARCH="arm64"
ENV CGO_ENABLED="1"
ENV CGO_LDFLAGS="-L/usr/lib/aarch64-linux-gnu -lSDL2_image -lSDL2_ttf -lSDL2 -ldl -lpthread -lm"
ENV CGO_CFLAGS="-I/usr/include/aarch64-linux-gnu -I/usr/include/aarch64-linux-gnu/SDL2 -D_REENTRANT"

# Criar diretório de trabalho
WORKDIR /workspace

COPY build-arm64.sh /usr/local/bin/build-arm64.sh
RUN chmod +x /usr/local/bin/build-arm64.sh

# Comando padrão
CMD ["/usr/local/bin/build-arm64.sh"]
