# Lab 6.1 — Static site + rsync deploy

## Что такое rsync
`rsync` — утилита для синхронизации файлов по локальным и удалённым путям. Она копирует только изменённые блоки файлов (алгоритм дельт-передачи), может сжимать трафик, сохранять права и владельцев, работать по SSH. Благодаря этому обновление статики на сервере выполняется быстро и повторяемо одной командой.

## Структура проекта
- `site/` — статический сайт, который будет отдавать nginx.
- `deploy.sh` — сценарий деплоя на удалённый сервер через `rsync` + SSH.
- `infra/nginx.static-site.conf` — пример server-блока для nginx.
- `Практическая работа №6.1.docx` — задание из методички.

## Быстрый старт (локально)
```bash
python -m http.server 8080 --directory site
# открыть http://localhost:8080
```

## Подготовка сервера
1. Подключитесь по SSH к серверу (тот же, что в лабе 5.2).
2. Установите nginx: `sudo apt update && sudo apt install nginx -y` (для Debian/Ubuntu).
3. Создайте директорию для статики и выдайте права пользователю деплоя:
   ```bash
   sudo mkdir -p /var/www/static-site
   sudo chown -R $USER:$USER /var/www/static-site
   ```
4. Скопируйте пример конфига `infra/nginx.static-site.conf` на сервер как `/etc/nginx/sites-available/static-site` и поправьте `server_name`/`root` при необходимости.
5. Активируйте сайт и перезапустите nginx:
   ```bash
   sudo ln -s /etc/nginx/sites-available/static-site /etc/nginx/sites-enabled/static-site
   sudo nginx -t
   sudo systemctl reload nginx
   ```

## Деплой статики через rsync
1. Клонируйте репозиторий локально, перейдите в каталог `lab5`.
2. Убедитесь, что на сервере разрешён ваш SSH-ключ (`~/.ssh/id_rsa.pub`).
3. Запустите деплой (пример):
   ```bash
   REMOTE_HOST=203.0.113.10 \
   REMOTE_USER=deploy \
   REMOTE_PATH=/var/www/static-site \
   ./deploy.sh
   ```
   Переменные:
   - `REMOTE_HOST` (обязательно) — IP или домен сервера.
   - `REMOTE_USER` — пользователь SSH (по умолчанию текущий локальный).
   - `REMOTE_PATH` — путь, куда сложить файлы (по умолчанию `/var/www/static-site`).
   - `SSH_KEY` — путь до приватного ключа, если не используется дефолтный.
   - `EXTRA_RSYNC_OPTS` — дополнительные опции (например `--dry-run`).
   - `RELOAD_NGINX` — `1` (по умолчанию) перезагрузить nginx через `sudo systemctl reload nginx`, `0` — пропустить.
4. Откройте сайт: `http://REMOTE_HOST` или ваш домен. Если нужно отдать сайт во внешнюю сеть с домашнего сервера — поднимите туннель (например, cloudpub.ru) и прикрепите домен через freedns.afraid.org.

## Что происходит в deploy.sh
- Проверяет обязательные переменные и наличие папки `site/`.
- Готовит SSH-команду (с ключом, если указан `SSH_KEY`).
- Создаёт директорию на сервере.
- Вызывает `rsync -avz --delete --checksum` для синхронизации `site/` → `REMOTE_PATH/`.
- При `RELOAD_NGINX=1` пытается выполнить `sudo systemctl reload nginx` на сервере.

## Скриншоты/отчёт
Для отчёта по практике сделайте скриншоты ключевых шагов: установка nginx, проверка `nginx -t`, запуск `./deploy.sh` с выводом, содержимое сайта в браузере.

# Лаба 6.2 — Ansible playbook (управление конфигурацией)

## Структура
- `ansible/ansible.cfg` — базовые настройки (инвентарь, роли, отключение проверки ключей).
- `ansible/inventory.ini` — пример инвентаря (замените `ansible_host`, `ansible_user`, путь к ключу).
- `ansible/setup.yml` — основной playbook с ролями `base`, `nginx`, `app`, `ssh` и тегами.
- `ansible/roles/base` — apt-update/upgrade, установка утилит + fail2ban.
- `ansible/roles/nginx` — установка nginx, выкладка конфига, отключение default-сайта.
- `ansible/roles/app` — загрузка архива статики (`static-site.tar.gz`) в `app_root` (`/var/www/static-site`).
- `ansible/roles/ssh` — добавление открытого ключа в `~/.ssh/authorized_keys` выбранного пользователя.
- `ansible/roles/app/files/static-site.tar.gz` — архив статики (тот же сайт из 6.1).
- `ansible/roles/ssh/files/authorized_key.pub` — заглушка ключа; замените на свой.

## Быстрый запуск
1. Установите Ansible локально (`sudo apt install ansible` или через `pipx install ansible`).
2. В `ansible/inventory.ini` подставьте свой IP/пользователя/ключ. Пример оставлен для `31.192.110.47`.
3. При необходимости замените `ansible/roles/ssh/files/authorized_key.pub` на ваш реальный публичный ключ.
4. Запуск всех ролей:
   ```bash
   cd ansible
   ansible-playbook setup.yml
   ```
   Запуск отдельных ролей:
   ```bash
   ansible-playbook setup.yml --tags "app"
   ansible-playbook setup.yml --tags "nginx,app"
   ```

## Переменные, которые можно переопределить (`--extra-vars`)
- `nginx_server_name` — по умолчанию `_`.
- `nginx_web_root` / `app_root` — по умолчанию `/var/www/static-site`.
- `app_archive` — путь к архиву со статикой (по умолчанию архив в роли).
- `app_owner` / `app_group` — по умолчанию `www-data`.
- `ssh_user` — по умолчанию `ansible_user`/`root`.
- `ssh_public_key_src` — путь к файлу с публичным ключом.

## Проверка
- После `--tags "nginx,app"`: `curl -I http://<host>` должен отдавать `200 OK`.
- Логи nginx: `/var/log/nginx/static-site.access.log`.
- SSH-ключ: на сервере в `~user/.ssh/authorized_keys` должен появиться ваш ключ.
