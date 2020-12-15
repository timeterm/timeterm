#!/usr/bin/env bash

set -ux

cd backend
gox 

cd ../nats-manager
gox

cd ../frontend-admin-web
yarn build
7z a build.7z build

cd ../docs
hugo
7z a build.7z public

cd ../os
toolbox run bash -c 'export MACHINE=raspberrypi4 && . ./setup-environment.sh && bitbake b2qt-embedded-qt5-image'
cd ..

7z a build.7z backend/backend_* nats-manager/nats-manager_* frontend-admin-web/build docs/public os/build-raspberrypi4/tmp/deploy/images/raspberrypi4/b2qt-embedded-qt5-image-raspberrypi4-*.7z

rm backend/backend_*
rm nats-manager/nats-manager_*

