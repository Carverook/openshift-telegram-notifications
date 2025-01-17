# OpenShift Telegram Notifications

A project to send OpenShift error messages to a telegram channel/group of your choice.

## Cluster Deployment

First create a [Telegram Bot](https://core.telegram.org/bots#6-botfather).

Then deploy this bot to OpenShift with permissions via:

```shell
$ oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:<current-project-here>:default --as=system:admin
$ oc new-app -f https://raw.githubusercontent.com/dinhnn/openshift-telegram-notifications/master/template.yaml \
             -p TELEGRAM_API_URL=https://api.telegram.org \
             -p TELEGRAM_BOT_TOKEN=<bot-token> \
             -p TELEGRAM_BOT_CHANNEL=<chat_id of channel or group> \
             -p OPENSHIFT_CONSOLE_URL=https://<openshift-host-here>:8443/console
$ oc start-build openshift-telegram-notifications
```

Once the app is built and deployed, it will start sending notifications to slack when there are `Warning` type events.

![Telegram Message](images/telegram-bot-message.png)

## Local Development

### Cluster Requirements

First you need a running minishift cluster. This can be installed via homebrew:

```shell
$ brew install socat openshift-cli docker-machine-driver-xhyve
$ brew tap caskroom/versions
$ brew cask install minishift-beta
```

The xhyve hypervisor requires superuser privileges. To enable, execute:

```shell
$ sudo chown root:wheel /usr/local/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve
$ sudo chmod u+s /usr/local/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve
```

Then start the cluster with:

```shell
$ minishift start --memory 4048
```

### App Requirements

First add the privileges to mount volumes and read cluster state to your service account:

```shell
$ oc login -u system:admin
$ oc adm policy add-scc-to-user hostmount-anyuid system:serviceaccount:myproject:default --as=system:admin
$ oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:myproject:default --as=system:admin
```
