#!/bin/bash
echo "🔧 TrimUI Smart Pro SDL2 environment setup (Simplified)"
echo "GOOS: $GOOS"
echo "GOARCH: $GOARCH"
echo "CC: $CC"
echo "CGO_ENABLED: $CGO_ENABLED"
echo "PKG_CONFIG_PATH: $PKG_CONFIG_PATH"
echo "CGO_CFLAGS: $CGO_CFLAGS"
echo "CGO_LDFLAGS: $CGO_LDFLAGS"
echo ""
echo "🔍 Checking cross-compilation toolchain"
which aarch64-linux-gnu-gcc
aarch64-linux-gnu-gcc --version | head -n1
echo ""
echo "🔍 Checking SDL2 setup"
ls -la /usr/include/aarch64-linux-gnu/SDL2/ 2>/dev/null && echo "✅ SDL2 headers found" || echo "❌ SDL2 headers not found"
pkg-config --exists sdl2 && echo "✅ SDL2 pkg-config OK" || echo "❌ SDL2 pkg-config failed"
pkg-config --cflags sdl2 2>/dev/null && echo "✅ SDL2 cflags OK" || echo "❌ SDL2 cflags failed"
pkg-config --libs sdl2 2>/dev/null && echo "✅ SDL2 libs OK" || echo "❌ SDL2 libs failed"
echo ""
echo "📦 Downloading Go dependencies"
go mod download
echo ""
echo "🏗️  Compiling for TrimUI Smart Pro (ARM64)"
mkdir -p bin
go build -v -ldflags="-s -w" -o bin/retroart-trimui-arm64 cmd/retroart/main.go
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ TrimUI compilation successful!"
    echo "📁 Binary created: bin/retroart-trimui-arm64"
    ls -la bin/retroart-trimui-arm64
    file bin/retroart-trimui-arm64
    echo ""
    echo "🎯 Binary specific to TrimUI Smart Pro ready!"
    echo "📊 File details:"
    echo "   Architecture: $(file bin/retroart-trimui-arm64 | grep -o 'ARM aarch64')"
    echo "   Size: $(ls -lh bin/retroart-trimui-arm64 | awk '{print $5}')"
else
    echo ""
    echo "❌ TrimUI compilation failed"
    echo "🔍 Debugging information:"
    echo ""
    echo "📋 Environment variables:"
    echo "   CC: $CC"
    echo "   CGO_CFLAGS: $CGO_CFLAGS"
    echo "   CGO_LDFLAGS: $CGO_LDFLAGS"
    echo "   PKG_CONFIG_PATH: $PKG_CONFIG_PATH"
    echo ""
    echo "📋 SDL2 pkg-config:"
    pkg-config --cflags sdl2 2>&1 || echo "   ❌ Error getting SDL2 cflags"
    pkg-config --libs sdl2 2>&1 || echo "   ❌ Error getting SDL2 libs"
    echo ""
    echo "📋 Headers check:"
    ls -la /usr/include/aarch64-linux-gnu/SDL2/ 2>/dev/null || echo "   ❌ No SDL2 headers found"
    echo ""
    echo "📋 Libraries check:"
    ls -la /usr/lib/aarch64-linux-gnu/ | grep -i sdl 2>/dev/null || echo "   ❌ No SDL2 libs found"
    echo ""
    echo "📋 Toolchain check:"
    which aarch64-linux-gnu-gcc
    aarch64-linux-gnu-gcc --version 2>/dev/null || echo "   ❌ Toolchain not working"
    exit 1
fi
