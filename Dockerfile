# Dockerfile para cross-compilation Go ARM64 com SDL2 TrimUI Smart Pro
FROM golang:1.23-bullseye

# Instalar dependências básicas para cross-compilation ARM64
RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    pkg-config \
    wget \
    tar \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Baixar e instalar SDL2 específico do TrimUI Smart Pro
RUN mkdir -p /opt/trimui-sdl2 && \
    cd /opt/trimui-sdl2 && \
    wget https://github.com/trimui/toolchain_sdk_smartpro/releases/download/20231018/SDL2-2.26.1.GE8300.tgz && \
    tar -xzf SDL2-2.26.1.GE8300.tgz && \
    rm SDL2-2.26.1.GE8300.tgz

# Instalar SDL2 TrimUI para ARM64
RUN cd /opt/trimui-sdl2 && \
    if [ -d "SDL2-2.26.1" ]; then \
        cd SDL2-2.26.1; \
    elif [ -d "SDL2" ]; then \
        cd SDL2; \
    else \
        cd $(ls -d */ | head -n1); \
    fi && \
    # Criar diretório e copiar headers TrimUI específicos
    mkdir -p /usr/aarch64-linux-gnu/include/SDL2 && \
    # Primeiro, copiar headers do sistema SDL2 ARM64
    cp -r /usr/include/aarch64-linux-gnu/SDL2/* /usr/aarch64-linux-gnu/include/SDL2/ 2>/dev/null || \
    cp -r /usr/include/SDL2/* /usr/aarch64-linux-gnu/include/SDL2/ 2>/dev/null || true && \
    # Depois, sobrescrever com headers TrimUI específicos se existirem
    cp -r include/* /usr/aarch64-linux-gnu/include/SDL2/ 2>/dev/null || \
    find . -name "*.h" -exec cp {} /usr/aarch64-linux-gnu/include/SDL2/ \; 2>/dev/null || true

# Configurar variáveis de ambiente para cross-compilation com SDL2 TrimUI
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64
ENV CC=aarch64-linux-gnu-gcc
ENV CXX=aarch64-linux-gnu-g++
ENV PKG_CONFIG_PATH=/usr/aarch64-linux-gnu/lib/pkgconfig
ENV CGO_CFLAGS="-I/usr/aarch64-linux-gnu/include/SDL2"
ENV CGO_LDFLAGS="-L/usr/aarch64-linux-gnu/lib"

# Criar arquivo de configuração pkg-config para SDL2 TrimUI (usando bibliotecas do sistema)
RUN mkdir -p /usr/aarch64-linux-gnu/lib/pkgconfig && \
    echo 'prefix=/usr/aarch64-linux-gnu' > /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'exec_prefix=${prefix}' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'libdir=${exec_prefix}/lib' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'includedir=${prefix}/include' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo '' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Name: sdl2' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Description: Simple DirectMedia Layer (TrimUI Smart Pro version)' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Version: 2.26.1' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Requires:' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Conflicts:' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Libs: -L${libdir} -lSDL2' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Libs.private: -lm -ldl -lpthread -lrt' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc && \
    echo 'Cflags: -I${includedir}/SDL2' >> /usr/aarch64-linux-gnu/lib/pkgconfig/sdl2.pc

# Instalar SDL2 para ARM64 (para bibliotecas de linkagem)
RUN apt-get update && apt-get install -y \
    libsdl2-dev:arm64 \
    libsdl2-gfx-dev:arm64 \
    libsdl2-image-dev:arm64 \
    libsdl2-mixer-dev:arm64 \
    libsdl2-net-dev:arm64 \
    libsdl2-ttf-dev:arm64 \
    libfreetype6-dev:arm64 \
    && rm -rf /var/lib/apt/lists/*

# Criar diretório de trabalho
WORKDIR /workspace

# Script de build personalizado para SDL2 TrimUI
RUN echo '#!/bin/bash\n\
echo "🔧 Configurando ambiente TrimUI Smart Pro SDL2..."\n\
echo "GOOS: $GOOS"\n\
echo "GOARCH: $GOARCH"\n\
echo "CC: $CC"\n\
echo "CGO_ENABLED: $CGO_ENABLED"\n\
echo "PKG_CONFIG_PATH: $PKG_CONFIG_PATH"\n\
echo "CGO_CFLAGS: $CGO_CFLAGS"\n\
echo "CGO_LDFLAGS: $CGO_LDFLAGS"\n\
echo ""\n\
echo "� Verificando SDL2 TrimUI..."\n\
ls -la /opt/trimui-sdl2/\n\
ls -la /usr/aarch64-linux-gnu/include/SDL2/ 2>/dev/null || echo "Headers SDL2 não encontrados"\n\
pkg-config --exists sdl2 && echo "✅ SDL2 pkg-config OK" || echo "❌ SDL2 pkg-config falhou"\n\
echo ""\n\
echo "�📦 Baixando dependências Go..."\n\
go mod download\n\
echo ""\n\
echo "🏗️  Compilando para TrimUI Smart Pro (ARM64)..."\n\
mkdir -p bin\n\
go build -v -ldflags="-s -w" -o bin/retroart-trimui-arm64 cmd/retroart/main.go\n\
if [ $? -eq 0 ]; then\n\
    echo ""\n\
    echo "✅ Compilação TrimUI bem-sucedida!"\n\
    echo "📁 Binário criado: bin/retroart-trimui-arm64"\n\
    ls -la bin/retroart-trimui-arm64\n\
    file bin/retroart-trimui-arm64\n\
    echo ""\n\
    echo "🎯 Binário específico para TrimUI Smart Pro pronto!"\n\
    echo "   Para transferir: docker cp container:/workspace/bin/retroart-trimui-arm64 ./bin/"\n\
else\n\
    echo ""\n\
    echo "❌ Falha na compilação TrimUI"\n\
    echo "Debugando..."\n\
    echo "SDL2 pkg-config:"\n\
    pkg-config --cflags sdl2 || echo "Erro no pkg-config SDL2"\n\
    pkg-config --libs sdl2 || echo "Erro no pkg-config SDL2 libs"\n\
    echo "Headers SDL2:"\n\
    ls -la /usr/aarch64-linux-gnu/include/SDL2/\n\
    echo "Bibliotecas ARM64:"\n\
    ls -la /usr/aarch64-linux-gnu/lib/ | grep -i sdl || echo "Nenhuma lib SDL2 encontrada"\n\
    echo "Bibliotecas do sistema:"\n\
    find /usr -name "*SDL2*" 2>/dev/null | head -10\n\
    exit 1\n\
fi' > /usr/local/bin/build-arm64.sh && \
    chmod +x /usr/local/bin/build-arm64.sh

# Comando padrão
CMD ["/usr/local/bin/build-arm64.sh"]
