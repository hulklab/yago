#!/bin/sh

# 路径配置
DIR_CONF="conf"
APPS="app"
ROOT=`pwd`
USR=$1
CONF_POSTFIX="toml"

# 用法
if test $# -lt 1
then
    echo USAGE: ./env.init.sh who
    echo    eg: ./env.init.sh yago
    exit
fi

echo start to generate environment for: \"$USR\"

cd $DIR_CONF
for APP in $APPS
do
    # 复制模板
    if test ! -e $APP.$USR.${CONF_POSTFIX}
    then
        if test ! -e $APP.${CONF_POSTFIX}
        then
            echo $APP.${CONF_POSTFIX} not exist.
            exit
        fi

        cp $APP.${CONF_POSTFIX} $APP.$USR.${CONF_POSTFIX}
        echo create $APP.$USR.${CONF_POSTFIX} ...... ok
    fi
done

cd $ROOT

for APP in $APPS
do
    if [[ ${APP} != "app" ]];then
        continue
    fi

    # 软连接
    if ! test -e $APP.${CONF_POSTFIX}
    then
        ln -sfn  ${DIR_CONF}/${APP}.${USR}.${CONF_POSTFIX} ${APP}.${CONF_POSTFIX}
        echo link -s ${DIR_CONF}/$APP.$USR.${CONF_POSTFIX} ${APP}.${CONF_POSTFIX} ...... ok
    fi

done

echo "${USR} env init done!"
