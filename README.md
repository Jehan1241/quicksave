# Quicksave
An upcoming game library manager. Inspired by and seeking to improve upon [Playnite](https://playnite.link/).

## Important
This is an early WIP project. While third party integration is complete, due to the non-standard and documented nature of external APIs certain games might cause unexpected behaviour.

## API Keys
Your API Keys are never logged, stored or exposed.\
To use certain features of the app, you currently need API keys, this will not be the case in the future.
* [IGDB API key](https://api-docs.igdb.com/#getting-started) : Create an IGDB account and obtain your API key and a Secret key.
* [Steam API key](https://steamcommunity.com/discussions/forum/1/3047235828269633221/) : To transfer you steam library, your account ID and API key is required.
* PlayStation Support : Visit the [PlayStation Homepage](https://www.playstation.com/) and login. Then visit [This Page](https://ca.account.sony.com/api/v1/ssocookie) and obtain your NPSSO code.

## Build Instructions
* Frontend
```
git clone https://github.com/Jehan1241/QuickSave
cd QuickSave/Electron/electron-leprechaun
npm i
npm run dev
```
* Backend
```
cd backend
go make tidy
go run .
```


