#!/usr/bin/env bash

IS_PUSH=""

function build() {
    rm -rf _book

    if [[ ! -d "node_modules" ]];then
        gitbook install
    fi

    gitbook build

    if [[ $? != 0 ]] || [[ ! -d "_book" ]];then
        echo "[ERROR]编译 book 失败"
        exit 1
    fi
}

function deploy() {

    echo "[INFO]发布到 gh-pages 分支，发布后直接生效"
    if  read -p "确认要发布吗[y|n]，默认为 n: > " IS_PUSH;then
        if [[ -z ${IS_PUSH} ]];then
            IS_PUSH="n"
        fi
    fi

    if [[ ${IS_PUSH} != "y" ]];then
        return
    fi

    git pull
    rm -rf /tmp/_book
    mv _book /tmp/_book
    git checkout gh-pages
    cp -r -f /tmp/_book/* ./ && rm -rf /tmp/_book
    git add .
    git commit -a -m "update"
    git push origin gh-pages
    git checkout doc
}

function main() {
    build
    deploy
}

main