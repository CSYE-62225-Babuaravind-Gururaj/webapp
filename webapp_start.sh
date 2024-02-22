#!/bin/sh

sudo systemctl daemon-reload

sudo systemctl start application.service

sudo systemctl enable application.service