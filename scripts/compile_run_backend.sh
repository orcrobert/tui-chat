#!/bin/bash

build_path="../build"

if [ -d "$build_path" ]; then
    echo "Dir exists, deleating it's contents..."
    rm -rf "$build_path"/*
else 
    echo "Dir does not exist, creating it..."
    mkdir -p "$build_path"
fi

echo "Compiling backend source files..."

clang++ -std=c++17 -o ../build/validate ../backend/validate.cpp -lsqlite3
clang++ -std=c++17 -o ../build/init_db ../backend/init_db.cpp -lsqlite3

echo "Running backend init executables..."
cd ../build/
./init_db