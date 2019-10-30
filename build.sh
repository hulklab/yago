#!/usr/bin/env bash

IS_PUSH=""

function build() {
    rm -rf _book
    gitbook install
    gitbook build

    if [[ ! -d "_book" ]];then
        echo "[ERROR]编译 book 失败"
        exit 1
    fi
}

function deploy() {

    echo "[INFO]发布到 gh-page 分支，发布后直接生效"
    if  read -p "确认要发布吗[y|n]，默认为 n: > " IS_PUSH;then
        if [[ -z ${IS_PUSH} ]];then
            IS_PUSH="n"
        fi
    fi

    if [[ ${IS_PUSH} != "y" ]];then
        return
    fi

    rm -rf /tmp/_book
    cp -r _book /tmp/_book
    git co gh-pages
    git branch
}

function main() {
    build
    deploy
}

main