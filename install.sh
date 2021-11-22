#!/bin/sh
cd front  && npm install && npm run build
cd ..
go install

