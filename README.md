# Order Book Recorder

Records Order Book snapshots tick data into CSV files.

Currently works for USDT pairs on Binance.

## Settings

Open settings.json and fill in desired symbols.
Please list only base pairs. Do not write "USDT", program adds it automatically.

(i) Api Keys are not required for current functionality, made it just for the future.
To fill in Api Key: Follow .env_example and create .env file in the same folder.

## Output file name example

**2022_11_30\_\_BTCUSDT.csv**

```
[timestamp, side, price, quantity],
[timestamp, side, price, quantity],
[...] 
```

</br>
