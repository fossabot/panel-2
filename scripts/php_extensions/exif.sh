#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

: '
Copyright 2022 HaoZi Technology Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
'

HR="+----------------------------------------------------"

action="$1"      # 操作
phpVersion="$2" # PHP版本

Install() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^extension=exif$')
    if [ "${isInstall}" != "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 已安装 exif"
        exit 1
    fi

    cd /www/server/php/${phpVersion}/src/ext/exif
    /www/server/php/${phpVersion}/bin/phpize
    ./configure --with-php-config=/www/server/php/${phpVersion}/bin/php-config
    make
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} exif 编译失败"
        exit 1
    fi
    make install
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} exif 安装失败"
        exit 1
    fi

    sed -i '/;haozi/a\extension=exif' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} exif 安装成功"
}

Uninstall() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^extension=exif$')
    if [ "${isInstall}" == "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 未安装 exif"
        exit 1
    fi

    sed -i '/extension=exif/d' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} exif 卸载成功"
}

if [ "$action" == 'install' ]; then
    Install
fi
if [ "$action" == 'uninstall' ]; then
    Uninstall
fi
