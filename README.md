# Order Book Recorder

Records Order Book snapshots tick data into CSV files.

Currently works for USDT pairs on Binance.

## Settings

Open settings.json and fill in desired symbols.
Do not write "USDT", program does that automatically. Please list only base pairs.

(i) Api Keys are not required for current functionality, made it just for the future.
To fill in Api Key: Open the .env file in the same folder as your executable is in.

## Output file example

**2022_11_30\_\_BTCUSDT.csv**

```
[timestamp, side, price, quantity],
[timestamp, side, price, quantity],
[...]
```

</br>

## Pipeline:

```
(once) ObResponse -> SymbolObSnapshot ->
                                        -> AllSymbolsObSnapshots -> csv
(tick)                   ObWsResponse ->
```

-   **ObResponse**
    </br>Response from http-request to Binance
    </br>

-   **SymbolObSnapshot** {
    </br> &nbsp;&nbsp;&nbsp;&nbsp;Symbol string
    </br> &nbsp;&nbsp;&nbsp;&nbsp;Date int (UnixMilli)
    </br> &nbsp;&nbsp;&nbsp;&nbsp;Data []strings: [timestamp, side, price, qty]
    </br>}
    </br>Formatted ObResponse, prepared for recording to csv
    </br>

-   **AllSymbolsObSnapshots** {
    </br>&nbsp;&nbsp;&nbsp;&nbsp;[]SymbolObSnapshot
    </br>}
    </br>Slice of SymbolObSnapshot

</br>

## Pseudo-code

```
generate empty AllSymbolsObResponse from listOfSymbols

for symbol in symbolsList {
    goroutines A:
    (once app starts)
        - get OrderbookTemplate from http api call

    goroutines B:
    (once go A finished)
        (at each tick, ws)
        - get ObWsResponse
            - parse the message
        - get previous SymbolObSnapshot
        - for each price level in ObWsResponse:
            if quantity != 0 {
                update SymbolObSnapshot:
                    the price level with new values
            } else {
                update SymbolObSnapshot:
                    remove price level
            }
        - WriteAll as new lines to .csv file
}
```
