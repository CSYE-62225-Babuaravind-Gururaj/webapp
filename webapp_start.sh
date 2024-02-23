#!/bin/sh

sudo systemctl daemon-reload

sudo systemctl start webapp.service

sudo systemctl enable webapp.service