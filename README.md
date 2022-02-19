# Binance Repo Function Purpose

This Repo will take duration of ETHUSDT and BTCUSDT as offense, ETHUSDC  and BTCUSDC as defense.
Duration will be set during runtime with m as minute, h as hour, d as day.
Data Taken starts at 00:00 and will take interval of set Duration.
The Return will be a spreadsheet with from 2018 to 2021

# TODO
Add Taker Fee of 0.07% for every transaction IN or OUT
Scan input for start year & interval string

# How to run, execute
make sure you have go version 1.13.8 and above
execute file ./binance