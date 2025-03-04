init `.env` file
```shell
cp .env.example .env
# modify the .env file as you wanted
```

for backend(api)
```shell
go run .
```
if the backend run successfully, then you can see the following:
![backend-api.png](./assets/backend-api.png)

for frontend(ui)
```shell
cd ui
pnpm install

ln -s ../.env .env
pnpm run dev
```
if the frontend run successfully, then you can see the following:
![frontend-ui.png](./assets/frontend-ui.png)

then open the browser, visit `http://localhost:5173/`, the following is the gif
![web-terminal.gif](./assets/web-terminal.gif)
