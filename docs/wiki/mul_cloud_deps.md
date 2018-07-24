

Normally you can do a single wr cloud deploy --deployment production, and a single wr cloud deploy --deployment development. This should be fine if you're a normal single user running workflows for yourself. Other people can do their own deployments from the same machine and you won't have any conflicts.

If, however, you have full control of a machine and want to run multiple deployments yourself (eg. you're running wr for other people, but want to keep each of their workflows in separate cloud networks), you can do something like:

export MY_UNIQUE_DEPLOYMENT_NAME="one"
export WR_MANAGERPORT=11320
export WR_MANAGERWEB=11321
export WR_MANAGERDIR="~/.wr_$MY_UNIQUE_DEPLOYMENT_NAME"
wr cloud deploy --resource_name $MY_UNIQUE_DEPLOYMENT_NAME
[use wr commands as normal]
wr cloud teardown --resource_name $MY_UNIQUE_DEPLOYMENT_NAME

You will have to arrange that the value of --resource_name is unique within your cloud, and that WR_MANAGERPORT and WR_MANAGERWEB are unique (and not used by anyone else) on your machine, for each deployment you want to do.

