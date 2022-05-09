#!/bin/sh
cd front  && pnpm install && pnpm run build
cd ..
go install

