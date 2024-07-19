# Git Authentication Tokens

## BitBucket

1. Login into your Bickbucket account.

2. Go to **Personal Settings** (Click the cog icon at the upper right corner).

3. Go to **ACCESS MANAGEMENT > App Password** on the left side menu.

4. Click on **Create app password**.

5. Add a label for the new password, and set the following permissions:

| Permission | Type |
| ---------- | ---- |
| Repositories | [x] Read |
| Repositories | [x] Write |
| Pull requests | [x] Read |
| Pull requests | [x] Write |

6. Click **Create**, and make sure to save the generated password on a safe place.

Run `gmmit -pr`, and provide your username and password when asked.

Check the official [BitBucket Docs](https://support.atlassian.com/bitbucket-cloud/docs/create-an-app-password/) for more details.

## Github

1. Login into to your Github account.

2. Click on you profile pic (upper right corner) and go to **Settings**.

3. On the left side menu go to **Developer Settings** (at the very bottom).

4. Go to **Personal Access Tokens > Tokens (classic)**

5. Click on **Generate new token**

6. Fill the form by adding a Note/Name to your token, selecting the experiation time, and make sure to select the following permissions:

| Permission | Description |
| ---------- | ----------- |
| [x] repo | Full control of private repositories |
| [x] repo:status | Access commit status |
| [x] repo_deployment | Access deployment status |
| [x] public_repo | Access public repositories |
| [x] repo:invite | Access repository invitations |
| [x] security_events | Read and write security events |

7. Click **Generate token**, and make sure to save the generated password on a safe place.

Run `gmmit -pr`, and provide your username and password when asked.

Check the official [Github Docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) for more details.
